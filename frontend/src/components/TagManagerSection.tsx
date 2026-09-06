import React, { useState, useEffect, useRef, useCallback } from "react";
import { createPortal } from "react-dom";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import {
  useTagsQuery,
  useCreateTagMutation,
  useAssociateTagMutation,
  useDisassociateTagMutation,
} from "~features/tags/hooks/useTags";
import type { TagSummaryResponse } from "~types/api";
import styles from "../routes/projects.$projectId.module.css";

interface TagManagerSectionProps {
  itemId: string;
  itemType: "SNIPPET" | "PROBLEM" | "NOTE" | "LINK" | "CREDENTIAL";
  itemTags?: TagSummaryResponse[];
  projectId: string;
  onOpenManageTagsModal: () => void;
  title?: string;
  noBorder?: boolean;
}

export const TagManagerSection: React.FC<TagManagerSectionProps> = ({
  itemId,
  itemType,
  itemTags = [],
  projectId,
  onOpenManageTagsModal,
  title = "Associated Tags",
  noBorder = false,
}) => {
  const { data: tagsData = [] } = useTagsQuery(projectId);
  const createTagMutation = useCreateTagMutation(projectId);
  const associateTagMutation = useAssociateTagMutation(projectId);
  const disassociateTagMutation = useDisassociateTagMutation(projectId);

  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const [tagSearchQuery, setTagSearchQuery] = useState("");
  const [popoverPos, setPopoverPos] = useState<{ top: number; left: number } | null>(null);

  const triggerRef = useRef<HTMLButtonElement>(null);
  const popoverRef = useRef<HTMLDivElement>(null);

  const updatePopoverPosition = useCallback(() => {
    if (!triggerRef.current) return;
    const rect = triggerRef.current.getBoundingClientRect();
    const popoverHeight = 260; // approximate max-height
    const spaceBelow = window.innerHeight - rect.bottom;
    const top = spaceBelow >= popoverHeight
      ? rect.bottom + 4
      : rect.top - popoverHeight - 4;
    setPopoverPos({ top, left: rect.left });
  }, []);

  const openPopover = useCallback(() => {
    updatePopoverPosition();
    setIsPopoverOpen(true);
  }, [updatePopoverPosition]);

  const closePopover = useCallback(() => {
    setIsPopoverOpen(false);
    setTagSearchQuery("");
    setPopoverPos(null);
  }, []);

  // Close on outside click and keep the popover aligned while its container scrolls
  useEffect(() => {
    if (!isPopoverOpen) return;

    const handleOutside = (e: MouseEvent) => {
      if (
        popoverRef.current &&
        !popoverRef.current.contains(e.target as Node) &&
        !triggerRef.current?.contains(e.target as Node)
      ) {
        closePopover();
      }
    };

    const handleScroll = (e: Event) => {
      if (popoverRef.current?.contains(e.target as Node)) return;
      updatePopoverPosition();
    };

    document.addEventListener("mousedown", handleOutside);
    window.addEventListener("scroll", handleScroll, true);
    window.addEventListener("resize", updatePopoverPosition);
    return () => {
      document.removeEventListener("mousedown", handleOutside);
      window.removeEventListener("scroll", handleScroll, true);
      window.removeEventListener("resize", updatePopoverPosition);
    };
  }, [isPopoverOpen, closePopover, updatePopoverPosition]);

  const handleAddTag = async (tagId: string) => {
    try {
      await associateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag associated successfully");
      closePopover();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to associate tag");
    }
  };
  const handleRemoveTag = async (tagId: string) => {
    try {
      await disassociateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag removed successfully");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to remove tag");
    }
  };

  const handleCreateAndAddTag = async () => {
    if (!tagSearchQuery.trim()) return;
    let newTag;
    try {
      const presetColors = ["#8b5cf6", "#10b981", "#f43f5e", "#f59e0b", "#0ea5e9"];
      const randomColor = presetColors[Math.floor(Math.random() * presetColors.length)];
      newTag = await createTagMutation.mutateAsync({
        name: tagSearchQuery.trim(),
        color: randomColor,
      });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create tag");
      return;
    }

    try {
      await associateTagMutation.mutateAsync({
        itemType,
        itemId,
        tagId: newTag.id,
      });
      closePopover();
      toast.success(`Tag "${newTag.name}" created and associated`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : `Tag "${newTag.name}" created, but failed to associate`);
    }
  };

  const unassociatedTags = tagsData.filter(
    (t) =>
      t.name.toLowerCase().includes(tagSearchQuery.toLowerCase()) &&
      !itemTags.some((st) => st.id === t.id)
  );

  return (
    <div className={`${styles.tagSection} ${noBorder ? styles.tagSectionNoBorder : ""}`}>
      {title && title.trim() !== "" && (
        <div className={styles.tagHeader}>
          <Icons.Tag size={12} className="text-muted-foreground" />
          <span className={styles.tagSectionTitle}>{title}</span>
        </div>
      )}

      <div className={styles.tagList}>
        {itemTags.map((tag) => (
          <span key={tag.id} className={styles.tagPill}>
            <span
              className={styles.tagDot}
              style={{ backgroundColor: tag.color || "var(--color-primary)" }}
            />
            <span>{tag.name}</span>
            <button
              type="button"
              className={styles.tagRemoveBtn}
              onClick={() => handleRemoveTag(tag.id)}
              title={`Remove tag ${tag.name}`}
            >
              <Icons.X size={10} />
            </button>
          </span>
        ))}

        <div className={styles.addTagContainer}>
          <button
            ref={triggerRef}
            type="button"
            className={styles.addTagBtn}
            onClick={() => isPopoverOpen ? closePopover() : openPopover()}
          >
            <Icons.Plus size={10} />
            <span>Add Tag</span>
          </button>
        </div>
      </div>

      {/* Popover rendered via portal so it escapes any overflow:hidden/auto ancestor */}
      {isPopoverOpen && popoverPos && createPortal(
        <div
          ref={popoverRef}
          className={styles.tagPopover}
          style={{
            position: "fixed",
            top: popoverPos.top,
            left: popoverPos.left,
            bottom: "auto",
            zIndex: 9999,
          }}
          onClick={(e) => e.stopPropagation()}
          onKeyDown={(e) => e.stopPropagation()}
        >
          <div className={styles.popoverHeader}>
            <input
              type="text"
              placeholder="Filter/create tag..."
              className={styles.tagSearchInput}
              value={tagSearchQuery}
              onChange={(e) => setTagSearchQuery(e.target.value)}
              onKeyDown={(e) => {
                e.stopPropagation();
                if (e.key === "Enter") handleCreateAndAddTag();
                if (e.key === "Escape") closePopover();
              }}
              autoFocus
            />
          </div>

          <div className={styles.popoverList}>
            {unassociatedTags.map((t) => (
              <button
                key={t.id}
                type="button"
                className={styles.popoverItem}
                onClick={() => handleAddTag(t.id)}
              >
                <span
                  className={styles.tagColorPreview}
                  style={{ backgroundColor: t.color || "var(--color-primary)" }}
                />
                <span>{t.name}</span>
              </button>
            ))}

            {tagSearchQuery.trim() &&
              !tagsData.some(
                (t) => t.name.toLowerCase() === tagSearchQuery.toLowerCase()
              ) && (
                <button
                  type="button"
                  className={styles.popoverItemCreate}
                  onClick={handleCreateAndAddTag}
                >
                  <Icons.Plus size={10} />
                  <span>Create "{tagSearchQuery}"</span>
                </button>
              )}
            <div className="border-t border-border mt-2 pt-2 px-1">
              <button
                type="button"
                onClick={() => {
                  closePopover();
                  onOpenManageTagsModal();
                }}
                className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
              >
                <Icons.Settings size={10} />
                <span>Manage Tags</span>
              </button>
            </div>
          </div>
        </div>,
        document.body
      )}
    </div>
  );
};
