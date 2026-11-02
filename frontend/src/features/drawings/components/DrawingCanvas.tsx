import React, { useMemo, useEffect, useRef, useCallback } from "react";
import { Excalidraw } from "@excalidraw/excalidraw";
import "@excalidraw/excalidraw/index.css";
import { toast } from "sonner";
import type {
  NonDeletedExcalidrawElement,
  FileId,
  ExcalidrawElement,
} from "@excalidraw/excalidraw/element/types";
import type {
  AppState,
  BinaryFiles,
  DataURL,
  ExcalidrawImperativeAPI,
  ExcalidrawInitialDataState,
} from "@excalidraw/excalidraw/types";
import type { Drawing } from "~types/api";

// Ensure local asset path for desktop / offline safety without using `any`
interface WindowWithExcalidraw extends Window {
  EXCALIDRAW_ASSET_PATH?: string;
}

if (typeof window !== "undefined") {
  const win = window as WindowWithExcalidraw;
  if (!win.EXCALIDRAW_ASSET_PATH) {
    win.EXCALIDRAW_ASSET_PATH = "/";
  }
}

interface DrawingCanvasProps {
  drawing: Drawing;
  onSceneChange: (
    elements: readonly NonDeletedExcalidrawElement[],
    appState: Partial<AppState>,
    files: BinaryFiles
  ) => void;
  excalidrawAPIRef: React.MutableRefObject<ExcalidrawImperativeAPI | null>;
}

let _isCanvasFilterBrokenCache: boolean | null = null;

/**
 * Detects if CanvasRenderingContext2D.filter is broken or unsupported in the current engine.
 * In WebKitGTK (WRY/Tauri on Linux) and older Safari, context.filter is ignored/unsupported,
 * so Excalidraw's dark mode CSS filter (invert(93%) hue-rotate(180deg)) renders images as
 * negative X-ray pictures without the internal counter-filter being applied.
 */
function isCanvasFilterBroken(): boolean {
  if (_isCanvasFilterBrokenCache !== null) {
    return _isCanvasFilterBrokenCache;
  }
  if (typeof document === "undefined") return false;
  try {
    const testCanvas = document.createElement("canvas");
    testCanvas.width = 2;
    testCanvas.height = 2;
    const ctx = testCanvas.getContext("2d");
    if (!ctx || typeof ctx.filter !== "string") {
      _isCanvasFilterBrokenCache = true;
      return true;
    }

    ctx.filter = "invert(100%)";
    ctx.fillStyle = "#ffffff";
    ctx.fillRect(0, 0, 1, 1);
    const pixel = ctx.getImageData(0, 0, 1, 1).data;
    // If filter works, white (255) inverted becomes black (0).
    // If pixel[0] is still 255, context.filter is a no-op!
    _isCanvasFilterBrokenCache = pixel[0] > 128;
    return _isCanvasFilterBrokenCache;
  } catch {
    _isCanvasFilterBrokenCache = true;
    return true;
  }
}

/**
 * Safely searches prototype chain for a property descriptor
 */
function findPropertyDescriptor(
  obj: unknown,
  prop: string
): { owner: unknown; desc: PropertyDescriptor } | null {
  let curr = obj;
  while (curr) {
    try {
      const desc = Object.getOwnPropertyDescriptor(curr, prop);
      if (desc) return { owner: curr, desc };
    } catch {
      break;
    }
    curr = Object.getPrototypeOf(curr);
  }
  return null;
}

