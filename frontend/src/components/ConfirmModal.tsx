import React from "react";
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
  if (!isOpen) return null;

  const handleConfirm = async () => {
    await onConfirm();
  };

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <div className={styles.body}>
          <div className={styles.iconWrap}>
            <AlertTriangle size={22} color="#ef4444" />
          </div>
          <h2 className={styles.title}>{title}</h2>
          <p className={styles.message}>
            {message}{" "}
            {itemName && <span className={styles.highlight}>"{itemName}"</span>}?
          </p>
          {warningText && <p className={styles.warning}>{warningText}</p>}
        </div>

        <div className={styles.footer}>
          <button
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
