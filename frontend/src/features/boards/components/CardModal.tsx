import React, { useState, useEffect, useRef, useMemo, useCallback } from "react";
import * as Icons from "lucide-react";
import { toast } from "sonner";
import { ConfirmModal } from "../../../components/ConfirmModal";
import {
  useCardDetailQuery,
  useCreateCardMutation,
  useUpdateCardMutation,
  useDeleteCardMutation,
} from "../hooks/useBoards";
import {
  useTagsQuery,
  useAssociateTagMutation,
  useDisassociateTagMutation,
} from "~features/tags/hooks/useTags";
import { LinkedItemPicker } from "./LinkedItemPicker";
import { MarkdownWithMentions } from "./MarkdownWithMentions";
import { MentionAutocomplete } from "./MentionAutocomplete";
import { formatMentionToken, type ResolvedMentionItem } from "../utils/boardUtils";
import { getCaretCoordinates } from "../utils/caretCoordinates";
import { useSnippetsQuery } from "~features/snippets/hooks/useSnippets";
import { useProblemsQuery } from "~features/problems/hooks/useProblems";
import { useCredentialsQuery } from "~features/credentials/hooks/useCredentials";
import { useNotesQuery } from "~features/notes/hooks/useNotes";
import { useLinksQuery } from "~features/links/hooks/useLinks";
import type {
  CardPriority,
  CreateCardRequest,
  UpdateCardRequest,
  CreateCardItemCommand,
  TagSummaryResponse,
  CardView,
} from "~types/api";

interface CardModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  boardId: string;
  columnId: string;
  cardId?: string;
  onOpenManageTagsModal?: () => void;
  projectColor?: string;
}

interface CardModalInnerProps {
  projectId: string;
  boardId: string;
  columnId: string;
  cardId?: string;
  cardDetail?: CardView;
  allTags: TagSummaryResponse[];
  onClose: () => void;
  onOpenManageTagsModal?: () => void;
  projectColor?: string;
}

interface MentionState {
  isOpen: boolean;
  query: string;
  caretIndex: number;
  coords: { top: number; left: number };
  selectedIndex: number;
}

