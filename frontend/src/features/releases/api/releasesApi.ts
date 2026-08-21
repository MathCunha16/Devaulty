import { invoke } from "@tauri-apps/api/core";
import { listen } from "@tauri-apps/api/event";
import { check, Update } from "@tauri-apps/plugin-updater";
import { relaunch } from "@tauri-apps/plugin-process";
import pkg from "../../../../package.json";
import type {
  CurrentVersionResponse,
  AppEnvironment,
  ReleaseAssetFormat,
  AppUpdateInfoResponse,
  UpdateDownloadProgressResponse,
} from "../../../types/api";

let cachedUpdate: Update | null = null;

interface DownloadProgressPayload {
  percentage: number;
  downloaded_bytes: number;
  total_bytes: number;
}

export const releasesApi = {
  getCurrentVersion: async (): Promise<CurrentVersionResponse> => {
    const version = pkg.version;
    return {
      currentVersion: version,
      actualVersion: version,
    };
  },

  getAppEnvironment: async (): Promise<AppEnvironment> => {
    try {
      const env = await invoke<AppEnvironment>("get_app_environment");
      return env;
    } catch {
      return {
        os: "linux",
        arch: "x86_64",
        package_type: "deb",
        supports_in_place_update: false,
      };
    }
  },

  checkUpdates: async (): Promise<AppUpdateInfoResponse> => {
    try {
      const currentVersion = (await releasesApi.getCurrentVersion()).currentVersion;
      const env = await releasesApi.getAppEnvironment();
      const update = await check();
      cachedUpdate = update;

      if (update && update.available) {
        const latestVersion = update.version;
        const availableFormats: ReleaseAssetFormat[] = [
          {
            label: "Debian / Ubuntu Package (.deb)",
            packageType: "deb",
            filename: `Devaulty_${latestVersion}_amd64.deb`,
            url: `https://github.com/MathCunha16/Devaulty/releases/download/v${latestVersion}/Devaulty_${latestVersion}_amd64.deb`,
          },
          {
            label: "RedHat / Fedora Package (.rpm)",
            packageType: "rpm",
            filename: `Devaulty-${latestVersion}-1.x86_64.rpm`,
            url: `https://github.com/MathCunha16/Devaulty/releases/download/v${latestVersion}/Devaulty-${latestVersion}-1.x86_64.rpm`,
          },
          {
            label: "Standalone Linux AppImage (.AppImage)",
            packageType: "appimage",
            filename: `Devaulty_${latestVersion}_amd64.AppImage`,
            url: `https://github.com/MathCunha16/Devaulty/releases/download/v${latestVersion}/Devaulty_${latestVersion}_amd64.AppImage`,
          },
          {
            label: "Windows Setup Installer (.exe)",
            packageType: "exe",
            filename: `Devaulty_${latestVersion}_x64-setup.exe`,
            url: `https://github.com/MathCunha16/Devaulty/releases/download/v${latestVersion}/Devaulty_${latestVersion}_x64-setup.exe`,
          },
          {
            label: "macOS Disk Image (.dmg)",
            packageType: "dmg",
            filename: `Devaulty_${latestVersion}_aarch64.dmg`,
            url: `https://github.com/MathCunha16/Devaulty/releases/download/v${latestVersion}/Devaulty_${latestVersion}_aarch64.dmg`,
          },
        ];

        const matchedFormat =
          availableFormats.find((f) => f.packageType === env.package_type) ||
          availableFormats[0];

        return {
          updateAvailable: true,
          currentVersion: update.currentVersion || currentVersion,
          latestVersion,
          releaseTitle: `Release v${latestVersion}`,
          releaseNotes: update.body || "",
          publishedAt: update.date,
          packageType: env.package_type,
          supportsInPlaceUpdate: env.supports_in_place_update,
          downloadUrl: matchedFormat.url,
          downloadFilename: matchedFormat.filename,
          availableFormats,
        };
      }

      return {
        updateAvailable: false,
        currentVersion,
        latestVersion: currentVersion,
        packageType: env.package_type,
        supportsInPlaceUpdate: env.supports_in_place_update,
      };
    } catch (err: unknown) {
      console.error("Error checking for updates:", err);
      const currentVersion = (await releasesApi.getCurrentVersion()).currentVersion;
      return {
        updateAvailable: false,
        currentVersion,
        latestVersion: currentVersion,
      };
    }
  },

  downloadAndInstall: (
    onProgress: (data: UpdateDownloadProgressResponse) => void,
    onError: (errorMessage: string) => void
  ): (() => void) => {
    let isCancelled = false;

    (async () => {
      try {
        let update = cachedUpdate;
        if (!update) {
          update = await check();
          cachedUpdate = update;
        }

        if (isCancelled) return;

        if (!update) {
          onError("No update available to download.");
          return;
        }

        let downloadedBytes = 0;
        let totalBytes = 0;

        await update.downloadAndInstall((event) => {
          if (isCancelled) return;

          switch (event.event) {
            case "Started": {
              totalBytes = event.data.contentLength || 0;
              onProgress({
                status: "DOWNLOADING",
                percentage: 0,
                downloadedBytes: 0,
                totalBytes,
              });
              break;
            }

            case "Progress": {
              downloadedBytes += event.data.chunkLength;
              const percentage =
                totalBytes > 0
                  ? Math.min(100, Math.round((downloadedBytes / totalBytes) * 100))
                  : 0;
              onProgress({
                status: "DOWNLOADING",
                percentage,
                downloadedBytes,
                totalBytes,
              });
              break;
            }

            case "Finished": {
              onProgress({
                status: "INSTALLING",
                percentage: 100,
                downloadedBytes: totalBytes || downloadedBytes,
                totalBytes: totalBytes || downloadedBytes,
              });
              break;
            }
          }
        });

        if (!isCancelled) {
          onProgress({
            status: "COMPLETED",
            percentage: 100,
            downloadedBytes: totalBytes || downloadedBytes,
            totalBytes: totalBytes || downloadedBytes,
          });
        }
      } catch (err: unknown) {
        if (!isCancelled) {
          console.error("Tauri update error:", err);
          let errorMsg = "Failed to download and install update.";
          if (typeof err === "string") {
            errorMsg = err;
          } else if (err instanceof Error) {
            errorMsg = err.message;
          } else if (err && typeof err === "object" && "message" in err) {
            errorMsg = String((err as { message: unknown }).message);
          }
          onError(errorMsg);
        }
      }
    })();

    return () => {
      isCancelled = true;
      if (cachedUpdate) {
        cachedUpdate.close().catch(() => {});
      }
    };
  },

  downloadStandaloneInstaller: (
    url: string,
    filename: string,
    onProgress: (data: UpdateDownloadProgressResponse) => void,
    onError: (errorMessage: string) => void
  ): (() => void) => {
    let isCancelled = false;
    let unlistenProgress: (() => void) | null = null;

    (async () => {
      try {
        unlistenProgress = await listen<DownloadProgressPayload>(
          "download-file-progress",
          (event) => {
            if (isCancelled) return;
            onProgress({
              status: "DOWNLOADING",
              percentage: event.payload.percentage,
              downloadedBytes: event.payload.downloaded_bytes,
              totalBytes: event.payload.total_bytes,
            });
          }
        );

        if (isCancelled) return;

        onProgress({
          status: "DOWNLOADING",
          percentage: 0,
          downloadedBytes: 0,
          totalBytes: 0,
        });

        const savedPath = await invoke<string>("download_release_file", {
          url,
          filename,
        });

        if (!isCancelled) {
          onProgress({
            status: "COMPLETED",
            percentage: 100,
            downloadedBytes: 0,
            totalBytes: 0,
            savedFilePath: savedPath,
          });
        }
      } catch (err: unknown) {
        if (!isCancelled) {
          console.error("Download installer error:", err);
          let errorMsg = "Failed to download installer file.";
          if (typeof err === "string") {
            errorMsg = err;
          } else if (err instanceof Error) {
            errorMsg = err.message;
          } else if (err && typeof err === "object" && "message" in err) {
            errorMsg = String((err as { message: unknown }).message);
          }
          onError(errorMsg);
        }
      } finally {
        if (unlistenProgress) {
          unlistenProgress();
        }
      }
    })();

    return () => {
      isCancelled = true;
      if (unlistenProgress) {
        unlistenProgress();
      }
    };
  },

  openDownloadedFile: async (path: string): Promise<void> => {
    await invoke("open_file_path", { path });
  },

  relaunchApp: async (): Promise<void> => {
    await relaunch();
  },
};
