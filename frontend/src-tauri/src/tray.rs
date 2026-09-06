use std::sync::atomic::Ordering;
use std::sync::Arc;
use tauri::{
  image::Image,
  menu::{Menu, MenuItem},
  tray::TrayIconBuilder,
  AppHandle, Manager, WebviewWindowBuilder,
};

use crate::session::SessionState;

// Builds and registers the system tray icon with a minimal "Open" / "Quit" menu.
// Closing the main window hides it (and, after a grace period, destroys its
// WebView) instead of exiting the process (wired in lib.rs's on_window_event
// + run() exit-request handler), so the backend and its resolved resource
// paths (including the AppImage FUSE mount, when applicable) stay alive
// until the user explicitly quits from here.
pub fn setup_tray(app: &AppHandle) -> tauri::Result<()> {
  let open_item = MenuItem::with_id(app, "open", "Open Devaulty", true, None::<&str>)?;
  let quit_item = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
  let menu = Menu::with_items(app, &[&open_item, &quit_item])?;

  let icon_bytes = include_bytes!("../icons/64x64.png");
  let icon = Image::from_bytes(icon_bytes)?;

  TrayIconBuilder::new()
    .icon(icon)
    .menu(&menu)
    .on_menu_event(|app, event| match event.id.as_ref() {
      "open" => show_main_window(app),
      "quit" => quit_app(app),
      _ => {}
    })
    .build(app)?;

  Ok(())
}

// Rebuilds a window exactly as declared in tauri.conf.json, by label. This
// is the same construction path Tauri itself uses at boot, so it never
// drifts from the config file and never needs manual attribute-by-attribute
// duplication (or the "unstable" cargo feature that comes with calling
// individual builder methods like .transparent()).
fn rebuild_from_config(app: &AppHandle, label: &str) -> tauri::Result<()> {
  let config = app
    .config()
    .app
    .windows
    .iter()
    .find(|w| w.label == label)
    .unwrap_or_else(|| panic!("window '{}' missing from tauri.conf.json", label))
    .clone();

  WebviewWindowBuilder::from_config(app, &config)?.build()?;
  Ok(())
}

// Shows the main window, recreating it from scratch if it was previously
// destroyed (grace period elapsed after a tray minimize) or doesn't exist
// yet (second-instance relaunch). The Go backend is already running at this
// point (port/token resolved on first launch), so get_backend_info on the
// frontend resolves almost immediately here.
pub fn show_main_window(app: &AppHandle) {
  let state = app.state::<Arc<SessionState>>();

  // Hold the same lock the destroy timer uses. Bumping the epoch here is
  // what invalidates a pending timer, but only holding the lock for the
  // *entire* reopen/recreate sequence (not just the bump) guarantees no
  // destroy() can run concurrently and undo this reopen.
  let mut epoch = state.hide_epoch.lock().unwrap();
  *epoch += 1;

  if let Some(window) = app.get_webview_window("main") {
    let _ = window.show();
    let _ = window.set_focus();
    return;
  }

  // Cold path: the window was actually destroyed. Recreate the splash
  // window too (it was closed/destroyed by commands::close_splash right
  // after the last boot) and rebuild "main" from its tauri.conf.json entry,
  // which already declares visible: false - the frontend's normal boot
  // sequence (get_backend_info, then invoking close_splash once ready)
  // shows "main" and tears the splash back down, so reopening after a long
  // idle period looks the same as first launch instead of a blank window
  // mid-render.
  if app.get_webview_window("splash").is_none() {
    if let Err(e) = rebuild_from_config(app, "splash") {
      log::error!("Failed to recreate splash window: {}", e);
    }
  }

  if let Err(e) = rebuild_from_config(app, "main") {
    log::error!("Failed to recreate main window: {}", e);
  }
}

fn quit_app(app: &AppHandle) {
  let state = app.state::<Arc<SessionState>>();
  // Mark this exit as intentional before requesting it, so lib.rs's
  // ExitRequested handler lets it through instead of preventing it.
  state.quitting.store(true, Ordering::SeqCst);
  if let Some(mut child) = state.child_process.lock().unwrap().take() {
    let _ = child.kill();
  }
  app.exit(0);
}