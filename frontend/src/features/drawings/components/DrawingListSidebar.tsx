import React from "react";
import * as Icons from "lucide-react";
import type { DrawingSummary } from "~types/api";
import { DrawingCard } from "./DrawingCard";
import styles from "../../../routes/projects.$projectId.module.css";

interface DrawingListSidebarProps {
  drawings: DrawingSummary[];
  selectedId?: string;
  searchQuery: string;
  projectColor?: string;
  isCreating?: boolean;
  onSelect: (id: string) => void;
  onCreate: () => void;
  onDelete: (id: string) => void;
  onSearchChange: (query: string) => void;
  toggleCollapse: () => void;
}

export const DrawingListSidebar: React.FC<DrawingListSidebarProps> = ({
  drawings,
  selectedId,
  searchQuery,
  projectColor = "#10b981",
  isCreating = false,
  onSelect,
  onCreate,
  onDelete,
  onSearchChange,
  toggleCollapse,
}) => {
  const filteredDrawings = drawings.filter((drawing) => {
    const q = searchQuery.toLowerCase().trim();
    if (!q) return true;
    const matchesTitle = drawing.title.toLowerCase().includes(q);
    const matchesTag = drawing.tags?.some((t) => t.name.toLowerCase().includes(q));
    return matchesTitle || matchesTag;
  });

  return (
    <div className="flex flex-col h-full gap-2.5 p-3">
      {/* Top Action Bar */}
      <div className="flex items-center gap-1.5 shrink-0">
        <button
          type="button"
          className={`${styles.newSnippetBtn} flex-1 shadow-sm`}
          onClick={onCreate}
          disabled={isCreating}
        >
          <Icons.Plus size={14} />
          <span>{isCreating ? "Creating..." : "New Diagram"}</span>
        </button>
        <button
          type="button"
          onClick={toggleCollapse}
          className={styles.sidebarToggleIconBtn}
          title="Collapse diagrams list"
        >
          <Icons.PanelLeftClose size={15} />
        </button>
      </div>

      {/* Search Input */}
      <div className={styles.searchBar}>
        <Icons.Search className={styles.searchIcon} size={14} />
        <input
          type="text"
          placeholder="Search diagrams or tags..."
          className={styles.searchInput}
          value={searchQuery}
          onChange={(e) => onSearchChange(e.target.value)}
        />
        {searchQuery && (
          <button
            type="button"
            onClick={() => onSearchChange("")}
            className="text-muted-foreground hover:text-foreground pr-1"
          >
            <Icons.X size={12} />
          </button>
        )}
      </div>

      {/* Counter */}
      <div className="flex items-center justify-between text-[11px] text-muted-foreground px-1 shrink-0">
        <span>Diagrams</span>
        <span className="font-semibold text-foreground/75">
          {filteredDrawings.length}
        </span>
      </div>

      {/* Diagrams Scrollable List */}
      <div className="flex-1 overflow-y-auto pr-0.5 space-y-2 min-h-0">
        {filteredDrawings.length === 0 ? (
          <div className="flex flex-col items-center justify-center p-6 text-center text-muted-foreground gap-2 mt-4">
            <div className="w-10 h-10 rounded-full bg-secondary/50 flex items-center justify-center text-muted-foreground/60">
              <Icons.Shapes size={18} />
            </div>
            <span className="text-xs font-medium">
              {searchQuery ? "No diagrams match your search" : "No diagrams yet"}
            </span>
            <span className="text-[11px] text-muted-foreground/80 max-w-[180px]">
              {searchQuery
                ? "Try searching for another keyword or tag"
                : "Create your first Excalidraw whiteboard to start visualizing"}
            </span>
            {!searchQuery && (
              <button
                type="button"
                onClick={onCreate}
                className="mt-2 text-xs font-medium text-primary hover:underline"
              >
                + Create diagram
              </button>
            )}
          </div>
        ) : (
          filteredDrawings.map((drawing) => (
            <DrawingCard
              key={drawing.id}
              drawing={drawing}
              isSelected={drawing.id === selectedId}
              projectColor={projectColor}
              onSelect={() => onSelect(drawing.id)}
              onDelete={() => onDelete(drawing.id)}
            />
          ))
        )}
      </div>
    </div>
  );
};
