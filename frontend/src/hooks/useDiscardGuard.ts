import { useState, useCallback } from "react";

interface UseDiscardGuardOptions {
  isDirty: boolean;
  onClose: () => void;
}

export function useDiscardGuard({ isDirty, onClose }: UseDiscardGuardOptions) {
  const [isConfirmDiscardOpen, setIsConfirmDiscardOpen] = useState(false);

  const handleRequestClose = useCallback(() => {
    if (isDirty) {
      setIsConfirmDiscardOpen(true);
    } else {
      onClose();
    }
  }, [isDirty, onClose]);

  const handleConfirmDiscard = useCallback(() => {
    setIsConfirmDiscardOpen(false);
    onClose();
  }, [onClose]);

  const handleCancelDiscard = useCallback(() => {
    setIsConfirmDiscardOpen(false);
  }, []);

  return {
    isConfirmDiscardOpen,
    setIsConfirmDiscardOpen,
    handleRequestClose,
    handleConfirmDiscard,
    handleCancelDiscard,
  };
}
