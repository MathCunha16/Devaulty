use std::process::Command;
use std::sync::Arc;
use std::time::{Duration, Instant};
use tauri::Manager;

use crate::session::{BackendInfo, SessionState};

// Minimum duration the splash screen is displayed, ensuring a polished
// startup experience even when the Go backend initializes instantly.
const MIN_SPLASH_DURATION: Duration = Duration::from_secs(2);

fn is_allowed_download_host(host: &str) -> bool {
  host == "github.com"
    || host == "objects.githubusercontent.com"
    || host.ends_with(".github.com")
    || host.ends_with(".githubusercontent.com")
}

async fn cleanup_partial_download(target_path: &std::path::Path) {
  let _ = tokio::fs::remove_file(target_path).await;
}

// Async IPC command:
//   - Bundled mode (production): waits up to 30s for Go backend to print [DEVAULTY_SESSION],
//     then performs a health check, ensuring a minimum 2s splash display.
//   - Dev mode (no bundled binary): returns dev fallback immediately
#[tauri::command]
pub async fn get_backend_info(state: tauri::State<'_, Arc<SessionState>>) -> Result<BackendInfo, String> {
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
pub fn close_splash(app_handle: tauri::AppHandle) {
  if let Some(splash_window) = app_handle.get_webview_window("splash") {
    let _ = splash_window.close();
  }

  if let Some(main_window) = app_handle.get_webview_window("main") {
    let _ = main_window.maximize();
    let _ = main_window.show();
    let _ = main_window.set_focus();
  }
}

#[derive(serde::Serialize, Clone)]
pub struct AppEnvironment {
  pub os: String,
  pub arch: String,
  pub package_type: String,
  pub supports_in_place_update: bool,
}

#[tauri::command]
pub fn get_app_environment() -> AppEnvironment {
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

struct DownloadCleanupGuard<'a> {
  state: &'a Arc<SessionState>,
  download_id: &'a str,
}

impl<'a> Drop for DownloadCleanupGuard<'a> {
  fn drop(&mut self) {
    let mut lock = self.state.active_download_cancel.lock().unwrap();
    if let Some((ref id, _)) = *lock {
      if id == self.download_id {
        *lock = None;
      }
    }
  }
}

#[tauri::command]
pub async fn download_release_file(
  app: tauri::AppHandle,
  state: tauri::State<'_, Arc<SessionState>>,
  download_id: String,
  url: String,
  filename: String,
) -> Result<String, String> {
  use futures_util::StreamExt;
  use tauri::Emitter;
  use tokio::io::AsyncWriteExt;

  // Validate filename to prevent path traversal
  let raw_path = std::path::Path::new(&filename);
  let safe_name = raw_path
    .file_name()
    .and_then(|n| n.to_str())
    .ok_or_else(|| "Invalid filename".to_string())?;

  if safe_name.is_empty() || safe_name == "." || safe_name == ".." {
    return Err("Invalid filename".to_string());
  }

  // Validate URL scheme and host
  let parsed_url = reqwest::Url::parse(&url).map_err(|e| format!("Invalid URL: {}", e))?;
  if parsed_url.scheme() != "https" {
    return Err("Only HTTPS URLs are permitted".to_string());
  }

  let host = parsed_url.host_str().unwrap_or("");
  if !is_allowed_download_host(host) {
    return Err(format!("Unauthorized download host: {}", host));
  }

  let downloads_dir = dirs::download_dir()
    .or_else(dirs::home_dir)
    .unwrap_or_else(std::env::temp_dir);

  let target_path = downloads_dir.join(safe_name);

  let cancel_flag = Arc::new(std::sync::atomic::AtomicBool::new(false));
  {
    *state.active_download_cancel.lock().unwrap() =
      Some((download_id.clone(), Arc::clone(&cancel_flag)));
  }

  // RAII Guard ensures state.active_download_cancel is cleared on all exit paths
  let _guard = DownloadCleanupGuard {
    state: &state,
    download_id: &download_id,
  };

  let client = reqwest::Client::builder()
    .user_agent("Devaulty-Updater")
    .connect_timeout(Duration::from_secs(15))
    .read_timeout(Duration::from_secs(60))
    .redirect(reqwest::redirect::Policy::custom(|attempt| {
      let next_url = attempt.url();
      if next_url.scheme() != "https" {
        return attempt.stop();
      }

      let next_host = next_url.host_str().unwrap_or("");
      if !is_allowed_download_host(next_host) {
        return attempt.stop();
      }

      if attempt.previous().len() >= 5 {
        return attempt.stop();
      }

      attempt.follow()
    }))
    .build()
    .map_err(|e| format!("Failed to build HTTP client: {}", e))?;

  let response = client
    .get(parsed_url)
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

  while let Some(chunk_result) = stream.next().await {
    if cancel_flag.load(std::sync::atomic::Ordering::Relaxed) {
      drop(file);
      cleanup_partial_download(&target_path).await;
      return Err("Download cancelled by user".to_string());
    }

    let chunk = match chunk_result {
      Ok(chunk) => chunk,
      Err(e) => {
        drop(file);
        cleanup_partial_download(&target_path).await;
        return Err(format!("Error downloading chunk: {}", e));
      }
    };

    if let Err(e) = file.write_all(&chunk).await {
      drop(file);
      cleanup_partial_download(&target_path).await;
      return Err(format!("Error writing chunk to file: {}", e));
    }

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

  if let Err(e) = file.flush().await {
    drop(file);
    cleanup_partial_download(&target_path).await;
    return Err(format!("Failed to flush file: {}", e));
  }

  #[cfg(unix)]
  if safe_name.ends_with(".AppImage") {
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
pub fn cancel_download_release_file(state: tauri::State<'_, Arc<SessionState>>, download_id: String) {
  let lock = state.active_download_cancel.lock().unwrap();
  if let Some((ref id, ref flag)) = *lock {
    if id == &download_id {
      flag.store(true, std::sync::atomic::Ordering::Relaxed);
    }
  }
}

#[tauri::command]
pub fn open_file_path(path: String) -> Result<(), String> {
  let target = std::path::Path::new(&path);
  if !target.exists() {
    return Err(format!("Path does not exist: {}", path));
  }

  let canonical_target = target
    .canonicalize()
    .map_err(|e| format!("Failed to resolve path: {}", e))?;

  let downloads_dir = dirs::download_dir()
    .or_else(dirs::home_dir)
    .unwrap_or_else(std::env::temp_dir);

  let canonical_downloads = downloads_dir
    .canonicalize()
    .map_err(|e| format!("Failed to resolve downloads directory: {}", e))?;

  if !canonical_target.starts_with(&canonical_downloads) {
    return Err("Access denied: Path is outside the downloads directory".to_string());
  }

  #[cfg(target_os = "linux")]
  {
    let folder_to_open = canonical_target
      .parent()
      .unwrap_or(&canonical_downloads);

    Command::new("xdg-open")
      .arg(folder_to_open)
      .spawn()
      .map_err(|e| format!("Failed to open file manager: {}", e))?;
  }

  #[cfg(target_os = "windows")]
  {
    Command::new("explorer")
      .arg(format!("/select,{}", canonical_target.to_string_lossy()))
      .spawn()
      .map_err(|e| format!("Failed to open file manager: {}", e))?;
  }

  #[cfg(target_os = "macos")]
  {
    Command::new("open")
      .arg("-R")
      .arg(&canonical_target)
      .spawn()
      .map_err(|e| format!("Failed to open file manager: {}", e))?;
  }

  Ok(())
}

#[tauri::command]
pub fn open_external_url(url: String) -> Result<(), String> {
  let lower = url.to_lowercase();
  let safe_url = if !lower.starts_with("http://") && !lower.starts_with("https://") {
    format!("https://{}", url)
  } else {
    url
  };

  let parsed = reqwest::Url::parse(&safe_url).map_err(|e| format!("Invalid URL: {}", e))?;
  if parsed.scheme() != "http" && parsed.scheme() != "https" {
    return Err("Only HTTP and HTTPS URLs are permitted".to_string());
  }

  #[cfg(target_os = "linux")]
  {
    Command::new("xdg-open")
      .arg(&safe_url)
      .spawn()
      .map_err(|e| format!("Failed to open URL in browser: {}", e))?;
  }

  #[cfg(target_os = "windows")]
  {
    use std::ffi::OsStr;
    use std::os::windows::ffi::OsStrExt;

    #[link(name = "shell32")]
    extern "system" {
      fn ShellExecuteW(
        hwnd: *mut std::ffi::c_void,
        operation: *const u16,
        file: *const u16,
        parameters: *const u16,
        directory: *const u16,
        show_cmd: i32,
      ) -> isize;
    }

    let operation: Vec<u16> = OsStr::new("open").encode_wide().chain(Some(0)).collect();
    let file: Vec<u16> = OsStr::new(&safe_url).encode_wide().chain(Some(0)).collect();

    let ret = unsafe {
      ShellExecuteW(
        std::ptr::null_mut(),
        operation.as_ptr(),
        file.as_ptr(),
        std::ptr::null(),
        std::ptr::null(),
        1, // SW_SHOWNORMAL
      )
    };

    if ret <= 32 {
      return Err(format!("Failed to open URL in browser, error code: {}", ret));
    }
  }

  #[cfg(target_os = "macos")]
  {
    Command::new("open")
      .arg(&safe_url)
      .spawn()
      .map_err(|e| format!("Failed to open URL in browser: {}", e))?;
  }

  Ok(())
}

#[cfg(target_os = "linux")]
fn try_kdialog(default_path: &std::path::Path, extension: &str) -> Result<Option<std::path::PathBuf>, String> {
  let filter = format!("*.{}", extension);
  let res = std::process::Command::new("kdialog")
    .args(&[
      "--getsavefilename",
      &default_path.to_string_lossy(),
      &filter,
    ])
    .output();

  match res {
    Ok(output) if output.status.success() => {
      let path_str = String::from_utf8_lossy(&output.stdout).trim().to_string();
      if path_str.is_empty() {
        Ok(None)
      } else {
        Ok(Some(std::path::PathBuf::from(path_str)))
      }
    }
    Ok(output) if output.status.code() == Some(1) => Ok(None),
    Ok(output) => Err(format!("kdialog exited with code {:?}", output.status.code())),
    Err(e) => Err(format!("Failed to execute kdialog: {}", e)),
  }
}

#[tauri::command]
pub async fn save_diagram_file(
  default_filename: String,
  extension: String,
  file_bytes: Option<Vec<u8>>,
  file_text: Option<String>,
) -> Result<Option<String>, String> {
  let ext_clean = extension.trim().trim_start_matches('.').to_lowercase();
  if !["png", "svg", "excalidraw"].contains(&ext_clean.as_str()) {
    return Err(format!("Unsupported file extension: {}", extension));
  }

  let downloads_dir = dirs::download_dir()
    .or_else(dirs::home_dir)
    .unwrap_or_else(std::env::temp_dir);

  let default_path = downloads_dir.join(&default_filename);

  let selected_path: Option<std::path::PathBuf> = tokio::task::spawn_blocking({
    let default_path = default_path.clone();
    let ext_clean = ext_clean.clone();
    #[allow(unused_variables)]
    let default_filename = default_filename.clone();
    move || -> Result<Option<std::path::PathBuf>, String> {
      #[cfg(target_os = "linux")]
      {
        let zenity_res = std::process::Command::new("zenity")
          .args(&[
            "--file-selection",
            "--save",
            "--confirm-overwrite",
            &format!("--filename={}", default_path.to_string_lossy()),
            "--title=Save Diagram",
          ])
          .output();

        match zenity_res {
          Ok(output) if output.status.success() => {
            let path_str = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if path_str.is_empty() {
              Ok(None)
            } else {
              Ok(Some(std::path::PathBuf::from(path_str)))
            }
          }
          Ok(output) if output.status.code() == Some(1) => {
            // User cancelled
            Ok(None)
          }
          _ => {
            // Fallback to kdialog if zenity is not installed or failed with non-cancel status
            match try_kdialog(&default_path, &ext_clean) {
              Ok(kdialog_path) => Ok(kdialog_path),
              Err(_) => Err("Failed to open file save dialog (neither zenity nor kdialog is available)".to_string()),
            }
          }
        }
      }

      #[cfg(target_os = "windows")]
      {
        let ps_script = "Add-Type -AssemblyName System.Windows.Forms; \
           $d = New-Object System.Windows.Forms.SaveFileDialog; \
           $d.FileName = $env:DEVAULTY_DEFAULT_FILENAME; \
           $ext = $env:DEVAULTY_FILE_EXT; \
           $d.Filter = \"Files (*.$ext)|*.$ext|All Files (*.*)|*.*\"; \
           if ($d.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) { Write-Host $d.FileName }";

        let mut cmd = std::process::Command::new("powershell");
        cmd.args(&["-NoProfile", "-Command", ps_script])
          .env("DEVAULTY_DEFAULT_FILENAME", &default_filename)
          .env("DEVAULTY_FILE_EXT", &ext_clean);

        match cmd.output() {
          Ok(output) if output.status.success() => {
            let path_str = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if path_str.is_empty() {
              Ok(None)
            } else {
              Ok(Some(std::path::PathBuf::from(path_str)))
            }
          }
          Ok(output) => Err(format!("PowerShell dialog exited with code {:?}", output.status.code())),
          Err(e) => Err(format!("Failed to execute PowerShell dialog: {}", e)),
        }
      }

      #[cfg(target_os = "macos")]
      {
        let script = "set defName to (system attribute \"DEVAULTY_DEFAULT_FILENAME\")\n\
                      POSIX path of (choose file name with prompt \"Save Diagram\" default name defName)";
        let mut cmd = std::process::Command::new("osascript");
        cmd.args(&["-e", script])
          .env("DEVAULTY_DEFAULT_FILENAME", &default_filename);

        match cmd.output() {
          Ok(output) if output.status.success() => {
            let path_str = String::from_utf8_lossy(&output.stdout).trim().to_string();
            if path_str.is_empty() {
              Ok(None)
            } else {
              Ok(Some(std::path::PathBuf::from(path_str)))
            }
          }
          Ok(output) => {
            let stderr = String::from_utf8_lossy(&output.stderr);
            if stderr.contains("User canceled") || stderr.contains("-128") {
              Ok(None)
            } else {
              Err(format!("osascript failed: {}", stderr.trim()))
            }
          }
          Err(e) => Err(format!("Failed to execute osascript dialog: {}", e)),
        }
      }

      #[cfg(not(any(target_os = "linux", target_os = "windows", target_os = "macos")))]
      {
        Err("File save dialog is not supported on this platform".to_string())
      }
    }
  })
  .await
  .map_err(|e| format!("Dialog task error: {}", e))??;

  let target = match selected_path {
    Some(p) => p,
    None => return Ok(None), // User cancelled dialog
  };

  // Ensure extension is retained if user omitted it
  let final_path = if target.extension().is_none() {
    target.with_extension(&extension)
  } else {
    target
  };

  let data_to_write = if let Some(ref text) = file_text {
    text.as_bytes().to_vec()
  } else if let Some(bytes) = file_bytes {
    bytes
  } else {
    return Err("No content provided to save".to_string());
  };

  tokio::fs::write(&final_path, &data_to_write)
    .await
    .map_err(|e| format!("Failed to write file {:?}: {}", final_path, e))?;

  Ok(Some(final_path.to_string_lossy().to_string()))
}