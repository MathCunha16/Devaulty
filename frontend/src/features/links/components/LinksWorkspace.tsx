import React, { useState } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { TagManagerSection } from "../../../components/TagManagerSection";
import { formatUpdatedDate } from "../../../utils/dateUtils";
import { LinkForm } from "../components/LinkForm";
import {
  useLinksQuery,
  useLinkQuery,
  useDeleteLinkMutation,
} from "../hooks/useLinks";
import styles from "../../../routes/projects.$projectId.module.css";

interface LinksWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal: () => void;
  initialSelectedId?: string;
}

export const LinksWorkspace: React.FC<LinksWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
  initialSelectedId,
}) => {
  const { data: linksData } = useLinksQuery(projectId);
  const deleteLinkMutation = useDeleteLinkMutation(projectId);

  const [selectedLinkId, setSelectedLinkId] = useState<string | undefined>(initialSelectedId);
  const [linkSearchQuery, setLinkSearchQuery] = useState("");

  React.useEffect(() => {
    if (initialSelectedId) {
      setSelectedLinkId(initialSelectedId);
    }
  }, [initialSelectedId]);

  const [isLinkFormOpen, setIsLinkFormOpen] = useState(false);
  const [editingLinkId, setEditingLinkId] = useState<string | undefined>(undefined);

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

  const links = linksData?.content || [];

  const filteredLinks = links.filter((l) => {
    return (
      l.title.toLowerCase().includes(linkSearchQuery.toLowerCase()) ||
      l.url.toLowerCase().includes(linkSearchQuery.toLowerCase()) ||
      (l.tags && l.tags.some((t) => t.name.toLowerCase().includes(linkSearchQuery.toLowerCase())))
    );
  });

  const { data: linkDetail } = useLinkQuery(projectId, selectedLinkId || "");

  const handleDeleteLink = (linkId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Web Link",
      message: "Are you sure you want to delete the web link",
      itemName: title,
      warningText: "This action cannot be undone.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteLinkMutation.mutateAsync(linkId);
          toast.success("Link deleted successfully");
          if (selectedLinkId === linkId) setSelectedLinkId(undefined);
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete link");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  return (
    <>
      {/* Left Side: Links navigation list */}
      <div className={styles.leftPanel}>
        <button
          type="button"
          className={styles.newSnippetBtn}
          onClick={() => {
            setEditingLinkId(undefined);
            setIsLinkFormOpen(true);
          }}
        >
          <Icons.Plus size={14} />
          <span>Add Link</span>
        </button>

        <div className={styles.searchBar}>
          <Icons.Search className={styles.searchIcon} size={14} />
          <input
            type="text"
            placeholder="Search links..."
            className={styles.searchInput}
            value={linkSearchQuery}
            onChange={(e) => setLinkSearchQuery(e.target.value)}
          />
        </div>

        <div className={styles.snippetList}>
          {filteredLinks.length === 0 ? (
            <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border font-mono">
              No links found
            </div>
          ) : (
            filteredLinks.map((l) => (
              <button
                key={l.id}
                className={`${styles.snippetItem} ${selectedLinkId === l.id ? styles.snippetItemActive : ""}`}
                onClick={() => setSelectedLinkId(l.id)}
              >
                <div className={styles.snippetItemHeader}>
                  <span className={styles.snippetItemTitle}>{l.title}</span>
                  <div className="flex items-center gap-1.5 shrink-0 ml-auto pt-0.5">
                    <span className="text-[10px] text-muted-foreground font-mono">
                      {new Date(l.createdAt).toLocaleDateString()}
                    </span>
                    <Icons.ExternalLink size={12} className="text-muted-foreground" />
                  </div>
                </div>


                <span className="text-[10px] text-muted-foreground font-mono truncate max-w-[220px] block">
                  {l.url}
                </span>
                {l.tags && l.tags.length > 0 && (
                  <div className="flex flex-wrap gap-1 mt-1">
                    {l.tags.map((tag) => (
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

      {/* Right Side: Link Workspace details panel */}
      <div className={styles.rightPanel}>
        {linkDetail ? (
          <div className={styles.problemDetailScroll}>
            <div className={styles.problemDetailContainer}>
              <div className={styles.detailHeader}>
                <div className={styles.detailTitleSection}>
                  <h2 className={styles.detailTitle}>{linkDetail.title}</h2>
                  <div className={styles.problemMetadataRow}>
                    <span>Created: {new Date(linkDetail.createdAt).toLocaleDateString()}</span>
                    {formatUpdatedDate(linkDetail.updatedAt, linkDetail.createdAt) && (
                      <span>
                        Updated: {formatUpdatedDate(linkDetail.updatedAt, linkDetail.createdAt)}
                      </span>
                    )}
                  </div>
                </div>

                <div className="flex gap-2">
                  <button
                    type="button"
                    onClick={() => {
                      setEditingLinkId(linkDetail.id);
                      setIsLinkFormOpen(true);
                    }}
                    className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                  >
                    <Icons.Edit3 size={12} />
                    <span>Edit</span>
                  </button>
                  <button
                    type="button"
                    onClick={() => handleDeleteLink(linkDetail.id, linkDetail.title)}
                    className="flex items-center gap-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 border border-red-500/20 px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                  >
                    <Icons.Trash2 size={12} />
                    <span>Delete</span>
                  </button>
                </div>
              </div>

              <div className="p-6 flex flex-col gap-6">
                {/* URL card block */}
                <div className="bg-background/50 border border-border rounded p-4 flex flex-col gap-3">
                  <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Web Address</span>
                  <div className="flex items-center justify-between gap-4">
                    <span className="font-mono text-sm text-primary truncate flex-grow">{linkDetail.url}</span>
                    <a
                      href={linkDetail.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="flex items-center gap-1.5 text-xs text-primary-foreground bg-primary hover:bg-primary/95 border border-primary px-3.5 py-1.5 rounded transition-colors font-bold decoration-none cursor-pointer"
                    >
                      <span>Open Link</span>
                      <Icons.ExternalLink size={12} />
                    </a>
                  </div>
                </div>

                {/* Tag Section */}
                <TagManagerSection
                  itemId={linkDetail.id}
                  itemType="LINK"
                  itemTags={linkDetail.tags}
                  projectId={projectId}
                  onOpenManageTagsModal={onOpenManageTagsModal}
                  title="Tags"
                />

                {/* Description Panel */}
                <div className="flex-grow flex flex-col gap-2">
                  <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Description</span>
                  <div className="bg-background/50 border border-border rounded p-4 font-mono text-sm whitespace-pre-wrap leading-relaxed overflow-y-auto min-h-[150px]">
                    {linkDetail.description || (
                      <span className="text-muted-foreground italic">
                        No description documented. Click Edit to add details.
                      </span>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className={styles.placeholder}>
            <Icons.Link2 size={48} className="text-muted-foreground animate-pulse" />
            <div className={styles.placeholderText}>
              No link selected. Select a web link from the navigator or click "Add Link" to register a new destination.
            </div>
          </div>
        )}
      </div>

      <LinkForm
        isOpen={isLinkFormOpen}
        onClose={() => {
          setIsLinkFormOpen(false);
          setEditingLinkId(undefined);
        }}
        projectId={projectId}
        linkId={editingLinkId}
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
