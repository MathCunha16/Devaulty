import React, { useState, useEffect, useRef, Suspense } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  useCreateLinkMutation,
  useUpdateLinkMutation,
  useLinkQuery,
} from "../hooks/useLinks";
import { useAutoResize } from "../../../hooks/useAutoResize";
import styles from "./LinkForm.module.css";

interface LinkFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  linkId?: string;
}

interface LinkFormValues {
  title: string;
  url: string;
  description: string;
}

interface LinkFormInnerProps {
  title: string;
  initialValues?: LinkFormValues;
  onSubmit: (values: LinkFormValues) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
}

const LinkFormInner: React.FC<LinkFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [url, setUrl] = useState(initialValues?.url || "");
  const [description, setDescription] = useState(initialValues?.description || "");

  const descRef = useAutoResize(description, 100);

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
    if (!url.trim()) {
      toast.error("URL is required");
      return;
    }
    onSubmit({
      title: formTitle,
      url,
      description,
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
        aria-labelledby="link-form-title"
      >
        <div className={styles.header}>
          <h2 id="link-form-title" className={styles.title}>{title}</h2>
          <button className={styles.closeBtn} onClick={onClose} disabled={isSubmitting} aria-label="Close modal">
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="link-title" className={styles.label}>Title</label>
            <input
              ref={firstInputRef}
              id="link-title"
              type="text"
              className={styles.input}
              placeholder="e.g., API documentation, Staging environment, Figma board..."
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
              required
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="link-url" className={styles.label}>URL</label>
            <input
              id="link-url"
              type="url"
              className={styles.input}
              placeholder="https://example.com"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              disabled={isSubmitting}
              required
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="link-description" className={styles.label}>Description</label>
            <textarea
              id="link-description"
              ref={descRef}
              className={styles.textarea}
              placeholder="Brief notes about this link, credentials, or instructions..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
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
              {title.includes("EDIT") ? "Save Changes" : "Create Link"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const CreateLinkFormModal: React.FC<{ projectId: string; onClose: () => void }> = ({
  projectId,
  onClose,
}) => {
  const createMutation = useCreateLinkMutation(projectId);

  const handleSubmit = async (values: LinkFormValues) => {
    try {
      await createMutation.mutateAsync({
        title: values.title,
        url: values.url,
        description: values.description || undefined,
      });
      toast.success("Link created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create link");
    }
  };

  return (
    <LinkFormInner
      title="NEW WEB LINK"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
    />
  );
};

const EditLinkFormModal: React.FC<{
  projectId: string;
  linkId: string;
  onClose: () => void;
}> = ({ projectId, linkId, onClose }) => {
  const { data: link, isLoading, isError } = useLinkQuery(projectId, linkId);
  const updateMutation = useUpdateLinkMutation(projectId, linkId);

  const handleSubmit = async (values: LinkFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        url: values.url,
        description: values.description || undefined,
      });
      toast.success("Link updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update link");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3">
            <Loader2 className="animate-spin text-primary" size={28} />
            <span className="text-xs text-muted-foreground font-mono">LOADING LINK...</span>
          </div>
        </div>
      </div>
    );
  }

  if (isError || !link) {
    return (
      <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3 text-destructive font-mono text-xs">
            <span>FAILED TO LOAD LINK.</span>
          </div>
        </div>
      </div>
    );
  }

  const initialValues: LinkFormValues = {
    title: link.title || "",
    url: link.url || "",
    description: link.description || "",
  };

  return (
    <LinkFormInner
      title="EDIT WEB LINK"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
    />
  );
};

export const LinkForm: React.FC<LinkFormProps> = ({
  isOpen,
  onClose,
  projectId,
  linkId,
}) => {
  if (!isOpen) return null;

  if (linkId) {
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
        <EditLinkFormModal projectId={projectId} linkId={linkId} onClose={onClose} />
      </Suspense>
    );
  }

  return <CreateLinkFormModal projectId={projectId} onClose={onClose} />;
};
