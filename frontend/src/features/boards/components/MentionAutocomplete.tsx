import React, { useEffect, useRef } from "react";
import * as Icons from "lucide-react";
import { getItemTypeMeta, type ResolvedMentionItem } from "../utils/boardUtils";

interface MentionAutocompleteProps {
  isOpen: boolean;
  position: { top: number; left: number };
  items: ResolvedMentionItem[];
  hasAttachedItems: boolean;
  selectedIndex: number;
  query: string;
  onSelect: (item: ResolvedMentionItem) => void;
  onClose: () => void;
}

export const MentionAutocomplete: React.FC<MentionAutocompleteProps> = ({
  isOpen,
  position,
  items,
  hasAttachedItems,
  selectedIndex,
  query,
  onSelect,
  onClose,
}) => {
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) return;

    // Scroll active item into view
    const activeEl = listRef.current?.children[selectedIndex] as HTMLElement | undefined;
    if (activeEl) {
      activeEl.scrollIntoView({ block: "nearest" });
    }
  }, [isOpen, selectedIndex]);

  if (!isOpen) return null;

  return (
    <div
      className="absolute z-50 w-80 rounded-lg bg-card/95 backdrop-blur-xl border border-border shadow-2xl overflow-hidden animate-in fade-in zoom-in-95 duration-100 flex flex-col"
      style={{
        top: position.top + 24, // 24px below caret line
        left: Math.max(10, Math.min(position.left, 240)), // keep within reasonable bounds
      }}
    >
      {/* Header / Query indicator */}
      <div className="flex items-center justify-between px-3 py-1.5 border-b border-border/80 bg-secondary/40 text-[10px] font-mono text-muted-foreground">
        <div className="flex items-center gap-1">
          <Icons.AtSign size={11} className="text-primary" />
          <span>Mention Item</span>
          {query && <span className="text-foreground font-semibold">"{query}"</span>}
        </div>
        <div className="flex items-center gap-1.5">
          {hasAttachedItems && <span className="text-[9px] opacity-70">↑↓ ↵</span>}
          <button
            type="button"
            onMouseDown={(e) => {
              e.preventDefault();
              onClose();
            }}
            className="hover:text-foreground p-0.5 rounded cursor-pointer"
            title="Close"
          >
            <Icons.X size={10} />
          </button>
        </div>
      </div>

      {/* Items list / Empty states */}
      <div ref={listRef} className="max-h-48 overflow-y-auto p-1.5 flex flex-col gap-0.5">
        {!hasAttachedItems ? (
          <div className="p-3 text-center text-xs text-muted-foreground font-mono flex flex-col items-center gap-1.5 bg-secondary/20 rounded">
            <div className="p-1.5 rounded-full bg-amber-500/10 text-amber-500 border border-amber-500/30">
              <Icons.Paperclip size={14} />
            </div>
            <span className="font-semibold text-foreground text-[11px]">No items attached to card</span>
            <p className="text-[10px] text-muted-foreground leading-relaxed">
              Attach items in the <strong>"Linked Devaulty Items"</strong> section below before mentioning them with <strong>@</strong>.
            </p>
          </div>
        ) : items.length === 0 ? (
          <div className="px-3 py-4 text-center text-xs text-muted-foreground font-mono flex flex-col items-center gap-1.5">
            <Icons.SearchX size={16} className="text-muted-foreground/60" />
            <span>No attached items match "{query}"</span>
            <span className="text-[10px] text-muted-foreground/70">
              Only items attached to this card can be mentioned.
            </span>
          </div>
        ) : (
          items.map((item, idx) => {
            const meta = getItemTypeMeta(item.type);
            const Icon = meta.icon;
            const isSelected = idx === selectedIndex;

            return (
              <button
                key={`${item.type}-${item.id}`}
                type="button"
                onMouseDown={(e) => {
                  // Prevent textarea from losing focus before selection completes
                  e.preventDefault();
                  onSelect(item);
                }}
                className={`w-full flex items-center justify-between px-2.5 py-1.5 rounded text-left transition-colors cursor-pointer text-xs ${
                  isSelected
                    ? "bg-primary/15 text-foreground font-medium ring-1 ring-primary/40"
                    : "hover:bg-secondary/70 text-foreground"
                }`}
              >
                <div className="flex items-center gap-2 overflow-hidden min-w-0">
                  <span
                    className="p-1 rounded shrink-0"
                    style={{
                      backgroundColor: `color-mix(in srgb, ${meta.color} 15%, transparent)`,
                      color: meta.color,
                    }}
                  >
                    <Icon size={12} />
                  </span>
                  <span className="truncate text-xs">{item.title}</span>
                </div>

                <span
                  className="text-[9px] font-mono font-semibold uppercase px-1.5 py-0.5 rounded shrink-0 ml-2"
                  style={{
                    backgroundColor: `color-mix(in srgb, ${meta.color} 10%, transparent)`,
                    color: meta.color,
                    border: `1px solid color-mix(in srgb, ${meta.color} 25%, transparent)`,
                  }}
                >
                  {meta.label}
                </span>
              </button>
            );
          })
        )}
      </div>
    </div>
  );
};
