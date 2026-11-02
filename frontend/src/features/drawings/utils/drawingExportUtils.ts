import {
  exportToSvg,
  exportToBlob,
  exportToClipboard,
  serializeAsJSON,
} from "@excalidraw/excalidraw";
import { invoke } from "@tauri-apps/api/core";
import type { NonDeletedExcalidrawElement } from "@excalidraw/excalidraw/element/types";
import type { AppState, BinaryFiles } from "@excalidraw/excalidraw/types";

/**
 * Fallback download for browser environment without Tauri IPC.
 */
function fallbackBrowserDownload(blob: Blob, filename: string): string {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  setTimeout(() => URL.revokeObjectURL(url), 2000);
  return filename;
}

/**
 * Generate a compact SVG data URI thumbnail for the diagram.
 * Returns null if the scene has no elements.
 */
export async function generateDrawingThumbnail(
  elements: readonly NonDeletedExcalidrawElement[],
  appState: Partial<AppState>,
  files: BinaryFiles
): Promise<string | null> {
  if (!elements || elements.length === 0) {
    return null;
  }

  try {
    const svg = await exportToSvg({
      elements,
      appState: {
        ...appState,
        exportBackground: true,
        viewBackgroundColor: appState.viewBackgroundColor || "#121316",
        exportWithDarkMode: true,
        exportEmbedScene: false,
      },
      files,
      exportPadding: 10,
    });

    if (!svg) return null;

    // Optional width/height capping for thumbnail performance
    svg.setAttribute("width", "320");
    svg.setAttribute("height", "200");
    svg.style.maxWidth = "100%";
    svg.style.height = "auto";

    const svgString = new XMLSerializer().serializeToString(svg);
    return `data:image/svg+xml;utf8,${encodeURIComponent(svgString)}`;
  } catch (err) {
    console.warn("[drawingExportUtils] Failed to generate thumbnail:", err);
    return null;
  }
}

/**
 * Copy diagram directly to system clipboard as PNG or SVG.
 */
export async function copyDrawingToClipboard(
  format: "png" | "svg",
  elements: readonly NonDeletedExcalidrawElement[],
  appState: Partial<AppState>,
  files: BinaryFiles
): Promise<void> {
  await exportToClipboard({
    elements,
    appState: {
      ...appState,
      exportBackground: true,
      viewBackgroundColor: appState.viewBackgroundColor || "#121316",
      exportWithDarkMode: appState.theme === "dark",
    },
    files,
    type: format,
    exportPadding: 20,
  });
}

/**
 * Prompts the user with the native desktop Save As dialog to choose the destination folder and filename.
 * Falls back to standard browser download if running outside the Tauri environment.
 * Returns the final path/filename saved, or null if the user cancelled the dialog.
 */
async function saveFileWithNativePicker(
  defaultFilename: string,
  extension: string,
  blob: Blob,
  textString?: string
): Promise<string | null> {
  // 1. Desktop (Tauri) environment: use native OS Save As dialog via Tauri IPC
  try {
    let fileBytes: number[] | null = null;
    if (!textString) {
      const buffer = await blob.arrayBuffer();
      fileBytes = Array.from(new Uint8Array(buffer));
    }

    const savedPath = await invoke<string | null>("save_diagram_file", {
      defaultFilename,
      extension,
      fileBytes,
      fileText: textString ?? null,
    });

    return savedPath;
  } catch (err) {
    console.warn("[drawingExportUtils] Tauri native dialog fallback:", err);
    // 2. Fallback to standard browser download (e.g. if running in web)
    return fallbackBrowserDownload(blob, defaultFilename);
  }
}

/**
 * Export diagram as PNG, SVG or Excalidraw JSON file with native OS Save As dialog
 * letting the user select their preferred destination folder and filename.
 * Returns the saved path/filename or null if cancelled.
 */
export async function exportDrawingToFile(
  format: "png" | "svg" | "json",
  elements: readonly NonDeletedExcalidrawElement[],
  appState: Partial<AppState>,
  files: BinaryFiles,
  title: string
): Promise<string | null> {
  const sanitizedTitle = (title.trim() || "untitled").replace(/[/\\?%*:|"<>]/g, "-");

  if (format === "json") {
    const jsonStr = serializeAsJSON(elements, appState, files, "local");
    const blob = new Blob([jsonStr], { type: "application/json" });
    return await saveFileWithNativePicker(
      `${sanitizedTitle}.excalidraw`,
      "excalidraw",
      blob,
      jsonStr
    );
  }

  if (format === "svg") {
    const svg = await exportToSvg({
      elements,
      appState: {
        ...appState,
        exportBackground: true,
        viewBackgroundColor: appState.viewBackgroundColor || "#121316",
        exportWithDarkMode: appState.theme === "dark",
      },
      files,
      exportPadding: 20,
    });
    const svgString = new XMLSerializer().serializeToString(svg);
    const blob = new Blob([svgString], { type: "image/svg+xml" });
    return await saveFileWithNativePicker(
      `${sanitizedTitle}.svg`,
      "svg",
      blob,
      svgString
    );
  }

  if (format === "png") {
    const blob = await exportToBlob({
      elements,
      appState: {
        ...appState,
        exportBackground: true,
        viewBackgroundColor: appState.viewBackgroundColor || "#121316",
        exportWithDarkMode: appState.theme === "dark",
      },
      files,
      mimeType: "image/png",
      exportPadding: 20,
    });

    if (blob) {
      return await saveFileWithNativePicker(
        `${sanitizedTitle}.png`,
        "png",
        blob
      );
    }
  }

  return null;
}