// Enforce continuous anti-aliased rendering and minimum physical stroke thickness
// on WebKitGTK (WRY/Tauri on Linux) and Safari to prevent extreme zoom-out pixelation and broken lines
if (typeof window !== "undefined") {
  try {
    const testCanvas = document.createElement("canvas");
    const testCtx = testCanvas.getContext("2d");
    const searchTarget =
      testCtx ||
      (typeof CanvasRenderingContext2D !== "undefined"
        ? CanvasRenderingContext2D.prototype
        : null);

    if (searchTarget) {
      // 1. Force imageSmoothingEnabled = true
      // Excalidraw sets imageSmoothingEnabled = false when rendering cached element canvases,
      // which causes nearest-neighbor staircased pixelation during zoom out in WebKitGTK.
      const smoothInfo =
        findPropertyDescriptor(searchTarget, "imageSmoothingEnabled") ||
        (typeof CanvasRenderingContext2D !== "undefined"
          ? findPropertyDescriptor(CanvasRenderingContext2D.prototype, "imageSmoothingEnabled")
          : null);

      if (smoothInfo && smoothInfo.desc.set) {
        const originalSet = smoothInfo.desc.set;
        Object.defineProperty(smoothInfo.owner, "imageSmoothingEnabled", {
          set(_val: boolean) {
            void _val;
            // Always keep smoothing enabled for clean bilinear interpolation
            originalSet.call(this, true);
          },
          get: smoothInfo.desc.get,
          configurable: true,
          enumerable: true,
        });
      }

      // 2. Force imageSmoothingQuality = "high"
      const qualityInfo =
        findPropertyDescriptor(searchTarget, "imageSmoothingQuality") ||
        (typeof CanvasRenderingContext2D !== "undefined"
          ? findPropertyDescriptor(CanvasRenderingContext2D.prototype, "imageSmoothingQuality")
          : null);

      if (qualityInfo && qualityInfo.desc.set) {
        const originalSetQuality = qualityInfo.desc.set;
        Object.defineProperty(qualityInfo.owner, "imageSmoothingQuality", {
          set(_val: string) {
            void _val;
            originalSetQuality.call(this, "high");
          },
          get: qualityInfo.desc.get || (() => "high"),
          configurable: true,
          enumerable: true,
        });
      }

      // 3. Prevent lines from cutting in half / dropping out at extreme zoom out
      // In WebKitGTK (Cairo backend), when currentScale * lineWidth < 0.85px, subpixel alpha
      // coverage falls below the rasterization discard threshold, causing curved / hand-drawn
      // strokes to disappear in segments ("cut in half"). We clamp the physical stroke width to >= 0.85px.
      const lineWidthInfo =
        findPropertyDescriptor(searchTarget, "lineWidth") ||
        (typeof CanvasRenderingContext2D !== "undefined"
          ? findPropertyDescriptor(CanvasRenderingContext2D.prototype, "lineWidth")
          : null);

      if (lineWidthInfo && lineWidthInfo.desc.set) {
        const originalLineWidthSet = lineWidthInfo.desc.set;
        const originalLineWidthGet = lineWidthInfo.desc.get;
        Object.defineProperty(lineWidthInfo.owner, "lineWidth", {
          set(val: number) {
            try {
              if (val > 0 && typeof this.getTransform === "function") {
                const matrix = this.getTransform();
                const currentScale = Math.hypot(matrix.a, matrix.b);
                if (currentScale > 0) {
                  const physicalWidth = val * currentScale;
                  const minPhysicalWidth = 0.85;
                  if (physicalWidth < minPhysicalWidth) {
                    originalLineWidthSet.call(this, minPhysicalWidth / currentScale);
                    return;
                  }
                }
              }
            } catch {
              // Ignore matrix read errors
            }
            originalLineWidthSet.call(this, val);
          },
          get: originalLineWidthGet,
          configurable: true,
          enumerable: true,
        });
      }
    }
  } catch (err) {
    console.warn("[DrawingCanvas] Prototype optimization patch error:", err);
  }
}

/**
 * Pre-counter-inverts image pixels for WebKitGTK / Safari engines where context.filter is unsupported.
 * Matches Excalidraw's IMAGE_INVERT_FILTER ("invert(100%) hue-rotate(180deg) saturate(1.25)").
 * When Excalidraw's canvas CSS filter (invert(93%) hue-rotate(180deg)) is applied, this
 * cancels out, rendering photos with completely natural, true-to-life colors.
 */
