import React, { useState, useEffect, useRef, useMemo } from "react";
import { X, Loader2, Code2, Eye } from "lucide-react";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { toast } from "sonner";
import {
  useCreateNoteMutation,
  useUpdateNoteMutation,
  useNoteQuery,
} from "../hooks/useNotes";
import { useAutoResize } from "../../../hooks/useAutoResize";
import { useDiscardGuard } from "../../../hooks/useDiscardGuard";
import { DiscardConfirmModal } from "../../../components/DiscardConfirmModal";
import styles from "./NoteForm.module.css";

interface NoteFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  noteId?: string;
  projectColor?: string;
}

interface NoteFormValues {
  title: string;
  content: string;
}

interface NoteFormInnerProps {
  title: string;
  initialValues?: NoteFormValues;
  onSubmit: (values: NoteFormValues) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
  projectColor?: string;
}

const NoteFormInner: React.FC<NoteFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
  projectColor,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [content, setContent] = useState(initialValues?.content || "");
  const [contentTab, setContentTab] = useState<"write" | "preview">(
    initialValues?.content ? "preview" : "write"
  );

  const contentRef = useAutoResize(content, 180);

  const isDirty = useMemo(() => {
    const initTitle = initialValues?.title || "";
    const initContent = initialValues?.content || "";

    return (
      formTitle !== initTitle ||
      content !== initContent
    );
  }, [formTitle, content, initialValues]);

  const {
    isConfirmDiscardOpen,
    handleRequestClose,
    handleConfirmDiscard,
    handleCancelDiscard,
  } = useDiscardGuard({ isDirty, onClose });

  const renderPreviewHtml = () => {
    if (!content.trim()) return "";
    try {
      const rawHtml = marked.parse(content, { breaks: true, gfm: true }) as string;
      return DOMPurify.sanitize(rawHtml);
    } catch {
      return DOMPurify.sanitize(content);
    }
  };

  // Focus trap refs
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    previousActiveElement.current = document.activeElement as HTMLElement;
    const timer = setTimeout(() => {
      firstInputRef.current?.focus();
    }, 50);

    return () => {
      clearTimeout(timer);
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, []);

  // Trap focus inside modal
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        if (isConfirmDiscardOpen) {
          handleCancelDiscard();
        } else if (!isSubmitting) {
          handleRequestClose();
        }
        return;
      }
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
  }, [handleCancelDiscard, handleRequestClose, isConfirmDiscardOpen, isSubmitting]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formTitle.trim()) {
      toast.error("Title is required");
      return;
    }
    onSubmit({
      title: formTitle,
      content,
    });
  };

  return (
    <div
      className={styles.overlay}
      onClick={() => !isSubmitting && handleRequestClose()}
      style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
    >
      <div
        ref={modalRef}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="note-form-title"
      >
        <div className={styles.header}>
          <h2 id="note-form-title" className={styles.title}>{title}</h2>
          <button className={styles.closeBtn} onClick={handleRequestClose} disabled={isSubmitting} aria-label="Close modal">
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="note-title" className={styles.label}>Title</label>
            <input
              ref={firstInputRef}
              id="note-title"
              type="text"
              className={styles.input}
              placeholder="e.g., Deploy guidelines, meeting notes, database setup..."
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
              required
            />
          </div>

          <div className={styles.field}>
            <div className={styles.fieldHeader}>
              <label htmlFor="note-content" className={styles.label}>Content / Body</label>
              <div className="flex items-center gap-1 p-0.5 rounded-md bg-secondary/80 border border-border/80">
                <button
                  type="button"
                  onClick={() => setContentTab("write")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    contentTab === "write"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Code2 size={12} />
                  <span>Write / Code</span>
                </button>
                <button
                  type="button"
                  onClick={() => setContentTab("preview")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    contentTab === "preview"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Eye size={12} />
                  <span>Preview</span>
                </button>
              </div>
            </div>

            {contentTab === "write" ? (
              <textarea
                id="note-content"
                ref={contentRef}
                className={styles.textarea}
                placeholder="Write your markdown notes, reminders, or document outlines here..."
                value={content}
                onChange={(e) => setContent(e.target.value)}
                disabled={isSubmitting}
              />
            ) : (
              <div className={styles.previewContainer}>
                {content.trim() ? (
                  <div dangerouslySetInnerHTML={{ __html: renderPreviewHtml() }} />
                ) : (
                  <span className="text-muted-foreground italic text-xs font-mono">
                    No content written yet. Switch to "Write / Code" to add markdown details.
                  </span>
                )}
              </div>
            )}
          </div>

          <div className={styles.footer}>
            <button
              type="button"
              className={styles.btn}
              onClick={handleRequestClose}
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.btnPrimary}
              disabled={isSubmitting}
            >
              {title.includes("EDIT") ? "Save Changes" : "Create Note"}
            </button>
          </div>
        </form>
      </div>

      <DiscardConfirmModal
        isOpen={isConfirmDiscardOpen}
        onClose={handleCancelDiscard}
        onDiscard={handleConfirmDiscard}
        itemName="note"
      />
    </div>
  );
};

const CreateNoteFormModal: React.FC<{
  projectId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, onClose, projectColor }) => {
  const createMutation = useCreateNoteMutation(projectId);

  const handleSubmit = async (values: NoteFormValues) => {
    try {
      await createMutation.mutateAsync({
        title: values.title,
        content: values.content || undefined,
      });
      toast.success("Note created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create note");
    }
  };

  return (
    <NoteFormInner
      title="CREATE NEW NOTE"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
      projectColor={projectColor}
    />
  );
};

const EditNoteFormModal: React.FC<{
  projectId: string;
  noteId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, noteId, onClose, projectColor }) => {
  const { data: note, isLoading, isError } = useNoteQuery(projectId, noteId);
  const updateMutation = useUpdateNoteMutation(projectId, noteId);

  const handleSubmit = async (values: NoteFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        content: values.content,
      });
      toast.success("Note updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update note");
    }
  };

  if (isLoading) {
    return (
      <div
        className={styles.overlay}
        onClick={onClose}
        style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      >
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3">
            <Loader2 className="animate-spin text-primary" size={28} />
            <span className="text-xs text-muted-foreground font-mono">LOADING NOTE...</span>
          </div>
        </div>
      </div>
    );
  }

  if (isError || !note) {
    return (
      <div
        className={styles.overlay}
        onClick={onClose}
        style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      >
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3 text-destructive font-mono text-xs">
            <span>FAILED TO LOAD NOTE.</span>
          </div>
        </div>
      </div>
    );
  }

  const initialValues: NoteFormValues = {
    title: note.title || "",
    content: note.content || "",
  };

  return (
    <NoteFormInner
      title="EDIT SYSTEM NOTE"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
      projectColor={projectColor}
    />
  );
};

export const NoteForm: React.FC<NoteFormProps> = ({
  isOpen,
  onClose,
  projectId,
  noteId,
  projectColor,
}) => {
  if (!isOpen) return null;

  if (noteId) {
    return (
      <EditNoteFormModal
        projectId={projectId}
        noteId={noteId}
        onClose={onClose}
        projectColor={projectColor}
      />
    );
  }

  return (
    <CreateNoteFormModal
      projectId={projectId}
      onClose={onClose}
      projectColor={projectColor}
    />
  );
};
