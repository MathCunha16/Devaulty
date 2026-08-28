import React from "react";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import * as Icons from "lucide-react";
import type { CardSummaryView } from "~types/api";
import styles from "./KanbanWorkspace.module.css";

interface KanbanCardProps {
  card: CardSummaryView;
  onOpenCard: (cardId: string) => void;
}

export const KanbanCard: React.FC<KanbanCardProps> = ({ card, onOpenCard }) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({
    id: card.id,
    data: {
      type: "CARD",
      card,
    },
  });

  const style: React.CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition,
  };

  const formatDueDate = (dateStr?: string) => {
    if (!dateStr) return null;
    const date = new Date(dateStr);
    const now = new Date();
    const isOverdue = date < now && date.toDateString() !== now.toDateString();
    return {
      formatted: date.toLocaleDateString(undefined, { month: "short", day: "numeric" }),
      isOverdue,
    };
  };

  const due = formatDueDate(card.dueDate);

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      {...listeners}
      onClick={() => onOpenCard(card.id)}
      className={`${styles.card} ${isDragging ? styles.cardDragging : ""}`}
    >
      {/* Header: Title & Priority */}
      <div className={styles.cardHeader}>
        <h4 className={styles.cardTitle}>{card.title}</h4>
        {card.priority && (
          <span
            className={`${styles.priorityPill} ${
              card.priority === "EXTREMELY_HIGH"
                ? styles.priorityExtremelyHigh
                : card.priority === "HIGH"
                ? styles.priorityHigh
                : card.priority === "MEDIUM"
                ? styles.priorityMedium
                : styles.priorityLow
            }`}
          >
            {card.priority === "EXTREMELY_HIGH" ? "Urgent" : card.priority}
          </span>
        )}
      </div>

      {/* Tags Row */}
      {card.tags && card.tags.length > 0 && (
        <div className="flex flex-wrap gap-1">
          {card.tags.map((t) => (
            <span
              key={t.id}
              className="text-[10px] px-1.5 py-0.5 rounded-full font-mono border"
              style={{
                borderColor: t.color || "var(--color-border)",
                color: t.color || "var(--color-muted-foreground)",
                backgroundColor: t.color ? `${t.color}15` : "var(--color-secondary)",
              }}
            >
              {t.name}
            </span>
          ))}
        </div>
      )}

      {/* Meta Row: Due Date & Details Hint */}
      <div className="flex items-center justify-between pt-1 mt-auto text-[11px] text-muted-foreground">
        {due ? (
          <span
            className={`${styles.dueDatePill} ${due.isOverdue ? styles.dueDateOverdue : ""}`}
          >
            <Icons.Calendar size={10} />
            <span>{due.formatted}</span>
          </span>
        ) : (
          <div />
        )}

        <span className="opacity-0 group-hover:opacity-100 hover:text-primary transition-opacity flex items-center gap-0.5 text-[10px] font-mono">
          <span>Edit</span>
          <Icons.ChevronRight size={10} />
        </span>
      </div>
    </div>
  );
};
