import { useEffect, useRef, useCallback } from "react";
import { toast } from "sonner";
import { useLockVaultMutation } from "~features/security/hooks/useSecurity";

const INACTIVITY_TIMEOUT_MS = 15 * 60 * 1000; // 15 minutes
const THROTTLE_MS = 5000; // Throttle timer reset to max once every 5 seconds

export const useInactivityAutoLock = (enabled: boolean, onLockTriggered?: () => void) => {
  const lockMutation = useLockVaultMutation();
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastResetRef = useRef<number>(0);

  const onLockTriggeredRef = useRef(onLockTriggered);
  useEffect(() => {
    onLockTriggeredRef.current = onLockTriggered;
  }, [onLockTriggered]);

  const mutateAsync = lockMutation.mutateAsync;

  const resetTimer = useCallback(() => {
    lastResetRef.current = Date.now();
    if (timerRef.current) {
      clearTimeout(timerRef.current);
    }
    if (!enabled) return;

    timerRef.current = setTimeout(async () => {
      try {
        await mutateAsync();
        toast.warning("Vault automatically locked after 15 minutes of inactivity.");
        if (onLockTriggeredRef.current) {
          onLockTriggeredRef.current();
        }
      } catch {
        // Ignore lock errors if session already ended
      }
    }, INACTIVITY_TIMEOUT_MS);
  }, [enabled, mutateAsync]);

  useEffect(() => {
    if (!enabled) {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      return;
    }

    // Start timer initially
    resetTimer();

    const activityEvents = ["mousemove", "keydown", "click", "scroll", "touchstart"];
    const handleActivity = () => {
      const now = Date.now();
      if (now - lastResetRef.current >= THROTTLE_MS) {
        resetTimer();
      }
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
