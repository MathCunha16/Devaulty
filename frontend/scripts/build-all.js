import fs from "fs";
import path from "path";
import { execSync } from "child_process";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const backendDir = path.resolve(__dirname, "../../backend");
const frontendDir = path.resolve(__dirname, "..");
const resourcesDir = path.resolve(__dirname, "../src-tauri/resources");

function runCommand(command, cwd) {
  console.log(`\nRunning: "${command}" in ${cwd}`);
  execSync(command, { cwd, stdio: "inherit" });
}

function buildAll() {
  const isWindows = process.platform === "win32";
  const gradlewCmd = isWindows ? "gradlew.bat" : "./gradlew";

  // 1. Build Spring Boot executable JAR in backend
  console.log("Step 1: Building Backend Spring Boot JAR...");
  runCommand(`${gradlewCmd} bootJar`, backendDir);

  // 2. Locate generated JAR in backend/build/libs/
  const backendLibsDir = path.join(backendDir, "build/libs");
  if (!fs.existsSync(backendLibsDir)) {
    console.error(`Backend build libs directory not found: ${backendLibsDir}`);
    process.exit(1);
  }

  const jarFiles = fs.readdirSync(backendLibsDir).filter((file) => file.endsWith(".jar") && !file.endsWith("-plain.jar"));
  if (jarFiles.length === 0) {
    console.error("No executable JAR file found in backend/build/libs/");
    process.exit(1);
  }

  const sourceJarPath = path.join(backendLibsDir, jarFiles[0]);

  // 3. Ensure src-tauri/resources exists and copy backend.jar
  if (!fs.existsSync(resourcesDir)) {
    fs.mkdirSync(resourcesDir, { recursive: true });
  }

  const targetJarPath = path.join(resourcesDir, "backend.jar");
  fs.copyFileSync(sourceJarPath, targetJarPath);
  console.log(`Copied ${jarFiles[0]} -> src-tauri/resources/backend.jar`);

  // 4. Synchronize Version across files
  console.log("\nStep 2: Synchronizing project versions...");
  runCommand("node scripts/sync-version.js", frontendDir);

  // 5. Build Frontend assets (Vite)
  console.log("\nStep 3: Building React/Vite web bundle...");
  runCommand("npx vite build", frontendDir);

  console.log("\nCross-Platform Build Prep Complete! Ready for Tauri packaging.\n");
}

buildAll();
