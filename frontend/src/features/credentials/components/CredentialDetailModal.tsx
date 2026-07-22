import React, { useState, useEffect, useRef } from "react";
import {
  X,
  Eye,
  EyeOff,
  Copy,
  ExternalLink,
  Loader2,
  ShieldCheck,
  UserCheck,
  KeyRound,
  Code2,
} from "lucide-react";
import { toast } from "sonner";
import { useCredentialQuery } from "../hooks/useCredentials";
import styles from "./CredentialDetailModal.module.css";

interface CredentialDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  credentialId: string;
}

export const CredentialDetailModal: React.FC<CredentialDetailModalProps> = ({
  isOpen,
  onClose,
  projectId,
  credentialId,
}) => {
  const [showPassword, setShowPassword] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);
  const [showRawText, setShowRawText] = useState(false);

  const { data: cred, isLoading, isError } = useCredentialQuery(projectId, credentialId, isOpen);

  // Focus trap refs
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
        return;
      }
      if (e.key !== "Tab") return;

      const modal = modalRef.current;
      if (!modal) return;

      const focusables = modal.querySelectorAll<HTMLElement>(
        'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
      );
      if (focusables.length === 0) return;

      const firstElement = focusables[0];
      const lastElement = focusables[focusables.length - 1];

      if (e.shiftKey) {
        if (document.activeElement === firstElement) {
          lastElement.focus();
          e.preventDefault();
        }
      } else {
        if (document.activeElement === lastElement) {
          firstElement.focus();
          e.preventDefault();
        }
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    toast.success(`${label} copied to clipboard!`);
  };

  const payload = cred?.decryptedPayload || {};

  const getTypeIcon = (type?: string) => {
    switch (type) {
      case "LOGIN":
        return <UserCheck size={20} />;
      case "API_KEY":
        return <KeyRound size={20} />;
      case "RAW_TEXT":
        return <Code2 size={20} />;
      default:
        return <KeyRound size={20} />;
    }
  };

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div
        ref={modalRef}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="credential-detail-title"
      >
        <div className={styles.header}>
          <div className={styles.titleGroup}>
            <div className={styles.iconBadge}>{getTypeIcon(cred?.secretType)}</div>
            <div className={styles.titleWrapper}>
              <span className={styles.typeBadge}>{cred?.secretType || "CREDENTIAL"}</span>
              <h2 id="credential-detail-title" className={styles.title}>
                {cred?.title || "CREDENTIAL DETAILS"}
              </h2>
            </div>
          </div>
          <button className={styles.closeBtn} onClick={onClose} aria-label="Close modal">
            <X size={18} />
          </button>
        </div>

        <div className={styles.content}>
          {isLoading ? (
            <div className="flex flex-col items-center justify-center py-16 gap-3">
              <Loader2 className="animate-spin text-primary" size={32} />
              <span className="text-xs text-muted-foreground font-mono">DECRYPTING VAULT PAYLOAD...</span>
            </div>
          ) : isError || !cred ? (
            <div className="text-center py-12 text-destructive font-mono text-xs">
              Failed to decrypt credential. Please verify if your vault session remains active.
            </div>
          ) : (
            <>
              {/* Decrypted Payload Hero Card */}
              <div className={styles.secretSection}>
                <div className={styles.secretHeader}>
                  <div className={styles.secretHeaderTitle}>
                    <ShieldCheck size={16} />
                    <span>DECRYPTED VAULT PAYLOAD</span>
                  </div>
                </div>

                {cred.secretType === "LOGIN" && (
                  <>
                    {(payload.username || payload.user || payload.email) && (
                      <div className={styles.secretField}>
                        <div className={styles.labelRow}>
                          <span className={styles.label}>Username / Email</span>
                        </div>
                        <div className={styles.valueRow}>
                          <span className={styles.valueText}>
                            {payload.username || payload.user || payload.email}
                          </span>
                          <button
                            type="button"
                            className={styles.actionBtn}
                            onClick={() =>
                              copyToClipboard(
                                payload.username || payload.user || payload.email!,
                                "Username"
                              )
                            }
                            title="Copy Username"
                          >
                            <Copy size={14} />
                          </button>
                        </div>
                      </div>
                    )}

                    {(payload.password || payload.pass || payload.secret) && (
                      <div className={styles.secretField}>
                        <div className={styles.labelRow}>
                          <span className={styles.label}>Password</span>
                        </div>
                        <div className={styles.valueRow}>
                          <span className={`${styles.valueText} ${!showPassword ? styles.valueTextMasked : ""}`}>
                            {showPassword
                              ? payload.password || payload.pass || payload.secret
                              : "••••••••••••••••"}
                          </span>
                          <button
                            type="button"
                            className={styles.actionBtn}
                            onClick={() => setShowPassword(!showPassword)}
                            title={showPassword ? "Hide Password" : "Reveal Password"}
                          >
                            {showPassword ? <EyeOff size={14} /> : <Eye size={14} />}
                          </button>
                          <button
                            type="button"
                            className={styles.actionBtn}
                            onClick={() =>
                              copyToClipboard(
                                payload.password || payload.pass || payload.secret!,
                                "Password"
                              )
                            }
                            title="Copy Password"
                          >
                            <Copy size={14} />
                          </button>
                        </div>
                      </div>
                    )}
                  </>
                )}

                {cred.secretType === "API_KEY" &&
                  (payload.apiKey || payload.api_key || payload.token || payload.key) && (
                    <div className={styles.secretField}>
                      <div className={styles.labelRow}>
                        <span className={styles.label}>API Key / Token</span>
                      </div>
                      <div className={styles.valueRow}>
                        <span className={`${styles.valueText} ${!showApiKey ? styles.valueTextMasked : ""}`}>
                          {showApiKey
                            ? payload.apiKey || payload.api_key || payload.token || payload.key
                            : "••••••••••••••••••••••••••••••••"}
                        </span>
                        <button
                          type="button"
                          className={styles.actionBtn}
                          onClick={() => setShowApiKey(!showApiKey)}
                          title={showApiKey ? "Hide Token" : "Reveal Token"}
                        >
                          {showApiKey ? <EyeOff size={14} /> : <Eye size={14} />}
                        </button>
                        <button
                          type="button"
                          className={styles.actionBtn}
                          onClick={() =>
                            copyToClipboard(
                              (payload.apiKey || payload.api_key || payload.token || payload.key)!,
                              "API Key"
                            )
                          }
                          title="Copy API Key"
                        >
                          <Copy size={14} />
                        </button>
                      </div>
                    </div>
                  )}

                {cred.secretType === "RAW_TEXT" &&
                  (() => {
                    const rawTextValue =
                      payload.rawTextContent ||
                      payload.rawText ||
                      payload.content ||
                      payload.text ||
                      (Object.values(payload)[0] as string) ||
                      "";

                    if (!rawTextValue) return null;

                    return (
                      <div className={styles.secretField}>
                        <div className={styles.labelRow}>
                          <span className={styles.label}>Secret Content (Text / Private Key)</span>
                          <div className="flex items-center gap-1">
                            <button
                              type="button"
                              className={styles.actionBtn}
                              onClick={() => setShowRawText(!showRawText)}
                              title={showRawText ? "Hide Content" : "Reveal Content"}
                            >
                              {showRawText ? <EyeOff size={14} /> : <Eye size={14} />}
                            </button>
                            <button
                              type="button"
                              className={styles.actionBtn}
                              onClick={() => copyToClipboard(rawTextValue, "Secret Content")}
                              title="Copy Content"
                            >
                              <Copy size={14} />
                            </button>
                          </div>
                        </div>
                        <div className={styles.rawTextDisplay}>
                          {showRawText
                            ? rawTextValue
                            : "••••••••••••••••••••••••••••••••\n••••••••••••••••••••••••••••••••\n••••••••••••••••••••••••••••••••"}
                        </div>
                      </div>
                    );
                  })()}
              </div>

              {/* Metadata Section */}
              <div className={styles.metaGroup}>
                {cred.relatedUrl && (
                  <div className={styles.metaItem}>
                    <span className={styles.label}>Related URL</span>
                    <a
                      href={cred.relatedUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={styles.urlLink}
                    >
                      <span>{cred.relatedUrl}</span>
                      <ExternalLink size={13} />
                    </a>
                  </div>
                )}

                {cred.notes && (
                  <div className={styles.metaItem}>
                    <span className={styles.label}>Notes</span>
                    <div className={styles.notesBox}>{cred.notes}</div>
                  </div>
                )}

                {cred.tags && cred.tags.length > 0 && (
                  <div className={styles.metaItem}>
                    <span className={styles.label}>Tags</span>
                    <div className={styles.tagList}>
                      {cred.tags.map((tag) => (
                        <span
                          key={tag.id}
                          className={styles.tagBadge}
                          style={{
                            backgroundColor: `${tag.color || "#8b5cf6"}15`,
                            color: tag.color || "#8b5cf6",
                            borderColor: `${tag.color || "#8b5cf6"}30`,
                          }}
                        >
                          <span
                            style={{
                              width: "6px",
                              height: "6px",
                              borderRadius: "50%",
                              backgroundColor: tag.color || "#8b5cf6",
                            }}
                          />
                          <span>{tag.name}</span>
                        </span>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </>
          )}
        </div>

        <div className={styles.footer}>
          <button type="button" className={styles.btnSecondary} onClick={onClose}>
            Close
          </button>
        </div>
      </div>
    </div>
  );
};
