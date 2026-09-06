use std::sync::Arc;
use tauri::{
  image::Image,
  menu::{Menu, MenuItem},
  tray::TrayIconBuilder,
  AppHandle, Manager, WebviewUrl, WebviewWindowBuilder,
};

use crate::session::SessionState;

// Builds and registers the system tray icon with a minimal "Open" / "Quit" menu.
// Closing the main window destroys its WebView instead of exiting the process
// (wired in lib.rs's on_window_event + run() exit-request handler), so the
// backend and its resolved resource paths (including the AppImage FUSE mount,
// when applicable) stay alive until the user explicitly quits from here.
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
// destroyed (tray minimize) or doesn't exist yet (second-instance relaunch).
// Mirrors the "main" window entry in tauri.conf.json so the rebuilt window
// is indistinguishable from the one created at boot. The Go backend is
// already running at this point (port/token resolved on first launch), so
// get_backend_info on the frontend resolves almost immediately here - no
// need to route back through the splash window for a reopen.
pub fn show_main_window(app: &AppHandle) {
  if let Some(window) = app.get_webview_window("main") {
    let _ = window.show();
    let _ = window.set_focus();
    return;
  }

  match WebviewWindowBuilder::new(app, "main", WebviewUrl::App("index.html".into()))
    .title("Devaulty")
    .min_inner_size(1000.0, 650.0)
    .resizable(true)
    .fullscreen(false)
    .visible(false)
    .build()
  {
    Ok(window) => {
      let _ = window.maximize();
      let _ = window.show();
      let _ = window.set_focus();
    }
    Err(e) => log::error!("Failed to recreate main window: {}", e),
  }
}

fn quit_app(app: &AppHandle) {
  let state = app.state::<Arc<SessionState>>();
  if let Some(mut child) = state.child_process.lock().unwrap().take() {
    let _ = child.kill();
  }
  app.exit(0);
}