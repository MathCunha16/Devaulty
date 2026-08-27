import React, { useState, useEffect } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { TagManagerSection } from "../../../components/TagManagerSection";
import {
  useMasterPasswordSetupStatusQuery,
  useVaultStatusQuery,
} from "~features/security/hooks/useSecurity";
import { MasterPasswordSetupCard } from "~features/security/components/MasterPasswordSetupCard";
import { UnlockVaultCard } from "~features/security/components/UnlockVaultCard";
import { VaultSecurityBanner } from "~features/security/components/VaultSecurityBanner";
import {
  useCredentialsQuery,
  useDeleteCredentialMutation,
} from "../hooks/useCredentials";
import { CredentialForm } from "../components/CredentialForm";
import { CredentialDetailModal } from "../components/CredentialDetailModal";
import { useInactivityAutoLock } from "../../../hooks/useInactivityAutoLock";
import type { CredentialSecretType } from "~types/api";
import styles from "../../../routes/projects.$projectId.module.css";

interface CredentialsWorkspaceProps {
  projectId: string;
  isActive: boolean;
  onNavigateTab?: (tab: "snippets" | "problems" | "credentials" | "notes" | "links") => void;
  onOpenManageTagsModal: () => void;
  initialSelectedId?: string;
}

export const CredentialsWorkspace: React.FC<CredentialsWorkspaceProps> = ({
  projectId,
  isActive,
  onNavigateTab,
  onOpenManageTagsModal,
  initialSelectedId,
}) => {
  // Security queries
  const { data: isSetupRequired, isLoading: isSetupLoading } =
    useMasterPasswordSetupStatusQuery(isActive);
  const { data: vaultStatus, isLoading: isVaultStatusLoading } =
    useVaultStatusQuery(isActive);

  const isSecurityLoading = isActive && (isSetupLoading || isVaultStatusLoading);
  const isVaultActive = vaultStatus?.active === true;
  const isVaultLocked = !isSecurityLoading && !isSetupRequired && !isVaultActive;

  // Credentials queries & mutations
  const { data: credentialsData, isLoading: isCredentialsLoading } =
    useCredentialsQuery(projectId, isActive && isVaultActive);
  const deleteCredentialMutation = useDeleteCredentialMutation(projectId);

  // Auto-lock after 15 minutes of inactivity in credentials workspace
  useInactivityAutoLock(isActive && isVaultActive, () => {
    onNavigateTab?.("snippets");
  });

  // Toggle data-vault-active on root document when viewing credentials workspace
  useEffect(() => {
    if (isActive) {
      document.documentElement.dataset.vaultActive = "true";
    } else {
      delete document.documentElement.dataset.vaultActive;
    }
    return () => {
      delete document.documentElement.dataset.vaultActive;
    };
  }, [isActive]);

  // Credentials UI states
  const [credentialSearchQuery, setCredentialSearchQuery] = useState("");
  const [credentialTypeFilter, setCredentialTypeFilter] = useState<
    "ALL" | CredentialSecretType
  >("ALL");
  const [isCredentialFormOpen, setIsCredentialFormOpen] = useState(false);
  const [editingCredentialId, setEditingCredentialId] = useState<string | undefined>(undefined);
  const [viewingCredentialId, setViewingCredentialId] = useState<string | undefined>(initialSelectedId);

  useEffect(() => {
    if (initialSelectedId && isVaultActive) {
      setViewingCredentialId(initialSelectedId);
    }
  }, [initialSelectedId, isVaultActive]);

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

  const credentials = credentialsData?.content || [];

  const filteredCredentials = credentials.filter((c) => {
    const matchesSearch =
      c.title.toLowerCase().includes(credentialSearchQuery.toLowerCase()) ||
      (c.relatedUrl && c.relatedUrl.toLowerCase().includes(credentialSearchQuery.toLowerCase())) ||
      (c.tags && c.tags.some((t) => t.name.toLowerCase().includes(credentialSearchQuery.toLowerCase())));
    const matchesType = credentialTypeFilter === "ALL" || c.secretType === credentialTypeFilter;
    return matchesSearch && matchesType;
  });

  const handleDeleteCredential = (credId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Credential",
      message: "Are you sure you want to delete the credential",
      itemName: title,
      warningText: "This action cannot be undone. The credential payload will be permanently wiped.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteCredentialMutation.mutateAsync(credId);
          toast.success("Credential deleted successfully");
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete credential");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  return (
    <div className="flex-1 flex flex-col h-full overflow-y-auto p-2">
      {isSecurityLoading ? (
        <div className="flex-1 flex items-center justify-center p-8">
          <Icons.Loader2 className="animate-spin text-muted-foreground" size={32} />
        </div>
      ) : isSetupRequired ? (
        <MasterPasswordSetupCard />
      ) : isVaultLocked ? (
        <UnlockVaultCard />
      ) : (
        <div className="flex-1 flex flex-col h-full overflow-hidden">
          <VaultSecurityBanner secondsLeft={vaultStatus?.secondsLeft} />

          <div className="flex-1 flex flex-col gap-4 overflow-hidden">
            <div className="flex items-center justify-between gap-4 flex-wrap">
              <div className="flex items-center gap-3">
                <div className={styles.searchBar} style={{ minWidth: "260px" }}>
                  <Icons.Search className={styles.searchIcon} size={14} />
                  <input
                    type="text"
                    className={styles.searchInput}
                    placeholder="Search credentials in vault..."
                    value={credentialSearchQuery}
                    onChange={(e) => setCredentialSearchQuery(e.target.value)}
                  />
                </div>

                {/* Type filter tabs */}
                <div className={styles.filterTabs}>
                  {(["ALL", "LOGIN", "API_KEY", "RAW_TEXT"] as const).map((type) => (
                    <button
                      key={type}
                      type="button"
                      className={`${styles.filterTab} ${credentialTypeFilter === type ? styles.filterTabActive : ""}`}
                      onClick={() => setCredentialTypeFilter(type)}
                    >
                      {type === "ALL" ? "ALL" : type === "RAW_TEXT" ? "TEXT" : type}
                    </button>
                  ))}
                </div>
              </div>

              <button
                type="button"
                className="px-4 py-2 text-xs font-mono font-bold bg-primary text-primary-foreground rounded border border-primary hover:brightness-110 transition-all cursor-pointer flex items-center gap-1.5"
                onClick={() => {
                  setEditingCredentialId(undefined);
                  setIsCredentialFormOpen(true);
                }}
              >
                <Icons.Plus size={14} />
                <span>New Credential</span>
              </button>
            </div>

            {/* Grid / List of Credential Cards */}
            {isCredentialsLoading ? (
              <div className="flex flex-col items-center justify-center py-16 gap-3">
                <Icons.Loader2 className="animate-spin text-primary" size={28} />
                <span className="text-xs text-muted-foreground font-mono">LOADING CREDENTIALS...</span>
              </div>
            ) : filteredCredentials.length === 0 ? (
              <div className={styles.placeholder}>
                <Icons.KeyRound size={48} className="text-muted-foreground animate-pulse" />
                <div className={styles.placeholderText}>
                  {credentialSearchQuery || credentialTypeFilter !== "ALL"
                    ? "No credentials match the active filters."
                    : "No credentials stored in this project yet. Click 'New Credential' to add one."}
                </div>
              </div>
            ) : (
              <div className={styles.credentialGrid}>
                {filteredCredentials.map((cred) => (
                  <div
                    key={cred.id}
                    role="button"
                    tabIndex={0}
                    className={styles.credentialCard}
                    onClick={() => setViewingCredentialId(cred.id)}
                    onKeyDown={(e) => {
                      if (e.target !== e.currentTarget) return;
                      if (e.key === "Enter" || e.key === " ") {
                        e.preventDefault();
                        setViewingCredentialId(cred.id);
                      }
                    }}

                  >

                    <div className={styles.credentialCardHeader}>
                      <div className={styles.credentialTitleGroup}>
                        <span className={styles.credentialTypeBadge}>
                          {cred.secretType === "LOGIN" && <Icons.UserCheck size={12} />}
                          {cred.secretType === "API_KEY" && <Icons.KeyRound size={12} />}
                          {cred.secretType === "RAW_TEXT" && <Icons.Code2 size={12} />}
                          <span>{cred.secretType}</span>
                        </span>
                        <h3 className={styles.credentialCardTitle}>{cred.title}</h3>
                        {cred.relatedUrl && (
                          <a
                            href={cred.relatedUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            className={styles.credentialUrl}
                            onClick={(e) => e.stopPropagation()}
                          >
                            <Icons.ExternalLink size={12} />
                            <span>{cred.relatedUrl.replace(/^https?:\/\//, "")}</span>
                          </a>
                        )}
                      </div>
                      <div
                        className={styles.credentialUnlockHint}
                        title="Click card to unlock vault payload"
                      >
                        <Icons.Lock size={13} />
                      </div>
                    </div>

                    {cred.notes && (
                      <p className={styles.credentialNotes}>{cred.notes}</p>
                    )}



                    <div className={styles.credentialCardFooter}>

                      <div onClick={(e) => e.stopPropagation()}>
                        <TagManagerSection
                          itemId={cred.id}
                          itemType="CREDENTIAL"
                          itemTags={cred.tags}
                          projectId={projectId}
                          onOpenManageTagsModal={onOpenManageTagsModal}
                          title=""
                          noBorder={true}
                        />
                      </div>

                      {/* Actions */}
                      <div className={styles.credentialActions} onClick={(e) => e.stopPropagation()}>
                        <button
                          type="button"
                          className={styles.actionBtn}
                          onClick={(e) => {
                            e.stopPropagation();
                            setEditingCredentialId(cred.id);
                            setIsCredentialFormOpen(true);
                          }}
                          title="Edit"
                        >
                          <Icons.Edit2 size={12} />
                        </button>
                        <button
                          type="button"
                          className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                          onClick={(e) => {
                            e.stopPropagation();
                            handleDeleteCredential(cred.id, cred.title);
                          }}
                          title="Delete"
                        >
                          <Icons.Trash2 size={12} />
                        </button>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}

      <CredentialForm
        isOpen={isCredentialFormOpen}
        onClose={() => {
          setIsCredentialFormOpen(false);
          setEditingCredentialId(undefined);
        }}
        projectId={projectId}
        credentialId={editingCredentialId}
      />

      {viewingCredentialId && (
        <CredentialDetailModal
          isOpen={!!viewingCredentialId}
          onClose={() => setViewingCredentialId(undefined)}
          projectId={projectId}
          credentialId={viewingCredentialId}
        />
      )}

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
    </div>
  );
};
