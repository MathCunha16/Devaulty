import fs from "fs";
import path from "path";
import { execSync } from "child_process";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const backendGoDir = path.resolve(__dirname, "../../backend-go");
const frontendDir = path.resolve(__dirname, "..");
const resourcesDir = path.resolve(__dirname, "../src-tauri/resources");

const isWindows = process.platform === "win32";
const binaryName = isWindows ? "devaulty-backend.exe" : "devaulty-backend";

function runCommand(command, cwd) {
  console.log(`\nRunning: "${command}" in ${cwd}`);
  execSync(command, { cwd, stdio: "inherit" });
}

function buildAll() {
  // 1. Compile Go backend into a production-optimized native binary.
  // SQL migrations are embedded into the binary via go:embed at compile time,
  // so no external migration files need to be shipped.
  console.log("Step 1: Building Go backend binary...");
  const targetBinaryPath = path.join(resourcesDir, binaryName);

  if (!fs.existsSync(resourcesDir)) {
    fs.mkdirSync(resourcesDir, { recursive: true });
  } else {
    // Clean any residual backend binaries from previous builds to prevent cross-platform file pollution
    const existingFiles = fs.readdirSync(resourcesDir);
    for (const file of existingFiles) {
      if (file.startsWith("devaulty-backend")) {
        try {
          fs.unlinkSync(path.join(resourcesDir, file));
        } catch {
          // Ignore deletion errors
        }
      }
    }
  }

  // -ldflags="-s -w" strips debug symbols and DWARF info for a smaller binary
  const goBuildCmd = `go build -ldflags="-s -w" -o "${targetBinaryPath}" ./cmd/api/`;
  runCommand(goBuildCmd, backendGoDir);

  // Set executable permissions on Unix systems
  if (!isWindows) {
    fs.chmodSync(targetBinaryPath, 0o755);
  }

  console.log(`Go backend compiled -> src-tauri/resources/${binaryName}`);

  // 2. Synchronize version across all project manifests
  console.log("\nStep 2: Synchronizing project versions...");
  runCommand("node scripts/sync-version.js", frontendDir);

  // 3. Build frontend web assets (Vite)
  console.log("\nStep 3: Building React/Vite web bundle...");
  runCommand("npx vite build", frontendDir);

  console.log("\nCross-Platform Build Prep Complete! Ready for Tauri packaging.\n");
}

buildAll();
