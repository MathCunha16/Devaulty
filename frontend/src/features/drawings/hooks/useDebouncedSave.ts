import { useState, useRef, useEffect, useCallback } from "react";
import type { NonDeletedExcalidrawElement } from "@excalidraw/excalidraw/element/types";
import type { AppState, BinaryFiles } from "@excalidraw/excalidraw/types";
import { generateDrawingThumbnail } from "../utils/drawingExportUtils";

export type SaveStatus = "saved" | "unsaved" | "saving" | "error";

interface UseDebouncedSaveOptions {
  drawingId: string;
  initialContent: string | null;
  onSave: (payload: { content: string; thumbnail?: string | null }) => Promise<void>;
  debounceMs?: number;
}

function normalizeSceneContent(raw: string | null): string {
  if (!raw || !raw.trim()) {
    return JSON.stringify({
      elements: [],
      appState: {
        viewBackgroundColor: "#121316",
      },
      files: {},
    });
  }
  try {
    const parsed = JSON.parse(raw);
    return JSON.stringify({
      elements: parsed.elements || [],
      appState: {
        viewBackgroundColor: parsed.appState?.viewBackgroundColor || "#121316",
        gridSize: parsed.appState?.gridSize,
      },
      files: parsed.files || {},
    });
  } catch {
    return raw;
  }
}

function hasElementsTrulyChanged(
  prev: readonly NonDeletedExcalidrawElement[] | null,
  next: readonly NonDeletedExcalidrawElement[]
): boolean {
  if (!prev) return true;
  if (prev === next) return false;
  if (prev.length !== next.length) return true;
  for (let i = 0; i < next.length; i++) {
    const p = prev[i];
    const n = next[i];
    if (p.id !== n.id || p.version !== n.version || p.versionNonce !== n.versionNonce) {
      return true;
    }
  }
  return false;
}

