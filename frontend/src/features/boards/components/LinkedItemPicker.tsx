import React, { useState, useMemo } from "react";
import * as Icons from "lucide-react";
import { useNavigate } from "@tanstack/react-router";
import { useSnippetsQuery } from "~features/snippets/hooks/useSnippets";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { useCredentialsQuery } from "~features/credentials/hooks/useCredentials";
import { useNotesQuery } from "~features/notes/hooks/useNotes";
import { useLinksQuery } from "~features/links/hooks/useLinks";
import { getItemTypeMeta, getTabForItemType, type ProjectTabType } from "../utils/boardUtils";
import type { ItemType, CreateCardItemCommand } from "~types/api";

interface LinkedItemPickerProps {
  projectId: string;
  linkedItems: CreateCardItemCommand[];
  onChange: (items: CreateCardItemCommand[]) => void;
}

interface AvailableItem {
  id: string;
  title: string;
  type: ItemType;
}

export const LinkedItemPicker: React.FC<LinkedItemPickerProps> = ({
  projectId,
  linkedItems,
  onChange,
}) => {
  const navigate = useNavigate();
  const [isOpen, setIsOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [filterType, setFilterType] = useState<ItemType | "ALL">("ALL");

  // Fetch project entities for the picker
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const { data: problemsData } = useProblemsQuery(projectId);
  const { data: credentialsData } = useCredentialsQuery(projectId);
  const { data: notesData } = useNotesQuery(projectId);
  const { data: linksData } = useLinksQuery(projectId);

  const allAvailableItems: AvailableItem[] = useMemo(() => {
    const list: AvailableItem[] = [];

    (snippetsData?.content || []).forEach((s) => {
      list.push({ id: s.id, title: s.title, type: "SNIPPET" });
    });

    (problemsData?.content || []).forEach((p) => {
      list.push({ id: p.id, title: p.title, type: "PROBLEM" });
    });

    (credentialsData?.content || []).forEach((c) => {
      list.push({ id: c.id, title: c.title, type: "CREDENTIAL" });
    });

    (notesData?.content || []).forEach((n) => {
      list.push({ id: n.id, title: n.title, type: "NOTE" });
    });

    (linksData?.content || []).forEach((l) => {
      list.push({ id: l.id, title: l.title, type: "LINK" });
    });

    return list;
  }, [snippetsData, problemsData, credentialsData, notesData, linksData]);

  // Title lookup map
  const itemMap = useMemo(() => {
    const map = new Map<string, AvailableItem>();
    allAvailableItems.forEach((it) => map.set(`${it.type}-${it.id}`, it));
    return map;
  }, [allAvailableItems]);

  // Filter items by type and search query (excluding already attached items)
  const filteredAvailable = useMemo(() => {
    const linkedSet = new Set(linkedItems.map((item) => `${item.itemType}-${item.itemId}`));
    return allAvailableItems.filter((item) => {
      if (linkedSet.has(`${item.type}-${item.id}`)) return false;
      if (filterType !== "ALL" && item.type !== filterType) return false;
      if (search.trim()) {
        return item.title.toLowerCase().includes(search.toLowerCase());
      }
      return true;
    });
  }, [allAvailableItems, linkedItems, filterType, search]);

  const handleAddItem = (item: AvailableItem) => {
    onChange([...linkedItems, { itemType: item.type, itemId: item.id }]);
    setSearch("");
  };

  const handleRemoveItem = (index: number) => {
    const updated = linkedItems.filter((_, i) => i !== index);
    onChange(updated);
  };

  const handleNavigateToItem = (itemType: ItemType, itemId: string) => {
    const tab = getTabForItemType(itemType);
    navigate({
      to: "/projects/$projectId",
      params: { projectId },
      search: { tab, itemId } as { tab: ProjectTabType; itemId?: string },
    });
  };

  return (
    <div className="flex flex-col gap-2">
      {/* Section Header with Single Toggle Button */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-1.5">
          <Icons.Paperclip size={13} className="text-muted-foreground" />
          <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
            Linked Devaulty Items ({linkedItems.length})
          </label>
        </div>
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="text-xs flex items-center gap-1 text-primary hover:underline font-mono cursor-pointer transition-colors"
        >
          {isOpen ? (
            <>
              <Icons.ChevronUp size={13} />
              <span>Done</span>
            </>
          ) : (
            <>
              <Icons.Plus size={13} />
              <span>Attach Item</span>
            </>
          )}
        </button>
      </div>

      {/* Unified Single Row for Linked Items */}
      {linkedItems.length === 0 ? (
        <div className="p-3 rounded-md bg-secondary/30 border border-border/60 text-center">
          <p className="text-xs text-muted-foreground font-mono">
            No items linked yet. Attach snippets, problems, vaults, or notes to mention them with <code className="px-1 py-0.5 rounded bg-secondary text-primary font-bold">@</code> in description.
          </p>
        </div>
      ) : (
        <div className="flex flex-wrap gap-1.5 p-2 rounded-md bg-secondary/40 border border-border/70">
          {linkedItems.map((item, index) => {
            const meta = getItemTypeMeta(item.itemType);
            const Icon = meta.icon;
            const itemObj = itemMap.get(`${item.itemType}-${item.itemId}`);
            const title = itemObj?.title || `${meta.label} (${item.itemId.slice(0, 6)}...)`;

            return (
              <div
                key={`${item.itemType}-${item.itemId}-${index}`}
                className={`group inline-flex items-center gap-1.5 px-2.5 py-1 rounded text-xs font-mono border transition-all ${meta.bgColor} ${meta.borderColor} ${meta.textColor}`}
              >
                {/* Clickable Icon & Title to Navigate to Item */}
                <button
                  type="button"
                  onClick={() => handleNavigateToItem(item.itemType, item.itemId)}
                  className="inline-flex items-center gap-1.5 hover:underline cursor-pointer text-left"
                  title={`Open ${meta.label}: ${title}`}
                >
                  <Icon size={12} className="shrink-0" />
                  <span className="max-w-[160px] truncate font-medium">{title}</span>
                  <Icons.ExternalLink size={10} className="opacity-60 group-hover:opacity-100 transition-opacity shrink-0" />
                </button>

                {/* Remove button */}
                <button
                  type="button"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemoveItem(index);
                  }}
                  className="hover:bg-destructive/20 hover:text-destructive rounded p-0.5 cursor-pointer ml-1 transition-colors"
                  title="Remove from card"
                >
                  <Icons.X size={11} />
                </button>
              </div>
            );
          })}
        </div>
      )}

      {/* Attach Item Drawer */}
      {isOpen && (
        <div className="flex flex-col gap-2 p-3 rounded-lg bg-card border border-border shadow-xl animate-in fade-in zoom-in-95 duration-150 mt-1">
          {/* Type Filter Pills */}
          <div className="flex items-center gap-1 overflow-x-auto pb-1 text-xs">
            {(["ALL", "SNIPPET", "CREDENTIAL", "PROBLEM", "NOTE", "LINK"] as const).map((type) => (
              <button
                key={type}
                type="button"
                onClick={() => setFilterType(type)}
                className={`px-2 py-0.5 rounded font-mono text-[11px] transition-colors cursor-pointer whitespace-nowrap ${
                  filterType === type
                    ? "bg-primary text-primary-foreground font-semibold"
                    : "bg-secondary/70 text-muted-foreground hover:text-foreground"
                }`}
              >
                {type === "ALL" ? "All Items" : getItemTypeMeta(type).label}
              </button>
            ))}
          </div>

          {/* Search Bar */}
          <div className="relative flex items-center">
            <Icons.Search size={13} className="absolute left-2.5 text-muted-foreground" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search items to attach..."
              autoFocus
              className="w-full pl-8 pr-3 py-1.5 bg-background border border-border rounded text-xs text-foreground placeholder:text-muted-foreground outline-none focus:border-primary"
            />
          </div>

          {/* List of Available Items */}
          <div className="max-h-40 overflow-y-auto flex flex-col gap-1 pr-1">
            {filteredAvailable.length === 0 ? (
              <div className="text-center py-4 text-xs text-muted-foreground font-mono">
                {allAvailableItems.length === 0
                  ? "No items in this project yet."
                  : "No matching unattached items found."}
              </div>
            ) : (
              filteredAvailable.map((item) => {
                const meta = getItemTypeMeta(item.type);
                const Icon = meta.icon;
                return (
                  <button
                    key={`${item.type}-${item.id}`}
                    type="button"
                    onClick={() => handleAddItem(item)}
                    className="flex items-center justify-between p-1.5 rounded hover:bg-secondary/80 text-left text-xs transition-colors cursor-pointer group"
                  >
                    <div className="flex items-center gap-2 overflow-hidden">
                      <span
                        className="p-1 rounded"
                        style={{
                          backgroundColor: `color-mix(in srgb, ${meta.color} 15%, transparent)`,
                          color: meta.color,
                        }}
                      >
                        <Icon size={12} />
                      </span>
                      <span className="font-medium text-foreground truncate">{item.title}</span>
                    </div>
                    <span className="text-[10px] font-mono text-muted-foreground group-hover:text-primary shrink-0 ml-2">
                      + Attach
                    </span>
                  </button>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
};
