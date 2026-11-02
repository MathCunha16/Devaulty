import React, { useState, useRef, useEffect, useCallback } from "react";
import { createPortal } from "react-dom";
import * as Icons from "lucide-react";
import { toast } from "sonner";
import type { Drawing } from "~types/api";
import { formatUpdatedDate } from "../../../utils/dateUtils";
import type { SaveStatus } from "../hooks/useDebouncedSave";
import {
  useTagsQuery,
  useCreateTagMutation,
  useAssociateTagMutation,
  useDisassociateTagMutation,
} from "~features/tags/hooks/useTags";

interface DrawingHeaderProps {
  drawing: Drawing;
  projectId: string;
  saveStatus: SaveStatus;
  projectColor?: string;
  onRename: (newTitle: string) => void;
  onSaveNow: () => void;
  onDelete: () => void;
  onExport: (format: "png" | "svg" | "json") => void;
  onCopyToClipboard: (format: "png" | "svg") => void;
  onOpenImageDialog: () => void;
  onOpenManageTagsModal: () => void;
}

export const DrawingHeader: React.FC<DrawingHeaderProps> = ({
  drawing,
  projectId,
  saveStatus,
  projectColor = "#10b981",
  onRename,
  onSaveNow,
  onDelete,
  onExport,
  onCopyToClipboard,
  onOpenImageDialog,
  onOpenManageTagsModal,
}) => {
  // Inline title rename state without cascading render effect
  const [isEditingTitle, setIsEditingTitle] = useState(false);
  const [titleValue, setTitleValue] = useState(drawing.title);
  const [prevTitle, setPrevTitle] = useState(drawing.title);

  if (drawing.title !== prevTitle) {
    setPrevTitle(drawing.title);
    setTitleValue(drawing.title);
  }

  // Export menu state
  const [isExportMenuOpen, setIsExportMenuOpen] = useState(false);

  // Compact Tag Selector state
  const [isTagPopoverOpen, setIsTagPopoverOpen] = useState(false);
  const [tagSearchQuery, setTagSearchQuery] = useState("");
  const [popoverPos, setPopoverPos] = useState<{ top: number; left: number } | null>(null);

  const addTagBtnRef = useRef<HTMLButtonElement>(null);
  const tagPopoverRef = useRef<HTMLDivElement>(null);

  // Queries & Mutations for tags
  const { data: allTags = [] } = useTagsQuery(projectId);
  const createTagMutation = useCreateTagMutation(projectId);
  const associateTagMutation = useAssociateTagMutation(projectId);
  const disassociateTagMutation = useDisassociateTagMutation(projectId);

  const currentTags = drawing.tags || [];

  const updatePopoverPosition = useCallback(() => {
    if (!addTagBtnRef.current) return;
    const rect = addTagBtnRef.current.getBoundingClientRect();
    const popoverHeight = 240;
    const popoverWidth = 220;
    const viewportMargin = 8;
    const spaceBelow = window.innerHeight - rect.bottom;
    const preferredTop =
      spaceBelow >= popoverHeight ? rect.bottom + 6 : rect.top - popoverHeight - 6;
    const maxTop = Math.max(viewportMargin, window.innerHeight - popoverHeight - viewportMargin);
    const maxLeft = Math.max(viewportMargin, window.innerWidth - popoverWidth - viewportMargin);
    const top = Math.min(Math.max(preferredTop, viewportMargin), maxTop);
    const left = Math.min(Math.max(rect.left, viewportMargin), maxLeft);

    setPopoverPos({ top, left });
  }, []);

  const openTagPopover = useCallback(() => {
    updatePopoverPosition();
    setIsTagPopoverOpen(true);
  }, [updatePopoverPosition]);

  const closeTagPopover = useCallback(() => {
    setIsTagPopoverOpen(false);
    setTagSearchQuery("");
    setPopoverPos(null);
  }, []);

  useEffect(() => {
    if (!isTagPopoverOpen) return;

    const handleOutsideClick = (e: MouseEvent) => {
      if (
        tagPopoverRef.current &&
        !tagPopoverRef.current.contains(e.target as Node) &&
        !addTagBtnRef.current?.contains(e.target as Node)
      ) {
        closeTagPopover();
      }
    };

    const handleScroll = (e: Event) => {
      if (tagPopoverRef.current?.contains(e.target as Node)) return;
      updatePopoverPosition();
    };

    document.addEventListener("mousedown", handleOutsideClick);
    window.addEventListener("scroll", handleScroll, true);
    window.addEventListener("resize", updatePopoverPosition);
    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
      window.removeEventListener("scroll", handleScroll, true);
      window.removeEventListener("resize", updatePopoverPosition);
    };
  }, [isTagPopoverOpen, closeTagPopover, updatePopoverPosition]);

  const handleAssociateTag = async (tagId: string) => {
    try {
      await associateTagMutation.mutateAsync({
        itemType: "DRAWING",
        itemId: drawing.id,
        tagId,
      });
      toast.success("Tag associated");
      closeTagPopover();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to associate tag");
    }
  };

  const handleDisassociateTag = async (tagId: string) => {
    try {
      await disassociateTagMutation.mutateAsync({
        itemType: "DRAWING",
        itemId: drawing.id,
        tagId,
      });
      toast.success("Tag removed");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to remove tag");
    }
  };

  const handleCreateAndAssociateTag = async () => {
    const trimmed = tagSearchQuery.trim();
    if (!trimmed) return;

    const presetColors = ["#8b5cf6", "#10b981", "#f43f5e", "#f59e0b", "#0ea5e9"];
    const randomColor = presetColors[Math.floor(Math.random() * presetColors.length)];

    try {
      const newTag = await createTagMutation.mutateAsync({
        name: trimmed,
        color: randomColor,
      });
      await associateTagMutation.mutateAsync({
        itemType: "DRAWING",
        itemId: drawing.id,
        tagId: newTag.id,
      });
      closeTagPopover();
      toast.success(`Tag "${newTag.name}" created and added`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create tag");
    }
  };

  const unassociatedTags = allTags.filter(
    (t) =>
      t.name.toLowerCase().includes(tagSearchQuery.toLowerCase().trim()) &&
      !currentTags.some((ct) => ct.id === t.id)
  );

  const handleTitleSubmit = () => {
    const trimmed = titleValue.trim();
    if (trimmed && trimmed !== drawing.title) {
      onRename(trimmed);
    } else {
      setTitleValue(drawing.title);
    }
    setIsEditingTitle(false);
  };

  const handleTitleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleTitleSubmit();
    } else if (e.key === "Escape") {
      setTitleValue(drawing.title);
      setIsEditingTitle(false);
    }
  };

  const formattedDate = formatUpdatedDate(drawing.updatedAt, drawing.createdAt);

  return (
    <div className="flex items-center justify-between gap-3 px-3.5 py-2 bg-card/75 border-b border-border/80 backdrop-blur-md shrink-0 z-20">
      {/* Left: Icon, Title & Meta */}
      <div className="flex items-center gap-2.5 min-w-0 flex-1">
        <div
          className="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 border"
          style={{
            backgroundColor: `${projectColor}15`,
            borderColor: `${projectColor}30`,
            color: projectColor,
          }}
        >
          <Icons.Shapes size={15} />
        </div>

        <div className="flex flex-col min-w-0">
          {isEditingTitle ? (
            <div className="flex items-center gap-1.5">
              <input
                type="text"
                autoFocus
                className="bg-secondary/90 border border-primary/60 text-foreground text-xs font-semibold rounded px-2 py-0.5 outline-none focus:ring-1 focus:ring-primary w-full max-w-xs shadow-xs"
                value={titleValue}
                onChange={(e) => setTitleValue(e.target.value)}
                onBlur={handleTitleSubmit}
                onKeyDown={handleTitleKeyDown}
                placeholder="Diagram title..."
              />
              <button
                type="button"
                onClick={handleTitleSubmit}
                className="p-1 rounded bg-primary/20 hover:bg-primary/30 text-primary transition-colors"
                title="Save title"
              >
                <Icons.Check size={11} />
              </button>
            </div>
          ) : (
            <div
              role="button"
              tabIndex={0}
              onClick={() => setIsEditingTitle(true)}
              onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") setIsEditingTitle(true);
              }}
              className="flex items-center gap-1.5 px-1.5 py-0.5 -ml-1 rounded-md hover:bg-secondary/60 border border-transparent hover:border-border/60 transition-all cursor-pointer group"
              title="Click to rename this diagram"
            >
              <h2 className="text-xs font-semibold text-foreground truncate max-w-xs sm:max-w-sm group-hover:text-primary transition-colors">
                {drawing.title}
              </h2>
              <span className="flex items-center gap-0.5 text-[10px] text-muted-foreground/70 group-hover:text-primary transition-colors bg-secondary/50 group-hover:bg-primary/10 px-1.5 py-0.2 rounded border border-border/50 group-hover:border-primary/30 shrink-0">
                <Icons.Pencil size={9} />
                <span className="text-[9px] font-medium hidden sm:inline">Rename</span>
              </span>
            </div>
          )}

          {formattedDate && (
            <span className="text-[10px] text-muted-foreground/70 truncate pl-0.5">
              {formattedDate}
            </span>
          )}
        </div>

        {/* Sleek, Compact Inline Tags Section */}
        <div className="hidden md:flex items-center gap-1 ml-2 max-w-sm overflow-hidden flex-wrap">
          {currentTags.map((tag) => (
            <span
              key={tag.id}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[10px] font-medium border transition-colors group"
              style={{
                backgroundColor: `${tag.color || projectColor}15`,
                borderColor: `${tag.color || projectColor}35`,
                color: tag.color || projectColor,
              }}
            >
              <span
                className="w-1.5 h-1.5 rounded-full shrink-0"
                style={{ backgroundColor: tag.color || projectColor }}
              />
              <span className="truncate max-w-[90px]">{tag.name}</span>
              <button
                type="button"
                onClick={(e) => {
                  e.stopPropagation();
                  handleDisassociateTag(tag.id);
                }}
                className="opacity-40 group-hover:opacity-100 hover:text-destructive transition-opacity ml-0.5"
                title={`Remove tag ${tag.name}`}
              >
                <Icons.X size={9} />
              </button>
            </span>
          ))}

          {/* Add Tag Button */}
          <button
            ref={addTagBtnRef}
            type="button"
            onClick={() => (isTagPopoverOpen ? closeTagPopover() : openTagPopover())}
            className="inline-flex items-center gap-1 px-1.5 py-0.5 rounded-full text-[10px] font-medium text-muted-foreground hover:text-foreground bg-secondary/50 hover:bg-secondary border border-border/60 transition-colors"
            title="Add tag to diagram"
          >
            <Icons.Plus size={9} />
            <span>Tag</span>
          </button>
        </div>
      </div>

      {/* Center: Save Status Pill */}
      <div className="flex items-center shrink-0">
        {saveStatus === "saving" && (
          <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-[11px] font-medium">
            <span className="w-1.5 h-1.5 rounded-full bg-blue-400 animate-ping" />
            <span>Saving...</span>
          </div>
        )}
        {saveStatus === "unsaved" && (
          <div className="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-amber-500/10 border border-amber-500/20 text-amber-400 text-[11px] font-medium">
            <span className="w-1.5 h-1.5 rounded-full bg-amber-400" />
            <span>Unsaved</span>
          </div>
        )}
        {saveStatus === "saved" && (
          <div className="flex items-center gap-1 px-2 py-0.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-[11px] font-medium">
            <Icons.Check size={11} />
            <span>Saved</span>
          </div>
        )}
        {saveStatus === "error" && (
          <div className="flex items-center gap-1 px-2 py-0.5 rounded-full bg-destructive/10 border border-destructive/20 text-destructive text-[11px] font-medium">
            <Icons.AlertTriangle size={11} />
            <span>Save error</span>
          </div>
        )}
      </div>

      {/* Right: Actions (Save, Export, Delete) */}
      <div className="flex items-center gap-1.5 shrink-0">
        {/* Save Now Button (Instant manual flush) */}
        <button
          type="button"
          onClick={onSaveNow}
          className={`flex items-center gap-1 px-2.5 py-1 rounded-md text-xs font-medium transition-all ${
            saveStatus === "unsaved"
              ? "bg-primary text-primary-foreground hover:bg-primary/90 shadow-xs"
              : "bg-secondary/60 text-muted-foreground hover:text-foreground hover:bg-secondary"
          }`}
          title="Save diagram (Ctrl+S / Cmd+S)"
        >
          <Icons.Save size={12} />
          <span className="hidden sm:inline">Save</span>
        </button>

        {/* Export Dropdown */}
        <div className="relative">
          <button
            type="button"
            onClick={() => setIsExportMenuOpen((prev) => !prev)}
            className="flex items-center gap-1 px-2.5 py-1 rounded-md text-xs font-medium bg-secondary/60 text-muted-foreground hover:text-foreground hover:bg-secondary transition-all"
            title="Export drawing"
          >
            <Icons.Download size={12} />
            <span className="hidden sm:inline">Export</span>
            <Icons.ChevronDown size={10} />
          </button>

          {isExportMenuOpen && (
            <>
              <div
                className="fixed inset-0 z-30"
                onClick={() => setIsExportMenuOpen(false)}
              />
              <div className="absolute right-0 top-full mt-1.5 w-52 bg-card/95 border border-border/80 rounded-lg shadow-xl backdrop-blur-md p-1 z-40 flex flex-col gap-0.5 animate-in fade-in zoom-in-95 duration-100">
                {/* Clipboard Section */}
                <div className="px-2 py-1 text-[10px] font-semibold text-muted-foreground/70 uppercase tracking-wider">
                  Clipboard
                </div>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onCopyToClipboard("png");
                  }}
                >
                  <Icons.ClipboardCopy size={13} className="text-emerald-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Copy PNG</span>
                    <span className="text-[10px] text-muted-foreground/80">Direct paste for docs, chat, etc.</span>
                  </div>
                </button>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onCopyToClipboard("svg");
                  }}
                >
                  <Icons.Copy size={13} className="text-cyan-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Copy SVG</span>
                    <span className="text-[10px] text-muted-foreground/80">Vector for Figma, code or web</span>
                  </div>
                </button>

                <div className="h-[1px] bg-border/60 my-1" />

                {/* File Download Section */}
                <div className="px-2 py-1 text-[10px] font-semibold text-muted-foreground/70 uppercase tracking-wider">
                  Save to Disk
                </div>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onOpenImageDialog();
                  }}
                >
                  <Icons.Sliders size={13} className="text-amber-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Export Image Dialog...</span>
                    <span className="text-[10px] text-muted-foreground/80">Configure 2x/3x scale & background</span>
                  </div>
                </button>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onExport("png");
                  }}
                >
                  <Icons.Image size={13} className="text-emerald-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Save PNG Image...</span>
                    <span className="text-[10px] text-muted-foreground/80">Choose destination folder</span>
                  </div>
                </button>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onExport("svg");
                  }}
                >
                  <Icons.FileCode size={13} className="text-cyan-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Save SVG File...</span>
                    <span className="text-[10px] text-muted-foreground/80">Choose destination folder</span>
                  </div>
                </button>
                <button
                  type="button"
                  className="flex items-center gap-2 px-2.5 py-1.5 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/80 text-left transition-colors"
                  onClick={() => {
                    setIsExportMenuOpen(false);
                    onExport("json");
                  }}
                >
                  <Icons.FileJson size={13} className="text-purple-400 shrink-0" />
                  <div className="flex flex-col min-w-0">
                    <span className="font-medium">Save JSON (.excalidraw)...</span>
                    <span className="text-[10px] text-muted-foreground/80">Diagram backup file</span>
                  </div>
                </button>
              </div>
            </>
          )}
        </div>

        {/* Delete Button */}
        <button
          type="button"
          onClick={onDelete}
          className="p-1 rounded-md text-muted-foreground hover:text-destructive hover:bg-destructive/10 transition-colors"
          title="Delete diagram"
        >
          <Icons.Trash2 size={13} />
        </button>
      </div>

      {/* Portal for Tag Selection Popover */}
      {isTagPopoverOpen &&
        popoverPos &&
        createPortal(
          <div
            ref={tagPopoverRef}
            className="fixed w-56 bg-card/95 border border-border/80 rounded-xl shadow-2xl backdrop-blur-md p-2 z-[9999] flex flex-col gap-1 animate-in fade-in zoom-in-95 duration-100"
            style={{
              top: popoverPos.top,
              left: popoverPos.left,
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <div className="relative mb-1">
              <input
                type="text"
                autoFocus
                placeholder="Search or create tag..."
                className="w-full bg-secondary/80 border border-border/70 rounded-md px-2 py-1 text-xs text-foreground placeholder:text-muted-foreground/60 outline-none focus:border-primary/60"
                value={tagSearchQuery}
                onChange={(e) => setTagSearchQuery(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === "Enter") handleCreateAndAssociateTag();
                  if (e.key === "Escape") closeTagPopover();
                }}
              />
            </div>

            <div className="max-h-40 overflow-y-auto space-y-0.5">
              {unassociatedTags.length === 0 && !tagSearchQuery.trim() && (
                <div className="text-[11px] text-muted-foreground/70 py-2 text-center">
                  No other tags available
                </div>
              )}

              {unassociatedTags.map((tag) => (
                <button
                  key={tag.id}
                  type="button"
                  className="w-full flex items-center gap-2 px-2 py-1 rounded text-xs text-foreground/80 hover:text-foreground hover:bg-secondary/70 transition-colors text-left"
                  onClick={() => handleAssociateTag(tag.id)}
                >
                  <span
                    className="w-2 h-2 rounded-full shrink-0"
                    style={{ backgroundColor: tag.color || projectColor }}
                  />
                  <span className="truncate">{tag.name}</span>
                </button>
              ))}

              {tagSearchQuery.trim() &&
                !allTags.some(
                  (t) => t.name.toLowerCase() === tagSearchQuery.toLowerCase().trim()
                ) && (
                  <button
                    type="button"
                    className="w-full flex items-center gap-1.5 px-2 py-1 rounded text-xs text-primary hover:bg-primary/10 transition-colors text-left"
                    onClick={handleCreateAndAssociateTag}
                  >
                    <Icons.Plus size={11} />
                    <span>Create "{tagSearchQuery.trim()}"</span>
                  </button>
                )}
            </div>

            <div className="border-t border-border/60 pt-1 mt-1">
              <button
                type="button"
                onClick={() => {
                  closeTagPopover();
                  onOpenManageTagsModal();
                }}
                className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 uppercase tracking-wider font-medium"
              >
                <Icons.Settings size={10} />
                <span>Manage Tags</span>
              </button>
            </div>
          </div>,
          document.body
        )}
    </div>
  );
};
