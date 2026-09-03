use std::io::{BufRead, BufReader};
use std::process::{Command, Stdio};
use std::sync::Arc;
use tauri::Manager;
use uuid::Uuid;

use crate::session::SessionState;

// Resolves the native Go backend binary name based on the current platform.
pub fn backend_binary_name() -> &'static str {
  if cfg!(target_os = "windows") {
    "devaulty-backend.exe"
  } else {
    "devaulty-backend"
  }
}

// Locates the bundled backend binary inside the Tauri resource directory.
// Tries both `resources/<name>` (Tauri v2 default) and `<name>` (flat layout).
pub fn find_backend_binary(resource_dir: &std::path::Path) -> Option<std::path::PathBuf> {
  let name = backend_binary_name();

  let candidates = [
    resource_dir.join("resources").join(name),
    resource_dir.join(name),
  ];

  candidates.into_iter().find(|p| p.exists())
}

// Writes a tiny wrapper script at a fixed, predictable path that execs the
// real backend binary. This avoids duplicating the binary on disk while
// giving external tools (MCP clients) a stable path across all package
// formats (.deb, .rpm, .exe, .dmg, .AppImage).
pub fn install_standalone_cli(
  resource_binary: &std::path::Path,
  data_dir: &std::path::Path,
) -> std::io::Result<std::path::PathBuf> {
  let bin_dir = data_dir.join("bin");
  std::fs::create_dir_all(&bin_dir)?;

  #[cfg(unix)]
  {
    let dest = bin_dir.join("devaulty-backend");

    // Remove any stale file/symlink from a previous version before relinking.
    if dest.symlink_metadata().is_ok() {
      let _ = std::fs::remove_file(&dest);
    }

    std::os::unix::fs::symlink(resource_binary, &dest)?;
    Ok(dest)
  }

  #[cfg(windows)]
  {
    let dest = bin_dir.join("devaulty-backend.bat");
    let script = format!("@echo off\r\n\"{}\" %*\r\n", resource_binary.display());
    std::fs::write(&dest, script)?;
    Ok(dest)
  }
}

// Sets executable permission on Unix systems (Linux & macOS).
// On Windows this is a no-op since executability is determined by file extension.
#[cfg(unix)]
pub fn ensure_executable(path: &std::path::Path) -> std::io::Result<()> {
  use std::os::unix::fs::PermissionsExt;
  let mut perms = std::fs::metadata(path)?.permissions();
  perms.set_mode(0o755);
  std::fs::set_permissions(path, perms)
}

#[cfg(not(unix))]
pub fn ensure_executable(_path: &std::path::Path) -> std::io::Result<()> {
  Ok(())
}

// Locates, prepares, and spawns the bundled Go backend as a child process,
// wiring its stdout session handshake into the shared SessionState.
// Straight extraction of the logic that used to live inline in run()'s
// `.setup()` closure — same behavior, just callable on its own.
pub fn spawn_backend(app: &tauri::AppHandle, devaulty_data_dir: &std::path::Path, state: &Arc<SessionState>) {
  let resource_dir = app.path().resource_dir().ok();
  if let Some(ref res_path) = resource_dir {
    if let Some(binary_path) = find_backend_binary(res_path) {
      // Ensure executable permission on Unix platforms
      if let Err(e) = ensure_executable(&binary_path) {
        log::error!("Failed to set executable permission on backend binary: {}", e);
      }

      // Install a stable wrapper script pointing to the backend binary so external
      // tools (MCP clients) can reference a fixed path across all package formats.
      if let Err(e) = install_standalone_cli(&binary_path, devaulty_data_dir) {
        log::error!("Failed to install standalone CLI wrapper: {}", e);
      }

      // Generate a cryptographically secure random session token.
      // This token is injected exclusively via the child process environment
      // and never exposed to disk, logs, or other OS processes.
      let session_token = Uuid::new_v4().to_string();

      let data_dir_str = devaulty_data_dir.to_string_lossy().to_string();

      // Spawn the native Go backend process.
      // Migrations are embedded in the Go binary via go:embed,
      // so no external migration files need to be shipped.
      let mut command = Command::new(&binary_path);
      command
        .env("APP_ENV", "prod")
        .env("DEVAULTY_INTERNAL_TOKEN", &session_token)
        .env("DEVAULTY_DATA_DIR", &data_dir_str)
        .stdout(Stdio::piped())
        .stderr(Stdio::inherit());

      // Prevents a visible console window from flashing up when spawning
      // the Go backend on Windows (Go binaries default to console subsystem,
      // independent of the Tauri process's own windows_subsystem attribute).
      #[cfg(windows)]
      {
        use std::os::windows::process::CommandExt;
        const CREATE_NO_WINDOW: u32 = 0x08000000;
        command.creation_flags(CREATE_NO_WINDOW);
      }

      match command.spawn() {
        Ok(mut child) => {
          // Mark as bundled production mode only after successful process spawn
          *state.is_bundled_mode.lock().unwrap() = true;

          if let Some(stdout) = child.stdout.take() {
            let state_inner = Arc::clone(state);
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
          *state.child_process.lock().unwrap() = Some(child);
        }
        Err(e) => {
          *state.is_bundled_mode.lock().unwrap() = false;
          log::error!("Failed to spawn Go backend process: {}", e);
        }
      }
    }
  }
}