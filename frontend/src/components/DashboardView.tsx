import React, { useState } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
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

type TabFilter = "ACTIVE" | "ARCHIVED" | "ALL";

export const DashboardView: React.FC = () => {
  const navigate = useNavigate();
  const { data: projectsData, isLoading } = useProjectsQuery();

  const archiveMutation = useArchiveProjectMutation();
  const unarchiveMutation = useUnarchiveProjectMutation();
  const deleteMutation = useDeleteProjectMutation();

  const [searchQuery, setSearchQuery] = useState("");
  const [tabFilter, setTabFilter] = useState<TabFilter>("ACTIVE");
  const [editingProjectId, setEditingProjectId] = useState<string | undefined>(undefined);
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    itemId: string;
    itemName: string;
    isLoading: boolean;
  }>({ isOpen: false, itemId: "", itemName: "", isLoading: false });

  const projects = projectsData?.content || [];
  const activeProjectsCount = projects.filter((p) => !p.archived).length;
  const archivedProjectsCount = projects.filter((p) => p.archived).length;

  const effectiveTab: TabFilter =
    tabFilter === "ARCHIVED" && archivedProjectsCount === 0 ? "ACTIVE" : tabFilter;

  const isMutationPending =
    archiveMutation.isPending || unarchiveMutation.isPending || deleteMutation.isPending;

  const filteredProjects = projects.filter((project) => {
    // Filter by effective tab
    if (effectiveTab === "ACTIVE" && project.archived) return false;
    if (effectiveTab === "ARCHIVED" && !project.archived) return false;

    // Filter by search query
    if (!searchQuery.trim()) return true;
    const q = searchQuery.toLowerCase();
    return (
      project.name.toLowerCase().includes(q) ||
      (project.description && project.description.toLowerCase().includes(q))
    );
  });


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
      {/* Header section */}
      <div className={styles.header}>
        <div className={styles.headerTitleGroup}>
          <div className={styles.headerBadge}>
            <Icons.FolderGit2 size={14} />
            <span>WORKSPACES</span>
          </div>
          <h1 className={styles.title}>PROJECTS</h1>
          <p className={styles.subtitle}>
            Select a project workspace to manage snippets, vault credentials, notes, and links.
          </p>
        </div>

        <button
          type="button"
          className={styles.createBtn}
          onClick={() => setIsCreateOpen(true)}
          disabled={isMutationPending}
        >
          <Icons.Plus size={16} />
          <span>New Project</span>
        </button>
      </div>

      {/* Control Bar: Search & Tab Filters */}
      <div className={styles.controlBar}>
        <div className={styles.searchBox}>
          <Icons.Search size={14} className={styles.searchIcon} />
          <input
            type="text"
            className={styles.searchInput}
            placeholder="Filter projects by name or description..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button
              type="button"
              className={styles.clearSearchBtn}
              onClick={() => setSearchQuery("")}
            >
              <Icons.X size={12} />
            </button>
          )}
        </div>

        <div className={styles.filterTabs}>
          <button
            type="button"
            className={`${styles.filterTab} ${effectiveTab === "ACTIVE" ? styles.filterTabActive : ""}`}
            onClick={() => setTabFilter("ACTIVE")}
          >
            Active ({activeProjectsCount})
          </button>
          {archivedProjectsCount > 0 && (
            <button
              type="button"
              className={`${styles.filterTab} ${effectiveTab === "ARCHIVED" ? styles.filterTabActive : ""}`}
              onClick={() => setTabFilter("ARCHIVED")}
            >
              Archived ({archivedProjectsCount})
            </button>
          )}
          <button
            type="button"
            className={`${styles.filterTab} ${effectiveTab === "ALL" ? styles.filterTabActive : ""}`}
            onClick={() => setTabFilter("ALL")}
          >
            All ({projects.length})
          </button>
        </div>
      </div>

      {/* Projects Grid */}
      {isLoading ? (
        <div className="flex flex-col items-center justify-center py-20 gap-3">
          <Icons.Loader2 size={32} className="animate-spin text-primary" />
          <span className="text-xs font-mono text-muted-foreground">LOADING WORKSPACES...</span>
        </div>
      ) : filteredProjects.length === 0 ? (
        <div className={styles.emptyState}>
          <Icons.FolderSearch size={44} className="text-muted-foreground animate-pulse" />
          <h3 className={styles.emptyTitle}>
            {searchQuery
              ? `No projects matching "${searchQuery}"`
              : effectiveTab === "ARCHIVED"
              ? "No archived projects found"
              : "No projects created yet"}
          </h3>
          <p className={styles.emptySubtitle}>
            {searchQuery
              ? "Try adjusting your search filter."
              : "Get started by creating your first dev project workspace."}
          </p>
          {searchQuery ? (
            <button
              type="button"
              className={styles.secondaryBtn}
              onClick={() => setSearchQuery("")}
            >
              Clear Search
            </button>
          ) : (
            <button
              type="button"
              className={styles.primaryBtn}
              onClick={() => setIsCreateOpen(true)}
            >
              <Icons.Plus size={14} />
              <span>Create Project</span>
            </button>
          )}
        </div>
      ) : (
        <div className={styles.projectsGrid}>
          {filteredProjects.map((project) => {
            const ProjectIcon = getIconComponent(project.icon);
            const cardColor = project.color || "var(--color-primary, #6366f1)";

            return (
              <div
                key={project.id}
                role="button"
                tabIndex={0}
                className={`${styles.projectCard} ${project.archived ? styles.projectCardArchived : ""}`}
                onClick={() =>
                  navigate({
                    to: "/projects/$projectId",
                    params: { projectId: project.id },
                  })
                }
                onKeyDown={(e) => {
                  if (e.target !== e.currentTarget) return;
                  if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    navigate({
                      to: "/projects/$projectId",
                      params: { projectId: project.id },
                    });
                  }
                }}
              >
                {/* Top Accent Strip */}
                <div
                  className={styles.cardAccentBar}
                  style={{ backgroundColor: cardColor }}
                />

                {/* Card Header */}
                <div className={styles.cardHeader}>
                  <div
                    className={styles.cardIconBox}
                    style={{
                      borderColor: `color-mix(in srgb, ${cardColor} 25%, transparent)`,
                      backgroundColor: `color-mix(in srgb, ${cardColor} 10%, transparent)`,
                    }}
                  >
                    <ProjectIcon size={20} style={{ color: cardColor }} />
                  </div>

                  <div className={styles.cardTitleBox}>
                    <h3 className={styles.cardTitle}>{project.name}</h3>
                    {project.archived && (
                      <span className={styles.archivedBadge}>ARCHIVED</span>
                    )}
                  </div>
                </div>

                {/* Card Body */}
                <p className={styles.cardDesc}>
                  {project.description || "No description provided."}
                </p>

                {/* Direct Shortcut Feature Badges */}
                <div className={styles.shortcutRow} onClick={(e) => e.stopPropagation()}>
                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    search={{ tab: "snippets" }}
                    className={styles.shortcutTag}
                    title="Code Snippets"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Icons.Code2 size={11} />
                    <span>Snippets</span>
                  </Link>
                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    search={{ tab: "credentials" }}
                    className={styles.shortcutTag}
                    title="Vault Credentials"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Icons.KeyRound size={11} />
                    <span>Vault</span>
                  </Link>
                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    search={{ tab: "problems" }}
                    className={styles.shortcutTag}
                    title="Problems & Solutions"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Icons.AlertCircle size={11} />
                    <span>Problems</span>
                  </Link>
                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    search={{ tab: "notes" }}
                    className={styles.shortcutTag}
                    title="Notes"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Icons.FileText size={11} />
                    <span>Notes</span>
                  </Link>
                  <Link
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    search={{ tab: "links" }}
                    className={styles.shortcutTag}
                    title="Links"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Icons.Link2 size={11} />
                    <span>Links</span>
                  </Link>
                </div>

                {/* Card Footer */}
                <div className={styles.cardFooter}>
                  <div className={styles.cardActions} onClick={(e) => e.stopPropagation()}>
                    {project.archived ? (
                      <button
                        type="button"
                        className={styles.actionBtn}
                        onClick={(e) => {
                          e.stopPropagation();
                          handleUnarchive(project.id, project.name);
                        }}
                        title="Restore Project"
                        disabled={isMutationPending}
                      >
                        <Icons.ArchiveRestore size={13} />
                      </button>
                    ) : (
                      <button
                        type="button"
                        className={styles.actionBtn}
                        onClick={(e) => {
                          e.stopPropagation();
                          handleArchive(project.id, project.name);
                        }}
                        title="Archive Project"
                        disabled={isMutationPending}
                      >
                        <Icons.Archive size={13} />
                      </button>
                    )}
                    <button
                      type="button"
                      className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDelete(project.id, project.name);
                      }}
                      title="Delete Project"
                      disabled={isMutationPending}
                    >
                      <Icons.Trash2 size={13} />
                    </button>
                  </div>

                  <button
                    type="button"
                    className={styles.enterLink}
                    onClick={(e) => {
                      e.stopPropagation();
                      setEditingProjectId(project.id);
                    }}
                    title="View & Edit Project Details"
                  >
                    <span>View Details</span>
                    <Icons.Settings size={13} />
                  </button>
                </div>
              </div>
            );
          })}


          {/* Create Project Card in Grid */}
          <button
            type="button"
            className={styles.createCard}
            onClick={() => setIsCreateOpen(true)}
            disabled={isMutationPending}
          >
            <div className={styles.createCardBadge}>
              <Icons.Plus size={20} />
            </div>
            <div className={styles.createCardContent}>
              <h4 className={styles.createCardTitle}>Create Project</h4>
              <p className={styles.createCardSub}>Add a new environment workspace</p>
            </div>
          </button>
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
        warningText="This cannot be undone. All snippets, credentials, problems, notes, and links in this project will be permanently deleted."
        isLoading={confirmModal.isLoading}
      />
    </div>
  );
};
