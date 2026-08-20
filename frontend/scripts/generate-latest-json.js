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

  for (const file of files) {
    if (file.endsWith(".sig")) {
      const targetFileName = file.slice(0, -4); // remove .sig
      const sigPath = path.join(assetsDir, file);
      const signature = fs.readFileSync(sigPath, "utf-8").trim();

      if (
        targetFileName.includes("AppImage") ||
        targetFileName.includes("linux") ||
        targetFileName.endsWith(".AppImage.tar.gz")
      ) {
        manifest.platforms["linux-x86_64"] = {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        };
      } else if (
        targetFileName.includes("nsis") ||
        targetFileName.endsWith(".nsis.zip") ||
        targetFileName.includes("setup.exe")
      ) {
        manifest.platforms["windows-x86_64"] = {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        };
      } else if (
        targetFileName.includes("darwin") ||
        targetFileName.includes("aarch64") ||
        targetFileName.includes("app.tar.gz")
      ) {
        manifest.platforms["darwin-aarch64"] = {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        };
        manifest.platforms["darwin-x86_64"] = {
          signature,
          url: `${baseUrl}/${targetFileName}`,
        };
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
