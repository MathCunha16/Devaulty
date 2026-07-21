import React, { useState, Suspense } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  useCreateProblemMutation,
  useUpdateProblemMutation,
  useProblemQuery,
} from "../hooks/useProblems";
import type { ProblemStatus, ProblemSeverity } from "~types/api";
import { useAutoResize } from "../../../hooks/useAutoResize";
import styles from "./ProblemForm.module.css";

interface ProblemFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  problemId?: string;
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
}

const ProblemFormInner: React.FC<ProblemFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [errorDescription, setErrorDescription] = useState(initialValues?.errorDescription || "");
  const [solution, setSolution] = useState(initialValues?.solution || "");
  const [status, setStatus] = useState<ProblemStatus>(initialValues?.status || "OPEN");
  const [severity, setSeverity] = useState<ProblemSeverity>(initialValues?.severity || "MEDIUM");

  const errorDescRef = useAutoResize(errorDescription, 100);
  const solutionRef = useAutoResize(solution, 100);

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
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.header}>
          <h2 className={styles.title}>{title}</h2>
          <button className={styles.closeBtn} onClick={onClose} disabled={isSubmitting}>
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label className={styles.label}>Title</label>
            <input
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
              <label className={styles.label}>Severity</label>
              <select
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
              <label className={styles.label}>Status</label>
              <select
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
            <label className={styles.label}>Error / Stack Trace Logs</label>
            <textarea
              ref={errorDescRef}
              className={styles.textarea}
              placeholder="Paste stack traces, logs, or error details here..."
              value={errorDescription}
              onChange={(e) => setErrorDescription(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Solution Code / Script</label>
            <textarea
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
              {title.includes("EDIT") ? "Save Changes" : "Create Problem Node"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const CreateProblemFormModal: React.FC<{ projectId: string; onClose: () => void }> = ({
  projectId,
  onClose,
}) => {
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
      toast.success("Problem node created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create problem node");
    }
  };

  return (
    <ProblemFormInner
      title="NEW DIAGNOSTIC PROBLEM"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
    />
  );
};

const EditProblemFormModal: React.FC<{
  projectId: string;
  problemId: string;
  onClose: () => void;
}> = ({ projectId, problemId, onClose }) => {
  const { data: problem } = useProblemQuery(projectId, problemId);
  const updateMutation = useUpdateProblemMutation(projectId, problemId);

  const handleSubmit = async (values: ProblemFormValues) => {
    try {
      await updateMutation.mutateAsync({
        title: values.title,
        errorDescription: values.errorDescription || undefined,
        solution: values.solution || undefined,
        severity: values.severity,
      });
      toast.success("Problem updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update problem");
    }
  };

  const initialValues: ProblemFormValues = {
    title: problem?.title || "",
    errorDescription: problem?.errorDescription || "",
    solution: problem?.solution || "",
    status: problem?.status || "OPEN",
    severity: problem?.severity || "MEDIUM",
  };

  return (
    <ProblemFormInner
      title="EDIT DIAGNOSTIC PROBLEM"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
    />
  );
};

export const ProblemForm: React.FC<ProblemFormProps> = ({
  isOpen,
  onClose,
  projectId,
  problemId,
}) => {
  if (!isOpen) return null;

  if (problemId) {
    return (
      <Suspense
        fallback={
          <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal}>
              <div className="flex flex-col items-center justify-center p-12 gap-3">
                <Loader2 className="animate-spin text-primary" size={28} />
                <span className="text-xs text-muted-foreground font-mono">LOADING DIAGNOSTICS...</span>
              </div>
            </div>
          </div>
        }
      >
        <EditProblemFormModal projectId={projectId} problemId={problemId} onClose={onClose} />
      </Suspense>
    );
  }

  return <CreateProblemFormModal projectId={projectId} onClose={onClose} />;
};
