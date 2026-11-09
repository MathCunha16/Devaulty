import React from "react";
import * as Icons from "lucide-react";
import type { DrawingSummary } from "~types/api";
import { formatUpdatedDate } from "../../../utils/dateUtils";

interface DrawingCardProps {
  drawing: DrawingSummary;
  isSelected: boolean;
  projectColor?: string;
  onSelect: () => void;
  onDelete: () => void;
}

export const DrawingCard: React.FC<DrawingCardProps> = ({
  drawing,
  isSelected,
  projectColor = "#10b981",
  onSelect,
  onDelete,
}) => {
  const formattedDate = formatUpdatedDate(drawing.updatedAt, drawing.createdAt);

  return (
    <div
      onClick={onSelect}
      className={`group relative flex flex-col gap-2 p-2.5 rounded-xl border transition-all cursor-pointer ${
        isSelected
          ? "bg-card/90 shadow-md ring-1"
          : "bg-card/40 border-border/70 hover:bg-card/70 hover:border-border"
      }`}
      style={
        isSelected
          ? {
              borderColor: `${projectColor}60`,
              boxShadow: `0 0 12px ${projectColor}15`,
            }
          : undefined
      }
    >
      {/* Thumbnail or Fallback */}
      <div className="w-full h-24 rounded-lg bg-secondary/40 border border-border/50 overflow-hidden relative flex items-center justify-center">
        {drawing.thumbnail ? (
          <img
            src={drawing.thumbnail}
            alt={drawing.title}
            className="w-full h-full object-contain p-1 opacity-90 group-hover:opacity-100 transition-opacity"
            loading="lazy"
          />
        ) : (
          <div className="flex flex-col items-center gap-1 text-muted-foreground/50">
            <Icons.Shapes size={22} />
            <span className="text-[10px] tracking-wide">Blank Diagram</span>
          </div>
        )}

        {/* Delete button overlay */}
        <button
          type="button"
          onClick={(e) => {
            e.stopPropagation();
            onDelete();
          }}
          className="absolute top-1.5 right-1.5 p-1 rounded-md bg-background/80 hover:bg-destructive/20 text-muted-foreground hover:text-destructive opacity-0 group-hover:opacity-100 transition-all backdrop-blur-xs shadow-xs"
          title="Delete diagram"
        >
          <Icons.Trash2 size={12} />
        </button>
      </div>

      {/* Info */}
      <div className="flex flex-col gap-1 min-w-0">
        <div className="flex items-center justify-between gap-1.5">
          <span
            className={`text-xs font-semibold truncate ${
              isSelected ? "" : "text-foreground group-hover:text-primary transition-colors"
            }`}
            style={isSelected ? { color: projectColor } : undefined}
          >
            {drawing.title}
          </span>
        </div>

        <div className="flex items-center justify-between gap-1 text-[10px] text-muted-foreground">
          <span className="truncate">{formattedDate || "Just created"}</span>

          {/* Tags counter or preview */}
          {drawing.tags && drawing.tags.length > 0 && (
            <div className="flex items-center gap-1 shrink-0">
              <span className="flex items-center gap-0.5 px-1.5 py-0.2 rounded-full bg-secondary text-muted-foreground font-medium">
                <Icons.Tag size={9} />
                {drawing.tags.length}
              </span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
