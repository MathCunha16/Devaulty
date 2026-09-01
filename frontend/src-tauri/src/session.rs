use std::process::Child;
use std::sync::{Arc, Mutex};

// In-memory state shared between Tauri commands and the backend child process.
#[derive(Default)]
pub struct SessionState {
  pub port: Mutex<Option<u16>>,
  pub token: Mutex<Option<String>>,
  pub child_process: Mutex<Option<Child>>,
  pub is_bundled_mode: Mutex<bool>,
  pub active_download_cancel: Mutex<Option<(String, Arc<std::sync::atomic::AtomicBool>)>>,
}

// Response sent to React when it calls `invoke("get_backend_info")`
#[derive(serde::Serialize, Clone)]
pub struct BackendInfo {
  pub port: u16,
  pub token: String,
}