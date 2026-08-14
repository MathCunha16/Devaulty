import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Source of truth: package.json version field
const packageJsonPath = path.resolve(__dirname, "../package.json");
const cargoTomlPath = path.resolve(__dirname, "../src-tauri/Cargo.toml");
const tauriConfPath = path.resolve(__dirname, "../src-tauri/tauri.conf.json");

function syncVersion() {
  if (!fs.existsSync(packageJsonPath)) {
    console.error(`Source of truth file not found at: ${packageJsonPath}`);
    process.exit(1);
  }

  // 1. Read version from package.json (source of truth)
  const pkg = JSON.parse(fs.readFileSync(packageJsonPath, "utf-8"));
  const version = pkg.version;

  if (!version) {
    console.error("Could not find 'version' field in package.json");
    process.exit(1);
  }

  console.log(`Source version (package.json): "${version}"`);

  // 2. Update Cargo.toml
  if (fs.existsSync(cargoTomlPath)) {
    let cargoContent = fs.readFileSync(cargoTomlPath, "utf-8");
    cargoContent = cargoContent.replace(/^version = ".*?"/m, `version = "${version}"`);
    fs.writeFileSync(cargoTomlPath, cargoContent);
    console.log(`Updated src-tauri/Cargo.toml -> ${version}`);
  }

  // 3. Update tauri.conf.json
  if (fs.existsSync(tauriConfPath)) {
    const tauriConf = JSON.parse(fs.readFileSync(tauriConfPath, "utf-8"));
    tauriConf.version = version;
    fs.writeFileSync(tauriConfPath, JSON.stringify(tauriConf, null, 2) + "\n");
    console.log(`Updated src-tauri/tauri.conf.json -> ${version}`);
  }

  console.log("\nAll frontend & Tauri version fields are synchronized!");
}

syncVersion();
