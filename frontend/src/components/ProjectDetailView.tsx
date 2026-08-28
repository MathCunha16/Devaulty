import React, { useState, useEffect } from "react";
import { Link, useMatch, getRouteApi } from "@tanstack/react-router";
import * as Icons from "lucide-react";
import { useSidebar } from "../hooks/useSidebar";
import { useProjectQuery } from "~features/projects/hooks/useProjects";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { TagsManagerModal } from "./TagsManagerModal";
import { GooeyNav, type GooeyNavItem } from "./GooeyNav";
import { ProjectOverview } from "~features/projects/components/ProjectOverview";
import { SnippetsWorkspace } from "~features/snippets/components/SnippetsWorkspace";
import { ProblemsWorkspace } from "~features/problems/components/ProblemsWorkspace";
import { CredentialsWorkspace } from "~features/credentials/components/CredentialsWorkspace";
import { NotesWorkspace } from "~features/notes/components/NotesWorkspace";
import { LinksWorkspace } from "~features/links/components/LinksWorkspace";
import { KanbanWorkspace } from "~features/boards/components/KanbanWorkspace";
import styles from "../routes/projects.$projectId.module.css";
import type { ProjectTabType } from "../routes/projects.$projectId";

const routeApi = getRouteApi("/projects/$projectId");

export const ProjectDetailView: React.FC = () => {
  const { projectId } = useParamsHelper();
  const { isOpen, close: closeSidebar } = useSidebar();
  const search = routeApi.useSearch();
  const navigate = routeApi.useNavigate();

  // Auto-close sidebar when entering a project to maximize workspace space
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

  // Workspace sub-navigation state derived directly from URL search params
  const activeTab: ProjectTabType = search.tab || "overview";

  const [isTagsManagerOpen, setIsTagsManagerOpen] = useState(false);

  const handleTabChange = (tab: ProjectTabType) => {
    navigate({ search: { tab }, replace: true });
  };

  const navItems: GooeyNavItem[] = [
    { id: "overview", label: "Overview", icon: Icons.LayoutGrid },
    { id: "boards", label: "Kanban", icon: Icons.SquareKanban },
    { id: "snippets", label: "Snippets", icon: Icons.Code2 },
    {
      id: "problems",
      label: "Problems",
      icon: Icons.AlertCircle,
      badgeCount: openProblemsCount,
    },
    { id: "credentials", label: "Credentials", icon: Icons.KeyRound },
    { id: "notes", label: "Notes", icon: Icons.FileText },
    { id: "links", label: "Links", icon: Icons.Link2 },
  ];

  return (
    <div
      className={`flex flex-col gap-2.5 h-full p-4 overflow-hidden transition-all duration-300 ${
        isOpen ? "pt-3" : "pt-20"
      }`}
    >
      {/* Centered Unified Navigation Dock */}
      <div className="flex items-center justify-center shrink-0 w-full mb-1">
        <div className="flex items-center gap-1.5 p-1 rounded-full bg-card/60 border border-border/80 backdrop-blur-md shadow-md">
          {/* Back to Dashboard */}
          <Link
            to="/"
            className="flex items-center gap-1.5 py-1 px-3 rounded-full text-xs text-muted-foreground hover:text-foreground hover:bg-secondary/80 transition-all font-medium"
            title="Back to Dashboard"
          >
            <Icons.ArrowLeft size={13} />
            <span>Dashboard</span>
          </Link>

          <div className="w-[1px] h-4 bg-border/80 mx-0.5" />

          {/* Horizontal GooeyNav */}
          <GooeyNav
            items={navItems}
            activeId={activeTab}
            onChange={(id) => handleTabChange(id as ProjectTabType)}
            projectColor={project?.color}
          />

          <div className="w-[1px] h-4 bg-border/80 mx-0.5" />

          {/* Manage Tags */}
          <button
            type="button"
            onClick={() => setIsTagsManagerOpen(true)}
            className="flex items-center gap-1.5 py-1 px-3 rounded-full text-xs text-muted-foreground hover:text-foreground hover:bg-secondary/80 transition-all font-medium cursor-pointer border-0 bg-transparent"
            title="Manage Project Tags"
          >
            <Icons.Tags size={13} />
            <span>Tags</span>
          </button>
        </div>
      </div>

      {/* Main Workspace Layout */}
      <div className={styles.pageLayout}>
        {activeTab === "overview" && (
          <div className="flex-1 overflow-y-auto pr-1">
            <ProjectOverview
              projectId={projectId}
              project={project}
              onNavigateTab={handleTabChange}
              onOpenTagsManager={() => setIsTagsManagerOpen(true)}
            />
          </div>
        )}

        {activeTab === "boards" && (
          <KanbanWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            projectColor={project?.color}
          />
        )}

        {activeTab === "snippets" && (
          <SnippetsWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            initialSelectedId={search.itemId}
          />
        )}

        {activeTab === "problems" && (
          <ProblemsWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            initialSelectedId={search.itemId}
          />
        )}

        {activeTab === "credentials" && (
          <CredentialsWorkspace
            projectId={projectId}
            isActive={activeTab === "credentials"}
            onNavigateTab={handleTabChange}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            initialSelectedId={search.itemId}
          />
        )}

        {activeTab === "notes" && (
          <NotesWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            initialSelectedId={search.itemId}
          />
        )}

        {activeTab === "links" && (
          <LinksWorkspace
            projectId={projectId}
            onOpenManageTagsModal={() => setIsTagsManagerOpen(true)}
            initialSelectedId={search.itemId}
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
