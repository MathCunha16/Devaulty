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
 * Uses clean syntax: @[Item Title]
 */
export const formatMentionToken = (title: string): string => {
  return `@[${title}]`;
};

/**
 * Replace mention tokens (e.g. @[Item Title] or @Title) in HTML with interactive pill spans.
 * Matches against the list of available items linked to the card.
 */
export const replaceMentionsWithPills = (
  html: string,
  availableItems: ResolvedMentionItem[]
): string => {
  if (!availableItems.length || !html) return html;

  let result = html;

  // 1. First replace bracketed mentions: @[Item Title]
  availableItems.forEach((item) => {
    const escapedTitle = item.title.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
    const bracketRegex = new RegExp(`@\\[${escapedTitle}\\]`, "gi");

    const meta = getItemTypeMeta(item.type);
    const pillHtml = `<span class="devaulty-mention-pill" data-item-type="${item.type}" data-item-id="${item.id}" data-display-name="${item.title}" style="--mention-color: ${meta.color}; cursor: pointer;">${item.title}</span>`;

    result = result.replace(bracketRegex, pillHtml);
  });

  // 2. Also replace single-word direct mentions: @ItemTitle (if word boundaries match)
  availableItems.forEach((item) => {
    // Only if title has no spaces
    if (!item.title.includes(" ")) {
      const escapedTitle = item.title.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      const wordRegex = new RegExp(`\\B@${escapedTitle}\\b`, "gi");

      const meta = getItemTypeMeta(item.type);
      const pillHtml = `<span class="devaulty-mention-pill" data-item-type="${item.type}" data-item-id="${item.id}" data-display-name="${item.title}" style="--mention-color: ${meta.color}; cursor: pointer;">${item.title}</span>`;

      result = result.replace(wordRegex, pillHtml);
    }
  });

  return result;
};
