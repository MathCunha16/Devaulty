import { useEffect, useRef, useCallback } from "react";
import { toast } from "sonner";
import { useLockVaultMutation } from "~features/security/hooks/useSecurity";

const INACTIVITY_TIMEOUT_MS = 15 * 60 * 1000; // 15 minutes
const THROTTLE_MS = 5000; // Throttle timer reset to max once every 5 seconds

export const useInactivityAutoLock = (enabled: boolean, onLockTriggered?: () => void) => {
  const lockMutation = useLockVaultMutation();
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const lastResetRef = useRef<number>(0);
  const lastActivityRef = useRef<number>(0);
  const resetTimerRef = useRef<() => void>(() => {});

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

    const elapsed = Date.now() - (lastActivityRef.current || Date.now());
    const delay = Math.max(INACTIVITY_TIMEOUT_MS - elapsed, 0);

    timerRef.current = setTimeout(async () => {
      // Re-check exact elapsed time from actual last user activity
      const actualElapsed = Date.now() - (lastActivityRef.current || Date.now());
      if (actualElapsed < INACTIVITY_TIMEOUT_MS) {
        resetTimerRef.current();
        return;
      }

      try {
        await mutateAsync();
        toast.warning("Vault automatically locked after 15 minutes of inactivity.");
        if (onLockTriggeredRef.current) {
          onLockTriggeredRef.current();
        }
      } catch {
        // Ignore lock errors if session already ended
      }
    }, delay);
  }, [enabled, mutateAsync]);

  useEffect(() => {
    resetTimerRef.current = resetTimer;
  }, [resetTimer]);

  useEffect(() => {
    if (!enabled) {
      if (timerRef.current) {
        clearTimeout(timerRef.current);
      }
      return;
    }

    lastActivityRef.current = Date.now();
    resetTimer();

    const activityEvents = ["mousemove", "keydown", "click", "scroll", "touchstart"];
    const handleActivity = () => {
      const now = Date.now();
      lastActivityRef.current = now;
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
