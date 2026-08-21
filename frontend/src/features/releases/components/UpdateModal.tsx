import React, { useState, useEffect, useRef } from "react";
import {
  DownloadCloud,
  CheckCircle2,
  AlertCircle,
  X,
  Sparkles,
  ArrowRight,
  Loader2,
  FolderOpen,
  Info,
} from "lucide-react";
import { marked } from "marked";
import DOMPurify from "dompurify";
import type {
  AppUpdateInfoResponse,
  UpdateDownloadProgressResponse,
  ReleaseAssetFormat,
} from "../../../types/api";
import { releasesApi } from "../api/releasesApi";
import styles from "./UpdateModal.module.css";

interface UpdateModalProps {
  isOpen: boolean;
  onClose: () => void;
  updateInfo: AppUpdateInfoResponse | null;
}

interface UpdateModalContentProps {
  onClose: () => void;
  updateInfo: AppUpdateInfoResponse;
}

const UpdateModalContent: React.FC<UpdateModalContentProps> = ({
  onClose,
  updateInfo,
}) => {
  const defaultFormat =
    updateInfo.availableFormats?.find(
      (f) => f.packageType === updateInfo.packageType
    ) ||
    updateInfo.availableFormats?.[0] ||
    null;

  const isStandaloneMode = updateInfo.supportsInPlaceUpdate === false;

  const [isDownloading, setIsDownloading] = useState(false);
  const [progress, setProgress] = useState<UpdateDownloadProgressResponse | null>(null);
  const [streamError, setStreamError] = useState<string | null>(null);
  const [restartCountdown, setRestartCountdown] = useState<number | null>(null);
  const [selectedFormat, setSelectedFormat] = useState<ReleaseAssetFormat | null>(defaultFormat);
  const [savedFilePath, setSavedFilePath] = useState<string | null>(null);

  const cleanupRef = useRef<(() => void) | null>(null);
  const modalRef = useRef<HTMLDivElement | null>(null);

  // Handle countdown and auto-relaunch
  useEffect(() => {
    let timer: ReturnType<typeof setTimeout>;
    if (restartCountdown !== null && restartCountdown > 0) {
      timer = setTimeout(() => {
        setRestartCountdown((prev) => (prev !== null ? prev - 1 : null));
      }, 1000);
    } else if (restartCountdown === 0) {
      releasesApi.relaunchApp().catch((err) => {
        console.error("Failed to relaunch application:", err);
      });
    }
    return () => clearTimeout(timer);
  }, [restartCountdown]);

  // Clean up background update stream on unmount
  useEffect(() => {
    return () => {
      if (cleanupRef.current) {
        cleanupRef.current();
        cleanupRef.current = null;
      }
    };
  }, []);

  // Keyboard accessibility and focus management
  useEffect(() => {
    const previouslyFocused = document.activeElement as HTMLElement | null;
    modalRef.current?.focus();

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape" && !isDownloading && progress?.status !== "INSTALLING") {
        onClose();
      }
    };

    document.addEventListener("keydown", handleKeyDown);
    return () => {
      document.removeEventListener("keydown", handleKeyDown);
      previouslyFocused?.focus();
    };
  }, [isDownloading, progress?.status, onClose]);

  // Handler for In-Place Auto-Update (AppImage / Windows / macOS)
  const handleStartInPlaceUpdate = () => {
    setIsDownloading(true);
    setStreamError(null);
    setRestartCountdown(null);
    setProgress({
      status: "DOWNLOADING",
      percentage: 0,
      downloadedBytes: 0,
      totalBytes: updateInfo.downloadSizeInBytes || 0,
    });

    if (cleanupRef.current) {
      cleanupRef.current();
    }

    cleanupRef.current = releasesApi.downloadAndInstall(
      (data) => {
        setProgress(data);
        if (data.status === "COMPLETED") {
          setRestartCountdown(2);
        } else if (data.status === "FAILED") {
          setIsDownloading(false);
          setProgress(null);
          setStreamError(data.errorMessage || "Update download failed.");
        }
      },
      (errorMsg) => {
        setIsDownloading(false);
        setProgress(null);
        setStreamError(errorMsg);
      }
    );
  };

  // Handler for Standalone Installer Download (.deb, .rpm, .exe, .dmg)
  const handleStartStandaloneDownload = () => {
    if (!selectedFormat) return;

    setIsDownloading(true);
    setStreamError(null);
    setSavedFilePath(null);
    setProgress({
      status: "DOWNLOADING",
      percentage: 0,
      downloadedBytes: 0,
      totalBytes: 0,
    });

    if (cleanupRef.current) {
      cleanupRef.current();
    }

    cleanupRef.current = releasesApi.downloadStandaloneInstaller(
      selectedFormat.url,
      selectedFormat.filename,
      (data) => {
        setProgress(data);
        if (data.status === "COMPLETED") {
          setIsDownloading(false);
          if (data.savedFilePath) {
            setSavedFilePath(data.savedFilePath);
          }
        } else if (data.status === "FAILED") {
          setIsDownloading(false);
          setProgress(null);
          setStreamError(data.errorMessage || "Installer download failed.");
        }
      },
      (errorMsg) => {
        setIsDownloading(false);
        setProgress(null);
        setStreamError(errorMsg);
      }
    );
  };

  const handleOpenFile = () => {
    if (savedFilePath) {
      releasesApi.openDownloadedFile(savedFilePath).catch((err) => {
        console.error("Failed to open file path:", err);
      });
    }
  };

  const formatBytes = (bytes?: number) => {
    if (!bytes || isNaN(bytes) || bytes <= 0) return "0 MB";
    const mb = bytes / (1024 * 1024);
    return `${mb.toFixed(1)} MB`;
  };

  const formatVersionTag = (v: string) => {
    if (!v) return "";
    return v.startsWith("v") ? v : `v${v}`;
  };

  const renderReleaseNotes = (notes: string) => {
    try {
      const rawHtml = marked.parse(notes, { async: false }) as string;
      return DOMPurify.sanitize(rawHtml);
    } catch {
      return notes;
    }
  };

  const isInstallingOrCompleted =
    progress?.status === "INSTALLING" ||
    (progress?.status === "COMPLETED" && !savedFilePath);

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
                A new version of Devaulty is ready to download.
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
              Current: {formatVersionTag(updateInfo.currentVersion)}
            </span>
            <ArrowRight size={14} className={styles.versionArrow} />
            <span className={`${styles.versionPill} ${styles.versionNew}`}>
              New: {formatVersionTag(updateInfo.latestVersion)}
            </span>
          </div>

          {/* Linux system package banner */}
          {isStandaloneMode && !savedFilePath && (
            <div className={styles.infoBanner}>
              <Info size={18} className="flex-shrink-0 text-indigo-400 mt-0.5" />
              <div>
                <strong>System Package Detected (.deb / .rpm)</strong>
                <div className="text-xs opacity-90 font-normal mt-0.5">
                  Direct in-place replacement requires system administrator privileges.
                  Click below to download the updated package to your Downloads folder.
                </div>
              </div>
            </div>
          )}

          {/* Format selector */}
          {updateInfo.availableFormats &&
            updateInfo.availableFormats.length > 0 &&
            !isDownloading &&
            !savedFilePath && (
              <div className={styles.formatSelectRow}>
                <label className={styles.formatSelectLabel} htmlFor="package-format-select">
                  Installer Package Format
                </label>
                <select
                  id="package-format-select"
                  className={styles.formatSelect}
                  value={selectedFormat?.packageType || ""}
                  onChange={(e) => {
                    const matched = updateInfo.availableFormats?.find(
                      (f) => f.packageType === e.target.value
                    );
                    if (matched) setSelectedFormat(matched);
                  }}
                >
                  {updateInfo.availableFormats.map((f) => (
                    <option key={f.packageType} value={f.packageType}>
                      {f.label} ({f.filename})
                    </option>
                  ))}
                </select>
              </div>
            )}

          {/* Release Notes */}
          {updateInfo.releaseNotes && !isDownloading && !savedFilePath && (
            <div className={styles.releaseNotesBox}>
              <span className={styles.releaseNotesLabel}>
                {updateInfo.releaseTitle ? updateInfo.releaseTitle : "RELEASE NOTES"}
              </span>
              <div
                className={styles.releaseNotesContent}
                dangerouslySetInnerHTML={{ __html: renderReleaseNotes(updateInfo.releaseNotes) }}
              />
            </div>
          )}

          {/* Stream Error Box */}
          {streamError && (
            <div className={styles.errorBox}>
              <AlertCircle size={18} className="flex-shrink-0" />
              <div className="flex flex-col gap-1 text-xs">
                <span>{streamError}</span>
                {!isStandaloneMode && (
                  <span className="opacity-80">
                    You can switch to downloading the standalone package (.deb / .AppImage) below.
                  </span>
                )}
              </div>
            </div>
          )}

          {/* Download Progress Bar Section */}
          {isDownloading && (
            <div className={styles.progressSection}>
              <div className={styles.progressHeader}>
                <div className={styles.progressStatus}>
                  {progress?.status === "COMPLETED" ? (
                    <>
                      <Sparkles size={16} className="text-emerald-400 animate-pulse" />
                      <span>DOWNLOAD COMPLETE!</span>
                    </>
                  ) : progress?.status === "INSTALLING" ? (
                    <>
                      <Loader2 size={16} className="animate-spin text-amber-400" />
                      <span>INSTALLING UPDATE...</span>
                    </>
                  ) : (
                    <>
                      <DownloadCloud size={16} className="animate-pulse text-indigo-400" />
                      <span>DOWNLOADING {selectedFormat?.filename || "UPDATE"}...</span>
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
                  {formatBytes(progress?.downloadedBytes)}
                  {progress?.totalBytes ? ` / ${formatBytes(progress.totalBytes)}` : ""}
                </span>
                <span>
                  {progress?.status === "COMPLETED"
                    ? "Completed"
                    : progress?.status === "INSTALLING"
                    ? "Applying update..."
                    : `${progress?.percentage || 0}% Completed`}
                </span>
              </div>
            </div>
          )}

          {/* Saved to Downloads Success Banner */}
          {savedFilePath && (
            <div className={styles.restartAlert}>
              <CheckCircle2 size={20} className="flex-shrink-0 text-emerald-400" />
              <div>
                <strong>Installer downloaded successfully!</strong>
                <div className="text-xs opacity-90 font-normal mt-0.5 break-all">
                  Saved to: {savedFilePath}
                </div>
              </div>
            </div>
          )}

          {/* In-place In Progress Notice */}
          {progress?.status === "INSTALLING" && (
            <div className={styles.restartAlert}>
              <Loader2 size={20} className="flex-shrink-0 animate-spin text-amber-400" />
              <div>
                <strong>Installing update...</strong>
                <div className="text-xs opacity-90 font-normal mt-0.5">
                  Devaulty is applying the new version.
                </div>
              </div>
            </div>
          )}

          {/* In-place Restart Notification Notice */}
          {progress?.status === "COMPLETED" && !savedFilePath && (
            <div className={styles.restartAlert}>
              <CheckCircle2 size={20} className="flex-shrink-0 text-emerald-400" />
              <div>
                <strong>Update successfully installed!</strong>
                <div className="text-xs opacity-90 font-normal mt-0.5">
                  {restartCountdown !== null && restartCountdown > 0
                    ? `Devaulty will restart automatically in ${restartCountdown} second${restartCountdown > 1 ? "s" : ""}...`
                    : "Devaulty is restarting now..."}
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
                {savedFilePath ? "Close" : "Remind Me Later"}
              </button>

              {savedFilePath ? (
                <button
                  type="button"
                  className={styles.btnSuccess}
                  onClick={handleOpenFile}
                >
                  <FolderOpen size={16} />
                  Open in File Manager
                </button>
              ) : isStandaloneMode || streamError ? (
                <button
                  type="button"
                  className={styles.btnPrimary}
                  onClick={handleStartStandaloneDownload}
                >
                  <DownloadCloud size={16} />
                  Download {selectedFormat?.label.split("(")[1]?.replace(")", "") || "Package"}
                </button>
              ) : (
                <button
                  type="button"
                  className={styles.btnPrimary}
                  onClick={handleStartInPlaceUpdate}
                >
                  <Sparkles size={16} />
                  Download & Install Update
                </button>
              )}
            </>
          ) : isInstallingOrCompleted ? (
            <div className="flex items-center gap-2 text-xs font-mono text-muted-foreground py-1">
              <Loader2 size={14} className="animate-spin text-emerald-400" />
              <span>
                {progress?.status === "COMPLETED"
                  ? restartCountdown !== null && restartCountdown > 0
                    ? `Restarting in ${restartCountdown}s...`
                    : "Restarting application..."
                  : "Installing update..."}
              </span>
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

export const UpdateModal: React.FC<UpdateModalProps> = ({
  isOpen,
  onClose,
  updateInfo,
}) => {
  if (!isOpen || !updateInfo) return null;

  return (
    <UpdateModalContent
      key={updateInfo.latestVersion}
      onClose={onClose}
      updateInfo={updateInfo}
    />
  );
};