export function useDebouncedSave({
  drawingId,
  initialContent,
  onSave,
  debounceMs = 1500,
}: UseDebouncedSaveOptions) {
  const [status, setStatus] = useState<SaveStatus>("saved");
  const [prevSyncKey, setPrevSyncKey] = useState(`${drawingId}_${initialContent || ""}`);

  // Adjust state during render when switching drawing (official React pattern)
  const currentSyncKey = `${drawingId}_${initialContent || ""}`;
  if (currentSyncKey !== prevSyncKey) {
    setPrevSyncKey(currentSyncKey);
    setStatus("saved");
  }

  // Keep track of last saved serialization and fast element comparison refs
  const lastSavedContentRef = useRef<string>(normalizeSceneContent(initialContent));
  const lastElementsRef = useRef<readonly NonDeletedExcalidrawElement[] | null>(null);
  const lastFilesRef = useRef<BinaryFiles | null>(null);
  const lastBgRef = useRef<string | undefined>(undefined);
  const lastGridRef = useRef<number | undefined>(undefined);

  const pendingPayloadRef = useRef<{
    elements: readonly NonDeletedExcalidrawElement[];
    appState: Partial<AppState>;
    files: BinaryFiles;
  } | null>(null);

  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const draftTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isSavingRef = useRef(false);

  // Synchronize refs and timers when switching to a different drawing or when initialContent updates
  useEffect(() => {
    lastSavedContentRef.current = normalizeSceneContent(initialContent);
    lastElementsRef.current = null;
    lastFilesRef.current = null;
    lastBgRef.current = undefined;
    lastGridRef.current = undefined;
    pendingPayloadRef.current = null;

    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    if (draftTimerRef.current) {
      clearTimeout(draftTimerRef.current);
      draftTimerRef.current = null;
    }
  }, [drawingId, initialContent]);

  const commitSave = useCallback(async () => {
    if (!pendingPayloadRef.current || isSavingRef.current || !drawingId) {
      return;
    }

    if (draftTimerRef.current) {
      clearTimeout(draftTimerRef.current);
      draftTimerRef.current = null;
    }

    const { elements, appState, files } = pendingPayloadRef.current;
    const contentString = JSON.stringify({
      elements,
      appState: {
        viewBackgroundColor: appState.viewBackgroundColor,
        gridSize: appState.gridSize,
      },
      files,
    });

    // Check if truly modified
    if (contentString === lastSavedContentRef.current) {
      pendingPayloadRef.current = null;
      try {
        localStorage.removeItem(`devaulty_drawing_draft_${drawingId}`);
      } catch {
        // Ignore localStorage errors
      }
      setStatus("saved");
      return;
    }

    try {
      isSavingRef.current = true;
      setStatus("saving");

      // Generate compact thumbnail on save
      const thumbnail = await generateDrawingThumbnail(elements, appState, files);

      await onSave({
        content: contentString,
        thumbnail,
      });

      lastSavedContentRef.current = contentString;
      pendingPayloadRef.current = null;
      try {
        localStorage.removeItem(`devaulty_drawing_draft_${drawingId}`);
      } catch {
        // Ignore localStorage errors
      }
      setStatus("saved");
    } catch (err) {
      console.error("[useDebouncedSave] Auto-save error:", err);
      setStatus("error");
    } finally {
      isSavingRef.current = false;
    }
  }, [drawingId, onSave]);

  const triggerChange = useCallback(
    (
      elements: readonly NonDeletedExcalidrawElement[],
      appState: Partial<AppState>,
      files: BinaryFiles
    ) => {
      // Fast path: check if actual elements, files, background color, or grid changed
      const elementsChanged = hasElementsTrulyChanged(lastElementsRef.current, elements);
      const filesChanged = lastFilesRef.current !== files;
      const bgChanged = lastBgRef.current !== appState.viewBackgroundColor;
      const gridChanged = lastGridRef.current !== appState.gridSize;

      if (!elementsChanged && !filesChanged && !bgChanged && !gridChanged) {
        // Viewport-only changes (scrollX, scrollY, zoom, hover) - return immediately with ZERO overhead!
        return;
      }

      // Update change tracking references
      lastElementsRef.current = elements;
      lastFilesRef.current = files;
      lastBgRef.current = appState.viewBackgroundColor;
      lastGridRef.current = appState.gridSize;

      // Store in memory without serializing yet
      pendingPayloadRef.current = { elements, appState, files };
      setStatus((prev) => (prev !== "unsaved" ? "unsaved" : prev));

      // Tier 1: Debounced local draft snapshot (500ms) to avoid blocking main thread during drawing
      if (drawingId) {
        if (draftTimerRef.current) {
          clearTimeout(draftTimerRef.current);
        }
        draftTimerRef.current = setTimeout(() => {
          try {
            const draftJson = JSON.stringify({
              elements,
              appState: {
                viewBackgroundColor: appState.viewBackgroundColor,
                gridSize: appState.gridSize,
              },
              files,
            });
            localStorage.setItem(`devaulty_drawing_draft_${drawingId}`, draftJson);
          } catch {
            // Ignore storage errors silently
          }
        }, 500);
      }

      // Tier 2: Debounced backend persistence to SQLite
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }

      timerRef.current = setTimeout(() => {
        commitSave();
      }, debounceMs);
    },
    [commitSave, debounceMs, drawingId]
  );

  const saveNow = useCallback(async () => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
      timerRef.current = null;
    }
    if (draftTimerRef.current) {
      clearTimeout(draftTimerRef.current);
      draftTimerRef.current = null;
    }
    await commitSave();
  }, [commitSave]);

  // Clean-up and flush on unmount
  useEffect(() => {
    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      if (draftTimerRef.current) {
        clearTimeout(draftTimerRef.current);
      }
    };
  }, []);

  return {
    status,
    triggerChange,
    saveNow,
    isDirty: status === "unsaved" || status === "saving",
  };
}
