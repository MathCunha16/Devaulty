import { apiClient, getInternalToken, getApiBaseUrl } from "../../../api/client";
import type {
  CurrentVersionResponse,
  AppUpdateInfoResponse,
  UpdateDownloadProgressResponse,
} from "../../../types/api";

export const releasesApi = {
  getCurrentVersion: async (): Promise<CurrentVersionResponse> => {
    const response = await apiClient.get<CurrentVersionResponse>("/releases/current-app-version");
    return response.data;
  },

  checkUpdates: async (): Promise<AppUpdateInfoResponse> => {
    const response = await apiClient.get<AppUpdateInfoResponse>("/releases/check");
    return response.data;
  },

  streamDownloadAndInstall: (
    onProgress: (data: UpdateDownloadProgressResponse) => void,
    onError: (errorMessage: string) => void
  ): (() => void) => {
    const controller = new AbortController();
    const baseUrl = getApiBaseUrl();
    const url = `${baseUrl}/releases/download-and-install`;

    const headers: Record<string, string> = {
      Accept: "text/event-stream",
    };
    const internalToken = getInternalToken();
    if (internalToken) {
      headers["X-Devaulty-Internal-Token"] = internalToken;
    }

    fetch(url, {
      method: "POST",
      headers,
      signal: controller.signal,
    })
      .then(async (response) => {
        if (!response.ok) {
          let errorMsg = `HTTP Error ${response.status}`;
          try {
            const errData = await response.json();
            if (errData?.message) errorMsg = errData.message;
          } catch {
            // ignore json parse error
          }
          onError(errorMsg);
          return;
        }

        if (!response.body) {
          onError("No stream body received from server");
          return;
        }

        const reader = response.body.getReader();
        const decoder = new TextDecoder();
        let buffer = "";

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          buffer += decoder.decode(value, { stream: true });
          const lines = buffer.split("\n");
          buffer = lines.pop() || "";

          for (const line of lines) {
            const trimmed = line.trim();
            if (trimmed.startsWith("data:")) {
              const jsonStr = trimmed.slice(5).trim();
              if (jsonStr) {
                try {
                  const data: UpdateDownloadProgressResponse = JSON.parse(jsonStr);
                  onProgress(data);
                } catch {
                  // ignore JSON parse error for malformed lines
                }
              }
            }
          }
        }
      })
      .catch((err) => {
        if (err.name !== "AbortError") {
          onError(err.message || "Failed to download update stream");
        }
      });

    return () => controller.abort();
  },
};
