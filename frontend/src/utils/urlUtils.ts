import { invoke } from "@tauri-apps/api/core";

/**
 * Opens an external URL in the system's default browser.
 * In desktop (Tauri) environment, invokes the Rust native command `open_external_url`.
 * In browser/dev environment, falls back to window.open.
 */
export const openExternalUrl = async (rawUrl?: string | null): Promise<boolean> => {
  if (!rawUrl || !rawUrl.trim()) return false;

  const trimmed = rawUrl.trim();
  const match = trimmed.match(/^https?:\/\//i);
  const url = match
    ? `${match[0].toLowerCase()}${trimmed.slice(match[0].length)}`
    : `https://${trimmed}`;

  try {
    await invoke("open_external_url", { url });
    return true;
  } catch {
    try {
      window.open(url, "_blank", "noopener,noreferrer");
      return true;
    } catch (fallbackErr) {
      console.error("Failed to open external URL:", fallbackErr);
      return false;
    }
  }
};
