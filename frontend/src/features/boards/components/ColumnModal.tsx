import React, { useState, useEffect, useRef } from "react";
import * as Icons from "lucide-react";
import { toast } from "sonner";
import { ConfirmModal } from "../../../components/ConfirmModal";
import {
  useCreateColumnMutation,
  useUpdateColumnMutation,
  useDeleteColumnMutation,
} from "../hooks/useBoards";
import type { BoardColumnView, CreateBoardColumnRequest, UpdateBoardColumnRequest } from "~types/api";

interface ColumnModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  boardId: string;
  column?: BoardColumnView; // If provided, edit mode
  projectColor?: string;
}

interface ColumnModalInnerProps {
  projectId: string;
  boardId: string;
  column?: BoardColumnView;
  onClose: () => void;
}

const ColumnModalInner: React.FC<ColumnModalInnerProps> = ({
  projectId,
  boardId,
  column,
  onClose,
}) => {
  const isEditing = Boolean(column);
  const createMutation = useCreateColumnMutation(projectId, boardId);
  const updateMutation = useUpdateColumnMutation(projectId, boardId);
  const deleteMutation = useDeleteColumnMutation(projectId, boardId);

  const [name, setName] = useState(column?.name || "");
  const [hasWipLimit, setHasWipLimit] = useState(
    Boolean(column?.wipLimit != null && column.wipLimit > 0)
  );
  const [wipLimit, setWipLimit] = useState<number>(column?.wipLimit || 3);
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
      toast.error("Column name is required");
      return;
    }

    try {
      if (isEditing && column) {
        const payload: UpdateBoardColumnRequest = {
          name: name.trim(),
          wipLimit: hasWipLimit && wipLimit > 0 ? wipLimit : 0,
        };
        await updateMutation.mutateAsync({ columnId: column.id, payload });
        toast.success(`Column "${name}" updated`);
      } else {
        const payload: CreateBoardColumnRequest = {
          name: name.trim(),
          wipLimit: hasWipLimit && wipLimit > 0 ? wipLimit : undefined,
        };
        await createMutation.mutateAsync(payload);
        toast.success(`Column "${name}" created`);
      }
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save column");
    }
  };

  const handleExecuteDelete = async () => {
    if (!column) return;
    try {
      setIsDeleting(true);
      await deleteMutation.mutateAsync(column.id);
      toast.success(`Column "${column.name}" deleted`);
      setIsConfirmDeleteOpen(false);
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete column");
    } finally {
      setIsDeleting(false);
    }
  };

  const isSubmitting = createMutation.isPending || updateMutation.isPending || isDeleting;

  return (
    <form onSubmit={handleSubmit} className="p-5 flex flex-col gap-4">
      {/* Column Name */}
      <div className="flex flex-col gap-1.5">
        <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
          Column Name <span className="text-destructive">*</span>
        </label>
        <input
          ref={inputRef}
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="e.g. In Progress, Testing, Deployed"
          required
          className="w-full px-3 py-2 bg-background border border-border rounded-md text-sm font-medium text-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all"
        />
      </div>

      {/* WIP Limit Toggle */}
      <div className="flex flex-col gap-2 p-3 rounded-lg bg-secondary/40 border border-border/80">
        <div className="flex items-center justify-between">
          <div className="flex flex-col">
            <span className="text-xs font-semibold text-foreground font-mono">WIP Limit</span>
            <span className="text-[11px] text-muted-foreground">
              Maximum active cards allowed in this column
            </span>
          </div>
          <button
            type="button"
            role="switch"
            aria-checked={hasWipLimit}
            onClick={() => setHasWipLimit(!hasWipLimit)}
            className={`w-9 h-5 rounded-full p-0.5 transition-colors cursor-pointer border ${
              hasWipLimit
                ? "bg-primary border-primary justify-end"
                : "bg-muted border-border justify-start"
            } flex items-center`}
          >
            <div className="w-3.5 h-3.5 rounded-full bg-white shadow-sm" />
          </button>
        </div>

        {hasWipLimit && (
          <div className="flex items-center gap-2 pt-2 border-t border-border/60">
            <label className="text-xs font-mono text-muted-foreground">Max cards:</label>
            <input
              type="number"
              min={1}
              max={99}
              value={wipLimit}
              onChange={(e) => setWipLimit(Math.max(1, parseInt(e.target.value, 10) || 1))}
              className="w-20 px-2.5 py-1 bg-background border border-border rounded text-xs font-mono text-foreground outline-none focus:border-primary"
            />
          </div>
        )}
      </div>

      {/* Footer Actions */}
      <div className="flex items-center justify-between pt-3 border-t border-border mt-1">
        {isEditing ? (
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
      {isEditing && column && (
        <ConfirmModal
          isOpen={isConfirmDeleteOpen}
          onClose={() => setIsConfirmDeleteOpen(false)}
          onConfirm={handleExecuteDelete}
          title="Delete Column"
          message={`Are you sure you want to delete column "${column.name}"? All cards inside this column will also be deleted.`}
          warningText="This action cannot be undone."
          confirmLabel="Delete Column"
          isLoading={isDeleting}
        />
      )}
    </form>
  );
};

export const ColumnModal: React.FC<ColumnModalProps> = ({
  isOpen,
  onClose,
  projectId,
  boardId,
  column,
  projectColor,
}) => {
  const isEditing = Boolean(column);

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
              <Icons.Columns3 size={16} />
            </span>
            <h2 className="text-sm font-bold font-mono text-foreground uppercase tracking-wide">
              {isEditing ? "Edit Column" : "New Column"}
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
        <ColumnModalInner
          key={column?.id || "new-column"}
          projectId={projectId}
          boardId={boardId}
          column={column}
          onClose={onClose}
        />
      </div>
    </div>
  );
};
