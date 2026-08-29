import React, { useMemo } from "react";
import * as Icons from "lucide-react";
import { MagicBento, type BentoCardItem } from "../../../components/MagicBento";
import { useSnippetsQuery } from "~features/snippets/hooks/useSnippets";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { useCredentialsQuery } from "~features/credentials/hooks/useCredentials";
import { useVaultStatusQuery } from "~features/security/hooks/useSecurity";
import { useNotesQuery } from "~features/notes/hooks/useNotes";
import { useLinksQuery } from "~features/links/hooks/useLinks";
import { useTagsQuery } from "~features/tags/hooks/useTags";
import { useBoardsQuery } from "~features/boards/hooks/useBoards";
import type { ProjectViewResponse } from "~types/api";
import type { ProjectTabType } from "../../../routes/projects.$projectId";

interface ProjectOverviewProps {
  projectId: string;
  project?: ProjectViewResponse;
  onNavigateTab: (tab: ProjectTabType) => void;
  onOpenTagsManager: () => void;
}

function hexToRgb(hex?: string): string {
  if (!hex || !hex.startsWith("#")) return "16, 185, 129";
  let cleanHex = hex.replace("#", "");
  if (cleanHex.length === 3) {
    cleanHex = cleanHex
      .split("")
      .map((c) => c + c)
      .join("");
  }
  const num = parseInt(cleanHex, 16);
  if (isNaN(num)) return "16, 185, 129";
  const r = (num >> 16) & 255;
  const g = (num >> 8) & 255;
  const b = num & 255;
  return `${r}, ${g}, ${b}`;
}

