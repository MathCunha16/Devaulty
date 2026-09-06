import React, { useState, useEffect, useRef, useMemo, Suspense } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
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
            <label htmlFor="problem-error-description" className={styles.label}>Error / Stack Trace Logs</label>
            <textarea
              id="problem-error-description"
              ref={errorDescRef}
              className={styles.textarea}
              placeholder="Paste stack traces, logs, or error details here..."
              value={errorDescription}
              onChange={(e) => setErrorDescription(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="problem-solution" className={styles.label}>Solution Code / Script</label>
            <textarea
              id="problem-solution"
              ref={solutionRef}
              className={styles.textarea}
              placeholder="Paste your fix, resolution script, or notes here..."
              value={solution}
              onChange={(e) => setSolution(e.target.value)}
              disabled={isSubmitting}
            />
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
