import React, { useState } from "react";
import { LockKeyhole, Eye, EyeOff, KeyRound, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useUnlockVaultMutation } from "../hooks/useSecurity";
import styles from "./UnlockVaultCard.module.css";

export const UnlockVaultCard: React.FC = () => {
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const unlockMutation = useUnlockVaultMutation();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!password) return;

    try {
      await unlockMutation.mutateAsync({ masterPassword: password });
      toast.success("Vault unlocked successfully!");
      setPassword("");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Incorrect master password or unlock failed.");
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.header}>
          <div className={styles.iconWrapper}>
            <LockKeyhole size={28} />
          </div>
          <h2 className={styles.title}>VAULT LOCKED</h2>
          <p className={styles.subtitle}>
            Enter your Master Password to unlock the vault and access this project's credentials.
          </p>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label htmlFor="unlock-password" className={styles.label}>Master Password</label>
            <div className={styles.inputWrapper}>
              <input
                id="unlock-password"
                type={showPassword ? "text" : "password"}
                className={styles.input}
                placeholder="Enter Master Password..."
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={unlockMutation.isPending}
                autoFocus
                required
              />
              <button
                type="button"
                className={styles.eyeBtn}
                onClick={() => setShowPassword(!showPassword)}
                aria-label={showPassword ? "Hide master password" : "Show master password"}
              >
                {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
          </div>

          <button
            type="submit"
            className={styles.btnSubmit}
            disabled={unlockMutation.isPending || !password}
          >
            {unlockMutation.isPending ? (
              <>
                <Loader2 size={16} className="animate-spin" />
                Unlocking...
              </>
            ) : (
              <>
                <KeyRound size={16} />
                Unlock Vault
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
};
