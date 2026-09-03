import React, { Suspense, useState, useEffect } from "react";
import { Link, Outlet } from "@tanstack/react-router";
import { Toaster, toast } from "sonner";
import { ThemeProvider } from "./ThemeProvider";
import { useTheme } from "../hooks/useTheme";
import { useProjectsQuery } from "~features/projects/hooks/useProjects";
import { ProjectForm } from "~features/projects/components/ProjectForm";
import { SidebarProvider } from "../contexts/SidebarProvider";
import { useSidebar } from "../hooks/useSidebar";
import * as Icons from "lucide-react";
import styles from "../routes/__root.module.css";
import { HackerLogo } from "./HackerLogo";
import { getIconComponent } from "../utils/icons";
import {
  useCurrentVersionQuery,
  useCheckUpdatesQuery,
  useCheckUpdatesMutation,
} from "~features/releases/hooks/useReleases";
import { UpdateModal } from "~features/releases/components/UpdateModal";
import { formatVersionTag } from "../utils/versionUtils";
import { invoke } from "@tauri-apps/api/core";
import type { AppUpdateInfoResponse } from "~types/api";


// ──────────────────────────────────────────────────────────────
// Edge toggle button — floats half in / half out of the sidebar.
// Uses a CSS custom property so its `left` transitions in sync
// with the sidebar wrapper's width animation.
// ──────────────────────────────────────────────────────────────
const SidebarEdgeButton: React.FC = () => {
  const { isOpen, toggle } = useSidebar();
  const SIDEBAR_WIDTH = 300;
  const BUTTON_HALF = 18; // half of the 36px button

  const leftPx = isOpen ? SIDEBAR_WIDTH - BUTTON_HALF : -BUTTON_HALF;

  return (
    <button
      className={styles.sidebarEdgeBtn}
      style={{ "--sb-left": `${leftPx}px` } as React.CSSProperties}
      onClick={toggle}
      title={isOpen ? "Collapse sidebar" : "Expand sidebar"}
      aria-label={isOpen ? "Collapse sidebar" : "Expand sidebar"}
    >
      {isOpen ? <Icons.ChevronLeft size={15} /> : <Icons.ChevronRight size={15} />}
    </button>
  );
};

