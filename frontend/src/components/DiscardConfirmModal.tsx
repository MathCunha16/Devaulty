import React from "react";
import { ConfirmModal } from "./ConfirmModal";

interface DiscardConfirmModalProps {
  isOpen: boolean;
  onClose: () => void;
  onDiscard: () => void;
  itemName?: string;
  message?: string;
}

export const DiscardConfirmModal: React.FC<DiscardConfirmModalProps> = ({
  isOpen,
  onClose,
  onDiscard,
  itemName = "item",
  message,
}) => {
  return (
    <ConfirmModal
      isOpen={isOpen}
      onClose={onClose}
      onConfirm={onDiscard}
      title="Discard Unsaved Changes"
      message={
        message ||
        `You have unsaved modifications on this ${itemName}. Are you sure you want to discard your changes and close?`
      }
      warningText="Any unsaved changes will be permanently lost."
      confirmLabel="Discard Changes"
    />
  );
};
