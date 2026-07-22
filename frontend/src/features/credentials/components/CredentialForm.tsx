import React, { useState, useEffect, useRef } from "react";
import { X, Eye, EyeOff, KeyRound, UserCheck, Code2, Loader2 } from "lucide-react";
import { toast } from "sonner";
import {
  useCreateCredentialMutation,
  useUpdateCredentialMutation,
  useCredentialQuery,
} from "../hooks/useCredentials";
import type { CredentialSecretType, CreateCredentialRequest, UpdateCredentialRequest } from "~types/api";
import styles from "./CredentialForm.module.css";

interface CredentialFormProps {
  isOpen: boolean;
  onClose: () => void;
  projectId: string;
  credentialId?: string;
}

interface CredentialFormInnerProps {
  title: string;
  initialValues?: {
    title: string;
    secretType: CredentialSecretType;
    username?: string;
    password?: string;
    apiKey?: string;
    rawTextContent?: string;
    notes?: string;
    relatedUrl?: string;
  };
  onSubmit: (values: CreateCredentialRequest) => Promise<void>;
  onClose: () => void;
  isSubmitting: boolean;
}

const CredentialFormInner: React.FC<CredentialFormInnerProps> = ({
  title,
  initialValues,
  onSubmit,
  onClose,
  isSubmitting,
}) => {
  const [formTitle, setFormTitle] = useState(initialValues?.title || "");
  const [secretType, setSecretType] = useState<CredentialSecretType>(
    initialValues?.secretType || "LOGIN"
  );
  const [username, setUsername] = useState(initialValues?.username || "");
  const [password, setPassword] = useState(initialValues?.password || "");
  const [apiKey, setApiKey] = useState(initialValues?.apiKey || "");
  const [rawTextContent, setRawTextContent] = useState(initialValues?.rawTextContent || "");
  const [notes, setNotes] = useState(initialValues?.notes || "");
  const [relatedUrl, setRelatedUrl] = useState(initialValues?.relatedUrl || "");

  const [showPassword, setShowPassword] = useState(false);
  const [showApiKey, setShowApiKey] = useState(false);

  // Focus trap refs
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const firstInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    previousActiveElement.current = document.activeElement as HTMLElement;
    setTimeout(() => {
      firstInputRef.current?.focus();
    }, 50);

    return () => {
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, []);

  // Escape key & focus trap
  useEffect(() => {
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
  }, [onClose]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formTitle.trim()) {
      toast.error("Title is required");
      return;
    }

    const payload: CreateCredentialRequest = {
      title: formTitle.trim(),
      secretType,
      notes: notes.trim() || undefined,
      relatedUrl: relatedUrl.trim() || undefined,
    };

    if (secretType === "LOGIN") {
      payload.username = username.trim() || undefined;
      payload.password = password || undefined;
    } else if (secretType === "API_KEY") {
      payload.apiKey = apiKey.trim() || undefined;
    } else if (secretType === "RAW_TEXT") {
      payload.rawTextContent = rawTextContent || undefined;
    }

    onSubmit(payload);
  };

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div
        ref={modalRef}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="credential-form-title"
      >
        <div className={styles.header}>
          <h2 id="credential-form-title" className={styles.title}>{title}</h2>
          <button className={styles.closeBtn} onClick={onClose} disabled={isSubmitting} aria-label="Close modal">
            <X size={16} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="cred-title" className={styles.label}>Credential Title</label>
            <input
              ref={firstInputRef}
              id="cred-title"
              type="text"
              className={styles.input}
              placeholder="e.g., GitHub Machine Token, Database Prod, AWS Admin..."
              value={formTitle}
              onChange={(e) => setFormTitle(e.target.value)}
              disabled={isSubmitting}
              maxLength={255}
              required
            />
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Secret Type</label>
            <div className={styles.typeSelector}>
              <button
                type="button"
                className={`${styles.typeBtn} ${secretType === "LOGIN" ? styles.typeBtnActive : ""}`}
                onClick={() => setSecretType("LOGIN")}
                disabled={isSubmitting}
              >
                <UserCheck size={16} />
                <span>LOGIN</span>
              </button>
              <button
                type="button"
                className={`${styles.typeBtn} ${secretType === "API_KEY" ? styles.typeBtnActive : ""}`}
                onClick={() => setSecretType("API_KEY")}
                disabled={isSubmitting}
              >
                <KeyRound size={16} />
                <span>API KEY</span>
              </button>
              <button
                type="button"
                className={`${styles.typeBtn} ${secretType === "RAW_TEXT" ? styles.typeBtnActive : ""}`}
                onClick={() => setSecretType("RAW_TEXT")}
                disabled={isSubmitting}
              >
                <Code2 size={16} />
                <span>TEXT / KEY</span>
              </button>
            </div>
          </div>

          {/* Conditional Secret Fields */}
          {secretType === "LOGIN" && (
            <>
              <div className={styles.field}>
                <label htmlFor="cred-username" className={styles.label}>Username / Email</label>
                <input
                  id="cred-username"
                  type="text"
                  className={styles.input}
                  placeholder="e.g., admin@company.com, root, deployer..."
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  disabled={isSubmitting}
                />
              </div>

              <div className={styles.field}>
                <label htmlFor="cred-password" className={styles.label}>Password</label>
                <div className={styles.inputWrapper}>
                  <input
                    id="cred-password"
                    type={showPassword ? "text" : "password"}
                    className={styles.input}
                    placeholder="Enter password to be encrypted..."
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    disabled={isSubmitting}
                  />
                  <button
                    type="button"
                    className={styles.eyeBtn}
                    onClick={() => setShowPassword(!showPassword)}
                    tabIndex={-1}
                  >
                    {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
                  </button>
                </div>
              </div>
            </>
          )}

          {secretType === "API_KEY" && (
            <div className={styles.field}>
              <label htmlFor="cred-apikey" className={styles.label}>API Key / Token</label>
              <div className={styles.inputWrapper}>
                <input
                  id="cred-apikey"
                  type={showApiKey ? "text" : "password"}
                  className={styles.input}
                  placeholder="e.g., sk_live_51M..., ghp_xxxx..., secret_key..."
                  value={apiKey}
                  onChange={(e) => setApiKey(e.target.value)}
                  disabled={isSubmitting}
                />
                <button
                  type="button"
                  className={styles.eyeBtn}
                  onClick={() => setShowApiKey(!showApiKey)}
                  tabIndex={-1}
                >
                  {showApiKey ? <EyeOff size={16} /> : <Eye size={16} />}
                </button>
              </div>
            </div>
          )}

          {secretType === "RAW_TEXT" && (
            <div className={styles.field}>
              <label htmlFor="cred-rawtext" className={styles.label}>Secret Content (Text / Private Key)</label>
              <textarea
                id="cred-rawtext"
                className={styles.textarea}
                placeholder="---BEGIN RSA PRIVATE KEY---&#10;Paste your private key, certificate, or environment variables here..."
                value={rawTextContent}
                onChange={(e) => setRawTextContent(e.target.value)}
                disabled={isSubmitting}
              />
            </div>
          )}

          <div className={styles.field}>
            <label htmlFor="cred-url" className={styles.label}>Related URL (Optional)</label>
            <input
              id="cred-url"
              type="url"
              className={styles.input}
              placeholder="https://console.aws.amazon.com, https://github.com..."
              value={relatedUrl}
              onChange={(e) => setRelatedUrl(e.target.value)}
              disabled={isSubmitting}
            />
          </div>

          <div className={styles.field}>
            <label htmlFor="cred-notes" className={styles.label}>Additional Notes (Optional)</label>
            <textarea
              id="cred-notes"
              className={styles.textarea}
              placeholder="Rotation instructions, staging environment notes..."
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              disabled={isSubmitting}
              style={{ minHeight: "70px" }}
            />
          </div>

          <div className={styles.footer}>
            <button
              type="button"
              className={styles.btn}
              onClick={onClose}
              disabled={isSubmitting}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.btnPrimary}
              disabled={isSubmitting}
            >
              {title.includes("EDIT") ? "Save Changes" : "Create Credential"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

const CreateCredentialModal: React.FC<{ projectId: string; onClose: () => void }> = ({
  projectId,
  onClose,
}) => {
  const createMutation = useCreateCredentialMutation(projectId);

  const handleSubmit = async (values: CreateCredentialRequest) => {
    try {
      await createMutation.mutateAsync(values);
      toast.success("Credential created successfully!");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to create credential.");
    }
  };

  return (
    <CredentialFormInner
      title="NEW VAULT CREDENTIAL"
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={createMutation.isPending}
    />
  );
};

const EditCredentialModal: React.FC<{
  projectId: string;
  credentialId: string;
  onClose: () => void;
}> = ({ projectId, credentialId, onClose }) => {
  const { data: cred, isLoading, isError } = useCredentialQuery(projectId, credentialId);
  const updateMutation = useUpdateCredentialMutation(projectId, credentialId);

  const handleSubmit = async (values: CreateCredentialRequest) => {
    try {
      const updatePayload: UpdateCredentialRequest = {
        title: values.title,
        secretType: values.secretType,
        username: values.username,
        password: values.password,
        apiKey: values.apiKey,
        rawTextContent: values.rawTextContent,
        notes: values.notes,
        relatedUrl: values.relatedUrl,
      };

      await updateMutation.mutateAsync(updatePayload);
      toast.success("Credential updated successfully!");
      onClose();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : "Failed to update credential.");
    }
  };

  if (isLoading) {
    return (
      <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3">
            <Loader2 className="animate-spin text-primary" size={28} />
            <span className="text-xs text-muted-foreground font-mono">LOADING CREDENTIAL...</span>
          </div>
        </div>
      </div>
    );
  }

  if (isError || !cred) {
    return (
      <div className={styles.overlay} onClick={onClose}>
        <div className={styles.modal}>
          <div className="flex flex-col items-center justify-center p-12 gap-3 text-destructive font-mono text-xs">
            <span>FAILED TO LOAD VAULT CREDENTIAL.</span>
          </div>
        </div>
      </div>
    );
  }

  const payload = cred.decryptedPayload || {};

  const initialValues = {
    title: cred.title || "",
    secretType: cred.secretType || "LOGIN",
    username: payload.username || payload.user || payload.email || "",
    password: payload.password || payload.pass || payload.secret || "",
    apiKey: payload.apiKey || payload.api_key || payload.token || payload.key || "",
    rawTextContent:
      payload.rawTextContent ||
      payload.rawText ||
      payload.content ||
      payload.text ||
      (Object.values(payload)[0] as string) ||
      "",
    notes: cred.notes || "",
    relatedUrl: cred.relatedUrl || "",
  };

  return (
    <CredentialFormInner
      title="EDIT VAULT CREDENTIAL"
      initialValues={initialValues}
      onSubmit={handleSubmit}
      onClose={onClose}
      isSubmitting={updateMutation.isPending}
    />
  );
};

export const CredentialForm: React.FC<CredentialFormProps> = ({
  isOpen,
  onClose,
  projectId,
  credentialId,
}) => {
  if (!isOpen) return null;

  if (credentialId) {
    return <EditCredentialModal projectId={projectId} credentialId={credentialId} onClose={onClose} />;
  }

  return <CreateCredentialModal projectId={projectId} onClose={onClose} />;
};
