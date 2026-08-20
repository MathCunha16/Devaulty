import fs from "fs";
import path from "path";

export function generateLatestJson(assetsDir, tag, repo = "MathCunha16/Devaulty") {
  if (!fs.existsSync(assetsDir)) {
    console.error(`Assets directory not found: ${assetsDir}`);
    return;
  }

  const files = fs.readdirSync(assetsDir);
  const version = tag.replace(/^v/, "");
  const baseUrl = `https://github.com/${repo}/releases/download/${tag}`;

  const manifest = {
    version: version,
    notes: `Release ${tag}`,
    pub_date: new Date().toISOString(),
    platforms: {},
  };

  const setPlatform = (platform, entry) => {
    if (manifest.platforms[platform]) {
      throw new Error(
        `Conflict: platform "${platform}" was already assigned in manifest.`
      );
    }
    manifest.platforms[platform] = entry;
  };

  for (const file of files) {
    if (file.endsWith(".sig")) {
      const targetFileName = file.slice(0, -4); // remove .sig
      const sigPath = path.join(assetsDir, file);
      const signature = fs.readFileSync(sigPath, "utf-8").trim();

      // Linux Updater Artifacts (AppImage)
      if (
        targetFileName.endsWith(".AppImage.tar.gz") ||
        (targetFileName.includes("AppImage") && !targetFileName.endsWith(".deb") && !targetFileName.endsWith(".rpm"))
      ) {
        const arch =
          targetFileName.includes("aarch64") || targetFileName.includes("arm64")
            ? "linux-aarch64"
            : "linux-x86_64";
        setPlatform(arch, {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        });
      }
      // Windows Updater Artifacts (NSIS exe / zip)
      else if (
        targetFileName.endsWith(".nsis.zip") ||
        targetFileName.endsWith("-setup.exe") ||
        (targetFileName.endsWith(".exe") && !targetFileName.endsWith(".msi")) ||
        (targetFileName.includes("nsis") && !targetFileName.endsWith(".msi"))
      ) {
        const arch =
          targetFileName.includes("aarch64") || targetFileName.includes("arm64")
            ? "windows-aarch64"
            : "windows-x86_64";
        setPlatform(arch, {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        });
      }
      // macOS Updater Artifacts (strictly .app.tar.gz, excluding .dmg)
      else if (targetFileName.endsWith(".app.tar.gz")) {
        if (targetFileName.includes("universal")) {
          setPlatform("darwin-aarch64", {
            signature,
            url: `${baseUrl}/${targetFileName}`,
          });
          setPlatform("darwin-x86_64", {
            signature,
            url: `${baseUrl}/${targetFileName}`,
          });
        } else if (
          targetFileName.includes("x86_64") ||
          targetFileName.includes("intel") ||
          targetFileName.includes("x64")
        ) {
          setPlatform("darwin-x86_64", {
            signature,
            url: `${baseUrl}/${targetFileName}`,
          });
        } else {
          // Default macOS build on GitHub Actions macos-latest runner (Apple Silicon ARM64)
          setPlatform("darwin-aarch64", {
            signature,
            url: `${baseUrl}/${targetFileName}`,
          });
        }
      }
    }
  }

  const outputPath = path.join(assetsDir, "latest.json");
  fs.writeFileSync(outputPath, JSON.stringify(manifest, null, 2) + "\n");
  console.log(`Generated ${outputPath}:\n`, JSON.stringify(manifest, null, 2));
}

const isDirectRun =
  process.argv[1] &&
  (process.argv[1].endsWith("generate-latest-json.js") ||
    process.argv[1].endsWith("generate-latest-json"));

if (isDirectRun) {
  const assetsDir = process.argv[2] || "./release-assets";
  const tag = process.argv[3] || "v0.1.9";
  generateLatestJson(assetsDir, tag);
}
