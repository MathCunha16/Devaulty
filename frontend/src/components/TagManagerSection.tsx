import React, { useState, useEffect } from "react";
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

  // Close tag popover when clicking outside
  useEffect(() => {
    if (!isPopoverOpen) return;

    const handleOutsideClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      if (!target.closest(`.${styles.addTagContainer}`)) {
        setIsPopoverOpen(false);
        setTagSearchQuery("");
      }
    };

    document.addEventListener("click", handleOutsideClick);
    return () => document.removeEventListener("click", handleOutsideClick);
  }, [isPopoverOpen]);

  const handleAddTag = async (tagId: string) => {
    try {
      await associateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag associated successfully");
      setIsPopoverOpen(false);
      setTagSearchQuery("");
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
      setTagSearchQuery("");
      setIsPopoverOpen(false);
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
            type="button"
            className={styles.addTagBtn}
            onClick={() => setIsPopoverOpen((prev) => !prev)}
          >
            <Icons.Plus size={10} />
            <span>Add Tag</span>
          </button>

          {isPopoverOpen && (
            <div
              className={styles.tagPopover}
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
                  onKeyDown={(e) => e.stopPropagation()}
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
                      setIsPopoverOpen(false);
                      onOpenManageTagsModal();
                    }}
                    className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                  >
                    <Icons.Settings size={10} />
                    <span>Manage Tags</span>
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};
