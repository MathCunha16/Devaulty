import React, { useState, useEffect } from "react";
import { Link, useMatch } from "@tanstack/react-router";
import * as Icons from "lucide-react";
import { useSidebar } from "../hooks/useSidebar";
import { useProjectQuery } from "~features/projects/hooks/useProjects";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { TagsManagerModal } from "./TagsManagerModal";
import { SnippetsWorkspace } from "~features/snippets/components/SnippetsWorkspace";
import { ProblemsWorkspace } from "~features/problems/components/ProblemsWorkspace";
import { CredentialsWorkspace } from "~features/credentials/components/CredentialsWorkspace";
import { NotesWorkspace } from "~features/notes/components/NotesWorkspace";
import { LinksWorkspace } from "~features/links/components/LinksWorkspace";
import styles from "../routes/projects.$projectId.module.css";
import { getIconComponent } from "../utils/icons";

export const ProjectDetailView: React.FC = () => {
  const { projectId } = useParamsHelper();
  const { close: closeSidebar } = useSidebar();

  // Auto-close sidebar when entering a project
  useEffect(() => {
    closeSidebar();
  }, [projectId, closeSidebar]);

  // Load project core details
  const { data: project } = useProjectQuery(projectId);

  // Load problems for badge counter
  const { data: problemsData } = useProblemsQuery(projectId);
  const problems = problemsData?.content || [];
  const openProblemsCount = problems.filter(
    (p) => p.status === "OPEN" || p.status === "WORKING_ON"
  ).length;



  // Workspace sub-navigation state
  const [activeTab, setActiveTab] = useState<
    "snippets" | "problems" | "credentials" | "notes" | "links"
  >("snippets");

  const [isTagsManagerOpen, setIsTagsManagerOpen] = useState(false);

  const handleTabChange = (
    tab: "snippets" | "problems" | "credentials" | "notes" | "links"
  ) => {
    setActiveTab(tab);
  };

  const projectIcon = getIconComponent(project?.icon);

  return (
    <div className="flex flex-col gap-4 h-full">
      {/* Project Header */}
      <div className={styles.projectHeader}>
        <div className={styles.projectTitleSection}>
          <div className={styles.projectIcon}>
            {React.createElement(projectIcon, {
              size: 22,
              style: { color: project?.color || "var(--color-primary)" },
            })}
          </div>
          <div>
            <h1 className={styles.projectName}>{project?.name}</h1>
            <p className={styles.projectDesc}>{project?.description || "No description provided."}</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setIsTagsManagerOpen(true)}
            className="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
          >
            <Icons.Tags size={12} />
            <span>Manage Tags</span>
          </button>

          <Link
            to="/"
            className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors"
          >
            <Icons.ArrowLeft size={12} />
            <span>Back to Dashboard</span>
          </Link>
        </div>
      </div>

      {/* Main Workspace Layout */}
      <div className={styles.pageLayout}>
        {/* Sidebar Workspace Tabs Selector */}
        <div className={styles.workspaceSidebar}>
          <button
            className={`${styles.workspaceTab} ${activeTab === "snippets" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("snippets")}
            title="Code Snippets"
          >
            <Icons.Code size={18} />
            <span className={styles.workspaceTabLabel}>Snippets</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "problems" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("problems")}
            title="Problems & Diagnostics"
          >
            <Icons.AlertCircle size={18} />
            <span className={styles.workspaceTabLabel}>Problems</span>
            {openProblemsCount > 0 && (
              <span className={styles.badgeCount}>{openProblemsCount}</span>
            )}
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "credentials" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("credentials")}
            title="Secure Credentials Vault"
          >
            <Icons.KeyRound size={18} />
            <span className={styles.workspaceTabLabel}>Credentials</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "notes" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("notes")}
            title="System Notes"
          >
            <Icons.FileText size={18} />
            <span className={styles.workspaceTabLabel}>Notes</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "links" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("links")}
            title="Web Links"
          >
            <Icons.Link2 size={18} />
            <span className={styles.workspaceTabLabel}>Links</span>
          </button>
        </div>

        {/* Feature Workspace Component Execution */}
        {activeTab === "snippets" && (
          <SnippetsWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
          />
        )}

        {activeTab === "problems" && (
          <ProblemsWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
          />
        )}

        {activeTab === "credentials" && (
          <CredentialsWorkspace
            projectId={projectId}
            isActive={activeTab === "credentials"}
            onNavigateTab={setActiveTab}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
          />
        )}

        {activeTab === "notes" && (
          <NotesWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
          />
        )}

        {activeTab === "links" && (
          <LinksWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
          />
        )}
      </div>

      {/* Global Tags Manager Modal */}
      <TagsManagerModal
        isOpen={isTagsManagerOpen}
        onClose={() => setIsTagsManagerOpen(false)}
        projectId={projectId}
      />
    </div>
  );
};

const useParamsHelper = () => {
  const match = useMatch({ from: "/projects/$projectId" });
  return match.params;
};

export const ProjectDetailRouteComponent: React.FC = () => {
  return <ProjectDetailView />;
};
