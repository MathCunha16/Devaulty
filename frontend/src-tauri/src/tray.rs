use std::sync::Arc;
use tauri::{
  image::Image,
  menu::{Menu, MenuItem},
  tray::TrayIconBuilder,
  AppHandle, Manager,
};

use crate::session::SessionState;

// Builds and registers the system tray icon with a minimal "Open" / "Quit" menu.
// Closing the main window hides it instead of exiting the process (wired in
// lib.rs's on_window_event), so the backend and its resolved resource paths
// (including the AppImage FUSE mount, when applicable) stay alive until the
// user explicitly quits from here.
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

fn show_main_window(app: &AppHandle) {
  if let Some(window) = app.get_webview_window("main") {
    let _ = window.show();
    let _ = window.set_focus();
  }
}

fn quit_app(app: &AppHandle) {
  let state = app.state::<Arc<SessionState>>();
  if let Some(mut child) = state.child_process.lock().unwrap().take() {
    let _ = child.kill();
  }
  app.exit(0);
}