use std::io::{BufRead, BufReader};
use std::process::{Child, Command, Stdio};
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};
use tauri::Manager;
use uuid::Uuid;

// In-memory state shared between Tauri commands and the backend child process.
#[derive(Default)]
pub struct SessionState {
  pub port: Mutex<Option<u16>>,
  pub token: Mutex<Option<String>>,
  pub child_process: Mutex<Option<Child>>,
  pub is_bundled_mode: Mutex<bool>,
}

// Response sent to React when it calls `invoke("get_backend_info")`
#[derive(serde::Serialize, Clone)]
pub struct BackendInfo {
  pub port: u16,
  pub token: String,
}

// Minimum duration the splash screen is displayed, ensuring a polished
// startup experience even when the Go backend initializes instantly.
const MIN_SPLASH_DURATION: Duration = Duration::from_secs(2);

// Async IPC command:
//   - Bundled mode (production): waits up to 30s for Go backend to print [DEVAULTY_SESSION],
//     then performs a health check, ensuring a minimum 2s splash display.
//   - Dev mode (no bundled binary): returns dev fallback immediately
#[tauri::command]
async fn get_backend_info(state: tauri::State<'_, Arc<SessionState>>) -> Result<BackendInfo, String> {
  let is_bundled = *state.is_bundled_mode.lock().unwrap();

  // Dev mode: no bundled binary — return immediately so splash closes fast
  if !is_bundled {
    return Ok(BackendInfo {
      port: 8080,
      token: "dev-token".to_string(),
    });
  }

  let splash_start = Instant::now();
  let timeout = Duration::from_secs(30);

  // Phase 1: Wait for the Go backend to emit its session handshake via stdout
  let info = loop {
    {
      let port = state.port.lock().unwrap();
      let token = state.token.lock().unwrap();

      if let (Some(p), Some(ref t)) = (*port, token.as_ref()) {
        break BackendInfo {
          port: p,
          token: t.to_string(),
        };
      }
    }

    if splash_start.elapsed() > timeout {
      return Err("Backend initialization timeout after 30s".to_string());
    }

    tokio::time::sleep(Duration::from_millis(50)).await;
  };

  // Phase 2: Wait for the HTTP health check endpoint to respond
  let health_url = format!("http://127.0.0.1:{}/health", info.port);
  let health_timeout = Duration::from_secs(10);
  let health_start = Instant::now();

  loop {
    match reqwest::Client::new()
      .get(&health_url)
      .timeout(Duration::from_secs(2))
      .send()
      .await
    {
      Ok(resp) if resp.status().is_success() => break,
      _ => {
        if health_start.elapsed() > health_timeout {
          return Err("Backend health check timeout after 10s".to_string());
        }
        tokio::time::sleep(Duration::from_millis(100)).await;
      }
    }
  }

  // Phase 3: Ensure the splash screen is shown for at least MIN_SPLASH_DURATION
  let elapsed = splash_start.elapsed();
  if elapsed < MIN_SPLASH_DURATION {
    tokio::time::sleep(MIN_SPLASH_DURATION - elapsed).await;
  }

  Ok(info)
}

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

// Resolves the native Go backend binary name based on the current platform.
fn backend_binary_name() -> &'static str {
  if cfg!(target_os = "windows") {
    "devaulty-backend.exe"
  } else {
    "devaulty-backend"
  }
}

// Locates the bundled backend binary inside the Tauri resource directory.
// Tries both `resources/<name>` (Tauri v2 default) and `<name>` (flat layout).
fn find_backend_binary(resource_dir: &std::path::Path) -> Option<std::path::PathBuf> {
  let name = backend_binary_name();

  let candidates = [
    resource_dir.join("resources").join(name),
    resource_dir.join(name),
  ];

  candidates.into_iter().find(|p| p.exists())
}

// Sets executable permission on Unix systems (Linux & macOS).
// On Windows this is a no-op since executability is determined by file extension.
#[cfg(unix)]
fn ensure_executable(path: &std::path::Path) -> std::io::Result<()> {
  use std::os::unix::fs::PermissionsExt;
  let mut perms = std::fs::metadata(path)?.permissions();
  perms.set_mode(0o755);
  std::fs::set_permissions(path, perms)
}

