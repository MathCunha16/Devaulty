use std::sync::atomic::Ordering;
use std::sync::Arc;
use tauri::{
  image::Image,
  menu::{Menu, MenuItem},
  tray::TrayIconBuilder,
  AppHandle, Manager, WebviewUrl, WebviewWindowBuilder,
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

// Shows the main window, recreating it from scratch if it was previously
// destroyed (grace period elapsed after a tray minimize) or doesn't exist
// yet (second-instance relaunch). Mirrors the "main" window entry in
// tauri.conf.json so the rebuilt window is indistinguishable from the one
// created at boot. The Go backend is already running at this point
// (port/token resolved on first launch), so get_backend_info on the
// frontend resolves almost immediately here.
pub fn show_main_window(app: &AppHandle) {
  if let Some(window) = app.get_webview_window("main") {
    // Fast path: still within the grace period, window was only hidden.
    // Bump the epoch so lib.rs's pending destroy timer (if any) sees it no
    // longer matches and skips tearing the WebView down.
    app
      .state::<Arc<SessionState>>()
      .hide_epoch
      .fetch_add(1, Ordering::SeqCst);
    let _ = window.show();
    let _ = window.set_focus();
    return;
  }

  // Cold path: the window was actually destroyed. Recreate the splash
  // window too (it was closed/destroyed by commands::close_splash right
  // after the last boot) and rebuild "main" hidden, exactly like app
  // startup. The frontend's normal boot sequence - get_backend_info, then
  // invoking close_splash once it's ready - takes care of showing "main"
  // and tearing the splash back down, so reopening after a long idle period
  // looks the same as first launch instead of a blank window mid-render.
  if app.get_webview_window("splash").is_none() {
    if let Err(e) = recreate_splash_window(app) {
      log::error!("Failed to recreate splash window: {}", e);
    }
  }

  if let Err(e) = WebviewWindowBuilder::new(app, "main", WebviewUrl::App("index.html".into()))
    .title("Devaulty")
    .min_inner_size(1000.0, 650.0)
    .resizable(true)
    .fullscreen(false)
    .visible(false)
    .build()
  {
    log::error!("Failed to recreate main window: {}", e);
  }
}

// Rebuilds the "splash" window with the same attributes it has in
// tauri.conf.json (used only at first boot otherwise).
fn recreate_splash_window(app: &AppHandle) -> tauri::Result<()> {
  WebviewWindowBuilder::new(app, "splash", WebviewUrl::App("/splash.html".into()))
    .title("Devaulty")
    .inner_size(550.0, 400.0)
    .center()
    .decorations(false)
    .transparent(true)
    .always_on_top(true)
    .build()?;
  Ok(())
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