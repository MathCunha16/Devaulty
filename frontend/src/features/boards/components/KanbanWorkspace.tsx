import React, { useState, useMemo } from "react";
import {
  DndContext,
  DragOverlay,
  closestCorners,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  type DragStartEvent,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  horizontalListSortingStrategy,
  arrayMove,
} from "@dnd-kit/sortable";
import * as Icons from "lucide-react";
import { toast } from "sonner";
import {
  useBoardsQuery,
  useDefaultBoardQuery,
  useBoardColumnsQuery,
  useBoardCardsQuery,
  useMoveCardMutation,
  useReorderColumnsMutation,
} from "../hooks/useBoards";
import { BoardColumn } from "./BoardColumn";
import { DragOverlayCard } from "./DragOverlayCard";
import { CardModal } from "./CardModal";
import { ColumnModal } from "./ColumnModal";
import { BoardModal } from "./BoardModal";
import { useProjectQuery } from "~features/projects/hooks/useProjects";
import type { BoardColumnView, CardSummaryView, CardPriority } from "~types/api";
import styles from "./KanbanWorkspace.module.css";

interface KanbanWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal?: () => void;
  projectColor?: string;
}

export const KanbanWorkspace: React.FC<KanbanWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
  projectColor,
}) => {
  const { data: project } = useProjectQuery(projectId);
  const activeProjectColor = projectColor || project?.color || "#10b981";

  // Board List & Selection
  const { data: boardsData, isLoading: isLoadingBoards } = useBoardsQuery(projectId);
  const { data: defaultBoard } = useDefaultBoardQuery(projectId);

  const boards = useMemo(() => boardsData?.content || [], [boardsData]);

  const [selectedBoardId, setSelectedBoardId] = useState<string>("");

  const activeBoard = useMemo(() => {
    if (selectedBoardId) {
      const found = boards.find((b) => b.id === selectedBoardId);
      if (found) return found;
    }
    return defaultBoard || boards[0];
  }, [boards, selectedBoardId, defaultBoard]);

  const activeBoardId = activeBoard?.id || "";

  // Query Columns & Cards for active board
  const { data: columns = [], isLoading: isLoadingColumns } = useBoardColumnsQuery(
    projectId,
    activeBoardId
  );
  const { data: cards = [], isLoading: isLoadingCards } = useBoardCardsQuery(
    projectId,
    activeBoardId
  );

  const moveCardMutation = useMoveCardMutation(projectId, activeBoardId);
  const reorderColumnsMutation = useReorderColumnsMutation(projectId, activeBoardId);

  // Filters State
  const [search, setSearch] = useState("");
  const [priorityFilter, setPriorityFilter] = useState<CardPriority | "ALL">("ALL");

  // Modals State
  const [isCardModalOpen, setIsCardModalOpen] = useState(false);
  const [targetColumnIdForCard, setTargetColumnIdForCard] = useState<string>("");
  const [editingCardId, setEditingCardId] = useState<string | undefined>(undefined);

  const [isColumnModalOpen, setIsColumnModalOpen] = useState(false);
  const [editingColumn, setEditingColumn] = useState<BoardColumnView | undefined>(undefined);

  const [isBoardModalOpen, setIsBoardModalOpen] = useState(false);
  const [editingBoard, setEditingBoard] = useState(false);

  // Drag and Drop State
  const [activeCard, setActiveCard] = useState<CardSummaryView | null>(null);
  const [activeColumn, setActiveColumn] = useState<BoardColumnView | null>(null);

  // Pointer sensor with activation constraint so clicks don't accidentally trigger drags
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 5,
      },
    }),
    useSensor(KeyboardSensor)
  );

  // Sorted Columns
  const sortedColumns = useMemo(
    () => [...columns].sort((a, b) => a.position - b.position),
    [columns]
  );

  const columnIds = useMemo(
    () => sortedColumns.map((c) => c.id),
    [sortedColumns]
  );

  // Filtered Cards
  const filteredCards = useMemo(() => {
    return cards.filter((card) => {
      if (priorityFilter !== "ALL" && card.priority !== priorityFilter) return false;
      if (search.trim()) {
        const q = search.toLowerCase();
        const matchesTitle = card.title.toLowerCase().includes(q);
        const matchesTag = card.tags?.some((t) => t.name.toLowerCase().includes(q));
        if (!matchesTitle && !matchesTag) return false;
      }
      return true;
    });
  }, [cards, priorityFilter, search]);

  // Group cards by column ID and sort by position
  const cardsByColumn = useMemo(() => {
    const map = new Map<string, CardSummaryView[]>();
    sortedColumns.forEach((col) => {
      map.set(col.id, []);
    });

    filteredCards.forEach((card) => {
      const list = map.get(card.columnId);
      if (list) {
        list.push(card);
      } else {
        map.set(card.columnId, [card]);
      }
    });

    map.forEach((list) => {
      list.sort((a, b) => a.position - b.position);
    });

    return map;
  }, [sortedColumns, filteredCards]);

  // ── Drag & Drop Handlers ─────────────────────────────────────

  const handleDragStart = (event: DragStartEvent) => {
    const { active } = event;
    const type = active.data.current?.type;

    if (type === "COLUMN") {
      const col = sortedColumns.find((c) => c.id === active.id);
      if (col) {
        setActiveColumn(col);
      }
    } else {
      const card = cards.find((c) => c.id === active.id);
      if (card) {
        setActiveCard(card);
      }
    }
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    const isColumnDrag = Boolean(
      activeColumn || active.data.current?.type === "COLUMN"
    );

    setActiveCard(null);
    setActiveColumn(null);

    if (!over || !activeBoardId) return;

    // Handle Column Reordering via Drag and Drop
    if (isColumnDrag) {
      if (active.id !== over.id) {
        const oldIndex = sortedColumns.findIndex((c) => c.id === active.id);
        const newIndex = sortedColumns.findIndex((c) => c.id === over.id);

        if (oldIndex !== -1 && newIndex !== -1) {
          const reordered = arrayMove(sortedColumns, oldIndex, newIndex);
          const positions = reordered.map((c) => c.id);
          try {
            await reorderColumnsMutation.mutateAsync({ positions });
            toast.success("Columns reordered");
          } catch (err) {
            toast.error(err instanceof Error ? err.message : "Failed to reorder columns");
          }
        }
      }
      return;
    }

    // Handle Card Drag and Drop
    const activeCardId = active.id as string;
    const overId = over.id as string;

    const sourceCard = cards.find((c) => c.id === activeCardId);
    if (!sourceCard) return;

    // Check if dropped over a column directly or over another card
    let targetColumnId = "";
    let targetPosition = 0;

    const isOverColumn = sortedColumns.some((col) => col.id === overId);

    if (isOverColumn) {
      targetColumnId = overId;
      const columnCards = cardsByColumn.get(targetColumnId) || [];
      targetPosition = columnCards.length;
    } else {
      // Over another card
      const targetCard = cards.find((c) => c.id === overId);
      if (targetCard) {
        targetColumnId = targetCard.columnId;
        targetPosition = targetCard.position;
      }
    }

    if (!targetColumnId) return;

    // If dropped in the same column at the same position, do nothing
    if (sourceCard.columnId === targetColumnId && sourceCard.position === targetPosition) {
      return;
    }

    try {
      await moveCardMutation.mutateAsync({
        cardId: activeCardId,
        payload: {
          targetColumnId,
          position: targetPosition,
        },
      });
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to move card");
    }
  };

  // ── Modal Trigger Handlers ───────────────────────────────────

  const handleOpenAddCard = (columnId: string) => {
    setTargetColumnIdForCard(columnId);
    setEditingCardId(undefined);
    setIsCardModalOpen(true);
  };

  const handleOpenCardDetail = (cardId: string) => {
    const card = cards.find((c) => c.id === cardId);
    setTargetColumnIdForCard(card?.columnId || sortedColumns[0]?.id || "");
    setEditingCardId(cardId);
    setIsCardModalOpen(true);
  };

  const handleOpenNewColumn = () => {
    setEditingColumn(undefined);
    setIsColumnModalOpen(true);
  };

  const handleOpenEditColumn = (col: BoardColumnView) => {
    setEditingColumn(col);
    setIsColumnModalOpen(true);
  };

  const handleOpenNewBoard = () => {
    setEditingBoard(false);
    setIsBoardModalOpen(true);
  };

  const handleOpenEditBoard = () => {
    setEditingBoard(true);
    setIsBoardModalOpen(true);
  };

  const isLoading = isLoadingBoards || isLoadingColumns || isLoadingCards;

  return (
    <div
      className={styles.workspaceRoot}
      style={{ "--color-primary": activeProjectColor } as React.CSSProperties}
    >
      {/* Top Toolbar */}
      <div className={styles.topBar}>
        <div className={styles.topBarLeft}>
          {/* Board Selector Group */}
          <div className="flex items-center gap-1.5">
            <div className={styles.boardSelectorWrapper}>
              <Icons.LayoutDashboard size={14} className={styles.boardSelectorIcon} />
              <select
                value={activeBoardId}
                onChange={(e) => setSelectedBoardId(e.target.value)}
                className={styles.boardSelector}
              >
                {boards.map((b) => (
                  <option key={b.id} value={b.id}>
                    {b.name}
                  </option>
                ))}
              </select>
              <Icons.ChevronDown size={12} className={styles.boardSelectorChevron} />
            </div>

            {activeBoard?.isDefault && (
              <span className={styles.defaultBadge}>Default</span>
            )}

            {activeBoard && (
              <button
                type="button"
                onClick={handleOpenEditBoard}
                className={styles.actionBtn}
                title="Edit active board settings"
              >
                <Icons.Settings2 size={13} />
              </button>
            )}

            <button
              type="button"
              onClick={handleOpenNewBoard}
              className={styles.actionBtn}
              title="Create new board"
            >
              <Icons.Plus size={13} />
              <span className="hidden sm:inline">New Board</span>
            </button>
          </div>

          <div className="w-[1px] h-4 bg-border/80 mx-0.5 hidden md:block" />

          {/* Search Input */}
          <div className={styles.searchBox}>
            <Icons.Search size={13} className={styles.searchIcon} />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search cards by title or tag..."
              className={styles.searchInput}
            />
            {search && (
              <button
                type="button"
                onClick={() => setSearch("")}
                className="absolute right-2 text-muted-foreground hover:text-foreground text-xs cursor-pointer p-0.5"
                title="Clear search"
              >
                <Icons.X size={12} />
              </button>
            )}
          </div>

          {/* Priority Segmented Filter Bar */}
          <div className={`hidden xl:flex ${styles.prioritySegmentedGroup}`}>
            {(
              [
                { id: "ALL", label: "All", dotColor: "bg-slate-400" },
                { id: "EXTREMELY_HIGH", label: "Urgent", dotColor: "bg-red-500" },
                { id: "HIGH", label: "High", dotColor: "bg-amber-500" },
                { id: "MEDIUM", label: "Medium", dotColor: "bg-blue-500" },
                { id: "LOW", label: "Low", dotColor: "bg-slate-400" },
              ] as const
            ).map((p) => {
              const isActive = priorityFilter === p.id;
              return (
                <button
                  key={p.id}
                  type="button"
                  onClick={() => setPriorityFilter(p.id)}
                  className={`${styles.priorityPill} ${isActive ? styles.priorityPillActive : ""}`}
                >
                  <span className={`w-1.5 h-1.5 rounded-full ${p.dotColor}`} />
                  <span>{p.label}</span>
                </button>
              );
            })}
          </div>
        </div>

        <div className={styles.topBarRight}>
          <button
            type="button"
            onClick={handleOpenNewColumn}
            className={styles.primaryBtn}
            title="Add a new column to this board"
          >
            <Icons.Plus size={14} />
            <span>Add Column</span>
          </button>
        </div>
      </div>

      {/* Main Board Canvas */}
      {isLoading && boards.length === 0 ? (
        <div className="flex-1 flex items-center justify-center text-muted-foreground gap-2">
          <Icons.Loader2 size={20} className="animate-spin text-primary" />
          <span className="text-xs font-mono">Loading Kanban board...</span>
        </div>
      ) : boards.length === 0 ? (
        <div className={styles.emptyState}>
          <div className="p-3 rounded-full bg-primary/10 text-primary border border-primary/20">
            <Icons.SquareKanban size={28} />
          </div>
          <h3 className="text-sm font-bold font-mono text-foreground">No Boards Found</h3>
          <p className="text-xs max-w-sm text-muted-foreground">
            Create your first Kanban board to start organizing tasks, code snippets, and issues.
          </p>
          <button
            type="button"
            onClick={handleOpenNewBoard}
            className={styles.primaryBtn}
          >
            <Icons.Plus size={14} />
            <span>Create Board</span>
          </button>
        </div>
      ) : (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <div className={styles.boardCanvas}>
            <SortableContext
              items={columnIds}
              strategy={horizontalListSortingStrategy}
            >
              {sortedColumns.map((column) => (
                <BoardColumn
                  key={column.id}
                  column={column}
                  cards={cardsByColumn.get(column.id) || []}
                  onOpenCard={handleOpenCardDetail}
                  onAddCardToColumn={handleOpenAddCard}
                  onEditColumn={handleOpenEditColumn}
                />
              ))}
            </SortableContext>

            {sortedColumns.length === 0 && (
              <div className="flex items-center justify-center w-full py-16 text-muted-foreground gap-3">
                <span className="text-xs font-mono">This board has no columns yet.</span>
                <button
                  type="button"
                  onClick={handleOpenNewColumn}
                  className={styles.actionBtn}
                >
                  <Icons.Plus size={13} />
                  <span>Add First Column</span>
                </button>
              </div>
            )}
          </div>

          {/* Drag Overlay for Smooth 60fps Drag Feedback */}
          <DragOverlay>
            {activeCard ? (
              <DragOverlayCard card={activeCard} />
            ) : activeColumn ? (
              <div className={`${styles.column} opacity-90 shadow-2xl border-primary scale-[1.02]`}>
                <div className={styles.columnHeader}>
                  <div className={styles.columnTitleBox}>
                    <div className="p-0.5 text-primary">
                      <Icons.GripVertical size={14} />
                    </div>
                    <h3 className={styles.columnTitle}>{activeColumn.name}</h3>
                  </div>
                </div>
                <div className="p-8 text-center text-xs text-muted-foreground font-mono">
                  {(cardsByColumn.get(activeColumn.id) || []).length} cards
                </div>
              </div>
            ) : null}
          </DragOverlay>
        </DndContext>
      )}

      {/* Modals */}
      {isCardModalOpen && (
        <CardModal
          isOpen={isCardModalOpen}
          onClose={() => setIsCardModalOpen(false)}
          projectId={projectId}
          boardId={activeBoardId}
          columnId={targetColumnIdForCard}
          cardId={editingCardId}
          onOpenManageTagsModal={onOpenManageTagsModal}
          projectColor={activeProjectColor}
        />
      )}

      {isColumnModalOpen && (
        <ColumnModal
          isOpen={isColumnModalOpen}
          onClose={() => setIsColumnModalOpen(false)}
          projectId={projectId}
          boardId={activeBoardId}
          column={editingColumn}
          projectColor={activeProjectColor}
        />
      )}

      {isBoardModalOpen && (
        <BoardModal
          isOpen={isBoardModalOpen}
          onClose={() => setIsBoardModalOpen(false)}
          projectId={projectId}
          board={editingBoard ? activeBoard : undefined}
          onSelectBoard={(newBoardId) => setSelectedBoardId(newBoardId)}
          projectColor={activeProjectColor}
        />
      )}
    </div>
  );
};
