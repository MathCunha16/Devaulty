import React, { useState, useEffect, useRef, Suspense, useMemo } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import type { ExcalidrawImperativeAPI } from "@excalidraw/excalidraw/types";
import type { DrawingSummary } from "~types/api";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { useResizableSidebar } from "../../../hooks/useResizableSidebar";
import styles from "../../../routes/projects.$projectId.module.css";
import {
  useDrawingsQuery,
  useDrawingQuery,
  useCreateDrawingMutation,
  useUpdateDrawingMutation,
  useDeleteDrawingMutation,
} from "../hooks/useDrawings";
import { useDebouncedSave } from "../hooks/useDebouncedSave";
import { exportDrawingToFile, copyDrawingToClipboard } from "../utils/drawingExportUtils";
import { DrawingListSidebar } from "./DrawingListSidebar";
import { DrawingHeader } from "./DrawingHeader";
import { DrawingCanvasSkeleton } from "./DrawingCanvasSkeleton";

// Lazy-load the Excalidraw canvas component for optimal initial mount & memory
const DrawingCanvas = React.lazy(() => import("./DrawingCanvas"));

interface DrawingsWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal: () => void;
  initialSelectedId?: string;
  projectColor?: string;
  isActive?: boolean;
}

export const DrawingsWorkspace: React.FC<DrawingsWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
  initialSelectedId,
  projectColor = "#10b981",
  isActive = true,
}) => {
  const { data: drawingsData } = useDrawingsQuery(projectId);
  const drawingsList = useMemo(() => drawingsData?.content || [], [drawingsData?.content]);

  // Track initial selected ID from URL if provided
  const [selectedDrawingId, setSelectedDrawingId] = useState<string | undefined>(
    initialSelectedId
  );
  const [prevInitialId, setPrevInitialId] = useState(initialSelectedId);

  if (initialSelectedId !== prevInitialId) {
    setPrevInitialId(initialSelectedId);
    setSelectedDrawingId(initialSelectedId);
  }

  // Pure declarative derivation for active diagram: selects clicked diagram or first available
  const activeDrawingId =
    selectedDrawingId && drawingsList.some((d) => d.id === selectedDrawingId)
      ? selectedDrawingId
      : drawingsList.length > 0
        ? drawingsList[0].id
        : undefined;

  const [searchQuery, setSearchQuery] = useState("");
  const [drawingToDelete, setDrawingToDelete] = useState<DrawingSummary | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);

  const {
    sidebarWidth,
    isCollapsed,
    isDragging,
    minWidth,
    maxWidth,
    toggleCollapse,
    handleMouseDown,
    handleDoubleClick,
    handleKeyDown,
  } = useResizableSidebar({
    storageKey: "drawings",
    defaultWidth: 320,
    minWidth: 240,
    maxWidth: 520,
  });

  // Mutations
  const createMutation = useCreateDrawingMutation(projectId);
  const updateMutation = useUpdateDrawingMutation(projectId, activeDrawingId || "");
  const deleteMutation = useDeleteDrawingMutation(projectId);

  // Query selected drawing details
  const { data: currentDrawing, isLoading: isLoadingCurrent } = useDrawingQuery(
    projectId,
    activeDrawingId || ""
  );

  // Excalidraw API ref for direct exports & commands
  const excalidrawAPIRef = useRef<ExcalidrawImperativeAPI | null>(null);

  // Debounced auto-save hook with instant local crash recovery
  const {
    status: saveStatus,
    triggerChange,
    saveNow,
  } = useDebouncedSave({
    drawingId: activeDrawingId || "",
    initialContent: currentDrawing?.content || null,
    onSave: async (payload) => {
      if (!activeDrawingId) return;
      await updateMutation.mutateAsync(payload);
    },
    debounceMs: 1500,
  });

  // Keep laser pointer and cursor perfectly aligned when sidebar collapses, resizes, or tab becomes active
  useEffect(() => {
    if (isActive && excalidrawAPIRef.current) {
      excalidrawAPIRef.current.refresh();
    }
  }, [sidebarWidth, isCollapsed, isActive]);

  // Global Ctrl+S / Cmd+S shortcut to force immediate save
  useEffect(() => {
    const handleSaveShortcut = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "s") {
        e.preventDefault();
        saveNow();
      }
    };
    window.addEventListener("keydown", handleSaveShortcut);
    return () => window.removeEventListener("keydown", handleSaveShortcut);
  }, [saveNow]);

  // Create new drawing
  const handleCreateDrawing = async () => {
    try {
      const newDrawing = await createMutation.mutateAsync({
        title: "Untitled Diagram",
      });
      setSelectedDrawingId(newDrawing.id);
      toast.success("New diagram created");
    } catch {
      toast.error("Failed to create diagram");
    }
  };

  // Rename drawing
  const handleRenameDrawing = async (newTitle: string) => {
    if (!activeDrawingId) return;
    try {
      await updateMutation.mutateAsync({ title: newTitle });
      toast.success("Diagram renamed");
    } catch {
      toast.error("Failed to rename diagram");
    }
  };

  // Delete drawing flow
  const handleDeleteClick = (id: string) => {
    const target = drawingsList.find((d) => d.id === id);
    if (target) {
      setDrawingToDelete(target);
    }
  };

  const handleConfirmDelete = async () => {
    if (!drawingToDelete) return;
    try {
      setIsDeleting(true);
      await deleteMutation.mutateAsync(drawingToDelete.id);
      if (activeDrawingId === drawingToDelete.id) {
        const remaining = drawingsList.filter((d) => d.id !== drawingToDelete.id);
        setSelectedDrawingId(remaining.length > 0 ? remaining[0].id : undefined);
      }
      toast.success("Diagram deleted");
      setDrawingToDelete(null);
    } catch {
      toast.error("Failed to delete diagram");
    } finally {
      setIsDeleting(false);
    }
  };

  // Export diagram to file
  const handleExport = async (format: "png" | "svg" | "json") => {
    if (!excalidrawAPIRef.current || !currentDrawing) return;
    try {
      const elements = excalidrawAPIRef.current.getSceneElements();
      const appState = excalidrawAPIRef.current.getAppState();
      const files = excalidrawAPIRef.current.getFiles();

      const savedPath = await exportDrawingToFile(
        format,
        elements,
        appState,
        files,
        currentDrawing.title
      );

      // If user cancelled the dialog, return silently without error
      if (!savedPath) {
        return;
      }

      const filename = savedPath.split("/").pop() || savedPath;
      toast.success(`Diagram saved as "${filename}"`);
    } catch {
      toast.error(`Failed to export diagram as ${format.toUpperCase()}`);
    }
  };

  // Copy diagram to clipboard
  const handleCopyToClipboard = async (format: "png" | "svg") => {
    if (!excalidrawAPIRef.current || !currentDrawing) return;
    try {
      const elements = excalidrawAPIRef.current.getSceneElements();
      const appState = excalidrawAPIRef.current.getAppState();
      const files = excalidrawAPIRef.current.getFiles();

      await copyDrawingToClipboard(format, elements, appState, files);
      toast.success(`${format.toUpperCase()} copied to clipboard!`);
    } catch {
      toast.error(`Failed to copy ${format.toUpperCase()} to clipboard`);
    }
  };

  // Open native Excalidraw image export modal dialog
  const handleOpenImageDialog = () => {
    if (!excalidrawAPIRef.current) return;
    excalidrawAPIRef.current.updateScene({
      appState: {
        openDialog: { name: "imageExport" },
      },
    });
  };

  return (
    <>
      {/* Left Panel: Diagrams list navigator */}
      <div
        className={`${styles.leftPanel} ${isCollapsed ? styles.leftPanelCollapsed : ""}`}
        style={{ width: isCollapsed ? 0 : `${sidebarWidth}px` }}
      >
        <DrawingListSidebar
          drawings={drawingsList}
          selectedId={activeDrawingId}
          searchQuery={searchQuery}
          projectColor={projectColor}
          isCreating={createMutation.isPending}
          onSelect={setSelectedDrawingId}
          onCreate={handleCreateDrawing}
          onDelete={handleDeleteClick}
          onSearchChange={setSearchQuery}
          toggleCollapse={toggleCollapse}
        />
      </div>

      {/* Resizer Handle */}
      <div
        role="separator"
        tabIndex={0}
        aria-orientation="vertical"
        aria-label="Resize diagram sidebar"
        aria-valuenow={sidebarWidth}
        aria-valuemin={minWidth}
        aria-valuemax={maxWidth}
        className={`${styles.panelResizer} ${isDragging ? styles.panelResizerDragging : ""}`}
        onMouseDown={handleMouseDown}
        onDoubleClick={handleDoubleClick}
        onKeyDown={handleKeyDown}
      />

      {/* Floating expand button when collapsed */}
      {isCollapsed && (
        <button
          type="button"
          onClick={toggleCollapse}
          className={styles.sidebarToggleIconBtnFloating}
          title="Expand diagrams list"
        >
          <Icons.PanelLeftOpen size={16} />
        </button>
      )}

      {/* Right Panel: Active Excalidraw Canvas or Empty State */}
      <div className={styles.rightPanel}>
        {activeDrawingId && currentDrawing ? (
          <div className="flex flex-col h-full w-full overflow-hidden bg-card rounded-xl border border-border/70 shadow-xs relative">
            <DrawingHeader
              drawing={currentDrawing}
              projectId={projectId}
              saveStatus={saveStatus}
              projectColor={projectColor}
              onRename={handleRenameDrawing}
              onSaveNow={saveNow}
              onDelete={() => handleDeleteClick(currentDrawing.id)}
              onExport={handleExport}
              onCopyToClipboard={handleCopyToClipboard}
              onOpenImageDialog={handleOpenImageDialog}
              onOpenManageTagsModal={onOpenManageTagsModal}
            />

            <Suspense fallback={<DrawingCanvasSkeleton />}>
              <DrawingCanvas
                key={currentDrawing.id}
                drawing={currentDrawing}
                onSceneChange={triggerChange}
                excalidrawAPIRef={excalidrawAPIRef}
              />
            </Suspense>
          </div>
        ) : isLoadingCurrent ? (
          <div className="flex-1 h-full rounded-xl border border-border/70 overflow-hidden">
            <DrawingCanvasSkeleton />
          </div>
        ) : (
          /* Empty State when no diagram is selected or created */
          <div className="flex-1 h-full flex flex-col items-center justify-center p-8 text-center bg-card/60 rounded-xl border border-border/70">
            <div className="w-16 h-16 rounded-2xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary mb-4 shadow-xs">
              <Icons.Shapes size={32} />
            </div>
            <h3 className="text-base font-semibold text-foreground mb-1">
              No diagram selected
            </h3>
            <p className="text-xs text-muted-foreground max-w-sm mb-5 leading-relaxed">
              Create architectural diagrams, flows, wireframes, and whiteboard
              brainstorms right inside your project with Excalidraw.
            </p>
            <button
              type="button"
              onClick={handleCreateDrawing}
              disabled={createMutation.isPending}
              className="flex items-center gap-2 px-4 py-2 rounded-lg bg-primary text-primary-foreground text-xs font-semibold hover:bg-primary/90 shadow-md transition-all cursor-pointer"
            >
              <Icons.Plus size={14} />
              <span>Create New Diagram</span>
            </button>
          </div>
        )}
      </div>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={!!drawingToDelete}
        onClose={() => setDrawingToDelete(null)}
        onConfirm={handleConfirmDelete}
        title="Delete Diagram"
        message="Are you sure you want to delete this diagram? This action cannot be undone."
        itemName={drawingToDelete?.title}
        isLoading={isDeleting}
        confirmLabel="Delete Diagram"
      />
    </>
  );
};