export const ProjectOverview: React.FC<ProjectOverviewProps> = ({
  projectId,
  project,
  onNavigateTab,
  onOpenTagsManager,
}) => {
  // Check vault unlock status
  const { data: vaultStatus } = useVaultStatusQuery();
  const isVaultActive = vaultStatus?.active === true;

  // Query project asset metrics
  const { data: boardsData } = useBoardsQuery(projectId);
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const { data: problemsData } = useProblemsQuery(projectId);
  const { data: credentialsData } = useCredentialsQuery(projectId, isVaultActive);
  const { data: notesData } = useNotesQuery(projectId);
  const { data: linksData } = useLinksQuery(projectId);
  const { data: tagsData } = useTagsQuery(projectId);

  const boardsCount = boardsData?.page?.totalElements ?? boardsData?.content?.length ?? 0;
  const snippetsCount =
    snippetsData?.page?.totalElements ?? snippetsData?.content?.length ?? 0;
  
  const problems = problemsData?.content || [];
  const openProblemsCount = problems.filter(
    (p) => p.status === "OPEN" || p.status === "WORKING_ON"
  ).length;

  const credentialsCount =
    credentialsData?.page?.totalElements ?? credentialsData?.content?.length ?? 0;
  const notesCount =
    notesData?.page?.totalElements ?? notesData?.content?.length ?? 0;
  const linksCount =
    linksData?.page?.totalElements ?? linksData?.content?.length ?? 0;
  const tagsCount = tagsData?.length ?? 0;

  const glowRgb = useMemo(() => hexToRgb(project?.color), [project?.color]);

  const cards: BentoCardItem[] = useMemo(
    () => [
      {
        id: "boards",
        title: "Kanban Boards",
        description:
          "Manage sprint tasks, backlog priorities, and cross-cutting workflow cards with interactive Drag & Drop.",
        label: "Sprint & Tasks",
        icon: Icons.SquareKanban,
        stat: boardsCount,
        statLabel: "Boards",
        color: project?.color || "#10b981",
        onClick: () => onNavigateTab("boards"),
      },
      {
        id: "snippets",
        title: "Code Snippets",
        description:
          "Store and manage reusable multi-language code snippets with Monaco Editor syntax highlighting.",
        label: "Snippets Vault",
        icon: Icons.Code2,
        stat: snippetsCount,
        statLabel: "Snippets",
        color: project?.color || "#10b981",
        onClick: () => onNavigateTab("snippets"),
      },
      {
        id: "problems",
        title: "Problems & Diagnostics",
        description:
          "Track and diagnose active project blockers, bugs, and unresolved technical issues.",
        label: "Issue Tracker",
        icon: Icons.AlertCircle,
        stat: openProblemsCount,
        statLabel: "Open Issues",
        color: openProblemsCount > 0 ? "#ef4444" : "#10b981",
        onClick: () => onNavigateTab("problems"),
      },
      {
        id: "credentials",
        title: "Credentials & Secrets",
        description:
          "Securely store API keys, database credentials, and secret tokens protected with Argon2id + AES-256-GCM.",
        label: "Encrypted Vault",
        icon: isVaultActive ? Icons.KeyRound : Icons.Lock,
        stat: isVaultActive ? credentialsCount : "?",
        statLabel: isVaultActive ? "Secrets" : "Locked",
        color: "#8b5cf6",
        onClick: () => onNavigateTab("credentials"),
      },
      {
        id: "notes",
        title: "Project Notes",
        description:
          "Document technical decisions, architectures, and guidelines with rendered Markdown.",
        label: "Knowledge Base",
        icon: Icons.FileText,
        stat: notesCount,
        statLabel: "Notes",
        color: "#3b82f6",
        onClick: () => onNavigateTab("notes"),
      },
      {
        id: "links",
        title: "Links & Repositories",
        description:
          "Quick access to production dashboards, Git repositories, staging environments, and documentation.",
        label: "Bookmarks",
        icon: Icons.Link2,
        stat: linksCount,
        statLabel: "Links",
        color: "#06b6d4",
        onClick: () => onNavigateTab("links"),
      },
      {
        id: "tags",
        title: "Tags & Taxonomy",
        description:
          "Organize snippets, problems, credentials, and notes with customizable cross-cutting tags.",
        label: "Taxonomy & Labels",
        icon: Icons.Tags,
        stat: tagsCount,
        statLabel: "Tags",
        color: "#ec4899",
        onClick: onOpenTagsManager,
      },
    ],
    [
      boardsCount,
      snippetsCount,
      openProblemsCount,
      credentialsCount,
      isVaultActive,
      notesCount,
      linksCount,
      tagsCount,
      project?.color,
      onNavigateTab,
      onOpenTagsManager,
    ]
  );

  return (
    <div className="flex flex-col gap-2.5 w-full max-w-5xl mx-auto py-1 px-1 animate-in fade-in duration-300">
      {/* Compact Overview Intro Banner */}
      <div className="flex items-center justify-between gap-3 p-3 rounded-xl border border-border bg-card/40 backdrop-blur-md">
        <div className="flex items-center gap-2.5 min-w-0">
          <div
            className="w-8 h-8 rounded-lg flex items-center justify-center border border-border shrink-0 shadow-sm"
            style={{ backgroundColor: "var(--color-secondary)" }}
          >
            <Icons.FolderGit2
              size={17}
              style={{ color: project?.color || "var(--color-primary)" }}
            />
          </div>
          <div className="min-w-0">
            <h2 className="text-sm font-bold text-foreground flex items-center gap-2 font-mono truncate">
              <span>{project?.name || "Project Overview"}</span>
              {project?.archived && (
                <span className="text-[9px] uppercase font-mono px-1.5 py-0.5 rounded-full bg-muted text-muted-foreground border border-border">
                  Archived
                </span>
              )}
            </h2>
            <p className="text-[11px] text-muted-foreground truncate max-w-xl">
              {project?.description || "Select any workspace card below to access your project vault assets."}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-3 border-l border-border pl-4 shrink-0">
          <div className="text-right">
            <span className="block text-[9px] font-mono uppercase text-muted-foreground">
              Total Assets
            </span>
            <span className="text-sm font-bold font-mono text-foreground">
              {snippetsCount + credentialsCount + notesCount + linksCount}
            </span>
          </div>
        </div>
      </div>

      {/* MagicBento Grid */}
      <div className="w-full">
        <MagicBento
          cards={cards}
          glowColor={glowRgb}
          enableStars={true}
          enableSpotlight={true}
          enableBorderGlow={true}
          enableTilt={true}
          enableMagnetism={true}
          clickEffect={true}
        />
      </div>
    </div>
  );
};

export default ProjectOverview;
