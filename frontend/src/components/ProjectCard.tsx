import React, { useCallback, useMemo } from "react";
import { Link, useNavigate } from "@tanstack/react-router";
import * as Icons from "lucide-react";
import { getIconComponent } from "../utils/icons";
import type { ProjectViewResponse } from "~types/api";
import { BorderGlow } from "./BorderGlow";

interface ProjectCardProps {
  project: ProjectViewResponse;
  isArchived: boolean;
  isMutationPending: boolean;
  onArchive: (id: string, name: string) => void;
  onUnarchive: (id: string, name: string) => void;
  onDelete: (id: string, name: string) => void;
  onEdit: (id: string) => void;
}

function hexToHslString(hex: string): string {
  let c = hex.replace("#", "");
  if (c.length === 3) c = c.split("").map((x) => x + x).join("");
  const num = parseInt(c, 16);
  if (isNaN(num)) return "160 84 39";
  const r = ((num >> 16) & 255) / 255;
  const g = ((num >> 8) & 255) / 255;
  const b = (num & 255) / 255;
  const max = Math.max(r, g, b);
  const min = Math.min(r, g, b);
  let h = 0;
  let s = 0;
  const l = (max + min) / 2;

  if (max !== min) {
    const d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch (max) {
      case r:
        h = (g - b) / d + (g < b ? 6 : 0);
        break;
      case g:
        h = (b - r) / d + 2;
        break;
      case b:
        h = (r - g) / d + 4;
        break;
    }
    h /= 6;
  }
  return `${Math.round(h * 360)} ${Math.round(s * 100)} ${Math.round(l * 100)}`;
}

