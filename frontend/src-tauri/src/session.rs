use std::process::Child;
use std::sync::atomic::{AtomicBool, AtomicU64};
use std::sync::{Arc, Mutex};

// In-memory state shared between Tauri commands and the backend child process.
#[derive(Default)]
pub struct SessionState {
  pub port: Mutex<Option<u16>>,
  pub token: Mutex<Option<String>>,
  pub child_process: Mutex<Option<Child>>,
  pub is_bundled_mode: Mutex<bool>,
  pub active_download_cancel: Mutex<Option<(String, Arc<std::sync::atomic::AtomicBool>)>>,
  // Bumped every time the main window is hidden (tray minimize) or reshown.
  // Used by the tray-destroy grace-period timer in lib.rs: a timer captures
  // the epoch value at the moment it's scheduled, and only destroys the
  // window if the epoch is still the same after the grace period elapses -
  // i.e. nobody reopened the window in the meantime. Reopening bumps the
  // epoch, silently invalidating any pending timer from an earlier hide.
  pub hide_epoch: AtomicU64,
  // Set right before tray::quit_app calls app.exit(0), so the ExitRequested
  // handler in lib.rs's run() closure knows this exit is intentional and
  // should NOT be swallowed by the prevent_exit() used to survive the main
  // window being destroyed with no windows left open.
  pub quitting: AtomicBool,
}

// Response sent to React when it calls `invoke("get_backend_info")`
#[derive(serde::Serialize, Clone)]
pub struct BackendInfo {
  pub port: u16,
  pub token: String,
}