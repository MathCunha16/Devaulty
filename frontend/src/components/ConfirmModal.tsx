import React, { useEffect, useRef } from "react";
import { AlertTriangle, Trash2 } from "lucide-react";
import styles from "./ConfirmModal.module.css";

interface ConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void | Promise<void>;
  title: string;
  message: string;
  itemName?: string;
  warningText?: string;
  isLoading?: boolean;
  confirmLabel?: string;
}

export const ConfirmModal: React.FC<ConfirmModalProps> = ({
  isOpen,
  onClose,
  onConfirm,
  title,
  message,
  itemName,
  warningText,
  isLoading = false,
  confirmLabel = "Delete",
}) => {
  const previousActiveElement = useRef<HTMLElement | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);
  const firstButtonRef = useRef<HTMLButtonElement>(null);

  useEffect(() => {
    if (isOpen) {
      previousActiveElement.current = document.activeElement as HTMLElement;
      setTimeout(() => {
        firstButtonRef.current?.focus();
      }, 0);
    }

    return () => {
      if (previousActiveElement.current) {
        previousActiveElement.current.focus();
      }
    };
  }, [isOpen]);

  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
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
  }, [isOpen]);

  if (!isOpen) return null;

  const handleConfirm = async () => {
    await onConfirm();
  };

  return (
    <div
      className={styles.overlay}
      onClick={() => !isLoading && onClose()}
    >
      <div
        ref={modalRef}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="confirm-modal-title"
      >
        <div className={styles.body}>
          <div className={styles.iconWrap}>
            <AlertTriangle size={22} color="#ef4444" />
          </div>
          <h2 id="confirm-modal-title" className={styles.title}>{title}</h2>
          <p className={styles.message}>
            {message}{" "}
            {itemName && <span className={styles.highlight}>"{itemName}"</span>}?
          </p>
          {warningText && <p className={styles.warning}>{warningText}</p>}
        </div>

        <div className={styles.footer}>
          <button
            ref={firstButtonRef}
            className={styles.btnCancel}
            onClick={onClose}
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            className={styles.btnDelete}
            onClick={handleConfirm}
            disabled={isLoading}
          >
            <Trash2 size={13} />
            <span>{isLoading ? "Deleting..." : confirmLabel}</span>
          </button>
        </div>
      </div>
    </div>
  );
};