function counterInvertImage(img: HTMLImageElement, mimeType: string): string {
  const canvas = document.createElement("canvas");
  const width = img.naturalWidth || img.width;
  const height = img.naturalHeight || img.height;
  if (width <= 0 || height <= 0) return img.src;

  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext("2d");
  if (!ctx) return img.src;

  ctx.drawImage(img, 0, 0);
  const imgData = ctx.getImageData(0, 0, width, height);
  const data = imgData.data;

  for (let i = 0; i < data.length; i += 4) {
    // Preserve transparent pixels
    if (data[i + 3] === 0) continue;

    const r = data[i];
    const g = data[i + 1];
    const b = data[i + 2];

    // 1. Invert (100%)
    const ir = 255 - r;
    const ig = 255 - g;
    const ib = 255 - b;

    // 2. Hue rotate 180° matrix
    const hr = -0.574 * ir + 1.430 * ig + 0.144 * ib;
    const hg =  0.426 * ir + 0.430 * ig + 0.144 * ib;
    const hb =  0.426 * ir + 1.430 * ig - 0.856 * ib;

    // 3. Saturate 1.25
    const gray = 0.2126 * hr + 0.7152 * hg + 0.0722 * hb;
    const satR = gray + 1.25 * (hr - gray);
    const satG = gray + 1.25 * (hg - gray);
    const satB = gray + 1.25 * (hb - gray);

    data[i] = Math.max(0, Math.min(255, Math.round(satR)));
    data[i + 1] = Math.max(0, Math.min(255, Math.round(satG)));
    data[i + 2] = Math.max(0, Math.min(255, Math.round(satB)));
  }

  ctx.putImageData(imgData, 0, 0);
  return canvas.toDataURL(mimeType.startsWith("image/png") ? "image/png" : "image/jpeg", 0.95);
}

/**
 * Asynchronously converts raw dataURL into counter-inverted dataURL
 */
function counterInvertDataURL(dataURL: string, mimeType: string): Promise<string> {
  return new Promise((resolve) => {
    const img = new Image();
    img.crossOrigin = "anonymous";
    img.onload = () => {
      try {
        const inverted = counterInvertImage(img, mimeType);
        resolve(inverted);
      } catch (e) {
        console.warn("[DrawingCanvas] Error counter-inverting image:", e);
        resolve(dataURL);
      }
    };
    img.onerror = (e) => {
      console.warn("[DrawingCanvas] Failed to load image for counter-inversion:", e);
      resolve(dataURL);
    };
    img.src = dataURL;
  });
}

