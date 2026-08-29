import * as Icons from "lucide-react";
import type { ItemType } from "~types/api";

export type ProjectTabType =
  | "overview"
  | "boards"
  | "snippets"
  | "problems"
  | "credentials"
  | "notes"
  | "links";

export interface ItemTypeMeta {
  label: string;
  icon: typeof Icons.Box;
  color: string;
  textColor: string;
  bgColor: string;
  borderColor: string;
}

export const getItemTypeMeta = (type: ItemType): ItemTypeMeta => {
  switch (type) {
    case "SNIPPET":
      return {
        label: "Snippet",
        icon: Icons.Code2,
        color: "#10b981",
        textColor: "text-emerald-500",
        bgColor: "bg-emerald-500/10",
        borderColor: "border-emerald-500/30",
      };
    case "CREDENTIAL":
      return {
        label: "Vault",
        icon: Icons.KeyRound,
        color: "#a855f7",
        textColor: "text-purple-500",
        bgColor: "bg-purple-500/10",
        borderColor: "border-purple-500/30",
      };
    case "PROBLEM":
      return {
        label: "Problem",
        icon: Icons.AlertCircle,
        color: "#ef4444",
        textColor: "text-red-500",
        bgColor: "bg-red-500/10",
        borderColor: "border-red-500/30",
      };
    case "NOTE":
      return {
        label: "Note",
        icon: Icons.FileText,
        color: "#3b82f6",
        textColor: "text-blue-500",
        bgColor: "bg-blue-500/10",
        borderColor: "border-blue-500/30",
      };
    case "LINK":
      return {
        label: "Link",
        icon: Icons.Link2,
        color: "#06b6d4",
        textColor: "text-cyan-500",
        bgColor: "bg-cyan-500/10",
        borderColor: "border-cyan-500/30",
      };
    default:
      return {
        label: "Item",
        icon: Icons.Box,
        color: "#64748b",
        textColor: "text-slate-500",
        bgColor: "bg-slate-500/10",
        borderColor: "border-slate-500/30",
      };
  }
};

export const getTabForItemType = (type: ItemType): ProjectTabType => {
  switch (type) {
    case "SNIPPET":
      return "snippets";
    case "PROBLEM":
      return "problems";
    case "CREDENTIAL":
      return "credentials";
    case "NOTE":
      return "notes";
    case "LINK":
      return "links";
    default:
      return "overview";
  }
};

export interface ResolvedMentionItem {
  id: string;
  type: ItemType;
  title: string;
}

/**
 * Format a mention token for inserting into markdown.
 * Uses clean syntax with item metadata: @[Item Title](item:TYPE:ID)
 */
export const formatMentionToken = (
  title: string,
  type?: ItemType,
  id?: string
): string => {
  if (type && id) {
    return `@[${title}](item:${type}:${id})`;
  }
  return `@[${title}]`;
};

const escapeHtml = (str: string): string => {
  return str
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&#039;");
};

/**
 * Replaces mentions in HTML with interactive stylized pills.
 */
export const replaceMentionsWithPills = (
  html: string,
  availableItems: ResolvedMentionItem[]
): string => {
  if (!html) return html;

  let result = html;

  // 1. First replace ID-based mentions: @<a href="item:TYPE:ID">Title</a> or <a href="item:TYPE:ID">Title</a>
  const itemLinkRegex =
    /@?<a\s+[^>]*href=["']item:([A-Z_]+):([a-zA-Z0-9-]+)["'][^>]*>([\s\S]*?)<\/a>/gi;
  result = result.replace(itemLinkRegex, (_match, _type, id, fallbackTitle) => {
    const item = availableItems.find((i) => i.id === id);
    if (item) {
      const meta = getItemTypeMeta(item.type);
      const safeTitle = escapeHtml(item.title);
      return `<span class="devaulty-mention-pill" data-item-type="${item.type}" data-item-id="${item.id}" data-display-name="${safeTitle}" style="--mention-color: ${meta.color}; cursor: pointer;">${safeTitle}</span>`;
    }
    // If not found, detached, or deleted, render fallback text without active pill
    return `@${fallbackTitle}`;
  });

  if (!availableItems.length) return result;

  // 2. Backward compatibility: replace legacy bracketed mentions: @[Item Title]
  availableItems.forEach((item) => {
    const escapedTitle = item.title.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const bracketRegex = new RegExp(`@\\[${escapedTitle}\\]`, "gi");

    const meta = getItemTypeMeta(item.type);
    const safeTitle = escapeHtml(item.title);
    const pillHtml = `<span class="devaulty-mention-pill" data-item-type="${item.type}" data-item-id="${item.id}" data-display-name="${safeTitle}" style="--mention-color: ${meta.color}; cursor: pointer;">${safeTitle}</span>`;

    result = result.replace(bracketRegex, pillHtml);
  });

  // 3. Backward compatibility: replace legacy single-word direct mentions: @ItemTitle
  availableItems.forEach((item) => {
    // Only if title has no spaces
    if (!item.title.includes(" ")) {
      const escapedTitle = item.title.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      const wordRegex = new RegExp(`\\B@${escapedTitle}\\b`, "gi");

      const meta = getItemTypeMeta(item.type);
      const safeTitle = escapeHtml(item.title);
      const pillHtml = `<span class="devaulty-mention-pill" data-item-type="${item.type}" data-item-id="${item.id}" data-display-name="${safeTitle}" style="--mention-color: ${meta.color}; cursor: pointer;">${safeTitle}</span>`;

      result = result.replace(wordRegex, pillHtml);
    }
  });

  return result;
};

/**
 * Formats a due date string (ISO / YYYY-MM-DD) safely without timezone rollover.
 */
export const formatDueDate = (
  dateStr?: string
): { formatted: string; isOverdue: boolean } | null => {
  if (!dateStr) return null;
  const datePart = dateStr.split("T")[0];
  const [year, month, day] = datePart.split("-").map(Number);
  if (!year || !month || !day) return null;

  const date = new Date(year, month - 1, day);
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const isOverdue = date < today;
  return {
    formatted: date.toLocaleDateString(undefined, { month: "short", day: "numeric" }),
    isOverdue,
  };
};
