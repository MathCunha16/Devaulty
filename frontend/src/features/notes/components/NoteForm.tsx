import React, { useState, useEffect, useRef, Suspense } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  useCreateNoteMutation,
  useUpdateNoteMutation,
  useNoteQuery,
} from "../hooks/useNotes";
import { useAutoResize } from "../../../hooks/useAutoResize";
import styles from "./NoteForm.module.css";

interface NoteFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  noteId?: string;
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
}

const NoteFormInner: React.FC<NoteFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [content, setContent] = useState(initialValues?.content || "");

  const contentRef = useAutoResize(content, 180);

  // Focus trap refs
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    previousActiveElement.current = document.activeElement as HTMLElement;
    setTimeout(() => {
      firstInputRef.current?.focus();
    }, 50);

    return () => {
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, []);

  // Trap focus inside modal
  useEffect(() => {
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
  }, []);

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
    <div className={styles.overlay} onClick={onClose}>
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
          <button className={styles.closeBtn} onClick={onClose} disabled={isSubmitting} aria-label="Close modal">
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
            <label htmlFor="note-content" className={styles.label}>Content / Body</label>
            <textarea
              id="note-content"
              ref={contentRef}
              className={styles.textarea}
              placeholder="Write your markdown notes, reminders, or document outlines here..."
              value={content}
              onChange={(e) => setContent(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <div className={styles.footer}>
            <button
              type="button"
              className={styles.btn}
              onClick={onClose}
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
    </div>
  );
};

const CreateNoteFormModal: React.FC<{ projectId: string; onClose: () => void }> = ({
  projectId,
  onClose,
}) => {
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
      title="NEW SYSTEM NOTE"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
    />
  );
};

const EditNoteFormModal: React.FC<{
  projectId: string;
  noteId: string;
  onClose: () => void;
}> = ({ projectId, noteId, onClose }) => {
  const { data: note, isLoading, isError } = useNoteQuery(projectId, noteId);
  const updateMutation = useUpdateNoteMutation(projectId, noteId);

  const handleSubmit = async (values: NoteFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        content: values.content || undefined,
      });
      toast.success("Note updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update note");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.overlay} onClick={onClose}>
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
      <div className={styles.overlay} onClick={onClose}>
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
    />
  );
};

export const NoteForm: React.FC<NoteFormProps> = ({
  isOpen,
  onClose,
  projectId,
  noteId,
}) => {
  if (!isOpen) return null;

  if (noteId) {
    return (
      <Suspense
        fallback={
          <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal}>
              <div className="flex flex-col items-center justify-center p-12 gap-3">
                <Loader2 className="animate-spin text-primary" size={28} />
                <span className="text-xs text-muted-foreground font-mono">LOADING...</span>
              </div>
            </div>
          </div>
        }
      >
        <EditNoteFormModal projectId={projectId} noteId={noteId} onClose={onClose} />
      </Suspense>
    );
  }

  return <CreateNoteFormModal projectId={projectId} onClose={onClose} />;
};
