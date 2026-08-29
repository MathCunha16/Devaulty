import React, { useState, useMemo, useCallback } from "react";
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
import { ProjectCard } from "./ProjectCard";
import { useTheme } from "../hooks/useTheme";
import FloatingLines from "./FloatingLines";
import styles from "../routes/index.module.css";

// ── Theme-aware FloatingLines configuration ──────────────────
// Stable references — MUST be outside the component to avoid re-creating the WebGL
// context on every React render (e.g. sidebar toggle, search typing, etc.)
const DARK_GRADIENT = ["#10b981", "#34d399", "#059669", "#064e3b"];
const LIGHT_GRADIENT = ["#047857", "#059669", "#10b981", "#065f46"];
const ENABLED_WAVES: Array<'middle' | 'bottom'> = ['middle', 'bottom'];
const LINE_COUNT = [8, 12];
const LINE_DISTANCE = [6, 4];

type TabFilter = "ACTIVE" | "ARCHIVED" | "ALL";

export const DashboardView: React.FC = () => {
  const { theme } = useTheme();
  const { data: projectsData, isLoading } = useProjectsQuery();

  // Memoize FloatingLines config to avoid re-creating the WebGL context on every render
  const floatingLinesConfig = useMemo(
    () => ({
      linesGradient: theme === "dark" ? DARK_GRADIENT : LIGHT_GRADIENT,
      mixBlendMode: (theme === "dark" ? "screen" : "normal") as React.CSSProperties["mixBlendMode"],
    }),
    [theme],
  );

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


  const handleArchive = useCallback(
    async (id: string, name: string) => {
      try {
        await archiveMutation.mutateAsync(id);
        toast.success(`Project "${name}" archived`);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to archive project");
      }
    },
    [archiveMutation]
  );

  const handleUnarchive = useCallback(
    async (id: string, name: string) => {
      try {
        await unarchiveMutation.mutateAsync(id);
        toast.success(`Project "${name}" restored`);
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to restore project");
      }
    },
    [unarchiveMutation]
  );

  const handleDelete = useCallback((id: string, name: string) => {
    setConfirmModal({ isOpen: true, itemId: id, itemName: name, isLoading: false });
  }, []);

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
    <div className={styles.dashboardRoot}>
      {/* WebGL animated background */}
      <div className={styles.bgLayer}>
        <FloatingLines
          theme={theme}
          linesGradient={floatingLinesConfig.linesGradient}
          enabledWaves={ENABLED_WAVES}
          lineCount={LINE_COUNT}
          lineDistance={LINE_DISTANCE}
          animationSpeed={0.6}
          interactive={false}
          parallax={false}
          mixBlendMode={floatingLinesConfig.mixBlendMode}
        />
      </div>

      {/* Readability overlay */}
      <div className={styles.bgOverlay} />

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
          {filteredProjects.map((project) => (
            <ProjectCard
              key={project.id}
              project={project}
              isArchived={project.archived}
              isMutationPending={isMutationPending}
              onArchive={handleArchive}
              onUnarchive={handleUnarchive}
              onDelete={handleDelete}
              onEdit={setEditingProjectId}
            />
          ))}


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
    </div>
  );
};
