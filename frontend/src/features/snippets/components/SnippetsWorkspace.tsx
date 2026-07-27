import React, { useState } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { CodeViewer } from "../../../components/CodeViewer";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { TagManagerSection } from "../../../components/TagManagerSection";
import { copyToClipboard } from "../../../utils/clipboardUtils";
import { formatUpdatedDate } from "../../../utils/dateUtils";


const SnippetForm = React.lazy(() =>
  import("./SnippetForm").then((m) => ({ default: m.SnippetForm }))
);

import {
  useSnippetsQuery,
  useSnippetQuery,
  useDeleteSnippetMutation,
} from "../hooks/useSnippets";
import { mapLanguageToMonaco } from "../utils/languageUtils";
import type { SnippetType } from "~types/api";
import styles from "../../../routes/projects.$projectId.module.css";


interface SnippetsWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal: () => void;
}

export const SnippetsWorkspace: React.FC<SnippetsWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
}) => {
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const deleteSnippetMutation = useDeleteSnippetMutation(projectId);

  const [selectedSnippetId, setSelectedSnippetId] = useState<string | undefined>(undefined);
  const [snippetSearchQuery, setSnippetSearchQuery] = useState("");
  const [snippetTypeFilter, setSnippetTypeFilter] = useState<"ALL" | SnippetType>("ALL");

  const [isSnippetFormOpen, setIsSnippetFormOpen] = useState(false);
  const [editingSnippetId, setEditingSnippetId] = useState<string | undefined>(undefined);

  const [copiedId, setCopiedId] = useState<string | null>(null);

  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    itemName?: string;
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

  const snippets = snippetsData?.content || [];
  const { data: snippetDetail } = useSnippetQuery(projectId, selectedSnippetId);

  const filteredSnippets = snippets.filter((s) => {
    const matchesSearch =
      s.title.toLowerCase().includes(snippetSearchQuery.toLowerCase()) ||
      (s.description && s.description.toLowerCase().includes(snippetSearchQuery.toLowerCase())) ||
      (s.content && s.content.toLowerCase().includes(snippetSearchQuery.toLowerCase()));
    const matchesType = snippetTypeFilter === "ALL" || s.snippetType === snippetTypeFilter;
    return matchesSearch && matchesType;
  });

  const selectedSnippet = snippetDetail || snippets.find((s) => s.id === selectedSnippetId);

  const handleOpenCreateSnippet = () => {
    setEditingSnippetId(undefined);
    setIsSnippetFormOpen(true);
  };

  const handleOpenEditSnippet = () => {
    if (!selectedSnippet) return;
    setEditingSnippetId(selectedSnippet.id);
    setIsSnippetFormOpen(true);
  };

  const handleDeleteSnippet = (snippetId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Snippet",
      message: "Are you sure you want to delete the snippet",
      itemName: title,
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteSnippetMutation.mutateAsync(snippetId);
          toast.success("Snippet deleted successfully");
          if (selectedSnippetId === snippetId) setSelectedSnippetId(undefined);
          closeConfirmModal();
        } catch (err) {
          toast.error(err instanceof Error ? err.message : "Failed to delete snippet");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const handleCopy = async (content: string, id: string) => {

    const success = await copyToClipboard(content);
    if (success) {
      setCopiedId(id);
      toast.success("Copied to clipboard!");
      setTimeout(() => setCopiedId(null), 2000);
    } else {
      toast.error("Failed to copy content");
    }
  };


  return (
    <>
      {/* Left Side: Snippets navigation list */}
      <div className={styles.leftPanel}>
        <button className={styles.newSnippetBtn} onClick={handleOpenCreateSnippet}>
          <Icons.Plus size={14} />
          <span>Add Snippet</span>
        </button>

        <div className={styles.searchBar}>
          <Icons.Search className={styles.searchIcon} size={14} />
          <input
            type="text"
            placeholder="Search snippets..."
            className={styles.searchInput}
            value={snippetSearchQuery}
            onChange={(e) => setSnippetSearchQuery(e.target.value)}
          />
        </div>

        <div className={styles.filterTabs}>
          <button
            className={`${styles.filterTab} ${snippetTypeFilter === "ALL" ? styles.filterTabActive : ""}`}
            onClick={() => setSnippetTypeFilter("ALL")}
          >
            ALL
          </button>
          <button
            className={`${styles.filterTab} ${snippetTypeFilter === "CODE" ? styles.filterTabActive : ""}`}
            onClick={() => setSnippetTypeFilter("CODE")}
          >
            CODE
          </button>
          <button
            className={`${styles.filterTab} ${snippetTypeFilter === "COMMAND" ? styles.filterTabActive : ""}`}
            onClick={() => setSnippetTypeFilter("COMMAND")}
          >
            CMD
          </button>
        </div>

        <div className={styles.snippetList}>
          {filteredSnippets.length === 0 ? (
            <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border font-mono">
              No snippets found
            </div>
          ) : (
            filteredSnippets.map((s) => (
              <button
                key={s.id}
                className={`${styles.snippetItem} ${
                  selectedSnippetId === s.id ? styles.snippetItemActive : ""
                }`}
                onClick={() => setSelectedSnippetId(s.id)}
              >
                <div className={styles.snippetItemHeader}>
                  <span className={styles.snippetItemTitle}>{s.title}</span>
                  <span className="text-[10px] text-muted-foreground font-mono shrink-0 ml-auto pt-0.5">
                    {new Date(s.createdAt).toLocaleDateString()}
                  </span>
                </div>


                {s.description && <p className={styles.snippetItemDesc}>{s.description}</p>}
                <div className={styles.snippetItemFooter}>
                  <span className={styles.badgeType}>{s.snippetType}</span>
                  <span className={styles.badgeLang}>{s.language}</span>
                </div>
              </button>
            ))
          )}
        </div>
      </div>

      {/* Right Side: Snippets detail view */}
      <div className={styles.rightPanel}>
        {selectedSnippet ? (
          <div className={styles.snippetDetailScroll}>
            <div className={styles.snippetDetailContainer}>
              <div className={styles.detailHeader}>
                <div className={styles.detailTitleSection}>
                  <div className={styles.detailTitleRow}>
                    <h2 className={styles.detailTitle}>{selectedSnippet.title}</h2>
                    <div className="flex gap-1.5">
                      <span className={styles.badgeType}>{selectedSnippet.snippetType}</span>
                      <span className={styles.badgeLang}>{selectedSnippet.language}</span>
                    </div>
                  </div>
                  {selectedSnippet.description && (
                    <p className={styles.detailDesc}>{selectedSnippet.description}</p>
                  )}
                  <div className={styles.problemMetadataRow}>
                    <span>Created: {new Date(selectedSnippet.createdAt).toLocaleDateString()}</span>
                    {formatUpdatedDate(selectedSnippet.updatedAt, selectedSnippet.createdAt) && (
                      <span>
                        Updated: {formatUpdatedDate(selectedSnippet.updatedAt, selectedSnippet.createdAt)}
                      </span>
                    )}
                  </div>
                </div>


                <div className={styles.detailActions}>
                  <button
                    className={styles.btnIcon}
                    onClick={handleOpenEditSnippet}
                    title="Edit Snippet"
                  >
                    <Icons.Edit3 size={14} />
                  </button>
                  <button
                    className={`${styles.btnIcon} ${styles.btnIconDanger}`}
                    onClick={() =>
                      handleDeleteSnippet(selectedSnippet.id, selectedSnippet.title)
                    }
                    title="Delete Snippet"
                  >
                    <Icons.Trash2 size={14} />
                  </button>
                </div>
              </div>

              {/* Tag Section */}
              <TagManagerSection
                itemId={selectedSnippet.id}
                itemType="SNIPPET"
                itemTags={selectedSnippet.tags}
                projectId={projectId}
                onOpenManageTagsModal={onOpenManageTagsModal}
              />

              <div className={styles.codePanel}>
                <div className={styles.codePanelHeader}>
                  <span className={styles.codePanelLabel}>
                    Source Code ({selectedSnippet.language.toLowerCase()})
                  </span>
                  <button
                    className={styles.copyButton}
                    onClick={() => handleCopy(selectedSnippet.content, selectedSnippet.id)}
                  >
                    {copiedId === selectedSnippet.id ? (
                      <>
                        <Icons.Check size={12} className="text-emerald-500" />
                        <span className="text-emerald-500 font-bold">Copied!</span>
                      </>
                    ) : (
                      <>
                        <Icons.Copy size={12} />
                        <span>Copy code</span>
                      </>
                    )}
                  </button>
                </div>

                <CodeViewer
                  height="380px"
                  language={mapLanguageToMonaco(selectedSnippet.language)}
                  code={selectedSnippet.content}
                />
              </div>
            </div>
          </div>
        ) : (
          <div className={styles.placeholder}>
            <Icons.Terminal size={48} className="text-muted-foreground animate-pulse" />
            <div className={styles.placeholderText}>
              No snippet selected. Choose a snippet from the list or click "Add Snippet" to create
              a new one.
            </div>
          </div>
        )}
      </div>

      <React.Suspense fallback={null}>
        <SnippetForm
          isOpen={isSnippetFormOpen}
          onClose={() => {
            setIsSnippetFormOpen(false);
            setEditingSnippetId(undefined);
          }}
          projectId={projectId}
          snippetId={editingSnippetId}
        />
      </React.Suspense>

      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={closeConfirmModal}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        itemName={confirmModal.itemName}
        isLoading={confirmModal.isLoading}
      />
    </>
  );
};