// ──────────────────────────────────────────────────────────────
// Main navigation sidebar content
// ──────────────────────────────────────────────────────────────
const NavigationSidebar: React.FC<{
  onOpenUpdateModal: (info: AppUpdateInfoResponse) => void;
}> = ({ onOpenUpdateModal }) => {
  const { theme, toggleTheme } = useTheme();
  const { data: projectsData } = useProjectsQuery();
  const { data: versionData } = useCurrentVersionQuery();
  const checkUpdatesMutation = useCheckUpdatesMutation();

  const { isOpen } = useSidebar();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const projects = projectsData?.content || [];
  const activeProjects = projects.filter((p) => !p.archived);

  const currentAppVersion = versionData?.currentVersion || versionData?.actualVersion;

  const handleManualCheck = async () => {
    try {
      const res = await checkUpdatesMutation.mutateAsync();
      if (res.updateAvailable) {
        onOpenUpdateModal(res);
      } else {
        toast.success(
          `Devaulty is up to date (${formatVersionTag(res.currentVersion || currentAppVersion)})`
        );
      }
    } catch {
      toast.error("Failed to check for application updates.");
    }
  };

  return (
    <div
      className={`${styles.sidebarWrapper} ${isOpen ? "" : styles.sidebarWrapperClosed}`}
      inert={!isOpen ? true : undefined}
    >
      <aside
        className={`${styles.sidebar} ${isOpen ? styles.sidebarOpen : styles.sidebarClosed}`}
      >
        <div className={styles.sidebarHeader}>
          <div className="flex flex-col items-center justify-center gap-1">
            <Link to="/" className={styles.appLogo} title="Devaulty Home">
              <HackerLogo height={46} />
            </Link>
            {currentAppVersion && (
              <button
                type="button"
                onClick={handleManualCheck}
                disabled={checkUpdatesMutation.isPending}
                className="text-[10px] font-mono text-muted-foreground/70 hover:text-foreground transition-colors border border-border/40 px-2 py-0.5 rounded-full bg-secondary/20 cursor-pointer flex items-center gap-1"
                title="Click to check for updates"
              >
                <span>{formatVersionTag(currentAppVersion)}</span>

                {checkUpdatesMutation.isPending && (
                  <Icons.Loader2 size={10} className="animate-spin text-primary" />
                )}
              </button>
            )}
          </div>
        </div>


        <div className={styles.sidebarContent}>
          <div className={styles.navGroup}>
            <span className={styles.navLabel}>System</span>
            <Link
              to="/"
              activeProps={{ className: `${styles.navItem} ${styles.navItemActive}` }}
              inactiveProps={{ className: styles.navItem }}
            >
              <div className={styles.navIconText}>
                <Icons.LayoutDashboard size={16} />
                <span>Dashboard</span>
              </div>
            </Link>
          </div>

          <div className={styles.navGroup}>
            <div className="flex justify-between items-center mb-1 pr-1">
              <span className={styles.navLabel}>Projects</span>
              <button
                onClick={() => setIsModalOpen(true)}
                className="text-xs text-muted-foreground hover:text-foreground hover:bg-secondary p-1 rounded"
                title="Create New Project"
              >
                <Icons.Plus size={12} />
              </button>
            </div>

            {activeProjects.length === 0 ? (
              <div className="text-xs text-muted-foreground px-2 py-4 border border-dashed rounded text-center border-border">
                No active projects
              </div>
            ) : (
              activeProjects.map((project) => {
                const ProjectIcon = getIconComponent(project.icon);
                return (
                  <Link
                    key={project.id}
                    to="/projects/$projectId"
                    params={{ projectId: project.id }}
                    activeProps={{ className: `${styles.navItem} ${styles.navItemActive}` }}
                    inactiveProps={{ className: styles.navItem }}
                  >
                    <div className={styles.navIconText}>
                      <ProjectIcon
                        size={16}
                        style={{ color: project.color || "var(--color-primary)" }}
                      />
                      <span className="truncate max-w-[140px]">{project.name}</span>
                    </div>
                    <div
                      className={styles.projectBadge}
                      style={{ backgroundColor: project.color || "var(--color-primary)" }}
                    />
                  </Link>
                );
              })
            )}
          </div>

          <button className={styles.newProjectButton} onClick={() => setIsModalOpen(true)}>
            <Icons.Plus size={14} />
            <span>New Project</span>
          </button>
        </div>

        <div className={styles.sidebarFooter}>
          <button
            className={styles.themeBtn}
            onClick={toggleTheme}
            title={theme === "light" ? "Switch to Dark Mode" : "Switch to Light Mode"}
          >
            {theme === "light" ? <Icons.Moon size={14} /> : <Icons.Sun size={14} />}
            <span>{theme === "light" ? "Dark Mode" : "Light Mode"}</span>
          </button>
        </div>

        {isModalOpen && (
          <ProjectForm isOpen={true} onClose={() => setIsModalOpen(false)} />
        )}
      </aside>
    </div>
  );
};

