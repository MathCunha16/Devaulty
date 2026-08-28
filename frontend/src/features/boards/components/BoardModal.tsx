import React, { useState, useEffect, useRef } from "react";
import * as Icons from "lucide-react";
import { toast } from "sonner";
import { ConfirmModal } from "../../../components/ConfirmModal";
import {
  useCreateBoardMutation,
  useUpdateBoardMutation,
  useDeleteBoardMutation,
} from "../hooks/useBoards";
import type { BoardView, CreateBoardRequest, UpdateBoardRequest } from "~types/api";

interface BoardModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  board?: BoardView; // If provided, edit mode
  onSelectBoard?: (boardId: string) => void;
  projectColor?: string;
}

interface BoardModalInnerProps {
  projectId: string;
  board?: BoardView;
  onClose: () => void;
  onSelectBoard?: (boardId: string) => void;
}

const BoardModalInner: React.FC<BoardModalInnerProps> = ({
  projectId,
  board,
  onClose,
  onSelectBoard,
}) => {
  const isEditing = Boolean(board);
  const createMutation = useCreateBoardMutation(projectId);
  const updateMutation = useUpdateBoardMutation(projectId);
  const deleteMutation = useDeleteBoardMutation(projectId);

  const [name, setName] = useState(board?.name || "");
  const [description, setDescription] = useState(board?.description || "");
  const [isDefault, setIsDefault] = useState(board?.isDefault || false);
  const [isDeleting, setIsDeleting] = useState(false);
  const [isConfirmDeleteOpen, setIsConfirmDeleteOpen] = useState(false);

  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    const timer = setTimeout(() => {
      inputRef.current?.focus();
    }, 50);
    return () => clearTimeout(timer);
  }, []);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        if (isConfirmDeleteOpen) {
          setIsConfirmDeleteOpen(false);
        } else {
          onClose();
        }
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isConfirmDeleteOpen, onClose]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) {
      toast.error("Board name is required");
      return;
    }

    try {
      if (isEditing && board) {
        const payload: UpdateBoardRequest = {
          name: name.trim(),
          description: description.trim() ? description.trim() : "",
          isDefault,
        };
        await updateMutation.mutateAsync({ boardId: board.id, payload });
        toast.success(`Board "${name}" updated`);
      } else {
        const payload: CreateBoardRequest = {
          name: name.trim(),
          description: description.trim() ? description.trim() : undefined,
          isDefault,
        };
        const created = await createMutation.mutateAsync(payload);
        toast.success(`Board "${name}" created`);
        onSelectBoard?.(created.id);
      }
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save board");
    }
  };

  const handleExecuteDelete = async () => {
    if (!board) return;
    try {
      setIsDeleting(true);
      await deleteMutation.mutateAsync(board.id);
      toast.success(`Board "${board.name}" deleted`);
      setIsConfirmDeleteOpen(false);
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete board");
    } finally {
      setIsDeleting(false);
    }
  };

  const isSubmitting = createMutation.isPending || updateMutation.isPending || isDeleting;

  return (
    <form onSubmit={handleSubmit} className="p-5 flex flex-col gap-4">
      {/* Name */}
      <div className="flex flex-col gap-1.5">
        <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
          Board Name <span className="text-destructive">*</span>
        </label>
        <input
          ref={inputRef}
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. Sprint 14, Architecture Backlog, Release Prep"
          required
          className="w-full px-3 py-2 bg-background border border-border rounded-md text-sm font-medium text-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all"
        />
      </div>

      {/* Description */}
      <div className="flex flex-col gap-1.5">
        <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
          Description (Optional)
        </label>
        <textarea
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          placeholder="Brief purpose of this board..."
          rows={3}
          className="w-full p-2.5 bg-background border border-border rounded-md text-xs font-mono text-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary resize-none"
        />
      </div>

      {/* isDefault Toggle */}
      <div className="flex items-center justify-between p-3 rounded-lg bg-secondary/30 border border-border/70">
        <div className="flex flex-col">
          <span className="text-xs font-semibold font-mono text-foreground">
            Set as Default Board
          </span>
          <span className="text-[11px] text-muted-foreground">
            Opens automatically when navigating to this project
          </span>
        </div>
        <input
          type="checkbox"
          checked={isDefault}
          onChange={(e) => setIsDefault(e.target.checked)}
          className="w-4 h-4 rounded border-border text-primary focus:ring-primary cursor-pointer"
        />
      </div>

      {/* Footer Actions */}
      <div className="flex items-center justify-between pt-3 border-t border-border mt-1">
        {isEditing && !board?.isDefault ? (
          <button
            type="button"
            onClick={() => setIsConfirmDeleteOpen(true)}
            disabled={isSubmitting}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-semibold font-mono text-destructive hover:bg-destructive/10 transition-colors cursor-pointer border border-destructive/20"
          >
            <Icons.Trash2 size={13} />
            <span>Delete</span>
          </button>
        ) : (
          <div />
        )}

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={onClose}
            disabled={isSubmitting}
            className="px-3.5 py-1.5 rounded-md text-xs font-semibold font-mono text-muted-foreground hover:bg-secondary transition-colors cursor-pointer border border-border"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            className="flex items-center gap-1.5 px-4 py-1.5 rounded-md text-xs font-bold font-mono bg-primary text-primary-foreground hover:brightness-110 shadow-md transition-all cursor-pointer disabled:opacity-50"
          >
            {isSubmitting && <Icons.Loader2 size={13} className="animate-spin" />}
            <span>{isEditing ? "Save" : "Create"}</span>
          </button>
        </div>
      </div>

      {/* Delete Confirmation Modal */}
      {isEditing && board && (
        <ConfirmModal
          isOpen={isConfirmDeleteOpen}
          onClose={() => setIsConfirmDeleteOpen(false)}
          onConfirm={handleExecuteDelete}
          title="Delete Board"
          message={`Are you sure you want to delete board "${board.name}"? All columns and cards inside will be permanently deleted.`}
          warningText="This action cannot be undone."
          confirmLabel="Delete Board"
          isLoading={isDeleting}
        />
      )}
    </form>
  );
};

export const BoardModal: React.FC<BoardModalProps> = ({
  isOpen,
  onClose,
  projectId,
  board,
  onSelectBoard,
  projectColor,
}) => {
  const isEditing = Boolean(board);

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
      style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      onClick={onClose}
    >
      <div
        className="relative w-full max-w-md flex flex-col rounded-xl bg-card border border-border shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-border/80 bg-secondary/30">
          <div className="flex items-center gap-2">
            <span className="p-1.5 rounded-lg bg-primary/10 text-primary border border-primary/20">
              <Icons.LayoutDashboard size={16} />
            </span>
            <h2 className="text-sm font-bold font-mono text-foreground uppercase tracking-wide">
              {isEditing ? "Edit Board" : "New Board"}
            </h2>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="p-1.5 rounded-md hover:bg-secondary text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
          >
            <Icons.X size={16} />
          </button>
        </div>

        {/* Body */}
        <BoardModalInner
          key={board?.id || "new-board"}
          projectId={projectId}
          board={board}
          onClose={onClose}
          onSelectBoard={onSelectBoard}
        />
      </div>
    </div>
  );
};
