import React, { useState, useRef, useCallback, useEffect, useMemo } from "react";
import * as Icons from "lucide-react";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { useNavigate } from "@tanstack/react-router";
import { useSnippetsQuery } from "~features/snippets/hooks/useSnippets";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { useCredentialsQuery } from "~features/credentials/hooks/useCredentials";
import { useNotesQuery } from "~features/notes/hooks/useNotes";
import { useLinksQuery } from "~features/links/hooks/useLinks";
import {
  getItemTypeMeta,
  getTabForItemType,
  replaceMentionsWithPills,
  type ResolvedMentionItem,
  type ProjectTabType,
} from "../utils/boardUtils";
import type { ItemType, CreateCardItemCommand } from "~types/api";
import styles from "./KanbanWorkspace.module.css";

interface MarkdownWithMentionsProps {
  text: string;
  projectId: string;
  linkedItems?: CreateCardItemCommand[];
}

interface HoverState {
  visible: boolean;
  itemType: ItemType;
  itemId: string;
  displayName: string;
  x: number;
  y: number;
  bottom: number;
}

type ItemData = Record<string, string | boolean | undefined>;

export const MarkdownWithMentions: React.FC<MarkdownWithMentionsProps> = ({
  text,
  projectId,
  linkedItems,
}) => {
  const navigate = useNavigate();
  const containerRef = useRef<HTMLDivElement>(null);
  const hoverTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const [hover, setHover] = useState<HoverState | null>(null);

  // Fetch project data for hover preview
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const { data: problemsData } = useProblemsQuery(projectId);
  const { data: credentialsData } = useCredentialsQuery(projectId);
  const { data: notesData } = useNotesQuery(projectId);
  const { data: linksData } = useLinksQuery(projectId);

  // Build a lookup map for all items
  const itemLookup = useMemo(() => {
    const map = new Map<string, { title: string; type: ItemType; data: ItemData }>();

    (snippetsData?.content || []).forEach((s) => {
      map.set(`SNIPPET-${s.id}`, {
        title: s.title,
        type: "SNIPPET",
        data: { language: s.language, snippetType: s.snippetType, description: s.description, content: s.content },
      });
    });

    (problemsData?.content || []).forEach((p) => {
      map.set(`PROBLEM-${p.id}`, {
        title: p.title,
        type: "PROBLEM",
        data: { status: p.status, severity: p.severity },
      });
    });

    (credentialsData?.content || []).forEach((c) => {
      map.set(`CREDENTIAL-${c.id}`, {
        title: c.title,
        type: "CREDENTIAL",
        data: { secretType: c.secretType, notes: c.notes, relatedUrl: c.relatedUrl },
      });
    });

    (notesData?.content || []).forEach((n) => {
      map.set(`NOTE-${n.id}`, {
        title: n.title,
        type: "NOTE",
        data: { archived: n.archived },
      });
    });

    (linksData?.content || []).forEach((l) => {
      map.set(`LINK-${l.id}`, {
        title: l.title,
        type: "LINK",
        data: { url: l.url, description: l.description },
      });
    });

    return map;
  }, [snippetsData, problemsData, credentialsData, notesData, linksData]);

  // List of resolved mention items (if linkedItems provided, prioritize those; else all project items)
  const availableMentionItems: ResolvedMentionItem[] = useMemo(() => {
    const list: ResolvedMentionItem[] = [];

    if (linkedItems && linkedItems.length > 0) {
      linkedItems.forEach((li) => {
        const found = itemLookup.get(`${li.itemType}-${li.itemId}`);
        if (found) {
          list.push({ id: li.itemId, type: li.itemType, title: found.title });
        }
      });
    } else {
      // Fallback: all project items
      itemLookup.forEach((val, key) => {
        const [_, id] = key.split("-");
        list.push({ id, type: val.type, title: val.title });
      });
    }

    return list;
  }, [linkedItems, itemLookup]);

  // Render markdown with title-based mention pills
  const renderedHtml = useMemo(() => {
    if (!text.trim()) return "";
    try {
      const rawHtml = marked.parse(text, { breaks: true, gfm: true }) as string;
      const withMentions = replaceMentionsWithPills(rawHtml, availableMentionItems);
      return DOMPurify.sanitize(withMentions, {
        ADD_ATTR: ["data-item-type", "data-item-id", "data-display-name", "style"],
      });
    } catch {
      return text;
    }
  }, [text, availableMentionItems]);

  // Navigate to item's tab with itemId
  const navigateToItem = useCallback(
    (itemType: ItemType, itemId?: string) => {
      const tab = getTabForItemType(itemType);
      navigate({
        to: "/projects/$projectId",
        params: { projectId },
        search: { tab, itemId } as { tab: ProjectTabType; itemId?: string },
      });
    },
    [navigate, projectId]
  );

  // Mouse event handlers for pills
  const handleMouseOver = useCallback((e: MouseEvent) => {
    const target = (e.target as HTMLElement).closest(".devaulty-mention-pill") as HTMLElement | null;
    if (!target) return;

    if (hoverTimerRef.current) clearTimeout(hoverTimerRef.current);

    hoverTimerRef.current = setTimeout(() => {
      const rect = target.getBoundingClientRect();
      setHover({
        visible: true,
        itemType: target.dataset.itemType as ItemType,
        itemId: target.dataset.itemId || "",
        displayName: target.dataset.displayName || "",
        x: rect.left + rect.width / 2,
        y: rect.top,
        bottom: rect.bottom,
      });
    }, 150);
  }, []);

  const handleMouseOut = useCallback((e: MouseEvent) => {
    const target = e.target as HTMLElement;
    const relatedTarget = e.relatedTarget as HTMLElement | null;

    // Don't hide if moving to the hover card itself
    if (relatedTarget?.closest?.(".devaulty-hover-card")) return;
    if (target.closest?.(".devaulty-mention-pill") && relatedTarget?.closest?.(".devaulty-mention-pill")) return;

    if (hoverTimerRef.current) clearTimeout(hoverTimerRef.current);
    hoverTimerRef.current = setTimeout(() => {
      setHover(null);
    }, 150);
  }, []);

  const handleClick = useCallback(
    (e: MouseEvent) => {
      const target = (e.target as HTMLElement).closest(".devaulty-mention-pill") as HTMLElement | null;
      if (!target) return;
      e.preventDefault();
      const itemType = target.dataset.itemType as ItemType;
      const itemId = target.dataset.itemId;
      navigateToItem(itemType, itemId);
    },
    [navigateToItem]
  );

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return;

    container.addEventListener("mouseover", handleMouseOver);
    container.addEventListener("mouseout", handleMouseOut);
    container.addEventListener("click", handleClick);

    return () => {
      container.removeEventListener("mouseover", handleMouseOver);
      container.removeEventListener("mouseout", handleMouseOut);
      container.removeEventListener("click", handleClick);
      if (hoverTimerRef.current) clearTimeout(hoverTimerRef.current);
    };
  }, [handleMouseOver, handleMouseOut, handleClick]);

  if (!text.trim()) {
    return (
      <span className="text-muted-foreground italic text-xs">
        No description written yet. Switch to "Write / Code" to add markdown details.
      </span>
    );
  }

  // Resolve hover item data
  const hoverItem = hover
    ? itemLookup.get(`${hover.itemType}-${hover.itemId}`)
    : null;

  return (
    <>
      <div
        ref={containerRef}
        className={`${styles.markdownBody} overflow-x-auto select-text`}
        dangerouslySetInnerHTML={{ __html: renderedHtml }}
      />

      {/* Hover Card */}
      {hover?.visible && (
        <HoverCard
          hover={hover}
          item={hoverItem}
          onMouseEnter={() => {
            if (hoverTimerRef.current) clearTimeout(hoverTimerRef.current);
          }}
          onMouseLeave={() => {
            setHover(null);
          }}
          onClick={() => {
            navigateToItem(hover.itemType, hover.itemId);
            setHover(null);
          }}
        />
      )}
    </>
  );
};

