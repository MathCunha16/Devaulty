use std::process::Child;
use std::sync::atomic::AtomicBool;
use std::sync::{Arc, Mutex};

// In-memory state shared between Tauri commands and the backend child process.
#[derive(Default)]
pub struct SessionState {
  pub port: Mutex<Option<u16>>,
  pub token: Mutex<Option<String>>,
  pub child_process: Mutex<Option<Child>>,
  pub is_bundled_mode: Mutex<bool>,
  pub active_download_cancel: Mutex<Option<(String, Arc<std::sync::atomic::AtomicBool>)>>,
  // Guards the hide -> grace-period -> destroy sequence against reopen races.
  // Both the destroy timer (lib.rs) and any path that reopens/recreates the
  // window (tray::show_main_window) must hold this lock for their entire
  // check-then-act sequence, not just the epoch read - otherwise a reopen
  // landing between the epoch check and the actual destroy() call would let
  // the timer tear down a window the user just reopened.
  pub hide_epoch: Mutex<u64>,
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