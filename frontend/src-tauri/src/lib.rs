mod backend;
mod commands;
mod session;
mod tray;

use std::sync::atomic::Ordering;
use std::sync::Arc;
use std::time::Duration;

use tauri::Manager;

use session::SessionState;

// How long the main window is kept alive (hidden, not destroyed) after being
// sent to the tray before its WebView is actually torn down to reclaim RAM.
// Reopening within this window is instant (just an unhide); reopening after
// it has elapsed goes through the full splash + rebuild flow, same as boot.
const TRAY_DESTROY_GRACE_PERIOD: Duration = Duration::from_secs(180);

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
          // Real close, not a permanent one: hide immediately so reopening
          // right away (the common case) is instant, then schedule the
          // actual WebView teardown after a grace period so idle time in
          // the tray still gets its RAM back. See TRAY_DESTROY_GRACE_PERIOD.
          api.prevent_close();
          let _ = window.hide();

          let app_handle = window.app_handle().clone();
          let state = app_handle.state::<Arc<SessionState>>().inner().clone();

          // Stamp this hide with the current epoch under the shared lock.
          let my_epoch = {
            let mut epoch = state.hide_epoch.lock().unwrap();
            *epoch += 1;
            *epoch
          };

          tauri::async_runtime::spawn(async move {
            tokio::time::sleep(TRAY_DESTROY_GRACE_PERIOD).await;

            // Hold the lock across the whole check-and-destroy sequence so a
            // concurrent reopen (tray::show_main_window) can't slip in
            // between the epoch check and the destroy() call.
            let epoch = state.hide_epoch.lock().unwrap();
            if *epoch == my_epoch {
              if let Some(w) = app_handle.get_webview_window("main") {
                let _ = w.destroy();
              }
            }
          });
        }
      }
    })
    .build(tauri::generate_context!())
    .expect("error while building tauri application")
    .run(|app_handle, event| {
      // Destroying the "main" window above would otherwise make Tauri quit
      // the whole process once no windows are left. We only want to exit
      // when the user explicitly picks "Quit" from the tray (tray::quit_app
      // sets state.quitting before calling app.exit(0)), so only swallow
      // the exit request when it wasn't an intentional quit.
      if let tauri::RunEvent::ExitRequested { api, .. } = event {
        let state = app_handle.state::<Arc<SessionState>>();
        if !state.quitting.load(Ordering::SeqCst) {
          api.prevent_exit();
        }
      }
    });
}