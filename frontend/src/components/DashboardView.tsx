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
import styles from "../routes/index.module.css";

// Allowed icons map
const ICON_MAPPING: Record<string, React.ComponentType<{ size?: number; className?: string; style?: React.CSSProperties }>> = {
  Folder: Icons.Folder,
  Terminal: Icons.Terminal,
  Database: Icons.Database,
  Globe: Icons.Globe,
  Cpu: Icons.Cpu,
  Activity: Icons.Activity,
  BookOpen: Icons.BookOpen,
  Code: Icons.Code,
};

// Helper to resolve project icon components securely via allowlist lookup
const getIconComponent = (iconName?: string) => {
  if (iconName && iconName in ICON_MAPPING) {
    return ICON_MAPPING[iconName];
  }
  return Icons.Folder;
};

export const DashboardView: React.FC = () => {
  const { data: projectsData } = useProjectsQuery();
  const archiveMutation = useArchiveProjectMutation();
  const unarchiveMutation = useUnarchiveProjectMutation();
  const deleteMutation = useDeleteProjectMutation();

  const [editingProjectId, setEditingProjectId] = useState<string | undefined>(undefined);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [showArchived, setShowArchived] = useState(false);

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

  const handleDelete = async (id: string, name: string) => {
    if (window.confirm(`Are you absolutely sure you want to delete project "${name}"? This action CANNOT be undone and will delete all associated snippets.`)) {
      try {
        await deleteMutation.mutateAsync(id);
        toast.success(`Project "${name}" deleted`);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to delete project");
      }
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
        <div className={styles.statCard}>
          <span className={styles.statLabel}>System Persistence</span>
          <span className="text-sm font-semibold text-emerald-500 font-mono flex items-center gap-1.5 mt-2">
            <Icons.Database size={16} />
            LOCAL DEV API
          </span>
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
    </div>
  );
};
