import React from "react";
import { Lock, Unlock, Clock, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useLockVaultMutation } from "../hooks/useSecurity";
import styles from "./VaultSecurityBanner.module.css";

interface VaultSecurityBannerProps {
  secondsLeft?: number;
}

export const VaultSecurityBanner: React.FC<VaultSecurityBannerProps> = ({ secondsLeft }) => {
  const lockMutation = useLockVaultMutation();

  const handleLock = async () => {
    try {
      await lockMutation.mutateAsync();
      toast.info("Vault locked.");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to lock vault.");
    }
  };

  const formatTimer = (totalSeconds?: number) => {
    if (totalSeconds === undefined || totalSeconds < 0) return "--:--";
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes.toString().padStart(2, "0")}:${seconds.toString().padStart(2, "0")}`;
  };

  return (
    <div className={styles.banner}>
      <div className={styles.statusInfo}>
        <div className={styles.statusBadge}>
          <Unlock size={14} />
          <span>VAULT UNLOCKED</span>
        </div>
        <div className={styles.timer}>
          <Clock size={14} />
          <span>Session expires in: </span>
          <span className={styles.timerValue}>{formatTimer(secondsLeft)}</span>
        </div>
      </div>

      <button
        type="button"
        className={styles.lockBtn}
        onClick={handleLock}
        disabled={lockMutation.isPending}
        title="Lock vault immediately"
      >
        {lockMutation.isPending ? (
          <Loader2 size={14} className="animate-spin" />
        ) : (
          <Lock size={14} />
        )}
        <span>Lock Vault</span>
      </button>
    </div>
  );
};
