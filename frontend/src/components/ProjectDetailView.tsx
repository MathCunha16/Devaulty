import React, { useState } from "react";
import { Link, useMatch } from "@tanstack/react-router";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import Editor from "@monaco-editor/react";
import { useTheme } from "../hooks/useTheme";
import { useProjectQuery } from "~features/projects/hooks/useProjects";
import {
  useSnippetsQuery,
  useCreateSnippetMutation,
  useDeleteSnippetMutation,
} from "~features/snippets/hooks/useSnippets";
import type { SnippetLanguage, SnippetType } from "~types/api";
import styles from "../routes/projects.$projectId.module.css";

// Preset list of popular languages for the editor select dropdown
const POPULAR_LANGUAGES: SnippetLanguage[] = [
  "PLAIN_TEXT",
  "JAVASCRIPT",
  "TYPESCRIPT",
  "PYTHON",
  "GO",
  "RUST",
  "JAVA",
  "KOTLIN",
  "CSHARP",
  "CPP",
  "C",
  "BASH",
  "ZSH",
  "SH",
  "POWERSHELL",
  "HTML",
  "CSS",
  "SQL",
  "JSON",
  "YAML",
  "XML",
  "MARKDOWN",
  "DOCKERFILE",
  "DOCKER_COMPOSE",
  "NGINX",
  "TERRAFORM",
  "GRADLE",
  "MAVEN_POM",
  "REGEX",
  "DIFF",
  "LOG",
];

const mapLanguageToMonaco = (lang: SnippetLanguage): string => {
  switch (lang) {
    case "JAVASCRIPT":
    case "JSX":
    case "MONGODB":
      return "javascript";
    case "TYPESCRIPT":
    case "TSX":
      return "typescript";
    case "PYTHON":
      return "python";
    case "GO":
      return "go";
    case "RUST":
      return "rust";
    case "JAVA":
      return "java";
    case "KOTLIN":
      return "kotlin";
    case "CSHARP":
      return "csharp";
    case "CPP":
      return "cpp";
    case "C":
      return "c";
    case "BASH":
    case "FISH":
    case "ZSH":
    case "SH":
      return "shell";
    case "POWERSHELL":
      return "powershell";
    case "BATCH":
      return "bat";
    case "HTML":
    case "VUE":
    case "SVELTE":
      return "html";
    case "CSS":
      return "css";
    case "SCSS":
      return "scss";
    case "LESS":
      return "less";
    case "JSON":
      return "json";
    case "YAML":
    case "DOCKER_COMPOSE":
    case "KUBERNETES_YAML":
      return "yaml";
    case "XML":
    case "MAVEN_POM":
      return "xml";
    case "TOML":
    case "INI":
    case "ENV":
    case "PROPERTIES":
      return "ini";
    case "MARKDOWN":
      return "markdown";
    case "DOCKERFILE":
      return "dockerfile";
    case "SQL":
    case "PLSQL":
      return "sql";
    case "GRAPHQL":
      return "graphql";
    case "DIFF":
      return "diff";
    default:
      return "plaintext";
  }
};

// Helper to resolve project icon components
const getIconComponent = (iconName?: string) => {
  if (!iconName) return Icons.Folder;
  const IconComponent = (Icons as unknown as Record<string, React.ComponentType<{ size?: number; className?: string; style?: React.CSSProperties }>>)[iconName];
  return IconComponent || Icons.Folder;
};

