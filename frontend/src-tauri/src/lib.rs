use std::io::{BufRead, BufReader};
use std::process::{Child, Command, Stdio};
use std::sync::{Arc, Mutex};
use std::time::{Duration, Instant};
use tauri::Manager;

// In-memory structure to store port, token, child process, and jar mode status.
#[derive(Default)]
pub struct SessionState {
  pub port: Mutex<Option<u16>>,
  pub token: Mutex<Option<String>>,
  pub child_process: Mutex<Option<Child>>,
  pub is_jar_mode: Mutex<bool>,
}

// Response sent to React when it calls `invoke("get_backend_info")`
#[derive(serde::Serialize, Clone)]
pub struct BackendInfo {
  pub port: u16,
  pub token: String,
}

// Async IPC command:
//   - If bundled JAR exists (production): waits up to 60s for Spring Boot to print [DEVAULTY_SESSION]
//   - If no JAR (dev mode with tauri dev): returns fallback immediately
#[tauri::command]
async fn get_backend_info(state: tauri::State<'_, Arc<SessionState>>) -> Result<BackendInfo, String> {
  let is_jar = *state.is_jar_mode.lock().unwrap();

  // Dev mode: no bundled JAR — return immediately so splash closes fast
  if !is_jar {
    return Ok(BackendInfo {
      port: 8080,
      token: "dev-session-token".to_string(),
    });
  }

  // Production mode: poll until Spring Boot writes the session line to stdout
  let start_time = Instant::now();
  let timeout = Duration::from_secs(60);

  loop {
    {
      let port = state.port.lock().unwrap();
      let token = state.token.lock().unwrap();

      if let (Some(p), Some(ref t)) = (*port, token.as_ref()) {
        return Ok(BackendInfo {
          port: p,
          token: t.to_string(),
        });
      }
    }

    if start_time.elapsed() > timeout {
      return Err("Backend initialization timeout after 60s".to_string());
    }

    tokio::time::sleep(Duration::from_millis(100)).await;
  }
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

// Resolves the `java` executable path, checking common install locations on Linux/macOS
// to work around missing PATH when app is launched from the desktop launcher.
fn resolve_java_binary() -> Option<String> {
  // 1. Try java directly from PATH
  if Command::new("java").arg("-version").stdout(Stdio::null()).stderr(Stdio::null()).status().is_ok() {
    return Some("java".to_string());
  }

  // 2. Common explicit paths as fallback
  let candidates = [
    "/usr/bin/java",
    "/usr/local/bin/java",
    "/usr/lib/jvm/default/bin/java",
    "/usr/lib/jvm/default-java/bin/java",
    "/usr/lib/jvm/java-21-openjdk-amd64/bin/java",
    "/usr/lib/jvm/java-17-openjdk-amd64/bin/java",
  ];

  for path in candidates {
    if std::path::Path::new(path).exists() {
      return Some(path.to_string());
    }
  }

  None
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  if let Some(config_dir) = dirs::config_dir() {
    let devaulty_dir = config_dir.join("devaulty");
    if !devaulty_dir.exists() {
      let _ = std::fs::create_dir_all(devaulty_dir);
    }
  }

  let session_state = Arc::new(SessionState::default());
  let state_clone = Arc::clone(&session_state);

  tauri::Builder::default()
    .manage(session_state)
    .invoke_handler(tauri::generate_handler![close_splash, get_backend_info])
    .setup(move |app| {
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      let resource_dir = app.path().resource_dir().ok();
      if let Some(res_path) = resource_dir {
        // Try both possible paths: Tauri v2 preserves the resources/ subdir by default
        let jar_path = if res_path.join("resources/backend.jar").exists() {
          res_path.join("resources/backend.jar")
        } else {
          res_path.join("backend.jar")
        };
        if jar_path.exists() {
          // Mark as production JAR mode BEFORE spawning Java
          *state_clone.is_jar_mode.lock().unwrap() = true;

          if let Some(java_bin) = resolve_java_binary() {
            if let Ok(mut child) = Command::new(&java_bin)
              .env("SPRING_PROFILES_ACTIVE", "prod")
              .args([
                "-Xms64m",
                "-Xmx256m",
                "-XX:MetaspaceSize=96m",
                "-XX:MaxMetaspaceSize=192m",
                "-XX:ParallelGCThreads=2",
                "-XX:ConcGCThreads=1",
                "-XX:+UseG1GC",
                "-XX:MaxGCPauseMillis=100",
                "-jar",
                jar_path.to_str().unwrap(),
                "--spring.profiles.active=prod",
              ])
              .stdout(Stdio::piped())
              .spawn()
            {
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
                        if part.starts_with("PORT=") {
                          port_val = part.trim_start_matches("PORT=").parse::<u16>().ok();
                        } else if part.starts_with("TOKEN=") {
                          token_val = Some(part.trim_start_matches("TOKEN=").to_string());
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
