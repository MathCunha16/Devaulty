import React from "react";
import { useSortable, SortableContext, verticalListSortingStrategy } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import * as Icons from "lucide-react";
import type { BoardColumnView, CardSummaryView } from "~types/api";
import { KanbanCard } from "./KanbanCard";
import styles from "./KanbanWorkspace.module.css";

interface BoardColumnProps {
  column: BoardColumnView;
  cards: CardSummaryView[];
  onOpenCard: (cardId: string) => void;
  onAddCardToColumn: (columnId: string) => void;
  onEditColumn: (column: BoardColumnView) => void;
}

export const BoardColumn: React.FC<BoardColumnProps> = ({
  column,
  cards,
  onOpenCard,
  onAddCardToColumn,
  onEditColumn,
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
    isOver,
  } = useSortable({
    id: column.id,
    data: {
      type: "COLUMN",
      column,
    },
  });

  const style: React.CSSProperties = {
    transform: CSS.Translate.toString(transform),
    transition,
  };

  const cardIds = React.useMemo(() => cards.map((c) => c.id), [cards]);

  const isWipExceeded =
    column.wipLimit != null &&
    column.wipLimit > 0 &&
    cards.length > column.wipLimit;

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`${styles.column} ${isOver ? styles.columnDragOver : ""} ${
        isDragging ? "opacity-40 border-dashed" : ""
      }`}
    >
      {/* Column Header */}
      <div className={styles.columnHeader}>
        <div className={styles.columnTitleBox}>
          {/* Column Drag Handle */}
          <div
            {...attributes}
            {...listeners}
            className="cursor-grab active:cursor-grabbing p-0.5 text-muted-foreground hover:text-foreground transition-colors"
            title="Drag to reorder column"
          >
            <Icons.GripVertical size={14} />
          </div>

          <h3 className={styles.columnTitle} title={column.name}>
            {column.name}
          </h3>
          <span
            className={`${styles.columnCount} ${
              isWipExceeded ? styles.wipWarning : ""
            }`}
            title={
              column.wipLimit
                ? `WIP Limit: ${column.wipLimit} cards (Current: ${cards.length})`
                : `${cards.length} cards`
            }
          >
            {column.wipLimit ? `${cards.length}/${column.wipLimit}` : cards.length}
          </span>
        </div>

        <div className={styles.columnActions}>
          <button
            type="button"
            onClick={() => onAddCardToColumn(column.id)}
            className={styles.columnActionBtn}
            title="Add card to column"
          >
            <Icons.Plus size={14} />
          </button>
          <button
            type="button"
            onClick={() => onEditColumn(column)}
            className={styles.columnActionBtn}
            title="Edit column settings"
          >
            <Icons.MoreHorizontal size={14} />
          </button>
        </div>
      </div>

      {/* Cards List Container */}
      <div className={styles.cardList}>
        <SortableContext
          items={cardIds}
          strategy={verticalListSortingStrategy}
        >
          {cards.map((card) => (
            <KanbanCard
              key={card.id}
              card={card}
              onOpenCard={onOpenCard}
            />
          ))}
        </SortableContext>

        {cards.length === 0 && (
          <div className="flex flex-col items-center justify-center p-6 text-center text-muted-foreground border border-dashed border-border/60 rounded-md my-auto">
            <span className="text-[11px] font-mono">Empty Column</span>
          </div>
        )}
      </div>

      {/* Bottom Quick Add */}
      <div className={styles.columnFooter}>
        <button
          type="button"
          onClick={() => onAddCardToColumn(column.id)}
          className={styles.quickAddCardBtn}
        >
          <Icons.Plus size={13} />
          <span>New Card</span>
        </button>
      </div>
    </div>
  );
};
