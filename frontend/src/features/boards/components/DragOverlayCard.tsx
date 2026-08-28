import React from "react";
import * as Icons from "lucide-react";
import type { CardSummaryView } from "~types/api";
import styles from "./KanbanWorkspace.module.css";

interface DragOverlayCardProps {
  card: CardSummaryView;
}

export const DragOverlayCard: React.FC<DragOverlayCardProps> = ({ card }) => {
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
    <div className={styles.dragOverlayCard}>
      {/* Title & Priority */}
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

      {/* Tags */}
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

      {/* Due Date & Meta */}
      {due && (
        <div className="flex items-center gap-1">
          <span
            className={`${styles.dueDatePill} ${due.isOverdue ? styles.dueDateOverdue : ""}`}
          >
            <Icons.Calendar size={10} />
            <span>{due.formatted}</span>
          </span>
        </div>
      )}
    </div>
  );
};