#[cfg(not(unix))]
fn ensure_executable(_path: &std::path::Path) -> std::io::Result<()> {
  Ok(())
}

#[derive(serde::Serialize, Clone)]
pub struct AppEnvironment {
  pub os: String,
  pub arch: String,
  pub package_type: String,
  pub supports_in_place_update: bool,
}

#[tauri::command]
fn get_app_environment() -> AppEnvironment {
  let os = std::env::consts::OS.to_string();
  let arch = std::env::consts::ARCH.to_string();

  let mut package_type = "unknown".to_string();
  let mut supports_in_place_update = true;

  if os == "linux" {
    if std::env::var_os("APPIMAGE").is_some() {
      package_type = "appimage".to_string();
      supports_in_place_update = true;
    } else {
      supports_in_place_update = false;
      if std::path::Path::new("/etc/debian_version").exists()
        || std::path::Path::new("/var/lib/dpkg").exists()
      {
        package_type = "deb".to_string();
      } else if std::path::Path::new("/etc/redhat-release").exists()
        || std::path::Path::new("/etc/fedora-release").exists()
      {
        package_type = "rpm".to_string();
      } else {
        package_type = "deb".to_string();
      }
    }
  } else if os == "windows" {
    package_type = "exe".to_string();
    supports_in_place_update = true;
  } else if os == "macos" {
    package_type = "dmg".to_string();
    supports_in_place_update = true;
  }

  AppEnvironment {
    os,
    arch,
    package_type,
    supports_in_place_update,
  }
}

#[derive(serde::Serialize, Clone)]
pub struct DownloadProgressPayload {
  pub percentage: u32,
  pub downloaded_bytes: u64,
  pub total_bytes: u64,
}

#[tauri::command]
async fn download_release_file(
  app: tauri::AppHandle,
  url: String,
  filename: String,
) -> Result<String, String> {
  use futures_util::StreamExt;
  use tauri::Emitter;
  use tokio::io::AsyncWriteExt;

  let downloads_dir = dirs::download_dir()
    .or_else(dirs::home_dir)
    .unwrap_or_else(std::env::temp_dir);

  let target_path = downloads_dir.join(&filename);

  let client = reqwest::Client::builder()
    .user_agent("Devaulty-Updater")
    .build()
    .map_err(|e| format!("Failed to build HTTP client: {}", e))?;

  let response = client
    .get(&url)
    .send()
    .await
    .map_err(|e| format!("Failed to connect to download URL: {}", e))?;

  if !response.status().is_success() {
    return Err(format!(
      "Download failed with HTTP status: {}",
      response.status()
    ));
  }

  let total_size = response.content_length().unwrap_or(0);
  let mut downloaded: u64 = 0;

  let mut file = tokio::fs::File::create(&target_path)
    .await
    .map_err(|e| format!("Failed to create destination file {:?}: {}", target_path, e))?;

  let mut stream = response.bytes_stream();

  while let Some(chunk) = stream.next().await {
    let chunk = chunk.map_err(|e| format!("Error downloading chunk: {}", e))?;
    file
      .write_all(&chunk)
      .await
      .map_err(|e| format!("Error writing chunk to file: {}", e))?;

    downloaded += chunk.len() as u64;
    let percentage = if total_size > 0 {
      ((downloaded as f64 / total_size as f64) * 100.0).min(100.0) as u32
    } else {
      0
    };

    let _ = app.emit(
      "download-file-progress",
      DownloadProgressPayload {
        percentage,
        downloaded_bytes: downloaded,
        total_bytes: total_size,
      },
    );
  }

  file
    .flush()
    .await
    .map_err(|e| format!("Failed to flush file: {}", e))?;

  #[cfg(unix)]
  if filename.ends_with(".AppImage") {
    use std::os::unix::fs::PermissionsExt;
    if let Ok(metadata) = std::fs::metadata(&target_path) {
      let mut perms = metadata.permissions();
      perms.set_mode(0o755);
      let _ = std::fs::set_permissions(&target_path, perms);
    }
  }

  Ok(target_path.to_string_lossy().to_string())
}

