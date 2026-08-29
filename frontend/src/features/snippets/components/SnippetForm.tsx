import React, { useState, useEffect, useRef } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import Editor from "@monaco-editor/react";
import { useTheme } from "../../../hooks/useTheme";
import { useAutoResize } from "../../../hooks/useAutoResize";
import { LanguageSelect } from "./LanguageSelect";
import {
  useCreateSnippetMutation,
  useUpdateSnippetMutation,
  useSnippetQuery,
} from "../hooks/useSnippets";
import { mapLanguageToMonaco, ALL_LANGUAGES } from "../utils/languageUtils";
import type { SnippetLanguage, SnippetType } from "~types/api";
import styles from "./SnippetForm.module.css";

interface SnippetFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  snippetId?: string;
  projectColor?: string;
}

interface SnippetFormValues {
  title: string;
  description?: string;
  content: string;
  language: SnippetLanguage;
  snippetType: SnippetType;
}

interface SnippetFormInnerProps {
  title: string;
  initialValues?: SnippetFormValues;
  onSubmit: (values: SnippetFormValues) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
  projectColor?: string;
}

const SnippetFormInner: React.FC<SnippetFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
  projectColor,
}) => {
  const { theme } = useTheme();
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [formDescription, setFormDescription] = useState(initialValues?.description || "");
  const [formContent, setFormContent] = useState(initialValues?.content || "");
  const [formLanguage, setFormLanguage] = useState<SnippetLanguage>(
    initialValues?.language || "PLAIN_TEXT"
  );
  const [formSnippetType, setFormSnippetType] = useState<SnippetType>(
    initialValues?.snippetType || "CODE"
  );

  const descRef = useAutoResize(formDescription, 60);

  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const editorRef = useRef<any>(null);
  const layoutTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

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

  // Dispose Monaco editor and model on unmount to prevent WebKitGTK leaks
  useEffect(() => {
    return () => {
      if (layoutTimerRef.current) {
        clearTimeout(layoutTimerRef.current);
        layoutTimerRef.current = null;
      }
      if (editorRef.current) {
        try {
          editorRef.current.getModel()?.dispose();
          editorRef.current.dispose();
        } catch {
          // ignore if already disposed
        }
        editorRef.current = null;
      }
    };
  }, []);

  // Manual layout calculation for automaticLayout: false
  useEffect(() => {
    const handleResize = () => {
      editorRef.current?.layout();
    };
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        if (!isSubmitting) onClose();
        return;
      }
      if (e.key === "Tab") {
        const modal = modalRef.current;
        if (!modal) return;

        // If focus is currently inside Monaco editor, let Monaco handle Tab key natively
        if (document.activeElement?.closest(".monaco-editor")) return;

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
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [onClose, isSubmitting]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formTitle.trim() || !formContent.trim()) {
      toast.error("Title and Content are required");
      return;
    }
    onSubmit({
      title: formTitle,
      description: formDescription || undefined,
      content: formContent,
      language: formLanguage,
      snippetType: formSnippetType,
    });
  };

  return (
    <div
      className={styles.overlay}
      onClick={() => !isSubmitting && onClose()}
      style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
    >
      <div
        ref={modalRef}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="snippet-form-title"
      >
        <div className={styles.header}>
          <h2 id="snippet-form-title" className={styles.title}>
            {title}
          </h2>
          <button
            type="button"
            className={styles.closeBtn}
            onClick={onClose}
            disabled={isSubmitting}
            aria-label="Close modal"
          >
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="snippet-title" className={styles.label}>
              Title
            </label>
            <input
              ref={firstInputRef}
              id="snippet-title"
              type="text"
              className={styles.input}
              placeholder="e.g., Get user list API request, Database backup script..."
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              disabled={isSubmitting}
              required
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="snippet-desc" className={styles.label}>
              Description
            </label>
            <textarea
              id="snippet-desc"
              ref={descRef}
              className={styles.textarea}
              placeholder="Provide a brief context or description..."
              value={formDescription}
              onChange={(e) => setFormDescription(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <div className={styles.row}>
            <div className={styles.field}>
              <label htmlFor="snippet-type" className={styles.label}>
                Type
              </label>
              <select
                id="snippet-type"
                className={styles.select}
                value={formSnippetType}
                onChange={(e) => setFormSnippetType(e.target.value as SnippetType)}
                disabled={isSubmitting}
              >
                <option value="CODE">Code Snippet</option>
                <option value="COMMAND">Terminal Command</option>
              </select>
            </div>

            <div className={styles.field}>
              <label id="snippet-language-label" className={styles.label}>
                Language
              </label>
              <LanguageSelect
                triggerId="snippet-language-trigger"
                ariaLabelledBy="snippet-language-label"
                languages={ALL_LANGUAGES}
                value={formLanguage}
                onChange={setFormLanguage}
                disabled={isSubmitting}
              />
            </div>
          </div>

          <div className={styles.field}>
            <label id="snippet-content-label" className={styles.label}>
              Source Code Content
            </label>
            <div className={styles.editorWrapper}>
              <Editor
                height="280px"
                language={mapLanguageToMonaco(formLanguage)}
                theme={theme === "dark" ? "vs-dark" : "light"}
                value={formContent}
                onChange={(val) => setFormContent(val || "")}
                onMount={(editor) => {
                  editorRef.current = editor;
                  if (layoutTimerRef.current) clearTimeout(layoutTimerRef.current);
                  layoutTimerRef.current = setTimeout(() => {
                    if (editorRef.current) {
                      editor.layout();
                    }
                  }, 100);
                }}
                loading={
                  <div className="flex items-center justify-center h-48 text-xs text-muted-foreground font-mono">
                    Loading Editor Environment...
                  </div>
                }
                options={{
                  ariaLabel: "Source Code Content",
                  minimap: { enabled: false },
                  fontSize: 13,
                  lineNumbers: "on",
                  scrollBeyondLastLine: false,
                  automaticLayout: false,
                  readOnly: isSubmitting,
                  padding: { top: 8, bottom: 8 },
                }}
              />
            </div>
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
              {title.includes("EDIT") ? "Save Changes" : "Create Snippet"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const CreateSnippetFormModal: React.FC<{
  projectId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, onClose, projectColor }) => {
  const createMutation = useCreateSnippetMutation(projectId);

  const handleSubmit = async (values: SnippetFormValues) => {
    try {
      await createMutation.mutateAsync(values);
      toast.success("Snippet created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create snippet");
    }
  };

  return (
    <SnippetFormInner
      title="CREATE NEW SNIPPET"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
      projectColor={projectColor}
    />
  );
};

const EditSnippetFormModal: React.FC<{
  projectId: string;
  snippetId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, snippetId, onClose, projectColor }) => {
  const { data: snippet, isLoading, isError } = useSnippetQuery(projectId, snippetId);
  const updateMutation = useUpdateSnippetMutation(projectId, snippetId);

  const handleSubmit = async (values: SnippetFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        description: values.description,
        content: values.content,
        language: values.language,
        snippetType: values.snippetType,
      });
      toast.success("Snippet updated successfully");
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
            <span className="text-xs text-muted-foreground font-mono">LOADING SNIPPET...</span>
          </div>
        </div>
      </div>
    );
  }

  if (isError || !snippet) {
    return (
      <div
        className={styles.overlay}
        onClick={onClose}
        style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      >
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3 text-destructive font-mono text-xs">
            <span>FAILED TO LOAD SNIPPET.</span>
          </div>
        </div>
      </div>
    );
  }

  const initialValues: SnippetFormValues = {
    title: snippet.title || "",
    description: snippet.description || "",
    content: snippet.content || "",
    language: snippet.language || "PLAIN_TEXT",
    snippetType: snippet.snippetType || "CODE",
  };

  return (
    <SnippetFormInner
      title="EDIT SNIPPET"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
      projectColor={projectColor}
    />
  );
};

export const SnippetForm: React.FC<SnippetFormProps> = ({
  isOpen,
  onClose,
  projectId,
  snippetId,
  projectColor,
}) => {
  if (!isOpen) return null;

  if (snippetId) {
    return (
      <EditSnippetFormModal
        projectId={projectId}
        snippetId={snippetId}
        onClose={onClose}
        projectColor={projectColor}
      />
    );
  }

  return (
    <CreateSnippetFormModal
      projectId={projectId}
      onClose={onClose}
      projectColor={projectColor}
    />
  );
};
