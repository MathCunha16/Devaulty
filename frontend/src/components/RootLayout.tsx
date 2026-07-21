import React, { Suspense, useState } from "react";
import { Link, Outlet } from "@tanstack/react-router";
import { Toaster } from "sonner";
import { ThemeProvider } from "./ThemeProvider";
import { useTheme } from "../hooks/useTheme";
import { useProjectsQuery } from "~features/projects/hooks/useProjects";
import { ProjectForm } from "~features/projects/components/ProjectForm";
import * as Icons from "lucide-react";
import styles from "../routes/__root.module.css";
import { LogoDevaulty } from "./LogoDevaulty";
import { getIconComponent } from "../utils/icons";

const NavigationSidebar: React.FC = () => {
  const { theme, toggleTheme } = useTheme();
  const { data: projectsData } = useProjectsQuery();
  const [isModalOpen, setIsModalOpen] = useState(false);

  const projects = projectsData?.content || [];
  const activeProjects = projects.filter((p) => !p.archived);

  return (
    <aside className={styles.sidebar}>
      <div className={styles.sidebarHeader}>
        <Link to="/" className={styles.appLogo} title="Devaulty Home">
          <LogoDevaulty height={52} />
        </Link>
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
                    <ProjectIcon size={16} style={{ color: project.color || "var(--color-primary)" }} />
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
  );
};

const RootLayoutInner: React.FC = () => {
  const { theme } = useTheme();

  return (
    <div className={styles.appContainer}>
      {/* Render loading suspense wrapper inside sidebar to load initial projects */}
      <Suspense
        fallback={
          <aside className={styles.sidebar}>
            <div className={styles.sidebarHeader}>
              <div className={styles.appLogo}>
                <LogoDevaulty height={52} />
              </div>
            </div>
            <div className="flex-1 flex items-center justify-center p-4">
              <Icons.Loader2 className="animate-spin text-muted-foreground" size={24} />
            </div>
          </aside>
        }
      >
        <NavigationSidebar />
      </Suspense>

      <main className={styles.mainLayout}>
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
    </div>
  );
};

export const RootLayout: React.FC = () => {
  return (
    <ThemeProvider>
      <RootLayoutInner />
    </ThemeProvider>
  );
};
