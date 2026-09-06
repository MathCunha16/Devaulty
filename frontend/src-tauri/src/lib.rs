mod backend;
mod commands;
mod session;
mod tray;

use std::sync::Arc;

use session::SessionState;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  // Ensure the user config/data directory exists for database storage
  let devaulty_data_dir = dirs::config_dir()
    .or_else(dirs::data_local_dir)
    .or_else(dirs::home_dir)
    .map(|path| path.join("devaulty"))
    .unwrap_or_else(|| std::env::temp_dir().join("devaulty"));

  if !devaulty_data_dir.exists() {
    if let Err(e) = std::fs::create_dir_all(&devaulty_data_dir) {
      log::error!("Failed to create devaulty data directory {:?}: {}", devaulty_data_dir, e);
      eprintln!("Failed to create devaulty data directory {:?}: {}", devaulty_data_dir, e);
    }
  }

  let session_state = Arc::new(SessionState::default());
  let state_clone = Arc::clone(&session_state);

  tauri::Builder::default()
    // Must be registered before other plugins: if Devaulty is launched again
    // while an instance is already running (double-click, launcher, CLI),
    // this intercepts the second launch and just brings the existing window
    // forward instead of spawning a whole new process + Go backend.
    .plugin(tauri_plugin_single_instance::init(|app, _args, _cwd| {
      tray::show_main_window(app);
    }))
    .plugin(tauri_plugin_updater::Builder::new().build())
    .plugin(tauri_plugin_process::init())
    .manage(session_state)
    .invoke_handler(tauri::generate_handler![
      commands::close_splash,
      commands::get_backend_info,
      commands::get_app_environment,
      commands::download_release_file,
      commands::cancel_download_release_file,
      commands::open_file_path
    ])
    .setup(move |app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      #[cfg(target_os = "macos")]
      app.handle().set_activation_policy(tauri::ActivationPolicy::Accessory);

      backend::spawn_backend(app.handle(), &devaulty_data_dir, &state_clone);
      tray::setup_tray(app.handle())?;

      Ok(())
    })
    .on_window_event(|window, event| {
      if let tauri::WindowEvent::CloseRequested { api, .. } = event {
        if window.label() == "main" {
          // Real close, not a hide: destroys the WebView (WebKitGTK /
          // WebView2 / WKWebView) so it stops holding onto RAM while parked
          // in the tray. Rebuilt from scratch in tray::show_main_window()
          // when the user reopens it.
          api.prevent_close();
          let _ = window.destroy();
        }
      }
    })
    .build(tauri::generate_context!())
    .expect("error while building tauri application")
    .run(|_app_handle, event| {
      // Destroying the "main" window above would otherwise make Tauri quit
      // the whole process once no windows are left. We only want to exit
      // when the user explicitly picks "Quit" from the tray (which calls
      // app.exit(0) itself in tray::quit_app), so swallow the automatic
      // exit-on-last-window-closed request here.
      if let tauri::RunEvent::ExitRequested { api, .. } = event {
        api.prevent_exit();
      }
    });
}