const CardModalInner: React.FC<CardModalInnerProps> = ({
  projectId,
  boardId,
  columnId,
  cardId,
  cardDetail,
  allTags,
  onClose,
  onOpenManageTagsModal,
  projectColor,
}) => {
  const isEditing = Boolean(cardId && cardDetail);

  const createMutation = useCreateCardMutation(projectId, boardId, columnId);
  const updateMutation = useUpdateCardMutation(projectId, boardId);
  const deleteMutation = useDeleteCardMutation(projectId, boardId);

  const associateTagMutation = useAssociateTagMutation(projectId);
  const disassociateTagMutation = useDisassociateTagMutation(projectId);

  // Queries to resolve titles of linked items
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const { data: problemsData } = useProblemsQuery(projectId);
  const { data: credentialsData } = useCredentialsQuery(projectId);
  const { data: notesData } = useNotesQuery(projectId);
  const { data: linksData } = useLinksQuery(projectId);

  const allProjectItems = useMemo(() => {
    const list: ResolvedMentionItem[] = [];
    (snippetsData?.content || []).forEach((s) => list.push({ id: s.id, title: s.title, type: "SNIPPET" }));
    (problemsData?.content || []).forEach((p) => list.push({ id: p.id, title: p.title, type: "PROBLEM" }));
    (credentialsData?.content || []).forEach((c) => list.push({ id: c.id, title: c.title, type: "CREDENTIAL" }));
    (notesData?.content || []).forEach((n) => list.push({ id: n.id, title: n.title, type: "NOTE" }));
    (linksData?.content || []).forEach((l) => list.push({ id: l.id, title: l.title, type: "LINK" }));
    return list;
  }, [snippetsData, problemsData, credentialsData, notesData, linksData]);

  // Form State
  const [title, setTitle] = useState(cardDetail?.title || "");
  const [description, setDescription] = useState(cardDetail?.description || "");
  const [descTab, setDescTab] = useState<"write" | "preview">(
    cardDetail?.description ? "preview" : "write"
  );
  const [priority, setPriority] = useState<CardPriority | "">(
    cardDetail?.priority || ""
  );
  const [dueDate, setDueDate] = useState<string>(
    cardDetail?.dueDate ? cardDetail.dueDate.split("T")[0] : ""
  );
  const [linkedItems, setLinkedItems] = useState<CreateCardItemCommand[]>(
    (cardDetail?.linkedItems || []).map((li) => ({
      itemType: li.itemType,
      itemId: li.itemId,
    }))
  );
  const [cardTags, setCardTags] = useState<TagSummaryResponse[]>(
    cardDetail?.tags || []
  );

  const [isDeleting, setIsDeleting] = useState(false);
  const [isConfirmDeleteOpen, setIsConfirmDeleteOpen] = useState(false);
  const [isConfirmDiscardOpen, setIsConfirmDiscardOpen] = useState(false);

  // Mention Autocomplete State
  const [mentionState, setMentionState] = useState<MentionState>({
    isOpen: false,
    query: "",
    caretIndex: 0,
    coords: { top: 0, left: 0 },
    selectedIndex: 0,
  });

  const titleInputRef = useRef<HTMLInputElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    const timer = setTimeout(() => {
      titleInputRef.current?.focus();
    }, 50);
    return () => clearTimeout(timer);
  }, []);

  // Compute resolved items attached to this card ONLY
  const resolvedLinkedItems: ResolvedMentionItem[] = useMemo(() => {
    const map = new Map<string, ResolvedMentionItem>();
    allProjectItems.forEach((it) => map.set(`${it.type}-${it.id}`, it));

    const list: ResolvedMentionItem[] = [];
    linkedItems.forEach((li) => {
      const found = map.get(`${li.itemType}-${li.itemId}`);
      if (found) {
        list.push(found);
      }
    });

    return list;
  }, [allProjectItems, linkedItems]);

  // Filter items matching @query strictly from attached items
  const filteredMentionItems = useMemo(() => {
    if (!mentionState.isOpen || resolvedLinkedItems.length === 0) return [];
    if (!mentionState.query.trim()) return resolvedLinkedItems;
    return resolvedLinkedItems.filter((item) =>
      item.title.toLowerCase().includes(mentionState.query.toLowerCase())
    );
  }, [mentionState.isOpen, mentionState.query, resolvedLinkedItems]);

  // Dirty checking for unsaved changes protection
  const initialTitle = cardDetail?.title || "";
  const initialDescription = cardDetail?.description || "";
  const initialPriority = cardDetail?.priority || "";
  const initialDueDate = cardDetail?.dueDate ? cardDetail.dueDate.split("T")[0] : "";
  const initialLinked = useMemo(
    () => (cardDetail?.linkedItems || []).map((li) => `${li.itemType}-${li.itemId}`).sort().join(","),
    [cardDetail]
  );
  const currentLinked = useMemo(
    () => linkedItems.map((li) => `${li.itemType}-${li.itemId}`).sort().join(","),
    [linkedItems]
  );
  const initialTagIds = useMemo(
    () => (cardDetail?.tags || []).map((t) => t.id).sort().join(","),
    [cardDetail]
  );
  const currentTagIds = useMemo(
    () => cardTags.map((t) => t.id).sort().join(","),
    [cardTags]
  );

  const isDirty = useMemo(() => {
    return (
      title.trim() !== initialTitle.trim() ||
      description !== initialDescription ||
      priority !== initialPriority ||
      dueDate !== initialDueDate ||
      initialLinked !== currentLinked ||
      initialTagIds !== currentTagIds
    );
  }, [title, description, priority, dueDate, initialTitle, initialDescription, initialPriority, initialDueDate, initialLinked, currentLinked, initialTagIds, currentTagIds]);

  // Request close: if dirty, prompt confirmation; otherwise close immediately
  const handleRequestClose = useCallback(() => {
    if (isDirty) {
      setIsConfirmDiscardOpen(true);
    } else {
      onClose();
    }
  }, [isDirty, onClose]);

  // Escape key handler
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        if (mentionState.isOpen) {
          setMentionState((prev) => ({ ...prev, isOpen: false }));
        } else if (isConfirmDeleteOpen) {
          setIsConfirmDeleteOpen(false);
        } else if (isConfirmDiscardOpen) {
          setIsConfirmDiscardOpen(false);
        } else {
          handleRequestClose();
        }
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handleRequestClose, mentionState.isOpen, isConfirmDeleteOpen, isConfirmDiscardOpen]);

  // Check for @ trigger on typing in textarea
  const handleTextareaInput = useCallback(() => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const sel = textarea.selectionStart;
    const val = textarea.value;
    const textBeforeCaret = val.substring(0, sel);

    // Look for @ followed by optional characters right before cursor
    const match = textBeforeCaret.match(/@([a-zA-Z0-9_\-\s]{0,25})$/);

    if (match) {
      const query = match[1];
      const atPos = sel - query.length - 1;
      const coords = getCaretCoordinates(textarea, atPos);

      setMentionState({
        isOpen: true,
        query,
        caretIndex: atPos,
        coords: { top: coords.top, left: coords.left },
        selectedIndex: 0,
      });
    } else {
      if (mentionState.isOpen) {
        setMentionState((prev) => ({ ...prev, isOpen: false }));
      }
    }
  }, [mentionState.isOpen]);

  const handleSelectMention = useCallback(
    (item: ResolvedMentionItem) => {
      const textarea = textareaRef.current;
      if (!textarea) return;

      const val = description;
      const atPos = mentionState.caretIndex;
      const currentCaret = textarea.selectionStart;

      const before = val.substring(0, atPos);
      const after = val.substring(currentCaret);
      const mentionToken = formatMentionToken(item.title);
      const newText = before + mentionToken + " " + after;

      setDescription(newText);
      setMentionState((prev) => ({ ...prev, isOpen: false }));

      // Refocus textarea and place caret after inserted token
      setTimeout(() => {
        textarea.focus();
        const nextPos = (before + mentionToken + " ").length;
        textarea.setSelectionRange(nextPos, nextPos);
      }, 30);
    },
    [description, mentionState.caretIndex]
  );

  const handleTextareaKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (!mentionState.isOpen || filteredMentionItems.length === 0) return;

    if (e.key === "ArrowDown") {
      e.preventDefault();
      setMentionState((prev) => ({
        ...prev,
        selectedIndex: (prev.selectedIndex + 1) % filteredMentionItems.length,
      }));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setMentionState((prev) => ({
        ...prev,
        selectedIndex:
          (prev.selectedIndex - 1 + filteredMentionItems.length) %
          filteredMentionItems.length,
      }));
    } else if (e.key === "Enter" || e.key === "Tab") {
      e.preventDefault();
      const selected = filteredMentionItems[mentionState.selectedIndex];
      if (selected) {
        handleSelectMention(selected);
      }
    } else if (e.key === "Escape") {
      e.preventDefault();
      setMentionState((prev) => ({ ...prev, isOpen: false }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) {
      toast.error("Card title is required");
      return;
    }

    try {
      const payloadPriority = priority ? (priority as CardPriority) : undefined;
      const payloadDueDate = dueDate ? new Date(dueDate).toISOString() : undefined;

      if (isEditing && cardId) {
        const updatePayload: UpdateCardRequest = {
          title: title.trim(),
          description: description.trim() ? description : "",
          priority: payloadPriority,
          dueDate: payloadDueDate,
          linkedItems,
        };
        await updateMutation.mutateAsync({ cardId, payload: updatePayload });
        toast.success("Card updated successfully");
      } else {
        const createPayload: CreateCardRequest = {
          title: title.trim(),
          description: description.trim() ? description : undefined,
          priority: payloadPriority,
          dueDate: payloadDueDate,
          linkedItems,
        };
        const createdCard = await createMutation.mutateAsync(createPayload);
        if (cardTags.length > 0) {
          for (const tag of cardTags) {
            try {
              await associateTagMutation.mutateAsync({
                itemType: "card",
                itemId: createdCard.id,
                tagId: tag.id,
              });
            } catch {
              // Non-blocking tag association error
            }
          }
        }
        toast.success("Card created successfully");
      }
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to save card");
    }
  };

  const handleExecuteDelete = async () => {
    if (!cardId) return;
    try {
      setIsDeleting(true);
      await deleteMutation.mutateAsync(cardId);
      toast.success("Card deleted");
      setIsConfirmDeleteOpen(false);
      onClose();
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to delete card");
    } finally {
      setIsDeleting(false);
    }
  };

  const handleToggleTag = async (tag: TagSummaryResponse) => {
    if (!cardId) {
      const exists = cardTags.some((t) => t.id === tag.id);
      if (exists) {
        setCardTags(cardTags.filter((t) => t.id !== tag.id));
      } else {
        setCardTags([...cardTags, tag]);
      }
      return;
    }

    const isAssigned = cardTags.some((t) => t.id === tag.id);
    try {
      if (isAssigned) {
        await disassociateTagMutation.mutateAsync({
          itemType: "card",
          itemId: cardId,
          tagId: tag.id,
        });
        setCardTags((prev) => prev.filter((t) => t.id !== tag.id));
        toast.success(`Removed tag "${tag.name}"`);
      } else {
        await associateTagMutation.mutateAsync({
          itemType: "card",
          itemId: cardId,
          tagId: tag.id,
        });
        setCardTags((prev) => [...prev, tag]);
        toast.success(`Added tag "${tag.name}"`);
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to update tag");
    }
  };

  const isSubmitting = createMutation.isPending || updateMutation.isPending || isDeleting;

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200"
      style={{ "--color-primary": projectColor || "#10b981" } as React.CSSProperties}
      onClick={handleRequestClose}
    >
      <div
        className="relative w-full max-w-2xl max-h-[90vh] flex flex-col rounded-xl bg-card border border-border shadow-2xl overflow-hidden animate-in zoom-in-95 duration-150"
        onClick={(e) => e.stopPropagation()}
      >
        {/* Modal Header */}
        <div className="flex items-center justify-between p-4 border-b border-border/80 bg-secondary/30">
          <div className="flex items-center gap-2">
            <span className="p-1.5 rounded-lg bg-primary/10 text-primary border border-primary/20">
              <Icons.SquareKanban size={16} />
            </span>
            <h2 className="text-sm font-bold font-mono text-foreground uppercase tracking-wide">
              {isEditing ? "Card Details" : "Create Card"}
            </h2>
            {isDirty && (
              <span className="text-[10px] font-mono px-2 py-0.5 rounded-full bg-amber-500/15 text-amber-500 border border-amber-500/30">
                Unsaved Changes
              </span>
            )}
          </div>
          <button
            type="button"
            onClick={handleRequestClose}
            className="p-1.5 rounded-md hover:bg-secondary text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
            title="Close (Esc)"
          >
            <Icons.X size={16} />
          </button>
        </div>

        {/* Modal Form Body */}
        <form onSubmit={handleSubmit} className="flex-1 overflow-y-auto p-5 flex flex-col gap-4">
          {/* Title */}
          <div className="flex flex-col gap-1.5">
            <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
              Title <span className="text-destructive">*</span>
            </label>
            <input
              ref={titleInputRef}
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="What needs to be done?"
              required
              className="w-full px-3 py-2 bg-background border border-border rounded-md text-sm font-medium text-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary transition-all"
            />
          </div>

          {/* Markdown Description with Write/Preview Tabs */}
          <div className="flex flex-col gap-1.5">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
                  Description (Markdown)
                </label>
                <span className="text-[10px] font-mono text-muted-foreground bg-secondary/60 px-1.5 py-0.5 rounded border border-border/50">
                  Type <strong className="text-primary">@</strong> to mention attached items
                </span>
              </div>

              <div className="flex items-center gap-1 p-0.5 rounded-md bg-secondary/80 border border-border/80">
                <button
                  type="button"
                  onClick={() => setDescTab("write")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    descTab === "write"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Icons.Code2 size={12} />
                  <span>Write / Code</span>
                </button>
                <button
                  type="button"
                  onClick={() => setDescTab("preview")}
                  className={`flex items-center gap-1 px-2.5 py-1 rounded text-xs font-mono transition-colors cursor-pointer ${
                    descTab === "preview"
                      ? "bg-card text-foreground font-semibold shadow-sm border border-border/60"
                      : "text-muted-foreground hover:text-foreground"
                  }`}
                >
                  <Icons.Eye size={12} />
                  <span>Preview</span>
                </button>
              </div>
            </div>

            {descTab === "write" ? (
              <div className="relative">
                <textarea
                  ref={textareaRef}
                  value={description}
                  onChange={(e) => {
                    setDescription(e.target.value);
                    handleTextareaInput();
                  }}
                  onKeyUp={handleTextareaInput}
                  onClick={handleTextareaInput}
                  onKeyDown={handleTextareaKeyDown}
                  placeholder="Describe requirements, architectural notes, or markdown checklists...&#10;&#10;Type @ to link attached Devaulty items directly in text."
                  rows={6}
                  className="w-full p-3 bg-background border border-border rounded-md text-xs font-mono text-foreground placeholder:text-muted-foreground outline-none focus:border-primary focus:ring-1 focus:ring-primary resize-y min-h-[150px] leading-relaxed"
                />

                {/* In-Textarea Mention Autocomplete Popover (strictly attached items) */}
                <MentionAutocomplete
                  isOpen={mentionState.isOpen}
                  position={mentionState.coords}
                  items={filteredMentionItems}
                  hasAttachedItems={resolvedLinkedItems.length > 0}
                  selectedIndex={mentionState.selectedIndex}
                  query={mentionState.query}
                  onSelect={handleSelectMention}
                  onClose={() => setMentionState((prev) => ({ ...prev, isOpen: false }))}
                />
              </div>
            ) : (
              <div className="w-full p-4 bg-secondary/30 border border-border/70 rounded-md min-h-[150px] max-h-80 overflow-y-auto">
                <MarkdownWithMentions
                  text={description}
                  projectId={projectId}
                  linkedItems={linkedItems}
                />
              </div>
            )}
          </div>

          {/* Priority & Due Date Row */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Priority */}
            <div className="flex flex-col gap-1.5">
              <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
                Priority
              </label>
              <div className="flex items-center gap-1.5 flex-wrap">
                {(
                  [
                    { id: "LOW", label: "Low", color: "text-slate-400 border-slate-500/40 bg-slate-500/10" },
                    { id: "MEDIUM", label: "Medium", color: "text-blue-500 border-blue-500/40 bg-blue-500/10" },
                    { id: "HIGH", label: "High", color: "text-amber-500 border-amber-500/40 bg-amber-500/10" },
                    { id: "EXTREMELY_HIGH", label: "Urgent", color: "text-red-500 border-red-500/40 bg-red-500/10 font-bold" },
                  ] as const
                ).map((p) => {
                  const isSelected = priority === p.id;
                  return (
                    <button
                      key={p.id}
                      type="button"
                      onClick={() => setPriority(isSelected ? "" : p.id)}
                      className={`px-2.5 py-1 rounded text-xs font-mono border transition-all cursor-pointer ${
                        isSelected
                          ? `${p.color} ring-1 ring-primary/50 shadow-sm font-semibold`
                          : "border-border bg-secondary/40 text-muted-foreground hover:text-foreground"
                      }`}
                    >
                      {p.label}
                    </button>
                  );
                })}
              </div>
            </div>

            {/* Due Date */}
            <div className="flex flex-col gap-1.5">
              <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
                Due Date
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  className="w-full px-3 py-1.5 bg-background border border-border rounded-md text-xs font-mono text-foreground outline-none focus:border-primary"
                />
                {dueDate && (
                  <button
                    type="button"
                    onClick={() => setDueDate("")}
                    className="p-1.5 text-muted-foreground hover:text-destructive text-xs transition-colors cursor-pointer"
                    title="Clear due date"
                  >
                    <Icons.X size={14} />
                  </button>
                )}
              </div>
            </div>
          </div>

          {/* Tags Selection */}
          <div className="flex flex-col gap-1.5">
            <div className="flex items-center justify-between">
              <label className="text-xs font-semibold text-muted-foreground uppercase font-mono tracking-wider">
                Tags
              </label>
              {onOpenManageTagsModal && (
                <button
                  type="button"
                  onClick={onOpenManageTagsModal}
                  className="text-xs text-primary hover:underline font-mono cursor-pointer"
                >
                  + Manage Tags
                </button>
              )}
            </div>
            <div className="flex flex-wrap gap-1.5">
              {allTags.length === 0 ? (
                <span className="text-xs text-muted-foreground italic">No tags created in project yet.</span>
              ) : (
                allTags.map((tag) => {
                  const isSelected = cardTags.some((t) => t.id === tag.id);
                  return (
                    <button
                      key={tag.id}
                      type="button"
                      onClick={() => handleToggleTag(tag)}
                      className={`inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-xs font-mono border transition-all cursor-pointer ${
                        isSelected
                          ? "bg-primary/15 border-primary/40 text-primary font-semibold shadow-sm"
                          : "bg-secondary/40 border-border text-muted-foreground hover:text-foreground"
                      }`}
                      style={{
                        borderColor: isSelected && tag.color ? tag.color : undefined,
                        color: isSelected && tag.color ? tag.color : undefined,
                      }}
                    >
                      <span>{tag.name}</span>
                      {isSelected && <Icons.Check size={11} />}
                    </button>
                  );
                })
              )}
            </div>
          </div>

          {/* Single Unified Linked Devaulty Items Section */}
          <LinkedItemPicker
            projectId={projectId}
            linkedItems={linkedItems}
            onChange={setLinkedItems}
          />

          {/* Actions Bottom Bar */}
          <div className="flex items-center justify-between pt-4 border-t border-border mt-2">
            {isEditing ? (
              <button
                type="button"
                onClick={() => setIsConfirmDeleteOpen(true)}
                disabled={isSubmitting}
                className="flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-semibold font-mono text-destructive hover:bg-destructive/10 transition-colors cursor-pointer border border-destructive/20"
              >
                <Icons.Trash2 size={13} />
                <span>Delete Card</span>
              </button>
            ) : (
              <div />
            )}

            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={handleRequestClose}
                disabled={isSubmitting}
                className="px-4 py-1.5 rounded-md text-xs font-semibold font-mono text-muted-foreground hover:bg-secondary transition-colors cursor-pointer border border-border"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="flex items-center gap-1.5 px-4 py-1.5 rounded-md text-xs font-bold font-mono bg-primary text-primary-foreground hover:brightness-110 shadow-md transition-all cursor-pointer disabled:opacity-50"
              >
                {isSubmitting && <Icons.Loader2 size={13} className="animate-spin" />}
                <span>{isEditing ? "Save Changes" : "Create Card"}</span>
              </button>
            </div>
          </div>
        </form>

        {/* Delete Card Confirmation Modal */}
        {isEditing && (
          <ConfirmModal
            isOpen={isConfirmDeleteOpen}
            onClose={() => setIsConfirmDeleteOpen(false)}
            onConfirm={handleExecuteDelete}
            title="Delete Card"
            message={`Are you sure you want to delete card "${title || "Untitled"}"?`}
            warningText="This action cannot be undone."
            confirmLabel="Delete Card"
            isLoading={isDeleting}
          />
        )}

        {/* Discard Changes Confirmation Modal */}
        <ConfirmModal
          isOpen={isConfirmDiscardOpen}
          onClose={() => setIsConfirmDiscardOpen(false)}
          onConfirm={() => {
            setIsConfirmDiscardOpen(false);
            onClose();
          }}
          title="Discard Unsaved Changes"
          message="You have unsaved modifications on this card. Are you sure you want to discard your changes and close?"
          warningText="Any unsaved text, attachments, or tag updates will be lost."
          confirmLabel="Discard Changes"
        />
      </div>
    </div>
  );
};

export const CardModal: React.FC<CardModalProps> = ({
  isOpen,
  onClose,
  projectId,
  boardId,
  columnId,
  cardId,
  onOpenManageTagsModal,
  projectColor,
}) => {
  const isEditing = Boolean(cardId);

  const {
    data: cardDetail,
    isLoading: isLoadingCard,
    isError: isErrorCard,
    error: cardError,
  } = useCardDetailQuery(projectId, boardId, cardId);
  const { data: allTags = [] } = useTagsQuery(projectId);

  if (!isOpen) return null;

  if (isEditing && isLoadingCard) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
        <div className="relative w-full max-w-md p-8 flex flex-col items-center justify-center gap-3 rounded-xl bg-card border border-border shadow-2xl">
          <Icons.Loader2 size={24} className="animate-spin text-primary" />
          <span className="text-xs font-mono text-muted-foreground">Loading card details...</span>
        </div>
      </div>
    );
  }

  if (isEditing && (isErrorCard || !cardDetail)) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-200">
        <div className="relative w-full max-w-md p-6 flex flex-col items-center justify-center gap-3 text-center rounded-xl bg-card border border-destructive/40 shadow-2xl">
          <div className="p-2.5 rounded-full bg-destructive/10 text-destructive">
            <Icons.AlertCircle size={24} />
          </div>
          <h3 className="text-sm font-bold font-mono text-foreground">Failed to Load Card</h3>
          <p className="text-xs text-muted-foreground">
            {cardError instanceof Error ? cardError.message : "The requested card could not be found."}
          </p>
          <button
            type="button"
            onClick={onClose}
            className="mt-2 px-4 py-1.5 rounded-md bg-secondary text-xs font-mono font-medium hover:bg-secondary/80 cursor-pointer"
          >
            Close
          </button>
        </div>
      </div>
    );
  }

  return (
    <CardModalInner
      key={cardId || "new-card"}
      projectId={projectId}
      boardId={boardId}
      columnId={columnId}
      cardId={cardId}
      cardDetail={cardDetail}
      allTags={allTags}
      onClose={onClose}
      onOpenManageTagsModal={onOpenManageTagsModal}
      projectColor={projectColor}
    />
  );
};
