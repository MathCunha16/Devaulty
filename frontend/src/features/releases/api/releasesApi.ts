import { check, Update } from "@tauri-apps/plugin-updater";
import { relaunch } from "@tauri-apps/plugin-process";
import pkg from "../../../../package.json";
import type {
  CurrentVersionResponse,
  AppUpdateInfoResponse,
  UpdateDownloadProgressResponse,
} from "../../../types/api";

let cachedUpdate: Update | null = null;

export const releasesApi = {
  getCurrentVersion: async (): Promise<CurrentVersionResponse> => {
    const version = pkg.version;
    return {
      currentVersion: version,
      actualVersion: version,
    };
  },

  checkUpdates: async (): Promise<AppUpdateInfoResponse> => {
    try {
      const currentVersion = (await releasesApi.getCurrentVersion()).currentVersion;
      const update = await check();
      cachedUpdate = update;

      if (update && update.available) {
        return {
          updateAvailable: true,
          currentVersion: update.currentVersion || currentVersion,
          latestVersion: update.version,
          releaseTitle: `Release v${update.version}`,
          releaseNotes: update.body || "",
          publishedAt: update.date,
        };
      }

      return {
        updateAvailable: false,
        currentVersion,
        latestVersion: currentVersion,
      };
    } catch (err: unknown) {
      console.error("Error checking for updates via Tauri plugin:", err);
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
          const errorMsg =
            err instanceof Error ? err.message : "Failed to download and install update.";
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

  relaunchApp: async (): Promise<void> => {
    await relaunch();
  },
};
