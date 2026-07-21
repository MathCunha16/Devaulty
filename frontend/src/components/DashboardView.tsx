import React, { useState } from "react";
import { Link } from "@tanstack/react-router";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import {
  useProjectsQuery,
  useArchiveProjectMutation,
  useUnarchiveProjectMutation,
  useDeleteProjectMutation,
} from "~features/projects/hooks/useProjects";
import { ProjectForm } from "~features/projects/components/ProjectForm";
import { ConfirmModal } from "./ConfirmModal";
import { getIconComponent } from "../utils/icons";
import styles from "../routes/index.module.css";

export const DashboardView: React.FC = () => {
  const { data: projectsData } = useProjectsQuery();
  const archiveMutation = useArchiveProjectMutation();
  const unarchiveMutation = useUnarchiveProjectMutation();
  const deleteMutation = useDeleteProjectMutation();

  const [editingProjectId, setEditingProjectId] = useState<string | undefined>(undefined);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [showArchived, setShowArchived] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    itemId: string;
    itemName: string;
    isLoading: boolean;
  }>({ isOpen: false, itemId: "", itemName: "", isLoading: false });

  const projects = projectsData?.content || [];
  const activeProjects = projects.filter((p) => !p.archived);
  const archivedProjects = projects.filter((p) => p.archived);

  const isMutationPending = archiveMutation.isPending || unarchiveMutation.isPending || deleteMutation.isPending;

  const handleArchive = async (id: string, name: string) => {
    try {
      await archiveMutation.mutateAsync(id);
      toast.success(`Project "${name}" archived`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to archive project");
    }
  };

  const handleUnarchive = async (id: string, name: string) => {
    try {
      await unarchiveMutation.mutateAsync(id);
      toast.success(`Project "${name}" restored`);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to restore project");
    }
  };

  const handleDelete = (id: string, name: string) => {
    setConfirmModal({ isOpen: true, itemId: id, itemName: name, isLoading: false });
  };

  const handleConfirmDelete = async () => {
    setConfirmModal((prev) => ({ ...prev, isLoading: true }));
    try {
      await deleteMutation.mutateAsync(confirmModal.itemId);
      toast.success(`Project "${confirmModal.itemName}" deleted`);
      setConfirmModal({ isOpen: false, itemId: "", itemName: "", isLoading: false });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete project");
      setConfirmModal((prev) => ({ ...prev, isLoading: false }));
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>DASHBOARD</h1>
        <p className={styles.subtitle}>Manage your projects, snippets, and credentials</p>
      </div>

      <div className={styles.statsGrid}>
        <div className={styles.statCard}>
          <span className={styles.statLabel}>Active Projects</span>
          <span className={styles.statValue}>{activeProjects.length}</span>
        </div>
        <div className={styles.statCard}>
          <span className={styles.statLabel}>Archived Projects</span>
          <span className={styles.statValue}>{archivedProjects.length}</span>
        </div>
      </div>

      <div>
        <h2 className={styles.sectionTitle}>Active Projects</h2>
        <div className={styles.projectsGrid}>
          {activeProjects.map((project) => {
            const ProjectIcon = getIconComponent(project.icon);
            return (
              <div key={project.id} className={styles.projectCard}>
                <div
                  className="h-1 w-full absolute top-0 left-0"
                  style={{ backgroundColor: project.color || "var(--color-primary)" }}
                />
                <div className={styles.cardHeader}>
                  <div className={styles.cardTitleSection}>
                    <div className={styles.cardIcon}>
                      <ProjectIcon
                        size={18}
                        style={{ color: project.color || "var(--color-primary)" }}
                      />
                    </div>
                    <h3 className={styles.cardTitle}>{project.name}</h3>
                  </div>
                </div>

                <p className={styles.cardDesc}>
                  {project.description || "No description provided."}
                </p>

                <div className={styles.cardFooter}>
                  <div className={styles.cardActions}>
                    <button
                      className={styles.actionBtn}
                      onClick={() => setEditingProjectId(project.id)}
                      title="Edit Project"
                      disabled={isMutationPending}
                    >
                      <Icons.Edit3 size={12} />
                    </button>
                    <button
                      className={styles.actionBtn}
                      onClick={() => handleArchive(project.id, project.name)}
                      title="Archive Project"
                      disabled={isMutationPending}
                    >
                      <Icons.Archive size={12} />
                    </button>
                    <button
                      className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                      onClick={() => handleDelete(project.id, project.name)}
                      title="Delete Project"
                      disabled={isMutationPending}
                    >
                      <Icons.Trash2 size={12} />
                    </button>
                  </div>

                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    className={styles.openLink}
                  >
                    <span>Open snippets</span>
                    <Icons.ExternalLink size={12} />
                  </Link>
                </div>
              </div>
            );
          })}

          <button className={styles.emptyCard} onClick={() => setIsCreateOpen(true)} disabled={isMutationPending}>
            <div className="w-10 h-10 rounded-full border border-dashed border-border flex items-center justify-center text-muted-foreground">
              <Icons.Plus size={18} />
            </div>
            <div>
              <div className={styles.emptyCardTitle}>Create Project</div>
              <p className={styles.emptyCardText}>Add a new environment node</p>
            </div>
          </button>
        </div>
      </div>

      {archivedProjects.length > 0 && (
        <div className={styles.archiveSection}>
          <button
            type="button"
            className={styles.archiveHeader}
            style={{ background: "none", border: "none", width: "100%", padding: 0 }}
            onClick={() => setShowArchived(!showArchived)}
            aria-expanded={showArchived}
          >
            <span className={styles.archiveToggleText}>
              {showArchived ? "Hide Archived Projects" : `Show Archived Projects (${archivedProjects.length})`}
            </span>
            <span className={styles.actionBtn}>
              {showArchived ? <Icons.EyeOff size={12} /> : <Icons.Eye size={12} />}
            </span>
          </button>

          {showArchived && (
            <div className={`${styles.projectsGrid} mt-4`}>
              {archivedProjects.map((project) => {
                const ProjectIcon = getIconComponent(project.icon);
                return (
                  <div key={project.id} className={styles.projectCard} style={{ opacity: 0.7 }}>
                    <div className={styles.cardHeader}>
                      <div className={styles.cardTitleSection}>
                        <div className={styles.cardIcon}>
                          <ProjectIcon size={18} className="text-muted-foreground" />
                        </div>
                        <h3 className={styles.cardTitle}>{project.name}</h3>
                      </div>
                    </div>

                    <p className={styles.cardDesc}>
                      {project.description || "No description provided."}
                    </p>

                    <div className={styles.cardFooter}>
                      <div className={styles.cardActions}>
                        <button
                          className={styles.actionBtn}
                          onClick={() => handleUnarchive(project.id, project.name)}
                          title="Restore Project"
                          disabled={isMutationPending}
                        >
                          <Icons.ArchiveRestore size={12} />
                        </button>
                        <button
                          className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                          onClick={() => handleDelete(project.id, project.name)}
                          title="Delete Project"
                          disabled={isMutationPending}
                        >
                          <Icons.Trash2 size={12} />
                        </button>
                      </div>
                      <span className="text-xs text-muted-foreground font-mono">ARCHIVED</span>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}

      {/* Edit Form Modal */}
      {editingProjectId && (
        <ProjectForm
          isOpen={true}
          projectId={editingProjectId}
          onClose={() => setEditingProjectId(undefined)}
        />
      )}

      {/* Create Form Modal */}
      {isCreateOpen && (
        <ProjectForm isOpen={true} onClose={() => setIsCreateOpen(false)} />
      )}

      {/* Confirm Delete Modal */}
      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={() => setConfirmModal((prev) => ({ ...prev, isOpen: false }))}
        onConfirm={handleConfirmDelete}
        title="Delete Project"
        message="Are you sure you want to permanently delete the project"
        itemName={confirmModal.itemName}
        warningText="This cannot be undone. All snippets and diagnostics in this project will be permanently deleted."
        isLoading={confirmModal.isLoading}
      />
    </div>
  );
};
