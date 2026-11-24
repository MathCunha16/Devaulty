import React, { useState, useEffect, useRef, useMemo, Suspense } from "react";
import { X, Loader2, Code2, Eye } from "lucide-react";
import { toast } from "sonner";
import { MarkdownRenderer } from "../../../components/MarkdownRenderer";
import {
  useCreateProblemMutation,
  useUpdateProblemMutation,
  useProblemQuery,
} from "../hooks/useProblems";
import type { ProblemStatus, ProblemSeverity } from "~types/api";
import { useAutoResize } from "../../../hooks/useAutoResize";
import { useDiscardGuard } from "../../../hooks/useDiscardGuard";
import { DiscardConfirmModal } from "../../../components/DiscardConfirmModal";
import styles from "./ProblemForm.module.css";

interface ProblemFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  problemId?: string;
  projectColor?: string;
}

interface ProblemFormValues {
  title: string;
  errorDescription: string;
  solution: string;
  status: ProblemStatus;
  severity: ProblemSeverity;
}

interface ProblemFormInnerProps {
  title: string;
  initialValues?: ProblemFormValues;
  onSubmit: (values: ProblemFormValues) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
  projectColor?: string;
}

const ProblemFormInner: React.FC<ProblemFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
  projectColor,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [errorDescription, setErrorDescription] = useState(initialValues?.errorDescription || "");
  const [solution, setSolution] = useState(initialValues?.solution || "");
  const [status, setStatus] = useState<ProblemStatus>(initialValues?.status || "OPEN");
  const [severity, setSeverity] = useState<ProblemSeverity>(initialValues?.severity || "MEDIUM");
  const [errorTab, setErrorTab] = useState<"write" | "preview">("write");
  const [solutionTab, setSolutionTab] = useState<"write" | "preview">("write");

  const errorDescRef = useAutoResize(errorDescription, 100);
  const solutionRef = useAutoResize(solution, 100);

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

  const isDirty = useMemo(() => {
    const initTitle = initialValues?.title || "";
    const initErr = initialValues?.errorDescription || "";
    const initSol = initialValues?.solution || "";
    const initStatus = initialValues?.status || "OPEN";
    const initSev = initialValues?.severity || "MEDIUM";

    return (
      formTitle.trim() !== initTitle.trim() ||
      errorDescription.trim() !== initErr.trim() ||
      solution.trim() !== initSol.trim() ||
      status !== initStatus ||
      severity !== initSev
    );
  }, [formTitle, errorDescription, solution, status, severity, initialValues]);

  const {
    isConfirmDiscardOpen,
    handleRequestClose,
    handleConfirmDiscard,
    handleCancelDiscard,
  } = useDiscardGuard({ isDirty, onClose });

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
  }, [handleRequestClose, handleCancelDiscard, isConfirmDiscardOpen, isSubmitting]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formTitle.trim()) {
      toast.error("Title is required");
      return;
    }
    onSubmit({
      title: formTitle,
      errorDescription,
      solution,
      status,
      severity,
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
        aria-labelledby="problem-form-title"
      >
        <div className={styles.header}>
          <h2 id="problem-form-title" className={styles.title}>{title}</h2>
          <button
            type="button"
            className={styles.closeBtn}
            onClick={handleRequestClose}
            disabled={isSubmitting}
            aria-label="Close modal"
          >
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="problem-title" className={styles.label}>Title</label>
            <input
              ref={firstInputRef}
              id="problem-title"
              type="text"
              className={styles.input}
              placeholder="e.g., NullPointerException on Auth flow, Memory leak in cache"
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
              required
            />
          </div>

          <div className={styles.row}>
            <div className={styles.field}>
              <label htmlFor="problem-severity" className={styles.label}>Severity</label>
              <select
                id="problem-severity"
                className={styles.input}
                value={severity}
                onChange={(e) => setSeverity(e.target.value as ProblemSeverity)}
                disabled={isSubmitting}
              >
                <option value="LOW">Low</option>
                <option value="MEDIUM">Medium</option>
                <option value="HIGH">High</option>
                <option value="CRITICAL">Critical</option>
              </select>
            </div>

            <div className={styles.field}>
              <label htmlFor="problem-status" className={styles.label}>Status</label>
              <select
                id="problem-status"
                className={styles.input}
                value={status}
                onChange={(e) => setStatus(e.target.value as ProblemStatus)}
                disabled={isSubmitting}
              >
                <option value="OPEN">Open</option>
                <option value="WORKING_ON">Working On</option>
                <option value="RESOLVED">Resolved</option>
                <option value="WONT_FIX">Won't Fix</option>
              </select>
            </div>
          </div>

          <div className={styles.field}>
            <div className={styles.fieldHeader}>
              <label htmlFor="problem-error-description" className={styles.label}>
                Error / Stack Trace Logs
              </label>
              <div className="flex items-center gap-1 p-0.5 rounded-md bg-secondary/80 border border-border/80">
                <button
                  type="button"
                  onClick={() => setErrorTab("write")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    errorTab === "write"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Code2 size={12} />
                  <span>Write / Raw</span>
                </button>
                <button
                  type="button"
                  onClick={() => setErrorTab("preview")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    errorTab === "preview"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Eye size={12} />
                  <span>Preview</span>
                </button>
              </div>
            </div>

            {errorTab === "write" ? (
              <textarea
                id="problem-error-description"
                ref={errorDescRef}
                className={styles.textarea}
                placeholder="Paste raw stack traces, logs, or markdown with ```mermaid diagrams..."
                value={errorDescription}
                onChange={(e) => setErrorDescription(e.target.value)}
                disabled={isSubmitting}
              />
            ) : (
              <div className={styles.previewContainer}>
                {errorDescription.trim() ? (
                  <MarkdownRenderer content={errorDescription} />
                ) : (
                  <span className="text-muted-foreground italic text-xs font-mono">
                    No error details written yet. Switch to "Write / Raw" to add logs or markdown.
                  </span>
                )}
              </div>
            )}
          </div>

          <div className={styles.field}>
            <div className={styles.fieldHeader}>
              <label htmlFor="problem-solution" className={styles.label}>
                Solution Code / Script
              </label>
              <div className="flex items-center gap-1 p-0.5 rounded-md bg-secondary/80 border border-border/80">
                <button
                  type="button"
                  onClick={() => setSolutionTab("write")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    solutionTab === "write"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Code2 size={12} />
                  <span>Write / Raw</span>
                </button>
                <button
                  type="button"
                  onClick={() => setSolutionTab("preview")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    solutionTab === "preview"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Eye size={12} />
                  <span>Preview</span>
                </button>
              </div>
            </div>

            {solutionTab === "write" ? (
              <textarea
                id="problem-solution"
                ref={solutionRef}
                className={styles.textarea}
                placeholder="Paste your fix, resolution script, or markdown notes with ```mermaid..."
                value={solution}
                onChange={(e) => setSolution(e.target.value)}
                disabled={isSubmitting}
              />
            ) : (
              <div className={styles.previewContainer}>
                {solution.trim() ? (
                  <MarkdownRenderer content={solution} />
                ) : (
                  <span className="text-muted-foreground italic text-xs font-mono">
                    No solution written yet. Switch to "Write / Raw" to add code or markdown.
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
              {title.includes("EDIT") ? "Save Changes" : "Create Problem Node"}
            </button>
          </div>
        </form>
      </div>

      <DiscardConfirmModal
        isOpen={isConfirmDiscardOpen}
        onClose={handleCancelDiscard}
        onDiscard={handleConfirmDiscard}
        itemName="problem"
      />
    </div>
  );
};

const CreateProblemFormModal: React.FC<{
  projectId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, onClose, projectColor }) => {
  const createMutation = useCreateProblemMutation(projectId);

  const handleSubmit = async (values: ProblemFormValues) => {
    try {
      await createMutation.mutateAsync({
        title: values.title,
        errorDescription: values.errorDescription || undefined,
        solution: values.solution || undefined,
        status: values.status,
        severity: values.severity,
      });
      toast.success("Problem diagnostic created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create problem diagnostic");
    }
  };

  return (
    <ProblemFormInner
      title="LOG DIAGNOSTIC PROBLEM"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
      projectColor={projectColor}
    />
  );
};

const EditProblemFormModal: React.FC<{
  projectId: string;
  problemId: string;
  onClose: () => void;
  projectColor?: string;
}> = ({ projectId, problemId, onClose, projectColor }) => {
  const { data: problem, isLoading, isError } = useProblemQuery(projectId, problemId);
  const updateMutation = useUpdateProblemMutation(projectId, problemId);

  const handleSubmit = async (values: ProblemFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        errorDescription: values.errorDescription || undefined,
        solution: values.solution || undefined,
        severity: values.severity,
      });
      toast.success("Problem diagnostic updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update problem diagnostic");
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
            <span className="text-xs text-muted-foreground font-mono">LOADING DIAGNOSTICS...</span>
          </div>
        </div>
      </div>
    );
  }

  if (isError || !problem) {
    return (
      <div
        className={styles.overlay}
        onClick={onClose}
        style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      >
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3 text-destructive font-mono text-xs">
            <span>FAILED TO LOAD DIAGNOSTICS DATA.</span>
          </div>
        </div>
      </div>
    );
  }

  const initialValues: ProblemFormValues = {
    title: problem.title || "",
    errorDescription: problem.errorDescription || "",
    solution: problem.solution || "",
    status: problem.status || "OPEN",
    severity: problem.severity || "MEDIUM",
  };

  return (
    <ProblemFormInner
      title="EDIT DIAGNOSTIC PROBLEM"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
      projectColor={projectColor}
    />
  );
};

export const ProblemForm: React.FC<ProblemFormProps> = ({
  isOpen,
  onClose,
  projectId,
  problemId,
  projectColor,
}) => {
  if (!isOpen) return null;

  if (problemId) {
    return (
      <Suspense
        fallback={
          <div
            className={styles.overlay}
            onClick={onClose}
            style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
          >
            <div className={styles.modal}>
              <div className="flex flex-col items-center justify-center p-12 gap-3">
                <Loader2 className="animate-spin text-primary" size={28} />
                <span className="text-xs text-muted-foreground font-mono">LOADING DIAGNOSTICS...</span>
              </div>
            </div>
          </div>
        }
      >
        <EditProblemFormModal
          projectId={projectId}
          problemId={problemId}
          onClose={onClose}
          projectColor={projectColor}
        />
      </Suspense>
    );
  }

  return (
    <CreateProblemFormModal
      projectId={projectId}
      onClose={onClose}
      projectColor={projectColor}
    />
  );
};
