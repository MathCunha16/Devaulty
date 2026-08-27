import React, { useState } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { marked } from "marked";
import DOMPurify from "dompurify";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { TagManagerSection } from "../../../components/TagManagerSection";
import { formatUpdatedDate } from "../../../utils/dateUtils";
import { copyToClipboard } from "../../../utils/clipboardUtils";


import { NoteForm } from "../components/NoteForm";
import {
  useNotesQuery,
  useNoteQuery,
  useUpdateNoteMutation,
  useDeleteNoteMutation,
} from "../hooks/useNotes";
import styles from "../../../routes/projects.$projectId.module.css";

interface NotesWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal: () => void;
  initialSelectedId?: string;
}

export const NotesWorkspace: React.FC<NotesWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
  initialSelectedId,
}) => {
  const { data: notesData } = useNotesQuery(projectId);
  const deleteNoteMutation = useDeleteNoteMutation(projectId);

  const [selectedNoteId, setSelectedNoteId] = useState<string | undefined>(initialSelectedId);
  const [noteSearchQuery, setNoteSearchQuery] = useState("");
  const [noteArchivedFilter, setNoteArchivedFilter] = useState<"ACTIVE" | "ARCHIVED" | "ALL">("ACTIVE");

  React.useEffect(() => {
    if (initialSelectedId) {
      setSelectedNoteId(initialSelectedId);
    }
  }, [initialSelectedId]);

  const [isNoteFormOpen, setIsNoteFormOpen] = useState(false);
  const [editingNoteId, setEditingNoteId] = useState<string | undefined>(undefined);

  // Track which note is currently having its content edited inline
  const [editingNoteContentId, setEditingNoteContentId] = useState<string | null>(null);
  const [inlineContent, setInlineContent] = useState("");
  const isEditingContent = editingNoteContentId === selectedNoteId && !!selectedNoteId;

  const updateNoteMutation = useUpdateNoteMutation(projectId, selectedNoteId || "");
  const { data: noteDetail } = useNoteQuery(projectId, selectedNoteId || "");

  const [useMarkdown, setUseMarkdown] = useState<boolean>(() => {
    try {
      const saved = localStorage.getItem("devaulty_notes_markdown_preview");
      return saved !== "false";
    } catch {
      return true;
    }
  });

  const handleToggleMarkdown = (enabled: boolean) => {
    setUseMarkdown(enabled);
    try {
      localStorage.setItem("devaulty_notes_markdown_preview", String(enabled));
    } catch {
      // ignore storage errors
    }
  };

  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    itemName?: string;
    warningText?: string;
    onConfirm: () => Promise<void>;
    isLoading: boolean;
  }>({
    isOpen: false,
    title: "",
    message: "",
    onConfirm: async () => {},
    isLoading: false,
  });

  const closeConfirmModal = () =>
    setConfirmModal((prev) => ({ ...prev, isOpen: false, isLoading: false }));

  const notes = notesData?.content || [];

  const filteredNotes = notes.filter((n) => {
    const matchesSearch =
      n.title.toLowerCase().includes(noteSearchQuery.toLowerCase()) ||
      (n.tags && n.tags.some((t) => t.name.toLowerCase().includes(noteSearchQuery.toLowerCase())));

    let matchesArchived = true;
    if (noteArchivedFilter === "ACTIVE") {
      matchesArchived = !n.archived;
    } else if (noteArchivedFilter === "ARCHIVED") {
      matchesArchived = n.archived;
    }

    return matchesSearch && matchesArchived;
  });


  const handleDeleteNote = (noteId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete System Note",
      message: "Are you sure you want to permanently delete the note",
      itemName: title,
      warningText: "This action cannot be undone. The note contents will be permanently lost.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteNoteMutation.mutateAsync(noteId);
          toast.success("Note deleted successfully");
          if (selectedNoteId === noteId) setSelectedNoteId(undefined);
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete note");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const handleSaveInlineContent = async () => {
    if (!noteDetail) return;
    try {
      await updateNoteMutation.mutateAsync({
        title: noteDetail.title,
        content: inlineContent,
      });
      toast.success("Note content saved");
      setEditingNoteContentId(null);
    } catch {
      toast.error("Failed to save note content");
    }
  };

  // Interactive Checkbox Toggle in Markdown Preview Mode
  const handleMarkdownClick = async (e: React.MouseEvent<HTMLDivElement>) => {
    const target = e.target as HTMLElement;
    if (target && target.tagName === "INPUT" && (target as HTMLInputElement).type === "checkbox") {
      const checkbox = target as HTMLInputElement;
      const container = e.currentTarget;
      const checkboxes = Array.from(container.querySelectorAll('input[type="checkbox"]'));
      const index = checkboxes.indexOf(checkbox);
      if (index !== -1 && noteDetail?.content) {
        const lines = noteDetail.content.split("\n");
        let inCodeBlock = false;
        let taskListIndex = 0;
        const newLines = lines.map((line) => {
          if (line.trim().startsWith("```")) {
            inCodeBlock = !inCodeBlock;
            return line;
          }
          if (!inCodeBlock && /^\s*[-*+]\s*\[([ xX])\]/.test(line)) {
            if (taskListIndex === index) {
              taskListIndex++;
              return line.replace(/^(\s*[-*+]\s*\[)([ xX])(\])/, (_, p1, p2, p3) => {
                const isChecked = p2.toLowerCase() === "x";
                return `${p1}${isChecked ? " " : "x"}${p3}`;
              });
            }
            taskListIndex++;
          }
          return line;
        });

        const updatedContent = newLines.join("\n");

        if (updatedContent !== noteDetail.content) {
          try {
            await updateNoteMutation.mutateAsync({
              title: noteDetail.title,
              content: updatedContent,
            });
            toast.success("Checkbox state updated");
          } catch {
            toast.error("Failed to update checkbox");
          }
        }
      }
    }
  };

  const renderMarkdown = (content: string | undefined) => {
    if (!content) {
      return '<span class="text-muted-foreground italic">No content documented. Click Edit Content or click here to start writing.</span>';
    }
    try {
      const rawHtml = marked.parse(content, { breaks: true, gfm: true }) as string;
      const sanitized = DOMPurify.sanitize(rawHtml);
      return sanitized.replace(/<input([^>]*)\sdisabled(=["']?["']?)?/gi, '<input$1');
    } catch {
      return DOMPurify.sanitize(content);
    }
  };


  const formattedUpdatedDate = noteDetail
    ? formatUpdatedDate(noteDetail.updatedAt, noteDetail.createdAt)
    : null;

  return (
    <>
      {/* Left Side: Notes navigation list */}
      <div className={styles.leftPanel}>
        <button
          type="button"
          className={styles.newSnippetBtn}
          onClick={() => {
            setEditingNoteId(undefined);
            setIsNoteFormOpen(true);
          }}
        >
          <Icons.Plus size={14} />
          <span>Add Note</span>
        </button>

        <div className={styles.searchBar}>
          <Icons.Search className={styles.searchIcon} size={14} />
          <input
            type="text"
            placeholder="Search notes..."
            className={styles.searchInput}
            value={noteSearchQuery}
            onChange={(e) => setNoteSearchQuery(e.target.value)}
          />
        </div>

        <div className={styles.filterTabs}>
          <button
            type="button"
            className={`${styles.filterTab} ${noteArchivedFilter === "ACTIVE" ? styles.filterTabActive : ""}`}
            onClick={() => setNoteArchivedFilter("ACTIVE")}
          >
            ACTIVE
          </button>
          <button
            type="button"
            className={`${styles.filterTab} ${noteArchivedFilter === "ARCHIVED" ? styles.filterTabActive : ""}`}
            onClick={() => setNoteArchivedFilter("ARCHIVED")}
          >
            ARCHIVED
          </button>
          <button
            type="button"
            className={`${styles.filterTab} ${noteArchivedFilter === "ALL" ? styles.filterTabActive : ""}`}
            onClick={() => setNoteArchivedFilter("ALL")}
          >
            ALL
          </button>
        </div>

        <div className={styles.snippetList}>
          {filteredNotes.length === 0 ? (
            <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border font-mono">
              No notes found
            </div>
          ) : (
            filteredNotes.map((n) => (
              <button
                key={n.id}
                className={`${styles.snippetItem} ${selectedNoteId === n.id ? styles.snippetItemActive : ""}`}
                onClick={() => setSelectedNoteId(n.id)}
              >
                <div className={styles.snippetItemHeader}>
                  <span className={styles.snippetItemTitle}>{n.title}</span>
                  <div className="flex items-center gap-1.5 shrink-0 ml-auto pt-0.5">
                    {n.archived && <span className="text-[10px] text-amber-500 font-mono">ARCHIVED</span>}
                    <span className="text-[10px] text-muted-foreground font-mono">
                      {new Date(n.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                </div>


                {n.tags && n.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1">
                    {n.tags.map((tag) => (
                      <span
                        key={tag.id}
                        className={styles.badge}
                        style={{
                          backgroundColor: `${tag.color || "#8b5cf6"}15`,
                          color: tag.color || "#8b5cf6",
                          border: `1px solid ${tag.color || "#8b5cf6"}30`,
                          fontSize: "10px",
                        }}
                      >
                        {tag.name}
                      </span>
                    ))}
                  </div>
                )}
              </button>
            ))
          )}
        </div>
      </div>

      {/* Right Side: Note Workspace details panel */}
      <div className={styles.rightPanel}>
        {noteDetail ? (
          <div className={styles.problemDetailScroll}>
            <div className={styles.problemDetailContainer}>
              <div className={styles.detailHeader}>
                <div className={styles.detailTitleSection}>
                  <h2 className={styles.detailTitle}>{noteDetail.title}</h2>
                  <div className={styles.problemMetadataRow}>
                    <span>Created: {new Date(noteDetail.createdAt).toLocaleDateString()}</span>
                    {formattedUpdatedDate && <span>Updated: {formattedUpdatedDate}</span>}
                  </div>
                </div>

                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={async () => {
                      const contentToCopy =
                        editingNoteContentId === noteDetail.id
                          ? inlineContent
                          : noteDetail.content;
                      if (!contentToCopy) return;
                      const ok = await copyToClipboard(contentToCopy);
                      if (ok) {
                        toast.success("Note content copied to clipboard");
                      } else {
                        toast.error("Failed to copy note content");
                      }
                    }}
                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                    title="Copy Note Content"
                  >
                    <Icons.Copy size={12} />
                    <span>Copy</span>
                  </button>


                  <button
                    type="button"
                    onClick={() => {
                      setEditingNoteId(noteDetail.id);
                      setIsNoteFormOpen(true);
                    }}
                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                    title="Edit Note Metadata (Title)"
                  >
                    <Icons.Edit3 size={12} />
                    <span>Edit Metadata</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => handleDeleteNote(noteDetail.id, noteDetail.title)}
                    className="flex items-center gap-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 border border-red-500/20 px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                  >
                    <Icons.Trash2 size={12} />
                    <span>Delete</span>
                  </button>
                </div>

              </div>

              <div className="p-6 flex flex-col gap-6">
                {/* Tag Section */}
                <TagManagerSection
                  itemId={noteDetail.id}
                  itemType="NOTE"
                  itemTags={noteDetail.tags}
                  projectId={projectId}
                  onOpenManageTagsModal={onOpenManageTagsModal}
                  title="Tags"
                />

                {/* Content panel with inline editing & markdown checkbox toggles */}
                <div className="flex-grow flex flex-col gap-2">
                  <div className="flex items-center justify-between">
                    <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">
                      Note Content {isEditingContent ? "(Editing)" : ""}
                    </span>
                    <div className="flex items-center gap-2">
                      {isEditingContent ? (
                        <div className="flex items-center gap-1.5">
                          <button
                            type="button"
                            onClick={() => setEditingNoteContentId(null)}
                            className="px-2.5 py-1 text-[11px] font-mono rounded cursor-pointer transition-all border border-border bg-transparent text-muted-foreground hover:text-foreground"
                            disabled={updateNoteMutation.isPending}
                          >
                            Cancel
                          </button>
                          <button
                            type="button"
                            onClick={handleSaveInlineContent}
                            className="px-3 py-1 text-[11px] font-mono font-bold rounded cursor-pointer transition-all bg-primary text-primary-foreground shadow-sm flex items-center gap-1"
                            disabled={updateNoteMutation.isPending}
                          >
                            {updateNoteMutation.isPending ? (
                              <Icons.Loader2 size={12} className="animate-spin" />
                            ) : (
                              <Icons.Save size={12} />
                            )}
                            <span>Save Content</span>
                          </button>
                        </div>
                      ) : (
                        <button
                          type="button"
                          onClick={() => {
                            setInlineContent(noteDetail.content || "");
                            setEditingNoteContentId(noteDetail.id);
                          }}
                          className="flex items-center gap-1 text-[11px] font-mono text-primary hover:underline bg-transparent border-0 cursor-pointer"
                        >
                          <Icons.Edit2 size={11} />
                          <span>Edit Content</span>
                        </button>
                      )}

                      <div className="flex items-center gap-1 bg-secondary/40 p-0.5 rounded border border-border">
                        <button
                          type="button"
                          onClick={() => handleToggleMarkdown(false)}
                          className={`px-2 py-1 text-[10px] font-mono rounded cursor-pointer transition-all ${
                            !useMarkdown
                              ? "bg-primary text-primary-foreground shadow-sm font-bold"
                              : "text-muted-foreground hover:text-foreground bg-transparent"
                          }`}
                          style={{ border: "none" }}
                        >
                          RAW
                        </button>
                        <button
                          type="button"
                          onClick={() => handleToggleMarkdown(true)}
                          className={`px-2 py-1 text-[10px] font-mono rounded cursor-pointer transition-all ${
                            useMarkdown
                              ? "bg-primary text-primary-foreground shadow-sm font-bold"
                              : "text-muted-foreground hover:text-foreground bg-transparent"
                          }`}
                          style={{ border: "none" }}
                        >
                          MARKDOWN
                        </button>
                      </div>
                    </div>
                  </div>

                  {isEditingContent ? (
                    <div className="flex flex-col gap-2">
                      <textarea
                        className="bg-background/80 border border-primary/50 rounded p-4 font-mono text-sm leading-relaxed min-h-[320px] focus:outline-none focus:ring-1 focus:ring-primary"
                        value={inlineContent}
                        onChange={(e) => setInlineContent(e.target.value)}
                        placeholder="Write your markdown note content here..."
                        autoFocus
                      />
                    </div>
                  ) : useMarkdown ? (
                    <div
                      className={`bg-background/50 border border-border rounded p-6 text-sm leading-relaxed overflow-y-auto min-h-[300px] ${styles.markdownContainer}`}
                      onClick={handleMarkdownClick}
                      dangerouslySetInnerHTML={{ __html: renderMarkdown(noteDetail.content) }}
                    />
                  ) : (
                    <div
                      className="bg-background/50 border border-border rounded p-4 font-mono text-sm whitespace-pre-wrap leading-relaxed overflow-y-auto min-h-[300px] cursor-pointer hover:border-border/80 transition-colors"
                      onClick={() => {
                        setInlineContent(noteDetail.content || "");
                        setEditingNoteContentId(noteDetail.id);
                      }}
                      title="Click to edit content"
                    >

                      {noteDetail.content || (
                        <span className="text-muted-foreground italic">
                          No content documented. Click here to add details.
                        </span>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className={styles.placeholder}>
            <Icons.FileText size={48} className="text-muted-foreground animate-pulse" />
            <div className={styles.placeholderText}>
              No note selected. Select a note from the navigator or click "Add Note" to write a new note.
            </div>
          </div>
        )}
      </div>

      <NoteForm
        isOpen={isNoteFormOpen}
        onClose={() => {
          setIsNoteFormOpen(false);
          setEditingNoteId(undefined);
        }}
        projectId={projectId}
        noteId={editingNoteId}
      />

      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={closeConfirmModal}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        itemName={confirmModal.itemName}
        warningText={confirmModal.warningText}
        isLoading={confirmModal.isLoading}
      />
    </>
  );
};
