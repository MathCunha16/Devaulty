import React, { useState } from "react";
import { ShieldAlert, Eye, EyeOff, Lock, Check, X, Loader2 } from "lucide-react";
import { toast } from "sonner";
import { useSetupMasterPasswordMutation } from "../hooks/useSecurity";
import styles from "./MasterPasswordSetupCard.module.css";

export const MasterPasswordSetupCard: React.FC = () => {
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);

  const setupMutation = useSetupMasterPasswordMutation();

  // Strength validation checks
  const hasMinLength = password.length >= 8;
  const hasUpper = /[A-Z]/.exec(password) !== null;
  const hasLower = /[a-z]/.exec(password) !== null;
  const hasNumberOrSpecial = /[0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.exec(password) !== null;

  const isStrong = hasMinLength && hasUpper && hasLower && hasNumberOrSpecial;
  const match = password.length > 0 && password === confirmPassword;

  // Calculate strength percentage
  const checksPassed = [hasMinLength, hasUpper, hasLower, hasNumberOrSpecial].filter(Boolean).length;
  const strengthPercent = (checksPassed / 4) * 100;

  const getStrengthColor = () => {
    if (checksPassed <= 1) return "#ef4444";
    if (checksPassed <= 3) return "#eab308";
    return "#10b981";
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isStrong) {
      toast.error("Password must meet all security requirements.");
      return;
    }
    if (!match) {
      toast.error("The passwords entered do not match.");
      return;
    }

    try {
      await setupMutation.mutateAsync({ masterPassword: password });
      toast.success("Master password configured successfully!");
      setPassword("");
      setConfirmPassword("");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to configure master password.");
    }
  };

  return (
    <div className={styles.container}>
      <div className={styles.card}>
        <div className={styles.header}>
          <div className={styles.iconWrapper}>
            <Lock size={28} />
          </div>
          <h2 className={styles.title}>CONFIGURE VAULT MASTER PASSWORD</h2>
          <p className={styles.subtitle}>
            To store and manage encrypted credentials securely, you need to define a Master Password for your Vault.
          </p>
        </div>

        <div className={styles.warningBox}>
          <ShieldAlert className={styles.warningIcon} size={22} />
          <div className={styles.warningText}>
            <span className={styles.warningHighlight}>CRITICAL WARNING:</span> This Master Password <strong style={{ color: "#fff" }}>APPLIES TO ALL OF YOUR PROJECTS</strong>. Do not forget it under any circumstances! Due to zero-knowledge encryption, if you forget this password, <strong style={{ color: "#fff" }}>there is no way to recover your saved secrets</strong>.
          </div>
        </div>

        <form onSubmit={handleSubmit} className={styles.form}>
          <div className={styles.field}>
            <label className={styles.label}>New Master Password</label>
            <div className={styles.inputWrapper}>
              <input
                type={showPassword ? "text" : "password"}
                className={styles.input}
                placeholder="Enter a strong password (minimum 8 characters)..."
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={setupMutation.isPending}
                required
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

            {/* Strength meter bar */}
            {password.length > 0 && (
              <div className={styles.strengthMeter}>
                <div className={styles.strengthTrack}>
                  <div
                    className={styles.strengthBar}
                    style={{
                      width: `${strengthPercent}%`,
                      backgroundColor: getStrengthColor(),
                    }}
                  />
                </div>
              </div>
            )}

            {/* Password requirements */}
            <div className={styles.requirements}>
              <div className={`${styles.reqItem} ${hasMinLength ? styles.reqItemValid : ""}`}>
                {hasMinLength ? <Check size={12} /> : <X size={12} />}
                <span>Minimum 8 characters</span>
              </div>
              <div className={`${styles.reqItem} ${hasUpper ? styles.reqItemValid : ""}`}>
                {hasUpper ? <Check size={12} /> : <X size={12} />}
                <span>Uppercase letter (A-Z)</span>
              </div>
              <div className={`${styles.reqItem} ${hasLower ? styles.reqItemValid : ""}`}>
                {hasLower ? <Check size={12} /> : <X size={12} />}
                <span>Lowercase letter (a-z)</span>
              </div>
              <div className={`${styles.reqItem} ${hasNumberOrSpecial ? styles.reqItemValid : ""}`}>
                {hasNumberOrSpecial ? <Check size={12} /> : <X size={12} />}
                <span>Number or symbol</span>
              </div>
            </div>
          </div>

          <div className={styles.field}>
            <label className={styles.label}>Confirm Master Password</label>
            <div className={styles.inputWrapper}>
              <input
                type={showConfirm ? "text" : "password"}
                className={styles.input}
                placeholder="Repeat the master password..."
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
                disabled={setupMutation.isPending}
                required
              />
              <button
                type="button"
                className={styles.eyeBtn}
                onClick={() => setShowConfirm(!showConfirm)}
                tabIndex={-1}
              >
                {showConfirm ? <EyeOff size={16} /> : <Eye size={16} />}
              </button>
            </div>
            {confirmPassword.length > 0 && (
              <div className="text-xs font-mono mt-1">
                {match ? (
                  <span className="text-emerald-500 flex items-center gap-1">
                    <Check size={12} /> Passwords match
                  </span>
                ) : (
                  <span className="text-destructive flex items-center gap-1">
                    <X size={12} /> Passwords do not match
                  </span>
                )}
              </div>
            )}
          </div>

          <button
            type="submit"
            className={styles.btnSubmit}
            disabled={setupMutation.isPending || !isStrong || !match}
          >
            {setupMutation.isPending ? (
              <>
                <Loader2 size={16} className="animate-spin" />
                Configuring Vault...
              </>
            ) : (
              <>
                <Lock size={16} />
                Initialize Vault Master Password
              </>
            )}
          </button>
        </form>
      </div>
    </div>
  );
};