export const ProjectDetailView: React.FC = () => {
  const { projectId } = useParamsHelper();
  const { theme } = useTheme();

  // Load project details
  const { data: project } = useProjectQuery(projectId);

  // Load snippets details
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const createSnippetMutation = useCreateSnippetMutation(projectId);
  const deleteSnippetMutation = useDeleteSnippetMutation(projectId);

  // Local UI States
  const [selectedSnippetId, setSelectedSnippetId] = useState<string | undefined>(undefined);
  const [isEditing, setIsEditing] = useState(false);
  const [isCreating, setIsCreating] = useState(false);
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Search & Filter state
  const [searchQuery, setSearchQuery] = useState("");
  const [typeFilter, setTypeFilter] = useState<"ALL" | SnippetType>("ALL");

  // Form states
  const [formTitle, setFormTitle] = useState("");
  const [formDescription, setFormDescription] = useState("");
  const [formContent, setFormContent] = useState("");
  const [formLanguage, setFormLanguage] = useState<SnippetLanguage>("PLAIN_TEXT");
  const [formSnippetType, setFormSnippetType] = useState<SnippetType>("CODE");

  const snippets = snippetsData?.content || [];

  // Filter logic
  const filteredSnippets = snippets.filter((s) => {
    const matchesSearch =
      s.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      (s.description && s.description.toLowerCase().includes(searchQuery.toLowerCase())) ||
      s.content.toLowerCase().includes(searchQuery.toLowerCase());
    const matchesType = typeFilter === "ALL" || s.snippetType === typeFilter;
    return matchesSearch && matchesType;
  });

  const selectedSnippet = snippets.find((s) => s.id === selectedSnippetId);

  // Trigger copy
  const handleCopy = async (content: string, id: string) => {
    try {
      await navigator.clipboard.writeText(content);
      setCopiedId(id);
      toast.success("Code copied to clipboard!");
      setTimeout(() => setCopiedId(null), 2000);
    } catch {
      toast.error("Failed to copy code");
    }
  };

  // Open Form for creating
  const handleOpenCreate = () => {
    setIsCreating(true);
    setIsEditing(false);
    setFormTitle("");
    setFormDescription("");
    setFormContent("");
    setFormLanguage("PLAIN_TEXT");
    setFormSnippetType("CODE");
  };

  // Open Form for editing
  const handleOpenEdit = () => {
    if (!selectedSnippet) return;
    setIsEditing(true);
    setIsCreating(false);
    setFormTitle(selectedSnippet.title);
    setFormDescription(selectedSnippet.description || "");
    setFormContent(selectedSnippet.content);
    setFormLanguage(selectedSnippet.language);
    setFormSnippetType(selectedSnippet.snippetType);
  };

  // Delete snippet
  const handleDeleteSnippet = async (snippetId: string, title: string) => {
    if (window.confirm(`Delete snippet "${title}"?`)) {
      try {
        await deleteSnippetMutation.mutateAsync(snippetId);
        toast.success("Snippet deleted successfully");
        if (selectedSnippetId === snippetId) {
          setSelectedSnippetId(undefined);
        }
      } catch (err) {
        toast.error(err instanceof Error ? err.message : "Failed to delete snippet");
      }
    }
  };

  const projectIcon = getIconComponent(project?.icon);

  return (
    <div className="flex flex-col gap-4 h-full">
      <div className={styles.projectHeader}>
        <div className={styles.projectTitleSection}>
          <div className={styles.projectIcon}>
            {React.createElement(projectIcon, { size: 22, style: { color: project?.color || "var(--color-primary)" } })}
          </div>
          <div>
            <h1 className={styles.projectTitle}>{project?.name}</h1>
            {project?.description && <p className={styles.projectDesc}>{project.description}</p>}
          </div>
        </div>

        <Link
          to="/"
          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors"
        >
          <Icons.ArrowLeft size={12} />
          <span>Back to Dashboard</span>
        </Link>
      </div>

      <div className={styles.pageLayout}>
        {/* Left Side: Snippets navigation list */}
        <div className={styles.leftPanel}>
          <button className={styles.newSnippetBtn} onClick={handleOpenCreate}>
            <Icons.Plus size={14} />
            <span>Add Snippet</span>
          </button>

          <div className={styles.searchBar}>
            <Icons.Search className={styles.searchIcon} size={14} />
            <input
              type="text"
              placeholder="Search snippets..."
              className={styles.searchInput}
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
            />
          </div>

          <div className={styles.filterTabs}>
            <button
              className={`${styles.filterTab} ${typeFilter === "ALL" ? styles.filterTabActive : ""}`}
              onClick={() => setTypeFilter("ALL")}
            >
              ALL
            </button>
            <button
              className={`${styles.filterTab} ${typeFilter === "CODE" ? styles.filterTabActive : ""}`}
              onClick={() => setTypeFilter("CODE")}
            >
              CODE
            </button>
            <button
              className={`${styles.filterTab} ${typeFilter === "COMMAND" ? styles.filterTabActive : ""}`}
              onClick={() => setTypeFilter("COMMAND")}
            >
              CMD
            </button>
          </div>

          <div className={styles.snippetList}>
            {filteredSnippets.length === 0 ? (
              <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border">
                No snippets found
              </div>
            ) : (
              filteredSnippets.map((s) => (
                <button
                  key={s.id}
                  className={`${styles.snippetItem} ${
                    selectedSnippetId === s.id && !isCreating ? styles.snippetItemActive : ""
                  }`}
                  onClick={() => {
                    setSelectedSnippetId(s.id);
                    setIsCreating(false);
                    setIsEditing(false);
                  }}
                >
                  <div className={styles.snippetItemHeader}>
                    <span className={styles.snippetItemTitle}>{s.title}</span>
                    <span className="text-[10px] text-muted-foreground font-mono">
                      {new Date(s.createdAt).toLocaleDateString()}
                    </span>
                  </div>
                  {s.description && <p className={styles.snippetItemDesc}>{s.description}</p>}
                  <div className={styles.snippetItemFooter}>
                    <span className={styles.badgeType}>{s.snippetType}</span>
                    <span className={styles.badgeLang}>{s.language}</span>
                  </div>
                </button>
              ))
            )}
          </div>
        </div>

        {/* Right Side: Code detail or editor form */}
        <div className={styles.rightPanel}>
          {isCreating || isEditing ? (
            <form
              onSubmit={async (e) => {
                e.preventDefault();
                if (!formTitle.trim() || !formContent.trim()) {
                  toast.error("Title and Content are required");
                  return;
                }
                const payload = {
                  title: formTitle,
                  description: formDescription || undefined,
                  content: formContent,
                  language: formLanguage,
                  snippetType: formSnippetType,
                };

                try {
                  if (isCreating) {
                    const res = await createSnippetMutation.mutateAsync(payload);
                    toast.success("Snippet created successfully");
                    setSelectedSnippetId(res.id);
                  } else if (isEditing && selectedSnippetId) {
                    const { snippetsApi } = await import("~features/snippets/api/snippetsApi");
                    await snippetsApi.update(projectId, selectedSnippetId, payload);
                    toast.success("Snippet updated successfully");
                  }
                  setIsCreating(false);
                  setIsEditing(false);
                  window.dispatchEvent(new CustomEvent("snippet-saved"));
                } catch (err) {
                  toast.error(err instanceof Error ? err.message : "Failed to save snippet");
                }
              }}
              className={styles.formScroll}
            >
              <h2 className={styles.formTitle}>
                {isCreating ? "CREATE NEW SNIPPET" : "EDIT SNIPPET"}
              </h2>

              <div className={styles.formFields}>
                <div className={styles.formField}>
                  <label className={styles.formLabel}>Title</label>
                  <input
                    type="text"
                    className={styles.formInput}
                    placeholder="e.g., Get user list, Run database backup"
                    value={formTitle}
                    onChange={(e) => setFormTitle(e.target.value)}
                    required
                  />
                </div>

                <div className={styles.formField}>
                  <label className={styles.formLabel}>Description</label>
                  <input
                    type="text"
                    className={styles.formInput}
                    placeholder="Provide a brief context or description..."
                    value={formDescription}
                    onChange={(e) => setFormDescription(e.target.value)}
                  />
                </div>

                <div className={styles.formRow}>
                  <div className={styles.formField}>
                    <label className={styles.formLabel}>Type</label>
                    <select
                      className={styles.formInput}
                      value={formSnippetType}
                      onChange={(e) => setFormSnippetType(e.target.value as SnippetType)}
                    >
                      <option value="CODE">Code Snippet</option>
                      <option value="COMMAND">Terminal Command</option>
                    </select>
                  </div>

                  <div className={styles.formField}>
                    <label className={styles.formLabel}>Language</label>
                    <select
                      className={styles.formInput}
                      value={formLanguage}
                      onChange={(e) => setFormLanguage(e.target.value as SnippetLanguage)}
                    >
                      {POPULAR_LANGUAGES.map((lang) => (
                        <option key={lang} value={lang}>
                          {lang}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="border border-border rounded overflow-hidden bg-background">
                  <Editor
                    height="260px"
                    language={mapLanguageToMonaco(formLanguage)}
                    theme={theme === "dark" ? "vs-dark" : "light"}
                    value={formContent}
                    onChange={(value) => setFormContent(value || "")}
                    loading={
                      <div className="flex items-center justify-center h-48 text-xs text-muted-foreground font-mono">
                        Loading Editor Environment...
                      </div>
                    }
                    options={{
                      minimap: { enabled: false },
                      fontSize: 13,
                      lineNumbers: "on",
                      scrollBeyondLastLine: false,
                      automaticLayout: true,
                      padding: { top: 8, bottom: 8 },
                    }}
                  />
                </div>
              </div>

              <div className={styles.formActions}>
                <button
                  type="button"
                  className={styles.btn}
                  onClick={() => {
                    setIsCreating(false);
                    setIsEditing(false);
                  }}
                >
                  Cancel
                </button>
                <button type="submit" className={styles.btnPrimary}>
                  Save
                </button>
              </div>
            </form>
          ) : selectedSnippet ? (
            <div className={styles.snippetDetailScroll}>
              <div className={styles.snippetDetailContainer}>
                <div className={styles.detailHeader}>
                  <div className={styles.detailTitleSection}>
                    <div className={styles.detailTitleRow}>
                      <h2 className={styles.detailTitle}>{selectedSnippet.title}</h2>
                      <div className="flex gap-1.5">
                        <span className={styles.badgeType}>{selectedSnippet.snippetType}</span>
                        <span className={styles.badgeLang}>{selectedSnippet.language}</span>
                      </div>
                    </div>
                    {selectedSnippet.description && (
                      <p className={styles.detailDesc}>{selectedSnippet.description}</p>
                    )}
                  </div>

                  <div className={styles.detailActions}>
                    <button className={styles.btnIcon} onClick={handleOpenEdit} title="Edit Snippet">
                      <Icons.Edit3 size={14} />
                    </button>
                    <button
                      className={`${styles.btnIcon} ${styles.btnIconDanger}`}
                      onClick={() =>
                        handleDeleteSnippet(selectedSnippet.id, selectedSnippet.title)
                      }
                      title="Delete Snippet"
                    >
                      <Icons.Trash2 size={14} />
                    </button>
                  </div>
                </div>

                <div className={styles.codePanel}>
                  <div className={styles.codePanelHeader}>
                    <span className={styles.codePanelLabel}>
                      Source Code ({selectedSnippet.language.toLowerCase()})
                    </span>
                    <button
                      className={styles.copyButton}
                      onClick={() => handleCopy(selectedSnippet.content, selectedSnippet.id)}
                    >
                      {copiedId === selectedSnippet.id ? (
                        <>
                          <Icons.Check size={12} className="text-emerald-500" />
                          <span className="text-emerald-500 font-bold">Copied!</span>
                        </>
                      ) : (
                        <>
                          <Icons.Copy size={12} />
                          <span>Copy code</span>
                        </>
                      )}
                    </button>
                  </div>

                  <div className="border border-border rounded overflow-hidden bg-[#0b0b0c]">
                    <Editor
                      height="380px"
                      language={mapLanguageToMonaco(selectedSnippet.language)}
                      theme={theme === "dark" ? "vs-dark" : "light"}
                      value={selectedSnippet.content}
                      loading={
                        <div className="flex items-center justify-center h-48 text-xs text-muted-foreground font-mono">
                          Loading Code Preview...
                        </div>
                      }
                      options={{
                        readOnly: true,
                        minimap: { enabled: false },
                        fontSize: 13,
                        lineNumbers: "on",
                        scrollBeyondLastLine: false,
                        automaticLayout: true,
                        domReadOnly: true,
                        padding: { top: 8, bottom: 8 },
                      }}
                    />
                  </div>
                </div>
              </div>
            </div>
          ) : (
            <div className={styles.placeholder}>
              <Icons.Terminal size={48} className="text-muted-foreground animate-pulse" />
              <div className={styles.placeholderText}>
                No snippet selected. Choose a snippet from the list or click "Add Snippet" to
                create a new one.
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

// Helper to fetch route parameters because of versioning issues inside file route scope
const useParamsHelper = () => {
  const match = useMatch({ from: "/projects/$projectId" });
  return match.params;
};

// Query Client hook retriever
import { useQueryClient } from "@tanstack/react-query";
const useQueryClientHelper = () => {
  return useQueryClient();
};

export const ProjectDetailRouteComponent: React.FC = () => {
  const queryClient = useQueryClientHelper();
  const { projectId } = useParamsHelper();
  React.useEffect(() => {
    const handleSave = () => {
      queryClient.invalidateQueries({ queryKey: ["projects", projectId, "snippets"] });
    };
    window.addEventListener("snippet-saved", handleSave);
    return () => window.removeEventListener("snippet-saved", handleSave);
  }, [projectId, queryClient]);

  return <ProjectDetailView />;
};
