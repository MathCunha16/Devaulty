import React, { useState, Suspense } from "react";
import { X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  useCreateProjectMutation,
  useUpdateProjectMutation,
  useProjectQuery,
} from "../hooks/useProjects";
import styles from "./ProjectForm.module.css";

interface ProjectFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId?: string;
}

interface ProjectFormValues {
  name: string;
  description: string;
  color: string;
  icon: string;
}

const PRESET_COLORS = [
  "#ef4444", // red
  "#f97316", // orange
  "#facc15", // yellow
  "#22c55e", // green
  "#06b6d4", // cyan
  "#3b82f6", // blue
  "#8b5cf6", // violet
  "#ec4899", // pink
];

const PRESET_ICONS = ["Folder", "Terminal", "Database", "Globe", "Cpu", "Activity", "BookOpen", "Code"];

interface ProjectFormInnerProps {
  title: string;
  initialValues?: ProjectFormValues;
  onSubmit: (values: ProjectFormValues) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
}

const ProjectFormInner: React.FC<ProjectFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
}) => {
  const [name, setName] = useState(initialValues?.name || "");
  const [description, setDescription] = useState(initialValues?.description || "");
  const [color, setColor] = useState(initialValues?.color || PRESET_COLORS[5]);
  const [icon, setIcon] = useState(initialValues?.icon || PRESET_ICONS[0]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({ name, description, color, icon });
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
            <label className={styles.label}>Name</label>
            <input
              type="text"
              className={styles.input}
              placeholder="e.g., Core API, Frontend"
              value={name}
              onChange={(e) => setName(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
              required
            />
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Description</label>
            <textarea
              className={styles.textarea}
              placeholder="Brief summary of this project..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
            />
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Icon</label>
            <div className="flex gap-2 overflow-x-auto py-1">
              {PRESET_ICONS.map((ico) => (
                <button
                  key={ico}
                  type="button"
                  className={`px-3 py-1.5 border rounded text-xs transition-colors ${
                    icon === ico
                      ? "bg-primary text-primary-foreground border-primary"
                      : "border-border hover:bg-secondary text-muted-foreground hover:text-foreground"
                  }`}
                  onClick={() => setIcon(ico)}
                  disabled={isSubmitting}
                >
                  {ico}
                </button>
              ))}
            </div>
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Accent Color</label>
            <div className={styles.colorPicker}>
              {PRESET_COLORS.map((col) => (
                <button
                  key={col}
                  type="button"
                  className={`${styles.colorOption} ${color === col ? styles.colorOptionActive : ""}`}
                  style={{ backgroundColor: col }}
                  onClick={() => setColor(col)}
                  disabled={isSubmitting}
                />
              ))}
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
              {title.includes("EDIT") ? "Save Changes" : "Create"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const CreateProjectFormModal: React.FC<{ onClose: () => void }> = ({ onClose }) => {
  const createMutation = useCreateProjectMutation();

  const handleSubmit = async (values: ProjectFormValues) => {
    if (!values.name.trim()) {
      toast.error("Name is required");
      return;
    }
    try {
      await createMutation.mutateAsync(values);
      toast.success("Project created successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create project");
    }
  };

  return (
    <ProjectFormInner
      title="NEW PROJECT"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
    />
  );
};

const EditProjectFormModal: React.FC<{ projectId: string; onClose: () => void }> = ({
  projectId,
  onClose,
}) => {
  const { data: project } = useProjectQuery(projectId);
  const updateMutation = useUpdateProjectMutation(projectId);

  const handleSubmit = async (values: ProjectFormValues) => {
    if (!values.name.trim()) {
      toast.error("Name is required");
      return;
    }
    try {
      await updateMutation.mutateAsync(values);
      toast.success("Project updated successfully");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update project");
    }
  };

  const initialValues: ProjectFormValues = {
    name: project?.name || "",
    description: project?.description || "",
    color: project?.color || PRESET_COLORS[5],
    icon: project?.icon || PRESET_ICONS[0],
  };

  return (
    <ProjectFormInner
      title="EDIT PROJECT"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
    />
  );
};

export const ProjectForm: React.FC<ProjectFormProps> = ({ isOpen, onClose, projectId }) => {
  if (!isOpen) return null;

  if (projectId) {
    return (
      <Suspense
        fallback={
          <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal}>
              <div className="flex flex-col items-center justify-center p-12 gap-3">
                <Loader2 className="animate-spin text-primary" size={28} />
                <span className="text-xs text-muted-foreground font-mono">LOADING PROJECT...</span>
              </div>
            </div>
          </div>
        }
      >
        <EditProjectFormModal projectId={projectId} onClose={onClose} />
      </Suspense>
    );
  }

  return <CreateProjectFormModal onClose={onClose} />;
};
