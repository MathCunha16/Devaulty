import React, { useState, useEffect, useRef } from "react";
import {
  DownloadCloud,
  Sparkles,
  ArrowRight,
  X,
  AlertCircle,
  CheckCircle2,
  Loader2,
  RefreshCw,
} from "lucide-react";
import type {
  AppUpdateInfoResponse,
  UpdateDownloadProgressResponse,
} from "~types/api";
import { releasesApi } from "../api/releasesApi";
import styles from "./UpdateModal.module.css";

interface UpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
  updateInfo: AppUpdateInfoResponse | null;
}

const formatBytes = (bytes?: number): string => {
  if (!bytes || bytes <= 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
};

export const UpdateModal: React.FC<UpdateModalProps> = ({
  isOpen,
  onClose,
  updateInfo,
}) => {
  const [isDownloading, setIsDownloading] = useState(false);
  const [progress, setProgress] = useState<UpdateDownloadProgressResponse | null>(null);
  const [streamError, setStreamError] = useState<string | null>(null);

  const cleanupRef = useRef<(() => void) | null>(null);
  const modalRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
      }
    };
  }, []);

  // Trap focus & escape listener
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !isDownloading && progress?.status !== "INSTALLING") {
        onClose();
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => document.removeEventListener("keydown", handleKeyDown);
  }, [isOpen, isDownloading, progress?.status, onClose]);

  if (!isOpen || !updateInfo) return null;

  const handleStartUpdate = () => {
    setIsDownloading(true);
    setStreamError(null);
    setProgress({
      status: "DOWNLOADING",
      percentage: 0,
      downloadedBytes: 0,
      totalBytes: updateInfo.downloadSizeInBytes || 0,
    });

    if (cleanupRef.current) {
      cleanupRef.current();
    }

    cleanupRef.current = releasesApi.streamDownloadAndInstall(
      (data) => {
        setProgress(data);
        if (data.status === "FAILED") {
          setIsDownloading(false);
          setStreamError(data.errorMessage || "Update download failed.");
        }
      },
      (errorMsg) => {
        setIsDownloading(false);
        setStreamError(errorMsg);
      }
    );
  };

  const isInstallingOrCompleted =
    progress?.status === "INSTALLING" || progress?.status === "COMPLETED";

  return (
    <div
      className={styles.overlay}
      onClick={() => {
        if (!isDownloading && !isInstallingOrCompleted) {
          onClose();
        }
      }}
    >
      <div
        ref={modalRef}
        tabIndex={-1}
        className={styles.modal}
        onClick={(e) => e.stopPropagation()}
        role="dialog"
        aria-modal="true"
        aria-labelledby="update-modal-title"
      >
        {/* Header */}
        <div className={styles.header}>
          <div className={styles.headerTitleSection}>
            <div className={styles.iconBadge}>
              <DownloadCloud size={24} />
            </div>
            <div className={styles.titleGroup}>
              <h2 id="update-modal-title" className={styles.title}>
                SOFTWARE UPDATE AVAILABLE
              </h2>
              <p className={styles.subtitle}>
                A new version of Devaulty is ready to download and install.
              </p>
            </div>
          </div>
          {!isDownloading && !isInstallingOrCompleted && (
            <button
              className={styles.closeBtn}
              onClick={onClose}
              aria-label="Close modal"
            >
              <X size={18} />
            </button>
          )}
        </div>

        {/* Body */}
        <div className={styles.body}>
          {/* Version badge row */}
          <div className={styles.versionBadgeRow}>
            <span className={`${styles.versionPill} ${styles.versionCurrent}`}>
              Current: v{updateInfo.currentVersion}
            </span>
            <ArrowRight size={14} className={styles.versionArrow} />
            <span className={`${styles.versionPill} ${styles.versionNew}`}>
              New: v{updateInfo.latestVersion}
            </span>
          </div>

          {/* Release Notes */}
          {updateInfo.releaseNotes && !isDownloading && !isInstallingOrCompleted && (
            <div className={styles.releaseNotesBox}>
              <span className={styles.releaseNotesLabel}>
                {updateInfo.releaseTitle ? updateInfo.releaseTitle : "RELEASE NOTES"}
              </span>
              <div className={styles.releaseNotesContent}>
                {updateInfo.releaseNotes}
              </div>
            </div>
          )}

          {/* Stream Error Box */}
          {streamError && (
            <div className={styles.errorBox}>
              <AlertCircle size={18} className="flex-shrink-0" />
              <span>{streamError}</span>
            </div>
          )}

          {/* Download Progress Bar Section */}
          {isDownloading && (
            <div className={styles.progressSection}>
              <div className={styles.progressHeader}>
                <div className={styles.progressStatus}>
                  {isInstallingOrCompleted ? (
                    <>
                      <Loader2 size={16} className="animate-spin text-emerald-400" />
                      <span>INSTALLING UPDATE...</span>
                    </>
                  ) : (
                    <>
                      <DownloadCloud size={16} className="animate-pulse text-indigo-400" />
                      <span>DOWNLOADING INSTALLER...</span>
                    </>
                  )}
                </div>
                <span className={styles.progressPercent}>
                  {progress?.percentage || 0}%
                </span>
              </div>

              <div className={styles.track}>
                <div
                  className={styles.fill}
                  style={{ width: `${Math.min(progress?.percentage || 0, 100)}%` }}
                />
              </div>

              <div className={styles.progressSubInfo}>
                <span>
                  {formatBytes(progress?.downloadedBytes)} /{" "}
                  {formatBytes(progress?.totalBytes || updateInfo.downloadSizeInBytes)}
                </span>
                <span>
                  {isInstallingOrCompleted
                    ? "Executing Native Installer..."
                    : `${progress?.percentage || 0}% Completed`}
                </span>
              </div>
            </div>
          )}

          {/* Restart Notification Notice */}
          {isInstallingOrCompleted && (
            <div className={styles.restartAlert}>
              <CheckCircle2 size={20} className="flex-shrink-0" />
              <div>
                <strong>Update successfully downloaded!</strong>
                <div className="text-xs opacity-90 font-normal mt-0.5">
                  Devaulty is now applying the update and will restart automatically in a few seconds.
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Footer Actions */}
        <div className={styles.footer}>
          {!isDownloading && !isInstallingOrCompleted ? (
            <>
              <button
                type="button"
                className={styles.btnSecondary}
                onClick={onClose}
              >
                Remind Me Later
              </button>
              <button
                type="button"
                className={styles.btnPrimary}
                onClick={handleStartUpdate}
              >
                {streamError ? (
                  <>
                    <RefreshCw size={16} />
                    Retry Download & Install
                  </>
                ) : (
                  <>
                    <Sparkles size={16} />
                    Download & Install Update
                  </>
                )}
              </button>
            </>
          ) : isInstallingOrCompleted ? (
            <div className="flex items-center gap-2 text-xs font-mono text-muted-foreground py-1">
              <Loader2 size={14} className="animate-spin text-emerald-400" />
              <span>Preparing restart...</span>
            </div>
          ) : (
            <button
              type="button"
              className={styles.btnSecondary}
              onClick={() => {
                if (cleanupRef.current) cleanupRef.current();
                setIsDownloading(false);
                setProgress(null);
              }}
            >
              Cancel Download
            </button>
          )}
        </div>
      </div>
    </div>
  );
};