#[tauri::command]
fn open_file_path(path: String) -> Result<(), String> {
  let target = std::path::Path::new(&path);
  if !target.exists() {
    return Err(format!("Path does not exist: {}", path));
  }

  #[cfg(target_os = "linux")]
  {
    Command::new("xdg-open")
      .arg(&path)
      .spawn()
      .map_err(|e| format!("Failed to open file: {}", e))?;
  }

  #[cfg(target_os = "windows")]
  {
    Command::new("explorer")
      .arg(&path)
      .spawn()
      .map_err(|e| format!("Failed to open file: {}", e))?;
  }

  #[cfg(target_os = "macos")]
  {
    Command::new("open")
      .arg(&path)
      .spawn()
      .map_err(|e| format!("Failed to open file: {}", e))?;
  }

  Ok(())
}

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
    .plugin(tauri_plugin_updater::Builder::new().build())
    .plugin(tauri_plugin_process::init())
    .manage(session_state)
    .invoke_handler(tauri::generate_handler![
      close_splash,
      get_backend_info,
      get_app_environment,
      download_release_file,
      open_file_path
    ])
    .setup(move |app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      let resource_dir = app.path().resource_dir().ok();
      if let Some(ref res_path) = resource_dir {
        if let Some(binary_path) = find_backend_binary(res_path) {
          // Ensure executable permission on Unix platforms
          if let Err(e) = ensure_executable(&binary_path) {
            log::error!("Failed to set executable permission on backend binary: {}", e);
          }

          // Generate a cryptographically secure random session token.
          // This token is injected exclusively via the child process environment
          // and never exposed to disk, logs, or other OS processes.
          let session_token = Uuid::new_v4().to_string();

          let data_dir_str = devaulty_data_dir.to_string_lossy().to_string();

          // Spawn the native Go backend process.
          // Migrations are embedded in the Go binary via go:embed,
          // so no external migration files need to be shipped.
          match Command::new(&binary_path)
            .env("APP_ENV", "prod")
            .env("DEVAULTY_INTERNAL_TOKEN", &session_token)
            .env("DEVAULTY_DATA_DIR", &data_dir_str)
            .stdout(Stdio::piped())
            .stderr(Stdio::inherit())
            .spawn()
          {
            Ok(mut child) => {
              // Mark as bundled production mode only after successful process spawn
              *state_clone.is_bundled_mode.lock().unwrap() = true;

              if let Some(stdout) = child.stdout.take() {
                let state_inner = Arc::clone(&state_clone);
                std::thread::spawn(move || {
                  let reader = BufReader::new(stdout);
                  for line in reader.lines().flatten() {
                    if line.contains("[DEVAULTY_SESSION]") {
                      let parts: Vec<&str> = line.split_whitespace().collect();
                      let mut port_val = None;
                      let mut token_val = None;
                      for part in parts {
                        if let Some(stripped) = part.strip_prefix("PORT=") {
                          port_val = stripped.parse::<u16>().ok();
                        } else if let Some(stripped) = part.strip_prefix("TOKEN=") {
                          token_val = Some(stripped.to_string());
                        }
                      }
                      if let (Some(p), Some(t)) = (port_val, token_val) {
                        *state_inner.port.lock().unwrap() = Some(p);
                        *state_inner.token.lock().unwrap() = Some(t);
                      }
                    }
                  }
                });
              }
              *state_clone.child_process.lock().unwrap() = Some(child);
            }
            Err(e) => {
              *state_clone.is_bundled_mode.lock().unwrap() = false;
              log::error!("Failed to spawn Go backend process: {}", e);
            }
          }
        }
      }

      Ok(())
    })
    .on_window_event(|window, event| {
      if let tauri::WindowEvent::CloseRequested { .. } = event {
        if window.label() == "main" {
          let state = window.state::<Arc<SessionState>>();
          let mut guard = state.child_process.lock().unwrap();
          if let Some(mut child) = guard.take() {
            let _ = child.kill();
          }
        }
      }
    })
    .run(tauri::generate_context!())
    .expect("error while running tauri application");
}
