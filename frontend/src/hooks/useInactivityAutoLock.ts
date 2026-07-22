import { useEffect, useRef, useCallback } from "react";
import { toast } from "sonner";
import { useLockVaultMutation } from "~features/security/hooks/useSecurity";

const INACTIVITY_TIMEOUT_MS = 15 * 60 * 1000; // 15 minutes

export const useInactivityAutoLock = (enabled: boolean, onLockTriggered?: () => void) => {
  const lockMutation = useLockVaultMutation();
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const resetTimer = useCallback(() => {
    if (timerRef.current) {
      clearTimeout(timerRef.current);
    }
    if (!enabled) return;

    timerRef.current = setTimeout(async () => {
      try {
        await lockMutation.mutateAsync();
        toast.warning("Vault automatically locked after 15 minutes of inactivity.");
        if (onLockTriggered) {
          onLockTriggered();
        }
      } catch {
        // Ignore lock errors if session already ended
      }
    }, INACTIVITY_TIMEOUT_MS);
  }, [enabled, lockMutation, onLockTriggered]);

  useEffect(() => {
    if (!enabled) {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      return;
    }

    // Start timer
    resetTimer();

    const activityEvents = ["mousemove", "keydown", "click", "scroll", "touchstart"];
    const handleActivity = () => {
      resetTimer();
    };

    activityEvents.forEach((event) => {
      window.addEventListener(event, handleActivity, { passive: true });
    });

    return () => {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      activityEvents.forEach((event) => {
        window.removeEventListener(event, handleActivity);
      });
    };
  }, [enabled, resetTimer]);
};
