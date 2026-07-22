import React, { useState, useEffect, useRef } from "react";
import { X, Plus, Edit2, Trash2, Check, RotateCcw } from "lucide-react";
import { toast } from "sonner";
import {
  useTagsQuery,
  useCreateTagMutation,
  useUpdateTagMutation,
  useDeleteTagMutation,
} from "~features/tags/hooks/useTags";
import { ConfirmModal } from "./ConfirmModal";
import styles from "./TagsManagerModal.module.css";

interface TagsManagerModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
}

const PRESET_COLORS = [
  "#ef4444", // Red
  "#f97316", // Orange
  "#eab308", // Yellow
  "#22c55e", // Green
  "#10b981", // Emerald
  "#14b8a6", // Teal
  "#06b6d4", // Cyan
  "#3b82f6", // Blue
  "#6366f1", // Indigo
  "#8b5cf6", // Violet
  "#ec4899", // Pink
  "#f43f5e", // Rose
  "#64748b", // Slate
];

export const TagsManagerModal: React.FC<TagsManagerModalProps> = ({
  isOpen,
  onClose,
  projectId,
}) => {
  const { data: tags = [], isLoading } = useTagsQuery(projectId);
  const createMutation = useCreateTagMutation(projectId);
  const updateMutation = useUpdateTagMutation(projectId);
  const deleteMutation = useDeleteTagMutation(projectId);

  // Focus trap refs
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const titleInputRef = useRef<HTMLInputElement>(null);

  // Form states
  const [newTagName, setNewTagName] = useState("");
  const [newTagColor, setNewTagColor] = useState(PRESET_COLORS[0]);
  const [searchQuery, setSearchQuery] = useState("");

  // Edit states
  const [editingTagId, setEditingTagId] = useState<string | null>(null);
  const [editingTagName, setEditingTagName] = useState("");
  const [editingTagColor, setEditingTagColor] = useState("");

  // Delete confirm state
  const [tagToDelete, setTagToDelete] = useState<{ id: string; name: string } | null>(null);

  useEffect(() => {
    if (isOpen) {
      previousActiveElement.current = document.activeElement as HTMLElement;
      setTimeout(() => {
        titleInputRef.current?.focus();
      }, 50);
    }
    return () => {
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, [isOpen]);

  // Focus trap implementation
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== "Tab") return;

      const modal = modalRef.current;
      if (!modal) return;

      const focusables = modal.querySelectorAll<HTMLElement>(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      if (focusables.length === 0) return;

      const firstElement = focusables[0];
      const lastElement = focusables[focusables.length - 1];

      if (e.shiftKey) {
        if (document.activeElement === firstElement) {
          lastElement.focus();
          e.preventDefault();
        }
      } else {
        if (document.activeElement === lastElement) {
          firstElement.focus();
          e.preventDefault();
        }
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [isOpen]);

  if (!isOpen) return null;

  const resetModalState = () => {
    setNewTagName('');
    setNewTagColor(PRESET_COLORS[0]);
    setSearchQuery('');
    setEditingTagId(null);
    setEditingTagName('');
    setEditingTagColor('');
  };

  const handleClose = () => {
    resetModalState();
    onClose();
  };

  const handleCreateTag = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newTagName.trim()) return;

    try {
      await createMutation.mutateAsync({
        name: newTagName.trim(),
        color: newTagColor,
      });
      setNewTagName("");
      toast.success("Tag created successfully");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to create tag");
    }
  };

  const startEditing = (tagId: string, name: string, color: string) => {
    setEditingTagId(tagId);
    setEditingTagName(name);
    setEditingTagColor(color);
  };

  const cancelEditing = () => {
    setEditingTagId(null);
    setEditingTagName("");
    setEditingTagColor("");
  };

  const handleUpdateTag = async (tagId: string) => {
    if (!editingTagName.trim()) return;

    try {
      await updateMutation.mutateAsync({
        tagId,
        request: {
          name: editingTagName.trim(),
          color: editingTagColor,
        },
      });
      setEditingTagId(null);
      toast.success("Tag updated successfully");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update tag");
    }
  };

  const handleDeleteClick = (tagId: string, name: string) => {
    setTagToDelete({ id: tagId, name });
  };

  const confirmDeleteTag = async () => {
    if (!tagToDelete) return;
    try {
      await deleteMutation.mutateAsync(tagToDelete.id);
      toast.success(`Tag "${tagToDelete.name}" deleted successfully`);
      setTagToDelete(null);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete tag");
    }
  };

  const filteredTags = tags.filter((t) =>
    t.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <>
      <div className={styles.overlay} onClick={handleClose}>
        <div
          ref={modalRef}
          className={styles.modal}
          onClick={(e) => e.stopPropagation()}
          role="dialog"
          aria-modal="true"
          aria-labelledby="tags-manager-title"
        >
          <div className={styles.header}>
            <h2 id="tags-manager-title" className={styles.title}>
              MANAGE PROJECT TAGS
            </h2>
            <button className={styles.closeBtn} onClick={handleClose} aria-label="Close modal">
              <X size={16} />
            </button>
          </div>

          <div className={styles.content}>
            {/* Create Tag Section */}
            <form onSubmit={handleCreateTag} className={styles.createSection}>
              <span className={styles.sectionLabel}>Create New Tag</span>
              <div className={styles.inputGroup}>
                <input
                  ref={titleInputRef}
                  type="text"
                  placeholder="Tag name (e.g., frontend, bug)..."
                  className={styles.input}
                  value={newTagName}
                  onChange={(e) => setNewTagName(e.target.value)}
                  maxLength={40}
                  required
                />
                <button
                  type="submit"
                  className={styles.btnCreate}
                  disabled={createMutation.isPending || !newTagName.trim()}
                >
                  <Plus size={14} className="inline mr-1" />
                  Create
                </button>
              </div>

              {/* Color list */}
              <div className={styles.colorPalette}>
                {PRESET_COLORS.map((c) => (
                  <button
                    key={c}
                    type="button"
                    className={`${styles.colorBtn} ${newTagColor === c ? styles.colorBtnActive : ""}`}
                    style={{ backgroundColor: c }}
                    onClick={() => setNewTagColor(c)}
                    title={`Select color ${c}`}
                  />
                ))}
              </div>
            </form>

            {/* List Tags Section */}
            <div className={styles.listSection}>
              <span className={styles.sectionLabel}>All Tags ({tags.length})</span>
              <input
                type="text"
                placeholder="Search tags by name..."
                className={styles.searchBar}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />

              <div className={styles.tagList}>
                {isLoading ? (
                  <div className="text-center py-6 text-xs text-muted-foreground font-mono">
                    LOADING TAGS...
                  </div>
                ) : filteredTags.length === 0 ? (
                  <div className={styles.noTags}>
                    {searchQuery ? "No tags match your search query." : "No tags recorded for this project yet."}
                  </div>
                ) : (
                  filteredTags.map((tag) => {
                    const isEditing = editingTagId === tag.id;
                    const tagStyle = {
                      backgroundColor: `${tag.color || "#8b5cf6"}15`,
                      color: tag.color || "#8b5cf6",
                      borderColor: `${tag.color || "#8b5cf6"}30`,
                    };

                    return (
                      <div key={tag.id} className={styles.tagRow}>
                        {isEditing ? (
                          <div className={styles.editForm}>
                            <input
                              type="text"
                              className={styles.input}
                              value={editingTagName}
                              onChange={(e) => setEditingTagName(e.target.value)}
                              maxLength={40}
                              autoFocus
                            />
                            <div className={styles.inlineColors}>
                              {PRESET_COLORS.map((c) => (
                                <button
                                  key={c}
                                  type="button"
                                  className={`${styles.colorBtn} ${editingTagColor === c ? styles.colorBtnActive : ""}`}
                                  style={{ backgroundColor: c, width: "16px", height: "16px" }}
                                  onClick={() => setEditingTagColor(c)}
                                  title={`Select color ${c}`}
                                />
                              ))}
                            </div>
                            <button
                              type="button"
                              className={styles.actionBtn}
                              onClick={() => handleUpdateTag(tag.id)}
                              disabled={updateMutation.isPending || !editingTagName.trim()}
                              title="Save Changes"
                            >
                              <Check size={12} className="text-emerald-500" />
                            </button>
                            <button
                              type="button"
                              className={styles.actionBtn}
                              onClick={cancelEditing}
                              title="Cancel"
                            >
                              <RotateCcw size={12} />
                            </button>
                          </div>
                        ) : (
                          <>
                            <div className={styles.tagLabel}>
                              <span className={styles.tagBadge} style={tagStyle}>
                                <span
                                  style={{
                                    width: "6px",
                                    height: "6px",
                                    borderRadius: "50%",
                                    backgroundColor: tag.color || "#8b5cf6",
                                  }}
                                />
                                <span>{tag.name}</span>
                              </span>
                            </div>

                            <div className={styles.tagActions}>
                              <button
                                type="button"
                                className={styles.actionBtn}
                                onClick={() => startEditing(tag.id, tag.name, tag.color || "#8b5cf6")}
                                title="Edit Tag"
                              >
                                <Edit2 size={12} />
                              </button>
                              <button
                                type="button"
                                className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                                onClick={() => handleDeleteClick(tag.id, tag.name)}
                                title="Delete Tag"
                              >
                                <Trash2 size={12} />
                              </button>
                            </div>
                          </>
                        )}
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      <ConfirmModal
        isOpen={!!tagToDelete}
        onClose={() => setTagToDelete(null)}
        onConfirm={confirmDeleteTag}
        title="Delete Tag"
        message="Are you sure you want to permanently delete the tag"
        itemName={tagToDelete?.name}
        warningText="This will disassociate this tag from all snippets, problems, credentials, notes, and links across the project"
        isLoading={deleteMutation.isPending}
      />
    </>
  );
};
