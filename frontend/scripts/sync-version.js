import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Path to application.yaml (source of truth)
const appYamlPath = path.resolve(__dirname, "../../backend/src/main/resources/application.yaml");
const packageJsonPath = path.resolve(__dirname, "../package.json");
const cargoTomlPath = path.resolve(__dirname, "../src-tauri/Cargo.toml");
const tauriConfPath = path.resolve(__dirname, "../src-tauri/tauri.conf.json");

function syncVersion() {
  if (!fs.existsSync(appYamlPath)) {
    console.error(`Source of truth file not found at: ${appYamlPath}`);
    process.exit(1);
  }

  // 1. Read raw version from application.yaml
  const yamlContent = fs.readFileSync(appYamlPath, "utf-8");
  const versionLine = yamlContent
    .split("\n")
    .find((line) => line.trim().startsWith("version:"));

  if (!versionLine) {
    console.error("Could not find 'version:' line in application.yaml");
    process.exit(1);
  }

  // Extract raw version string (e.g., "0.1.6-alpha" or '0.1.6')
  const rawVersion = versionLine
    .split("version:")[1]
    .trim()
    .replace(/^["']|["']$/g, "");

  // Clean version without suffix for npm/semver/cargo (e.g., "0.1.6")
  const cleanVersion = rawVersion.replace(/-.*$/, "");

  console.log(`Extracted Backend Version: "${rawVersion}" (Clean SemVer: "${cleanVersion}")`);

  // 2. Update package.json
  if (fs.existsSync(packageJsonPath)) {
    const pkg = JSON.parse(fs.readFileSync(packageJsonPath, "utf-8"));
    pkg.version = cleanVersion;
    fs.writeFileSync(packageJsonPath, JSON.stringify(pkg, null, 2) + "\n");
    console.log(`Updated package.json -> ${cleanVersion}`);
  }

  // 3. Update Cargo.toml
  if (fs.existsSync(cargoTomlPath)) {
    let cargoContent = fs.readFileSync(cargoTomlPath, "utf-8");
    cargoContent = cargoContent.replace(/^version = ".*?"/m, `version = "${cleanVersion}"`);
    fs.writeFileSync(cargoTomlPath, cargoContent);
    console.log(`Updated src-tauri/Cargo.toml -> ${cleanVersion}`);
  }

  // 4. Update tauri.conf.json
  if (fs.existsSync(tauriConfPath)) {
    const tauriConf = JSON.parse(fs.readFileSync(tauriConfPath, "utf-8"));
    tauriConf.version = cleanVersion;
    fs.writeFileSync(tauriConfPath, JSON.stringify(tauriConf, null, 2) + "\n");
    console.log(`Updated src-tauri/tauri.conf.json -> ${cleanVersion}`);
  }

  console.log("\nAll frontend & Tauri version fields are synchronized with Backend!");
}

syncVersion();