// ──────────────────────────────────────────────────────────────
// Root layout inner — consumes theme & sidebar context.
// By the time this mounts, the native session is already ready
// (gated by RootLayout below), so it's safe for its children to
// fire API queries immediately.
// ──────────────────────────────────────────────────────────────
const RootLayoutInner: React.FC = () => {
  const { theme } = useTheme();
  const { isOpen } = useSidebar();
  const { data: autoUpdateData } = useCheckUpdatesQuery(true);

  const [dismissedVersion, setDismissedVersion] = useState<string | null>(null);
  const [manualUpdateInfo, setManualUpdateInfo] = useState<AppUpdateInfoResponse | null>(null);

  const currentUpdateInfo =
    manualUpdateInfo || (autoUpdateData?.updateAvailable ? autoUpdateData : null);

  const isModalOpen =
    !!currentUpdateInfo && currentUpdateInfo.latestVersion !== dismissedVersion;

  const handleOpenUpdateModal = (info: AppUpdateInfoResponse) => {
    setManualUpdateInfo(info);
    setDismissedVersion(null);
  };

  const handleCloseUpdateModal = () => {
    if (currentUpdateInfo) {
      setDismissedVersion(currentUpdateInfo.latestVersion);
    }
    setManualUpdateInfo(null);
  };

  return (
    <div className={styles.appContainer}>
      <Suspense
        fallback={
          <div className={styles.sidebarWrapper}>
            <aside className={`${styles.sidebar} ${styles.sidebarOpen}`}>
              <div className={styles.sidebarHeader}>
                <div className={styles.appLogo}>
                  <HackerLogo height={46} />
                </div>
              </div>
              <div className="flex-1 flex items-center justify-center p-4">
                <Icons.Loader2 className="animate-spin text-muted-foreground" size={24} />
              </div>
            </aside>
          </div>
        }
      >
        <NavigationSidebar onOpenUpdateModal={handleOpenUpdateModal} />
      </Suspense>

      <SidebarEdgeButton />

      <main className={styles.mainLayout}>
        <div className={`${styles.topBar} ${isOpen ? styles.topBarHidden : ""}`}>
          <Link to="/" className={styles.appLogo} title="Devaulty Home">
            <HackerLogo height={44} />
          </Link>
        </div>

        <div className={styles.contentWrapper}>
          <Suspense
            fallback={
              <div className="absolute inset-0 flex items-center justify-center bg-background/50">
                <Icons.Loader2 className="animate-spin text-primary" size={32} />
              </div>
            }
          >
            <Outlet />
          </Suspense>
        </div>
      </main>

      <Toaster position="bottom-right" theme={theme} closeButton />

      <UpdateModal
        isOpen={isModalOpen}
        onClose={handleCloseUpdateModal}
        updateInfo={currentUpdateInfo}
      />
    </div>
  );
};

// ──────────────────────────────────────────────────────────────
// Entry point — resolves the native backend session (port + token)
// BEFORE mounting any provider/component that fires API queries.
// This is the fix for the Windows startup race: previously,
// RootLayoutInner's children rendered immediately while the session
// handshake was still in flight, so the first queries hit the
// hardcoded fallback base URL (wrong port in production).
// ──────────────────────────────────────────────────────────────
export const RootLayout: React.FC = () => {
  const [sessionReady, setSessionReady] = useState(false);

  useEffect(() => {
    const initNativeSession = async () => {
      try {
        const info = await invoke<{ port: number; token: string }>("get_backend_info");
        if (info) {
          window.DEVAULTY_INTERNAL_TOKEN = info.token;
          // Use 127.0.0.1, not "localhost": the backend only binds IPv4,
          // and on Windows "localhost" can resolve to ::1 first.
          window.DEVAULTY_API_BASE_URL = `http://127.0.0.1:${info.port}/api/v1`;
        }
      } catch (err) {
        console.error("Failed to initialize backend native session:", err);
        // Fallback for non-Tauri or dev environment — proceed with
        // whatever default getApiBaseUrl() resolves to.
      } finally {
        setSessionReady(true);
        await invoke("close_splash").catch(() => {});
      }
    };

    void initNativeSession();
  }, []);

  // Splash window is still covering the screen at this point, so there's
  // nothing to render here — just wait until the session is resolved.
  if (!sessionReady) {
    return null;
  }

  return (
    <ThemeProvider>
      <SidebarProvider>
        <RootLayoutInner />
      </SidebarProvider>
    </ThemeProvider>
  );
};