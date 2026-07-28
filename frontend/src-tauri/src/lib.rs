use tauri::Manager;

#[tauri::command]
fn close_splash(app_handle: tauri::AppHandle) {
  if let Some(splash_window) = app_handle.get_webview_window("splash") {
    let _ = splash_window.close();
}

  if let Some(main_window) = app_handle.get_webview_window("main") {
    let _ = main_window.maximize();
    let _ = main_window.show();
    let _ = main_window.set_focus();
}
}

#[tauri::command]
fn get_internal_token() -> String {
  "devaulty-native-session-token".to_string()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  if let Some(config_dir) = dirs::config_dir() {
    let devaulty_dir = config_dir.join("devaulty");
    if !devaulty_dir.exists() {
      let _ = std::fs::create_dir_all(devaulty_dir);
    }
  }

  tauri::Builder::default()
    .invoke_handler(tauri::generate_handler![close_splash, get_internal_token])
    .setup(|app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }
      Ok(())
    })
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}