// ── Hover Card Component ──────────────────────────────────

interface HoverCardProps {
  hover: HoverState;
  item?: { title: string; type: ItemType; data: Record<string, string | boolean | undefined> } | null;
  onMouseEnter: () => void;
  onMouseLeave: () => void;
  onClick: () => void;
}

const HoverCard: React.FC<HoverCardProps> = ({
  hover,
  item,
  onMouseEnter,
  onMouseLeave,
  onClick,
}) => {
  const meta = getItemTypeMeta(hover.itemType);
  const Icon = meta.icon;
  const title = item?.title || hover.displayName;
  const data: Record<string, string | boolean | undefined> = item?.data || {};

  // Smart viewport clamping & flipping
  const cardWidth = 320;
  const estimatedHeight = 300;
  const placeBelow = hover.y < estimatedHeight + 30;

  const clampedX = Math.max(
    cardWidth / 2 + 12,
    Math.min(hover.x, window.innerWidth - cardWidth / 2 - 12)
  );

  const clampedTop = placeBelow
    ? Math.min(hover.bottom + 8, window.innerHeight - 120)
    : Math.max(12, hover.y - 8);

  const cardStyle: React.CSSProperties = {
    position: "fixed",
    left: clampedX,
    top: clampedTop,
    transform: placeBelow ? "translate(-50%, 0)" : "translate(-50%, -100%)",
    zIndex: 9999,
  };

  return (
    <div
      className="devaulty-hover-card animate-in fade-in zoom-in-95 duration-150"
      style={cardStyle}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
    >
      <div
        className="w-80 max-w-sm rounded-lg border shadow-2xl overflow-hidden backdrop-blur-xl flex flex-col"
        style={{
          background: "var(--color-card)",
          borderColor: `color-mix(in srgb, ${meta.color} 30%, var(--color-border))`,
          maxHeight: placeBelow
            ? `${Math.max(200, window.innerHeight - clampedTop - 20)}px`
            : "360px",
        }}
      >
        {/* Header */}
        <div
          className="flex items-center justify-between px-3 py-2 border-b shrink-0"
          style={{
            borderColor: `color-mix(in srgb, ${meta.color} 15%, var(--color-border))`,
            background: `color-mix(in srgb, ${meta.color} 6%, transparent)`,
          }}
        >
          <div className="flex items-center gap-2">
            <span
              className="p-1 rounded"
              style={{
                background: `color-mix(in srgb, ${meta.color} 15%, transparent)`,
                color: meta.color,
              }}
            >
              <Icon size={13} />
            </span>
            <span
              className="text-[10px] font-bold font-mono uppercase tracking-wider"
              style={{ color: meta.color }}
            >
              {meta.label} Preview
            </span>
          </div>

          <span className="text-[9px] font-mono text-muted-foreground opacity-70">
            Click to open
          </span>
        </div>

        {/* Scrollable Body */}
        <div className="p-3 flex flex-col gap-2.5 overflow-y-auto max-h-64 scrollbar-thin">
          <h4 className="text-xs font-semibold text-foreground leading-snug">
            {title}
          </h4>

          {/* Snippet Preview */}
          {hover.itemType === "SNIPPET" && (
            <div className="flex flex-col gap-1.5">
              <div className="flex items-center gap-2">
                {data.language && (
                  <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-secondary border border-border text-muted-foreground">
                    {String(data.language)}
                  </span>
                )}
                {data.snippetType && (
                  <span className="text-[10px] font-mono text-muted-foreground">
                    {String(data.snippetType).toLowerCase()}
                  </span>
                )}
              </div>
              {data.description && (
                <p className="text-[11px] text-muted-foreground leading-relaxed">
                  {String(data.description)}
                </p>
              )}
              {data.content && (
                <pre className="text-[10px] font-mono leading-relaxed p-2 rounded bg-secondary/80 border border-border text-foreground overflow-auto max-h-36 scrollbar-thin select-text">
                  {String(data.content)}
                </pre>
              )}
            </div>
          )}

          {/* Problem Preview */}
          {hover.itemType === "PROBLEM" && (
            <div className="flex flex-col gap-2">
              <div className="flex items-center gap-2">
                <span
                  className={`text-[10px] font-mono font-bold px-1.5 py-0.5 rounded ${
                    data.severity === "CRITICAL"
                      ? "bg-red-500/15 text-red-500 border border-red-500/30"
                      : data.severity === "HIGH"
                      ? "bg-amber-500/15 text-amber-500 border border-amber-500/30"
                      : "bg-slate-500/15 text-slate-500 border border-slate-500/30"
                  }`}
                >
                  {String(data.severity || "MEDIUM")}
                </span>
                <span
                  className={`text-[10px] font-mono px-1.5 py-0.5 rounded ${
                    data.status === "RESOLVED"
                      ? "bg-emerald-500/15 text-emerald-500 border border-emerald-500/30"
                      : data.status === "OPEN"
                      ? "bg-blue-500/15 text-blue-500 border border-blue-500/30"
                      : "bg-amber-500/15 text-amber-500 border border-amber-500/30"
                  }`}
                >
                  {String(data.status || "OPEN")}
                </span>
              </div>

              {data.errorDescription && (
                <div className="flex flex-col gap-0.5">
                  <span className="text-[9px] font-mono uppercase text-muted-foreground font-semibold">Error Details:</span>
                  <p className="text-[11px] text-foreground bg-secondary/40 p-1.5 rounded border border-border/60 font-mono text-[10px] leading-relaxed max-h-24 overflow-y-auto select-text">
                    {String(data.errorDescription)}
                  </p>
                </div>
              )}

              {data.solution && (
                <div className="flex flex-col gap-0.5">
                  <span className="text-[9px] font-mono uppercase text-emerald-500 font-semibold">Solution:</span>
                  <p className="text-[11px] text-foreground bg-emerald-500/5 p-1.5 rounded border border-emerald-500/20 font-mono text-[10px] leading-relaxed max-h-24 overflow-y-auto select-text">
                    {String(data.solution)}
                  </p>
                </div>
              )}
            </div>
          )}

          {/* Credential Preview */}
          {hover.itemType === "CREDENTIAL" && (
            <div className="flex flex-col gap-1.5">
              <div className="flex items-center gap-2">
                <Icons.Shield size={11} className="text-purple-400" />
                <span className="text-[10px] font-mono text-muted-foreground">
                  Type: {String(data.secretType || "LOGIN")}
                </span>
              </div>
              {data.relatedUrl && (
                <span className="text-[10px] font-mono text-primary truncate">
                  {String(data.relatedUrl)}
                </span>
              )}
              {data.notes && (
                <p className="text-[11px] text-muted-foreground bg-secondary/30 p-1.5 rounded border border-border/60 leading-relaxed max-h-24 overflow-y-auto select-text">
                  {String(data.notes)}
                </p>
              )}
            </div>
          )}

          {/* Note Preview */}
          {hover.itemType === "NOTE" && (
            <div className="flex flex-col gap-1.5">
              {data.archived && (
                <span className="text-[10px] font-mono px-1.5 py-0.5 rounded bg-amber-500/15 text-amber-500 border border-amber-500/30 w-fit">
                  Archived
                </span>
              )}
              {data.content && (
                <p className="text-[11px] text-foreground bg-secondary/30 p-2 rounded border border-border/60 leading-relaxed max-h-32 overflow-y-auto select-text font-sans">
                  {String(data.content)}
                </p>
              )}
            </div>
          )}

          {/* Link Preview */}
          {hover.itemType === "LINK" && (
            <div className="flex flex-col gap-1.5">
              {data.url && (
                <span className="text-[10px] font-mono text-cyan-400 truncate">
                  {String(data.url)}
                </span>
              )}
              {data.description && (
                <p className="text-[11px] text-muted-foreground bg-secondary/30 p-1.5 rounded border border-border/60 leading-relaxed max-h-24 overflow-y-auto select-text">
                  {String(data.description)}
                </p>
              )}
            </div>
          )}
        </div>

        {/* Footer — Go to item */}
        <button
          type="button"
          onClick={onClick}
          className="w-full flex items-center justify-center gap-1.5 px-3 py-1.5 border-t text-[10px] font-mono font-semibold transition-colors cursor-pointer hover:brightness-110 shrink-0"
          style={{
            borderColor: `color-mix(in srgb, ${meta.color} 15%, var(--color-border))`,
            color: meta.color,
            background: `color-mix(in srgb, ${meta.color} 6%, transparent)`,
          }}
        >
          <Icons.ExternalLink size={11} />
          <span>Open Full {meta.label}</span>
        </button>
      </div>
    </div>
  );
};