export const ProjectCard: React.FC<ProjectCardProps> = React.memo(
  ({
    project,
    isArchived,
    isMutationPending,
    onArchive,
    onUnarchive,
    onDelete,
    onEdit,
  }) => {
    const navigate = useNavigate();
    const ProjectIcon = getIconComponent(project.icon);
    const cardColor = project.color || "#10b981";

    const glowHsl = useMemo(() => hexToHslString(cardColor), [cardColor]);
    const glowColors = useMemo(
      () => [cardColor, `color-mix(in srgb, ${cardColor} 80%, #38bdf8)`, `color-mix(in srgb, ${cardColor} 70%, #c084fc)`],
      [cardColor]
    );

    const navigateToProject = useCallback(() => {
      navigate({
        to: "/projects/$projectId",
        params: { projectId: project.id },
      });
    }, [navigate, project.id]);

    const handleKeyDown = useCallback(
      (e: React.KeyboardEvent<HTMLDivElement>) => {
        if (e.target !== e.currentTarget) return;
        if (e.key === "Enter" || e.key === " ") {
          e.preventDefault();
          navigateToProject();
        }
      },
      [navigateToProject]
    );

    const shortcuts = [
      { tab: "boards", icon: Icons.SquareKanban, label: "Kanban", color: "#10b981" },
      { tab: "snippets", icon: Icons.Code2, label: "Snippets", color: "#f59e0b" },
      { tab: "credentials", icon: Icons.KeyRound, label: "Vault", color: "#a855f7" },
      { tab: "problems", icon: Icons.AlertCircle, label: "Issues", color: "#ef4444" },
      { tab: "notes", icon: Icons.FileText, label: "Notes", color: "#3b82f6" },
      { tab: "links", icon: Icons.Link2, label: "Links", color: "#06b6d4" },
    ] as const;

    return (
      <BorderGlow
        borderRadius={14}
        glowRadius={35}
        glowIntensity={0.9}
        edgeSensitivity={35}
        animated={true}
        glowColor={glowHsl}
        colors={glowColors}
        backgroundColor="color-mix(in srgb, var(--color-card) 75%, transparent)"
        className={`group cursor-pointer transition-transform duration-200 hover:-translate-y-1 ${
          isArchived ? "opacity-70 grayscale-[0.3] hover:opacity-95 hover:grayscale-0" : ""
        }`}
        onClick={navigateToProject}
        onKeyDown={handleKeyDown}
      >
        {/* Top Header: Icon + Title/Desc + Actions */}
        <div className="flex items-start gap-3 w-full">
          <div
            className="w-10 h-10 rounded-lg border flex items-center justify-center shrink-0 shadow-sm transition-transform duration-200 group-hover:scale-105"
            style={{
              borderColor: `color-mix(in srgb, ${cardColor} 30%, transparent)`,
              backgroundColor: `color-mix(in srgb, ${cardColor} 10%, var(--color-secondary))`,
              color: cardColor,
            }}
          >
            <ProjectIcon size={20} />
          </div>

          <div className="flex flex-col gap-0.5 min-w-0 flex-1">
            <div className="flex items-center gap-2">
              <h3 className="text-sm font-bold font-mono text-foreground truncate tracking-tight">
                {project.name}
              </h3>
              {isArchived ? (
                <span className="inline-flex items-center gap-1 text-[9px] font-mono font-bold text-muted-foreground bg-muted/60 border border-border px-1.5 py-0.5 rounded-full uppercase">
                  <Icons.Archive size={9} />
                  Archived
                </span>
              ) : (
                <span
                  className="w-1.5 h-1.5 rounded-full shrink-0 animate-pulse"
                  style={{ backgroundColor: cardColor, boxShadow: `0 0 6px ${cardColor}` }}
                  title="Active Workspace"
                />
              )}
            </div>
            <p className="text-xs text-muted-foreground line-clamp-2 leading-relaxed">
              {project.description || "No description provided."}
            </p>
          </div>

          {/* Header Action Buttons */}
          <div className="flex items-center gap-0.5 shrink-0 opacity-70 hover:opacity-100 transition-opacity" onClick={(e) => e.stopPropagation()}>
            <button
              type="button"
              className="w-7 h-7 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-secondary border border-transparent hover:border-border transition-all"
              onClick={(e) => {
                e.stopPropagation();
                onEdit(project.id);
              }}
              title="Project Settings"
            >
              <Icons.Settings size={13} />
            </button>

            {isArchived ? (
              <button
                type="button"
                className="w-7 h-7 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-secondary border border-transparent hover:border-border transition-all"
                onClick={(e) => {
                  e.stopPropagation();
                  onUnarchive(project.id, project.name);
                }}
                title="Restore Project"
                disabled={isMutationPending}
              >
                <Icons.ArchiveRestore size={13} />
              </button>
            ) : (
              <button
                type="button"
                className="w-7 h-7 rounded-md flex items-center justify-center text-muted-foreground hover:text-foreground hover:bg-secondary border border-transparent hover:border-border transition-all"
                onClick={(e) => {
                  e.stopPropagation();
                  onArchive(project.id, project.name);
                }}
                title="Archive Project"
                disabled={isMutationPending}
              >
                <Icons.Archive size={13} />
              </button>
            )}

            <button
              type="button"
              className="w-7 h-7 rounded-md flex items-center justify-center text-muted-foreground hover:text-destructive hover:bg-destructive/15 border border-transparent hover:border-destructive/30 transition-all"
              onClick={(e) => {
                e.stopPropagation();
                onDelete(project.id, project.name);
              }}
              title="Delete Project"
              disabled={isMutationPending}
            >
              <Icons.Trash2 size={13} />
            </button>
          </div>
        </div>

        {/* Integrated Module Quick-Access Toolbar */}
        <div
          className="grid grid-cols-6 w-full rounded-lg bg-background/40 border border-border/50 overflow-hidden divide-x divide-border/30"
          onClick={(e) => e.stopPropagation()}
        >
          {shortcuts.map(({ tab, icon: Icon, label, color }) => (
            <Link
              key={tab}
              to="/projects/$projectId"
              params={{ projectId: project.id }}
              search={{ tab }}
              className="flex flex-col items-center justify-center py-2 px-1 text-muted-foreground hover:text-foreground transition-all duration-150 group/item hover:bg-secondary/60"
              style={{ "--module-hover-color": color } as React.CSSProperties}
              onClick={(e) => e.stopPropagation()}
            >
              <Icon size={14} className="transition-transform duration-150 group-hover/item:-translate-y-0.5 group-hover/item:text-[var(--module-hover-color)]" />
              <span className="text-[10px] font-mono font-medium tracking-tight mt-1 truncate w-full text-center group-hover/item:text-[var(--module-hover-color)]">
                {label}
              </span>
            </Link>
          ))}
        </div>

        {/* Card Footer: Metadata & Open Link */}
        <div className="flex items-center justify-between pt-2.5 border-t border-border/40 text-xs font-mono w-full">
          <span className="text-[9px] uppercase tracking-wider text-muted-foreground font-bold opacity-60">
            Workspace
          </span>

          <Link
            to="/projects/$projectId"
            params={{ projectId: project.id }}
            className="inline-flex items-center gap-1.5 text-xs font-bold transition-all hover:brightness-110 focus:outline-none focus-visible:underline"
            style={{ color: cardColor }}
            onClick={(e) => e.stopPropagation()}
          >
            <span>Open</span>
            <Icons.ArrowRight size={13} className="transition-transform duration-150 group-hover:translate-x-1" />
          </Link>
        </div>
      </BorderGlow>
    );
  }
);

ProjectCard.displayName = "ProjectCard";