export const DrawingCanvas: React.FC<DrawingCanvasProps> = ({
  drawing,
  onSceneChange,
  excalidrawAPIRef,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const invertedMapRef = useRef<Map<string, string>>(new Map());
  const inFlightInversionsRef = useRef<Set<string>>(new Set());
  const isUpdatingSceneRef = useRef(false);

  // Parse initial scene data safely, recovering from local crash-draft if present
  const initialData = useMemo<ExcalidrawInitialDataState>(() => {
    let rawContent = drawing.content;
    try {
      const draft = localStorage.getItem(`devaulty_drawing_draft_${drawing.id}`);
      if (draft) {
        rawContent = draft;
      }
    } catch {
      // Ignore storage access errors
    }

    if (!rawContent) {
      return {
        elements: [],
        appState: {
          theme: "dark",
          viewBackgroundColor: "#121316",
          shouldCacheIgnoreZoom: true,
        },
      };
    }

    try {
      const parsed = JSON.parse(rawContent);
      return {
        elements: parsed.elements || [],
        appState: {
          ...(parsed.appState || {}),
          theme: "dark",
          viewBackgroundColor: parsed.appState?.viewBackgroundColor || "#121316",
          shouldCacheIgnoreZoom: true,
        },
        files: parsed.files || {},
      };
    } catch (err) {
      console.warn("[DrawingCanvas] Error parsing drawing content:", err);
      return {
        elements: [],
        appState: {
          theme: "dark",
          viewBackgroundColor: "#121316",
          shouldCacheIgnoreZoom: true,
        },
      };
    }
  }, [drawing.id, drawing.content]);

  // Clean-up API reference on unmount
  useEffect(() => {
    return () => {
      excalidrawAPIRef.current = null;
    };
  }, [excalidrawAPIRef]);

  // Fix laser pointer & cursor alignment: recalculate canvas rect when container resizes
  useEffect(() => {
    if (!containerRef.current) return;

    const ro = new ResizeObserver(() => {
      if (excalidrawAPIRef.current) {
        excalidrawAPIRef.current.refresh();
      }
    });

    ro.observe(containerRef.current);
    return () => ro.disconnect();
  }, [excalidrawAPIRef]);

  // Inverts an image file's dataURL and heals any matching elements on the Excalidraw canvas
  const invertFileAndHealScene = useCallback(
    (fileId: string, dataURL: string, mimeType: string) => {
      if (inFlightInversionsRef.current.has(fileId)) return;
      inFlightInversionsRef.current.add(fileId);

      counterInvertDataURL(dataURL, mimeType)
        .then((invertedDataURL) => {
          inFlightInversionsRef.current.delete(fileId);

          const invertedFileId = `devaulty_inv_${fileId}` as unknown as FileId;
          invertedMapRef.current.set(fileId, invertedFileId as unknown as string);

          const activeApi = excalidrawAPIRef.current;
          if (!activeApi) return;

          // Register inverted file AND update original fileId in Excalidraw store
          activeApi.addFiles([
            {
              id: invertedFileId,
              dataURL: invertedDataURL as unknown as DataURL,
              mimeType: mimeType as unknown as "image/png",
              created: Date.now(),
            },
            {
              id: fileId as unknown as FileId,
              dataURL: invertedDataURL as unknown as DataURL,
              mimeType: mimeType as unknown as "image/png",
              created: Date.now(),
            },
          ]);

          // Immediately heal any elements currently on the scene with this fileId
          const currentElements = activeApi.getSceneElements();
          let elementChanged = false;
          const updated = currentElements.map((el) => {
            if (el.type === "image" && (el.fileId === fileId || el.fileId === invertedFileId)) {
              if (el.fileId !== invertedFileId) {
                elementChanged = true;
                return {
                  ...el,
                  fileId: invertedFileId,
                  version: (el.version || 1) + 1,
                  versionNonce: Math.floor(Math.random() * 100000),
                };
              }
            }
            return el;
          });

          if (elementChanged) {
            isUpdatingSceneRef.current = true;
            try {
              activeApi.updateScene({ elements: updated });
            } finally {
              isUpdatingSceneRef.current = false;
            }
          }
        })
        .catch((err) => {
          inFlightInversionsRef.current.delete(fileId);
          console.warn("[DrawingCanvas] Error counter-inverting image:", err);
        });
    },
    [excalidrawAPIRef]
  );

  // Intercept and auto-heal imported images (from Excalidraw toolbar Image tool, drag-drop, or file load)
  const handleSceneChange = useCallback(
    (
      elements: readonly NonDeletedExcalidrawElement[],
      appState: Partial<AppState>,
      files: BinaryFiles
    ) => {
      // 1. Forward change event to parent debounced save hook
      onSceneChange(elements, appState, files);

      // 2. Fast check: if engine doesn't need counter-inversion or scene has no images/files, exit immediately
      if (!isCanvasFilterBroken()) {
        return;
      }

      const hasImages = elements.some((el) => el.type === "image");
      const hasFiles = Object.keys(files || {}).length > 0;
      if (!hasImages && !hasFiles && invertedMapRef.current.size === 0) {
        return;
      }

      const api = excalidrawAPIRef.current;
      if (!api) return;

      // Avoid re-entry loops while updateScene is executing
      if (isUpdatingSceneRef.current) {
        return;
      }

      const allFiles: BinaryFiles = {
        ...(api.getFiles() || {}),
        ...(files || {}),
      };

      // 3. Immediately heal any elements on the scene whose fileId already has an inverted counterpart
      const sceneElements = api.getSceneElements();
      let needsSceneUpdate = false;
      const healedElements = sceneElements.map((el) => {
        if (el.type === "image" && el.fileId) {
          const targetInvertedId = invertedMapRef.current.get(el.fileId);
          if (targetInvertedId && el.fileId !== targetInvertedId) {
            needsSceneUpdate = true;
            return {
              ...el,
              fileId: targetInvertedId as unknown as FileId,
              version: (el.version || 1) + 1,
              versionNonce: Math.floor(Math.random() * 100000),
            };
          }
        }
        return el;
      });

      if (needsSceneUpdate) {
        isUpdatingSceneRef.current = true;
        try {
          api.updateScene({ elements: healedElements });
        } finally {
          isUpdatingSceneRef.current = false;
        }
      }

      // 4. Identify any image file in allFiles that needs counter-inversion
      for (const [fileId, fileData] of Object.entries(allFiles)) {
        if (!fileData || !fileData.dataURL) continue;
        if (fileId.startsWith("devaulty_inv_")) continue;
        if (invertedMapRef.current.has(fileId)) continue;
        if (inFlightInversionsRef.current.has(fileId)) continue;

        const mime = fileData.mimeType || "";
        if (mime.startsWith("image/") && !mime.includes("svg")) {
          invertFileAndHealScene(fileId, fileData.dataURL, mime);
        }
      }

      // 5. Check if any image element currently in the scene references an un-inverted file
      for (const el of sceneElements) {
        if (el.type === "image" && el.fileId) {
          const fid = el.fileId as string;
          if (
            !fid.startsWith("devaulty_inv_") &&
            !invertedMapRef.current.has(fid) &&
            !inFlightInversionsRef.current.has(fid)
          ) {
            const rawFile = allFiles[fid];
            if (rawFile && rawFile.dataURL) {
              const mime = rawFile.mimeType || "image/png";
              invertFileAndHealScene(fid, rawFile.dataURL, mime);
            }
          }
        }
      }
    },
    [onSceneChange, excalidrawAPIRef, invertFileAndHealScene]
  );

  // Robust image paste & drag-drop handler: processes screenshots, files, and browser images
  useEffect(() => {
    const insertImageFile = (
      imageBlob: Blob,
      mimeType: string,
      clientPos?: { x: number; y: number }
    ) => {
      const api = excalidrawAPIRef.current;
      if (!api) return;

      const reader = new FileReader();
      reader.onload = () => {
        const rawDataURL = reader.result as string;
        const fileId = `devaulty_inv_${Date.now()}_${Math.random()
          .toString(36)
          .substring(2, 9)}` as unknown as FileId;

        invertedMapRef.current.set(fileId as unknown as string, fileId as unknown as string);

        const img = new Image();
        img.onload = () => {
          // Counter-invert image if WebKitGTK / Safari context.filter is broken in dark mode
          let finalDataURL = rawDataURL;
          if (isCanvasFilterBroken()) {
            finalDataURL = counterInvertImage(img, mimeType);
          }

          const currentAppState = api.getAppState();
          const naturalWidth = img.naturalWidth || img.width || 400;
          const naturalHeight = img.naturalHeight || img.height || 300;
          const maxWidth = Math.min(naturalWidth, 640);
          const width = maxWidth;
          const height = Math.round((naturalHeight / naturalWidth) * width);

          let sceneX: number;
          let sceneY: number;

          if (clientPos && containerRef.current) {
            const rect = containerRef.current.getBoundingClientRect();
            const viewportX = clientPos.x - rect.left;
            const viewportY = clientPos.y - rect.top;
            sceneX = Math.round(
              (viewportX - currentAppState.scrollX) / currentAppState.zoom.value - width / 2
            );
            sceneY = Math.round(
              (viewportY - currentAppState.scrollY) / currentAppState.zoom.value - height / 2
            );
          } else {
            sceneX = Math.round(
              -currentAppState.scrollX + currentAppState.width / 2 - width / 2
            );
            sceneY = Math.round(
              -currentAppState.scrollY + currentAppState.height / 2 - height / 2
            );
          }

          // Register binary image file in Excalidraw store
          api.addFiles([
            {
              id: fileId,
              dataURL: finalDataURL as unknown as DataURL,
              mimeType: mimeType as unknown as "image/png",
              created: Date.now(),
            },
          ]);

          const imageElement: ExcalidrawElement = {
            type: "image",
            id: `img_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`,
            status: "saved",
            fileId,
            x: sceneX,
            y: sceneY,
            width,
            height,
            angle: 0,
            strokeColor: "transparent",
            backgroundColor: "transparent",
            fillStyle: "hachure",
            strokeWidth: 1,
            strokeStyle: "solid",
            roughness: 0,
            opacity: 100,
            groupIds: [],
            frameId: null,
            roundness: null,
            seed: Math.floor(Math.random() * 100000),
            version: 1,
            versionNonce: Math.floor(Math.random() * 100000),
            isDeleted: false,
            boundElements: null,
            updated: Date.now(),
            link: null,
            locked: false,
            scale: [1, 1],
            crop: null,
          } as unknown as ExcalidrawElement;

          const existingElements = api.getSceneElements();
          api.updateScene({
            elements: [...existingElements, imageElement],
            appState: {
              selectedElementIds: { [imageElement.id]: true },
            },
          });
          toast.success("Image added to canvas");
        };
        img.src = rawDataURL;
      };
      reader.readAsDataURL(imageBlob);
    };

    const handlePaste = async (e: ClipboardEvent) => {
      // If user is typing in an input, textarea or contenteditable element, let normal paste happen
      const activeEl = document.activeElement;
      if (
        activeEl &&
        (activeEl.tagName === "INPUT" ||
          activeEl.tagName === "TEXTAREA" ||
          activeEl.getAttribute("contenteditable") === "true" ||
          activeEl.classList.contains("excalidraw-wysiwyg"))
      ) {
        return;
      }

      const api = excalidrawAPIRef.current;
      if (!api) return;

      // Don't intercept if user is editing a text element on canvas
      const appState = api.getAppState();
      if (appState.editingTextElement) {
        return;
      }

      let imageBlob: Blob | null = null;
      let mimeType = "image/png";

      // 1. Check clipboard items (standard for browser Copy Image and screenshot tools)
      if (e.clipboardData?.items) {
        for (const item of Array.from(e.clipboardData.items)) {
          if (item.type.startsWith("image/")) {
            imageBlob = item.getAsFile();
            mimeType = item.type;
            if (imageBlob) break;
          }
        }
      }

      // 2. Check clipboard files (standard for file manager copy)
      if (!imageBlob && e.clipboardData?.files && e.clipboardData.files.length > 0) {
        const file = Array.from(e.clipboardData.files).find((f) =>
          f.type.startsWith("image/")
        );
        if (file) {
          imageBlob = file;
          mimeType = file.type;
        }
      }

      // 3. Fallback to navigator.clipboard.read()
      if (!imageBlob && navigator.clipboard && navigator.clipboard.read) {
        try {
          const items = await navigator.clipboard.read();
          for (const item of items) {
            const foundType = item.types.find((t) => t.startsWith("image/"));
            if (foundType) {
              imageBlob = await item.getType(foundType);
              mimeType = foundType;
              break;
            }
          }
        } catch {
          // Ignore clipboard permission errors silently
        }
      }

      if (!imageBlob) return;

      e.preventDefault();
      e.stopPropagation();
      insertImageFile(imageBlob, mimeType);
    };

    const handleDrop = (e: DragEvent) => {
      if (!e.dataTransfer) return;
      const file = Array.from(e.dataTransfer.files).find((f) =>
        f.type.startsWith("image/")
      );
      if (file) {
        e.preventDefault();
        e.stopPropagation();
        insertImageFile(file, file.type, { x: e.clientX, y: e.clientY });
      }
    };

    const handleDragOver = (e: DragEvent) => {
      if (e.dataTransfer?.types.includes("Files")) {
        e.preventDefault();
      }
    };

    window.addEventListener("paste", handlePaste, true);
    window.addEventListener("drop", handleDrop, true);
    window.addEventListener("dragover", handleDragOver, true);
    return () => {
      window.removeEventListener("paste", handlePaste, true);
      window.removeEventListener("drop", handleDrop, true);
      window.removeEventListener("dragover", handleDragOver, true);
    };
  }, [excalidrawAPIRef]);

  return (
    <div
      ref={containerRef}
      tabIndex={-1}
      onClick={() => {
        // Ensure inner Excalidraw canvas container is focused
        const el = containerRef.current?.querySelector(".excalidraw") as HTMLElement | null;
        el?.focus();
      }}
      className="flex-1 w-full h-full relative overflow-hidden select-none bg-[#121316] outline-none"
    >
      <Excalidraw
        key={drawing.id}
        initialData={initialData}
        excalidrawAPI={(api) => {
          excalidrawAPIRef.current = api;
          api.updateScene({
            appState: {
              shouldCacheIgnoreZoom: true,
            },
          });
          // Refresh bounding box after mount
          setTimeout(() => api.refresh(), 50);
        }}
        theme="dark"
        autoFocus={true}
        handleKeyboardGlobally={true}
        onChange={handleSceneChange}
        UIOptions={{
          canvasActions: {
            changeViewBackgroundColor: true,
            clearCanvas: true,
            export: false,
            loadScene: false,
            saveToActiveFile: false,
            toggleTheme: false,
          },
        }}
      />
    </div>
  );
};

export default DrawingCanvas;
