import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Source of truth: package.json version field
const packageJsonPath = path.resolve(__dirname, "../package.json");
const cargoTomlPath = path.resolve(__dirname, "../src-tauri/Cargo.toml");
const tauriConfPath = path.resolve(__dirname, "../src-tauri/tauri.conf.json");
const backendVersionGoPath = path.resolve(__dirname, "../../backend/internal/domain/model/version.go");

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

  // 2. Update Cargo.toml (supports full SemVer including pre-release tags like 0.1.6-alpha.1)
  if (fs.existsSync(cargoTomlPath)) {
    let cargoContent = fs.readFileSync(cargoTomlPath, "utf-8");
    cargoContent = cargoContent.replace(/^version = ".*?"/m, `version = "${version}"`);
    fs.writeFileSync(cargoTomlPath, cargoContent);
    console.log(`Updated src-tauri/Cargo.toml -> ${version}`);
  }

  // 3. Update tauri.conf.json (full SemVer matching package.json)
  if (fs.existsSync(tauriConfPath)) {
    const tauriConf = JSON.parse(fs.readFileSync(tauriConfPath, "utf-8"));
    tauriConf.version = version;
    fs.writeFileSync(tauriConfPath, JSON.stringify(tauriConf, null, 2) + "\n");
    console.log(`Updated src-tauri/tauri.conf.json -> ${version}`);
  }

  // 4. Update backend version.go (Go backend compiled version)
  if (fs.existsSync(backendVersionGoPath)) {
    const formattedGoVersion = version.startsWith("v") ? version : `v${version}`;
    let goContent = fs.readFileSync(backendVersionGoPath, "utf-8");
    goContent = goContent.replace(/var AppVersion = ".*?"/m, `var AppVersion = "${formattedGoVersion}"`);
    fs.writeFileSync(backendVersionGoPath, goContent);
    console.log(`Updated backend/internal/domain/model/version.go -> ${formattedGoVersion}`);
  }

  console.log("\nAll frontend, backend & Tauri version fields are synchronized!");
}

syncVersion();

