import React, { useState, useEffect } from "react";
import { Link, useMatch } from "@tanstack/react-router";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import Editor from "@monaco-editor/react";
import { useTheme } from "../hooks/useTheme";
import { useAutoResize } from "../hooks/useAutoResize";
import { useSidebar } from "../hooks/useSidebar";
import { ConfirmModal } from "./ConfirmModal";
import { useProjectQuery } from "~features/projects/hooks/useProjects";
import {
  useSnippetsQuery,
  useCreateSnippetMutation,
  useUpdateSnippetMutation,
  useDeleteSnippetMutation,
} from "~features/snippets/hooks/useSnippets";
import {
  useProblemsQuery,
  useProblemQuery,
  useUpdateProblemStatusMutation,
  useDeleteProblemMutation,
} from "~features/problems/hooks/useProblems";
import {
  useTagsQuery,
  useCreateTagMutation,
  useAssociateTagMutation,
  useDisassociateTagMutation,
} from "~features/tags/hooks/useTags";
import { ProblemForm } from "~features/problems/components/ProblemForm";
import type { SnippetLanguage, SnippetType, ProblemStatus, ProblemSeverity, TagSummaryResponse } from "~types/api";
import styles from "../routes/projects.$projectId.module.css";
import { getIconComponent } from "../utils/icons";

const ALL_LANGUAGES: SnippetLanguage[] = [
  "PLAIN_TEXT",
  "BASH",
  "FISH",
  "ZSH",
  "SH",
  "POWERSHELL",
  "BATCH",
  "JAVA",
  "KOTLIN",
  "JAVASCRIPT",
  "TYPESCRIPT",
  "PYTHON",
  "GO",
  "RUST",
  "C",
  "CPP",
  "CSHARP",
  "PHP",
  "RUBY",
  "SWIFT",
  "DART",
  "SCALA",
  "LUA",
  "PERL",
  "R",
  "ELIXIR",
  "HASKELL",
  "CLOJURE",
  "GROOVY",
  "HTML",
  "CSS",
  "SCSS",
  "LESS",
  "JSX",
  "TSX",
  "VUE",
  "SVELTE",
  "JSON",
  "YAML",
  "XML",
  "TOML",
  "INI",
  "ENV",
  "CSV",
  "MARKDOWN",
  "PROPERTIES",
  "DOCKERFILE",
  "DOCKER_COMPOSE",
  "NGINX",
  "APACHE",
  "TERRAFORM",
  "ANSIBLE",
  "KUBERNETES_YAML",
  "HELM",
  "MAKEFILE",
  "CMAKE",
  "GRADLE",
  "MAVEN_POM",
  "SQL",
  "PLSQL",
  "GRAPHQL",
  "MONGODB",
  "GITHUB_ACTIONS",
  "GITLAB_CI",
  "JENKINSFILE",
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

export const ProjectDetailView: React.FC = () => {
  const { projectId } = useParamsHelper();
  const { theme } = useTheme();
  const { close: closeSidebar } = useSidebar();

  // Auto-close sidebar when entering a project
  useEffect(() => {
    closeSidebar();
  }, [projectId, closeSidebar]);

  // Load project core details
  const { data: project } = useProjectQuery(projectId);

  // Load snippets details
  const { data: snippetsData } = useSnippetsQuery(projectId);
  const createSnippetMutation = useCreateSnippetMutation(projectId);
  const deleteSnippetMutation = useDeleteSnippetMutation(projectId);

  // Load problems details
  const { data: problemsData } = useProblemsQuery(projectId);
  const updateProblemStatusMutation = useUpdateProblemStatusMutation(projectId);
  const deleteProblemMutation = useDeleteProblemMutation(projectId);

  // Load tags details
  const { data: tagsData = [] } = useTagsQuery(projectId);
  const createTagMutation = useCreateTagMutation(projectId);
  const associateTagMutation = useAssociateTagMutation(projectId);
  const disassociateTagMutation = useDisassociateTagMutation(projectId);

  // Workspace sub-navigation state
  const [activeTab, setActiveTab] = useState<"snippets" | "problems">("snippets");

  // Selected item states
  const [selectedSnippetId, setSelectedSnippetId] = useState<string | undefined>(undefined);
  const [selectedProblemId, setSelectedProblemId] = useState<string | undefined>(undefined);

  // Edit / Creation states
  const [isEditingSnippet, setIsEditingSnippet] = useState(false);
  const [isCreatingSnippet, setIsCreatingSnippet] = useState(false);
  const [isProblemFormOpen, setIsProblemFormOpen] = useState(false);
  const [editingProblemId, setEditingProblemId] = useState<string | undefined>(undefined);

  // ConfirmModal states
  const [confirmModal, setConfirmModal] = useState<{
    isOpen: boolean;
    title: string;
    message: string;
    itemName?: string;
    warningText?: string;
    onConfirm: () => Promise<void>;
    isLoading: boolean;
  }>({
    isOpen: false,
    title: "",
    message: "",
    onConfirm: async () => {},
    isLoading: false,
  });

  const closeConfirmModal = () =>
    setConfirmModal((prev) => ({ ...prev, isOpen: false, isLoading: false }));

  // Copy status state
  const [copiedId, setCopiedId] = useState<string | null>(null);

  // Tags filter UI state
  const [showTagPopoverId, setShowTagPopoverId] = useState<string | null>(null);
  const [tagSearchQuery, setTagSearchQuery] = useState("");

  // Snippet update mutation
  const updateSnippetMutation = useUpdateSnippetMutation(projectId, selectedSnippetId || "");

  // Auto-resize for snippet description textarea
  const [formDescription, setFormDescription] = useState("");
  const descTextareaRef = useAutoResize(formDescription, 60);

  // Search & Filters (Snippets)
  const [snippetSearchQuery, setSnippetSearchQuery] = useState("");
  const [snippetTypeFilter, setSnippetTypeFilter] = useState<"ALL" | SnippetType>("ALL");

  // Search & Filters (Problems)
  const [problemSearchQuery, setProblemSearchQuery] = useState("");
  const [problemSeverityFilter, setProblemSeverityFilter] = useState<"ALL" | ProblemSeverity>("ALL");
  const [problemStatusFilter, setProblemStatusFilter] = useState<"ALL" | "UNRESOLVED" | ProblemStatus>("UNRESOLVED");

  // Form states (Snippets)
  const [formTitle, setFormTitle] = useState("");
  const [formContent, setFormContent] = useState("");
  const [formLanguage, setFormLanguage] = useState<SnippetLanguage>("PLAIN_TEXT");
  const [formSnippetType, setFormSnippetType] = useState<SnippetType>("CODE");

  const snippets = snippetsData?.content || [];
  const problems = problemsData?.content || [];

  // Filter calculations
  const filteredSnippets = snippets.filter((s) => {
    const matchesSearch =
      s.title.toLowerCase().includes(snippetSearchQuery.toLowerCase()) ||
      (s.description && s.description.toLowerCase().includes(snippetSearchQuery.toLowerCase())) ||
      s.content.toLowerCase().includes(snippetSearchQuery.toLowerCase());
    const matchesType = snippetTypeFilter === "ALL" || s.snippetType === snippetTypeFilter;
    return matchesSearch && matchesType;
  });

  const filteredProblems = problems.filter((p) => {
    const matchesSearch =
      p.title.toLowerCase().includes(problemSearchQuery.toLowerCase()) ||
      (p.tags && p.tags.some(t => t.name.toLowerCase().includes(problemSearchQuery.toLowerCase())));
    const matchesSeverity = problemSeverityFilter === "ALL" || p.severity === problemSeverityFilter;
    
    let matchesStatus = true;
    if (problemStatusFilter === "UNRESOLVED") {
      matchesStatus = p.status === "OPEN" || p.status === "WORKING_ON";
    } else if (problemStatusFilter !== "ALL") {
      matchesStatus = p.status === problemStatusFilter;
    }
    
    return matchesSearch && matchesSeverity && matchesStatus;
  });

  const selectedSnippet = snippets.find((s) => s.id === selectedSnippetId);
  const selectedProblem = problems.find((p) => p.id === selectedProblemId);

  // Load detail view model for the selected problem (holds errorDescription and solution)
  const { data: problemDetail, isLoading: isLoadingProblemDetail } = useProblemQuery(
    projectId,
    selectedProblemId || ""
  );

  // Active / Open problems count for badge
  const openProblemsCount = problems.filter(p => p.status === "OPEN" || p.status === "WORKING_ON").length;

  // Copy handler
  const handleCopy = async (content: string, id: string) => {
    try {
      await navigator.clipboard.writeText(content);
      setCopiedId(id);
      toast.success("Copied to clipboard!");
      setTimeout(() => setCopiedId(null), 2000);
    } catch {
      toast.error("Failed to copy content");
    }
  };

  // Tags association operations
  const handleAddTag = async (itemId: string, itemType: "SNIPPET" | "PROBLEM", tagId: string) => {
    try {
      await associateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag associated successfully");
      setShowTagPopoverId(null);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to associate tag");
    }
  };

  const handleRemoveTag = async (itemId: string, itemType: "SNIPPET" | "PROBLEM", tagId: string) => {
    try {
      await disassociateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag removed successfully");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to remove tag");
    }
  };

  const handleCreateAndAddTag = async (itemId: string, itemType: "SNIPPET" | "PROBLEM") => {
    if (!tagSearchQuery.trim()) return;
    try {
      const presetColors = ["#8b5cf6", "#10b981", "#f43f5e", "#f59e0b", "#0ea5e9"];
      const randomColor = presetColors[Math.floor(Math.random() * presetColors.length)];
      const newTag = await createTagMutation.mutateAsync({
        name: tagSearchQuery.trim(),
        color: randomColor,
      });
      await associateTagMutation.mutateAsync({
        itemType,
        itemId,
        tagId: newTag.id,
      });
      setTagSearchQuery("");
      setShowTagPopoverId(null);
      toast.success(`Tag "${newTag.name}" created and associated`);
    } catch {
      toast.error("Failed to create tag");
    }
  };

  // Snippet Form open wrappers
  const handleOpenCreateSnippet = () => {
    setIsCreatingSnippet(true);
    setIsEditingSnippet(false);
    setFormTitle("");
    setFormDescription("");
    setFormContent("");
    setFormLanguage("PLAIN_TEXT");
    setFormSnippetType("CODE");
  };

  const handleOpenEditSnippet = () => {
    if (!selectedSnippet) return;
    setIsEditingSnippet(true);
    setIsCreatingSnippet(false);
    setFormTitle(selectedSnippet.title);
    setFormDescription(selectedSnippet.description || "");
    setFormContent(selectedSnippet.content);
    setFormLanguage(selectedSnippet.language);
    setFormSnippetType(selectedSnippet.snippetType);
  };

  const handleDeleteSnippet = (snippetId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Snippet",
      message: "Are you sure you want to delete the snippet",
      itemName: title,
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteSnippetMutation.mutateAsync(snippetId);
          toast.success("Snippet deleted successfully");
          if (selectedSnippetId === snippetId) setSelectedSnippetId(undefined);
          closeConfirmModal();
        } catch (err) {
          toast.error(err instanceof Error ? err.message : "Failed to delete snippet");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  // Problem Status patch
  const handleStatusChange = async (problemId: string, status: ProblemStatus) => {
    try {
      await updateProblemStatusMutation.mutateAsync({ problemId, status });
      toast.success(`Status updated to ${status.replace("_", " ")}`);
    } catch {
      toast.error("Failed to update status");
    }
  };

  const handleDeleteProblem = (problemId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Diagnostic Node",
      message: "Are you sure you want to delete the diagnostic node",
      itemName: title,
      warningText: "This action cannot be undone. All logs and solution data will be permanently lost.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteProblemMutation.mutateAsync(problemId);
          toast.success("Problem node deleted successfully");
          if (selectedProblemId === problemId) setSelectedProblemId(undefined);
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete problem");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const projectIcon = getIconComponent(project?.icon);
  const isSubmittingSnippet = createSnippetMutation.isPending || updateSnippetMutation.isPending;

  return (
    <div className="flex flex-col gap-4 h-full">
      <div className={styles.projectHeader}>
        <div className={styles.projectTitleSection}>
          <div className={styles.projectIcon}>
            {React.createElement(projectIcon, {
              size: 22,
              style: { color: project?.color || "var(--color-primary)" },
            })}
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
        {/* Workspace Sidebar Tabs Selector */}
        <div className={styles.workspaceSidebar}>
          <button
            className={`${styles.workspaceTab} ${activeTab === "snippets" ? styles.workspaceTabActive : ""}`}
            onClick={() => setActiveTab("snippets")}
            title="Code Snippets"
          >
            <Icons.Code size={18} />
            <span className={styles.workspaceTabLabel}>Snippets</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "problems" ? styles.workspaceTabActive : ""}`}
            onClick={() => setActiveTab("problems")}
            title="Problems & Diagnostics"
          >
            <Icons.AlertCircle size={18} />
            <span className={styles.workspaceTabLabel}>Problems</span>
            {openProblemsCount > 0 && (
              <span className={styles.badgeCount}>{openProblemsCount}</span>
            )}
          </button>

          <button className={`${styles.workspaceTab} ${styles.workspaceTabDisabled}`} title="Credentials (Coming Soon)" disabled>
            <Icons.KeyRound size={18} />
            <span className={styles.workspaceTabLabel}>Credentials</span>
          </button>

          <button className={`${styles.workspaceTab} ${styles.workspaceTabDisabled}`} title="System Notes (Coming Soon)" disabled>
            <Icons.FileText size={18} />
            <span className={styles.workspaceTabLabel}>Notes</span>
          </button>

          <button className={`${styles.workspaceTab} ${styles.workspaceTabDisabled}`} title="Web Links (Coming Soon)" disabled>
            <Icons.Link2 size={18} />
            <span className={styles.workspaceTabLabel}>Links</span>
          </button>
        </div>

        {/* Tab 1: Snippets Workspace */}
        {activeTab === "snippets" && (
          <>
            {/* Left Side: Snippets navigation list */}
            <div className={styles.leftPanel}>
              <button
                className={styles.newSnippetBtn}
                onClick={handleOpenCreateSnippet}
                disabled={isSubmittingSnippet}
              >
                <Icons.Plus size={14} />
                <span>Add Snippet</span>
              </button>

              <div className={styles.searchBar}>
                <Icons.Search className={styles.searchIcon} size={14} />
                <input
                  type="text"
                  placeholder="Search snippets..."
                  className={styles.searchInput}
                  value={snippetSearchQuery}
                  onChange={(e) => setSnippetSearchQuery(e.target.value)}
                  disabled={isSubmittingSnippet}
                />
              </div>

              <div className={styles.filterTabs}>
                <button
                  className={`${styles.filterTab} ${snippetTypeFilter === "ALL" ? styles.filterTabActive : ""}`}
                  onClick={() => setSnippetTypeFilter("ALL")}
                  disabled={isSubmittingSnippet}
                >
                  ALL
                </button>
                <button
                  className={`${styles.filterTab} ${snippetTypeFilter === "CODE" ? styles.filterTabActive : ""}`}
                  onClick={() => setSnippetTypeFilter("CODE")}
                  disabled={isSubmittingSnippet}
                >
                  CODE
                </button>
                <button
                  className={`${styles.filterTab} ${snippetTypeFilter === "COMMAND" ? styles.filterTabActive : ""}`}
                  onClick={() => setSnippetTypeFilter("COMMAND")}
                  disabled={isSubmittingSnippet}
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
                        selectedSnippetId === s.id && !isCreatingSnippet ? styles.snippetItemActive : ""
                      }`}
                      onClick={() => {
                        setSelectedSnippetId(s.id);
                        setIsCreatingSnippet(false);
                        setIsEditingSnippet(false);
                      }}
                      disabled={isSubmittingSnippet}
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

            {/* Right Side: Snippets details / form editor */}
            <div className={styles.rightPanel}>
              {isCreatingSnippet || isEditingSnippet ? (
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
                      if (isCreatingSnippet) {
                        const res = await createSnippetMutation.mutateAsync(payload);
                        toast.success("Snippet created successfully");
                        setSelectedSnippetId(res.id);
                      } else if (isEditingSnippet && selectedSnippetId) {
                        await updateSnippetMutation.mutateAsync(payload);
                        toast.success("Snippet updated successfully");
                      }
                      setIsCreatingSnippet(false);
                      setIsEditingSnippet(false);
                    } catch (err) {
                      toast.error(err instanceof Error ? err.message : "Failed to save snippet");
                    }
                  }}
                  className={styles.formScroll}
                >
                  <h2 className={styles.formTitle}>
                    {isCreatingSnippet ? "CREATE NEW SNIPPET" : "EDIT SNIPPET"}
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
                        disabled={isSubmittingSnippet}
                        required
                      />
                    </div>

                    <div className={styles.formField}>
                      <label className={styles.formLabel}>Description</label>
                      <textarea
                        ref={descTextareaRef}
                        className={styles.formTextarea}
                        placeholder="Provide a brief context or description..."
                        value={formDescription}
                        onChange={(e) => setFormDescription(e.target.value)}
                        disabled={isSubmittingSnippet}
                      />
                    </div>

                    <div className={styles.formRow}>
                      <div className={styles.formField}>
                        <label className={styles.formLabel}>Type</label>
                        <select
                          className={styles.formInput}
                          value={formSnippetType}
                          onChange={(e) => setFormSnippetType(e.target.value as SnippetType)}
                          disabled={isSubmittingSnippet}
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
                          disabled={isSubmittingSnippet}
                        >
                          {ALL_LANGUAGES.map((lang) => (
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
                          readOnly: isSubmittingSnippet,
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
                        setIsCreatingSnippet(false);
                        setIsEditingSnippet(false);
                      }}
                      disabled={isSubmittingSnippet}
                    >
                      Cancel
                    </button>
                    <button type="submit" className={styles.btnPrimary} disabled={isSubmittingSnippet}>
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
                        <button
                          className={styles.btnIcon}
                          onClick={handleOpenEditSnippet}
                          title="Edit Snippet"
                          disabled={isSubmittingSnippet}
                        >
                          <Icons.Edit3 size={14} />
                        </button>
                        <button
                          className={`${styles.btnIcon} ${styles.btnIconDanger}`}
                          onClick={() =>
                            handleDeleteSnippet(selectedSnippet.id, selectedSnippet.title)
                          }
                          title="Delete Snippet"
                          disabled={isSubmittingSnippet}
                        >
                          <Icons.Trash2 size={14} />
                        </button>
                      </div>
                    </div>

                    {/* Shared Tags Manager for Snippets */}
                    <div className={styles.tagSection}>
                      <div className={styles.tagHeader}>
                        <Icons.Tag size={12} className="text-muted-foreground" />
                        <span className={styles.tagSectionTitle}>Associated Tags</span>
                      </div>
                      <div className={styles.tagList}>
                        {selectedSnippet.tags && selectedSnippet.tags.map((tag: TagSummaryResponse) => (
                          <span key={tag.id} className={styles.tagPill}>
                            <span
                              className={styles.tagDot}
                              style={{ backgroundColor: tag.color || "var(--color-primary)" }}
                            />
                            <span>{tag.name}</span>
                            <button
                              type="button"
                              className={styles.tagRemoveBtn}
                              onClick={() => handleRemoveTag(selectedSnippet.id, "SNIPPET", tag.id)}
                              title={`Remove tag ${tag.name}`}
                            >
                              <Icons.X size={10} />
                            </button>
                          </span>
                        ))}

                        <div className={styles.addTagContainer}>
                          <button
                            type="button"
                            className={styles.addTagBtn}
                            onClick={() =>
                              setShowTagPopoverId(
                                showTagPopoverId === `snippet-${selectedSnippet.id}`
                                  ? null
                                  : `snippet-${selectedSnippet.id}`
                              )
                            }
                          >
                            <Icons.Plus size={10} />
                            <span>Add Tag</span>
                          </button>

                          {showTagPopoverId === `snippet-${selectedSnippet.id}` && (
                            <div className={styles.tagPopover}>
                              <div className={styles.popoverHeader}>
                                <input
                                  type="text"
                                  placeholder="Filter/create tag..."
                                  className={styles.tagSearchInput}
                                  value={tagSearchQuery}
                                  onChange={(e) => setTagSearchQuery(e.target.value)}
                                  autoFocus
                                />
                              </div>
                              <div className={styles.popoverList}>
                                {tagsData
                                  .filter(
                                    (t) =>
                                      t.name.toLowerCase().includes(tagSearchQuery.toLowerCase()) &&
                                      (!selectedSnippet.tags ||
                                        !selectedSnippet.tags.some((st: TagSummaryResponse) => st.id === t.id))
                                  )
                                  .map((t) => (
                                    <button
                                      key={t.id}
                                      type="button"
                                      className={styles.popoverItem}
                                      onClick={() =>
                                        handleAddTag(selectedSnippet.id, "SNIPPET", t.id)
                                      }
                                    >
                                      <span
                                        className={styles.tagColorPreview}
                                        style={{ backgroundColor: t.color || "var(--color-primary)" }}
                                      />
                                      <span>{t.name}</span>
                                    </button>
                                  ))}

                                {tagSearchQuery.trim() &&
                                  !tagsData.some(
                                    (t) => t.name.toLowerCase() === tagSearchQuery.toLowerCase()
                                  ) && (
                                    <button
                                      type="button"
                                      className={styles.popoverItemCreate}
                                      onClick={() =>
                                        handleCreateAndAddTag(selectedSnippet.id, "SNIPPET")
                                      }
                                    >
                                      <Icons.Plus size={10} />
                                      <span>Create "{tagSearchQuery}"</span>
                                    </button>
                                  )}
                              </div>
                            </div>
                          )}
                        </div>
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
          </>
        )}

        {/* Tab 2: Problems / Errors Diagnostic Workspace */}
        {activeTab === "problems" && (
          <>
            {/* Middle Panel: Problems Navigation List */}
            <div className={styles.leftPanel}>
              <button
                className={styles.newSnippetBtn}
                onClick={() => {
                  setEditingProblemId(undefined);
                  setIsProblemFormOpen(true);
                }}
              >
                <Icons.Plus size={14} />
                <span>Log Problem Node</span>
              </button>

              <div className={styles.searchBar}>
                <Icons.Search className={styles.searchIcon} size={14} />
                <input
                  type="text"
                  placeholder="Search errors or tags..."
                  className={styles.searchInput}
                  value={problemSearchQuery}
                  onChange={(e) => setProblemSearchQuery(e.target.value)}
                />
              </div>

              <div className={styles.problemsFilterArea}>
                {/* Status Toggle buttons */}
                <div className={styles.filterTabs}>
                  <button
                    className={`${styles.filterTab} ${problemStatusFilter === "UNRESOLVED" ? styles.filterTabActive : ""}`}
                    onClick={() => setProblemStatusFilter("UNRESOLVED")}
                  >
                    OPEN
                  </button>
                  <button
                    className={`${styles.filterTab} ${problemStatusFilter === "RESOLVED" ? styles.filterTabActive : ""}`}
                    onClick={() => setProblemStatusFilter("RESOLVED")}
                  >
                    RESOLVED
                  </button>
                  <button
                    className={`${styles.filterTab} ${problemStatusFilter === "ALL" ? styles.filterTabActive : ""}`}
                    onClick={() => setProblemStatusFilter("ALL")}
                  >
                    ALL
                  </button>
                </div>

                {/* Severity select dropdown */}
                <select
                  className={styles.searchInput}
                  value={problemSeverityFilter}
                  onChange={(e) => setProblemSeverityFilter(e.target.value as "ALL" | ProblemSeverity)}
                >
                  <option value="ALL">ALL SEVERITIES</option>
                  <option value="CRITICAL">CRITICAL ONLY</option>
                  <option value="HIGH">HIGH SEVERITY</option>
                  <option value="MEDIUM">MEDIUM SEVERITY</option>
                  <option value="LOW">LOW SEVERITY</option>
                </select>
              </div>

              <div className={styles.problemList}>
                {filteredProblems.length === 0 ? (
                  <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border">
                    No diagnostics logged
                  </div>
                ) : (
                  filteredProblems.map((p) => (
                    <button
                      key={p.id}
                      className={`${styles.problemItem} ${
                        selectedProblemId === p.id ? styles.problemItemActive : ""
                      }`}
                      onClick={() => setSelectedProblemId(p.id)}
                    >
                      {/* Left indicator bar matching severity */}
                      <div
                        className={`${styles.severityIndicator} ${
                          p.severity === "CRITICAL"
                            ? styles.severityCritical
                            : p.severity === "HIGH"
                              ? styles.severityHigh
                              : p.severity === "MEDIUM"
                                ? styles.severityMedium
                                : styles.severityLow
                        }`}
                      />

                      <div className={styles.snippetItemHeader}>
                        <span className={styles.snippetItemTitle}>{p.title}</span>
                        <span className="text-[10px] text-muted-foreground font-mono">
                          {new Date(p.createdAt).toLocaleDateString()}
                        </span>
                      </div>

                      <div className="flex gap-1.5 items-center mt-1 flex-wrap">
                        <span
                          className={`${styles.statusBadge} ${
                            p.status === "OPEN"
                              ? styles.statusOpen
                              : p.status === "WORKING_ON"
                                ? styles.statusWorking
                                : p.status === "RESOLVED"
                                  ? styles.statusResolved
                                  : styles.statusWontFix
                          }`}
                        >
                          {p.status.replace("_", " ")}
                        </span>

                        <span
                          className={`${styles.severityBadge} ${
                            p.severity === "CRITICAL"
                              ? styles.badgeCritical
                              : p.severity === "HIGH"
                                ? styles.badgeHigh
                                : p.severity === "MEDIUM"
                                  ? styles.badgeMedium
                                  : styles.badgeLow
                          }`}
                        >
                          {p.severity}
                        </span>
                      </div>

                      {p.tags && p.tags.length > 0 && (
                        <div className="flex gap-1 flex-wrap mt-1.5">
                          {p.tags.map((t) => (
                            <span
                              key={t.id}
                              className="text-[9px] font-mono px-1.5 py-0.5 rounded-full border border-border flex items-center gap-1"
                            >
                              <span
                                className="w-1 h-1 rounded-full"
                                style={{ backgroundColor: t.color || "var(--color-primary)" }}
                              />
                              {t.name}
                            </span>
                          ))}
                        </div>
                      )}
                    </button>
                  ))
                )}
              </div>
            </div>

            {/* Right Side: Problems Workspace Detail Console */}
            <div className={styles.rightPanel}>
              {isLoadingProblemDetail && selectedProblemId ? (
                <div className="flex-1 flex flex-col items-center justify-center p-12 gap-3">
                  <Icons.Loader2 className="animate-spin text-primary" size={32} />
                  <span className="text-xs text-muted-foreground font-mono">LOADING ERROR DETAILS...</span>
                </div>
              ) : selectedProblem && problemDetail ? (
                <div className={styles.problemDetailScroll} key={problemDetail.id}>
                  <div className={styles.problemDetailContainer}>
                    <div className={styles.detailHeader}>
                      <div className={styles.detailTitleSection}>
                        <div className={styles.detailTitleRow}>
                          <h2 className={styles.detailTitle}>{problemDetail.title}</h2>
                          <div className="flex gap-1.5">
                            <span
                              className={`${styles.severityBadge} ${
                                problemDetail.severity === "CRITICAL"
                                  ? styles.badgeCritical
                                  : problemDetail.severity === "HIGH"
                                    ? styles.badgeHigh
                                    : problemDetail.severity === "MEDIUM"
                                      ? styles.badgeMedium
                                      : styles.badgeLow
                              }`}
                            >
                              {problemDetail.severity} SEVERITY
                            </span>
                          </div>
                        </div>

                        <div className={styles.problemMetadataRow}>
                          <div className={styles.problemMetadataItem}>
                            <Icons.Calendar size={12} />
                            <span>Logged: {new Date(problemDetail.createdAt).toLocaleString()}</span>
                          </div>
                          <div className={styles.problemMetadataItem}>
                            <Icons.RefreshCw size={12} />
                            <span>Updated: {new Date(problemDetail.updatedAt).toLocaleString()}</span>
                          </div>
                        </div>
                      </div>

                      <div className={styles.detailActions}>
                        <button
                          className={styles.btnIcon}
                          onClick={() => {
                            setEditingProblemId(problemDetail.id);
                            setIsProblemFormOpen(true);
                          }}
                          title="Edit Diagnostics"
                        >
                          <Icons.Edit3 size={14} />
                        </button>
                        <button
                          className={`${styles.btnIcon} ${styles.btnIconDanger}`}
                          onClick={() =>
                            handleDeleteProblem(problemDetail.id, problemDetail.title)
                          }
                          title="Delete Diagnostic Node"
                        >
                          <Icons.Trash2 size={14} />
                        </button>
                      </div>
                    </div>

                    {/* Interactive Resolution Action Panel */}
                    <div className="px-5 pt-4">
                      <div className={styles.statusActionPanel}>
                        <span className={styles.statusActionHeader}>Resolution Workflow Controls</span>
                        <div className={styles.statusActionRow}>
                          <button
                            type="button"
                            className={`${styles.statusSwitchBtn} ${
                              problemDetail.status === "OPEN" ? styles.statusSwitchBtnActive : ""
                            }`}
                            onClick={() => handleStatusChange(problemDetail.id, "OPEN")}
                          >
                            <Icons.CircleAlert size={12} />
                            <span>Set Open</span>
                          </button>
                          <button
                            type="button"
                            className={`${styles.statusSwitchBtn} ${
                              problemDetail.status === "WORKING_ON" ? styles.statusSwitchBtnActive : ""
                            }`}
                            onClick={() => handleStatusChange(problemDetail.id, "WORKING_ON")}
                          >
                            <Icons.Play size={12} />
                            <span>Investigate (Work)</span>
                          </button>
                          <button
                            type="button"
                            className={`${styles.statusSwitchBtn} ${
                              problemDetail.status === "RESOLVED" ? styles.statusSwitchBtnActive : ""
                            }`}
                            onClick={() => handleStatusChange(problemDetail.id, "RESOLVED")}
                          >
                            <Icons.CheckCircle2 size={12} className="text-emerald-500" />
                            <span className="text-emerald-500">Resolve Error</span>
                          </button>
                          <button
                            type="button"
                            className={`${styles.statusSwitchBtn} ${
                              problemDetail.status === "WONT_FIX" ? styles.statusSwitchBtnActive : ""
                            }`}
                            onClick={() => handleStatusChange(problemDetail.id, "WONT_FIX")}
                          >
                            <Icons.EyeOff size={12} />
                            <span>Won't Fix</span>
                          </button>
                        </div>
                      </div>
                    </div>

                    {/* Tags Manager for Problems */}
                    <div className={styles.tagSection}>
                      <div className={styles.tagHeader}>
                        <Icons.Tag size={12} className="text-muted-foreground" />
                        <span className={styles.tagSectionTitle}>Diagnostic Labels & Tags</span>
                      </div>
                      <div className={styles.tagList}>
                        {problemDetail.tags && problemDetail.tags.map((tag: TagSummaryResponse) => (
                          <span key={tag.id} className={styles.tagPill}>
                            <span
                              className={styles.tagDot}
                              style={{ backgroundColor: tag.color || "var(--color-primary)" }}
                            />
                            <span>{tag.name}</span>
                            <button
                              type="button"
                              className={styles.tagRemoveBtn}
                              onClick={() => handleRemoveTag(problemDetail.id, "PROBLEM", tag.id)}
                              title={`Remove tag ${tag.name}`}
                            >
                              <Icons.X size={10} />
                            </button>
                          </span>
                        ))}

                        <div className={styles.addTagContainer}>
                          <button
                            type="button"
                            className={styles.addTagBtn}
                            onClick={() =>
                              setShowTagPopoverId(
                                showTagPopoverId === `problem-${problemDetail.id}`
                                  ? null
                                  : `problem-${problemDetail.id}`
                              )
                            }
                          >
                            <Icons.Plus size={10} />
                            <span>Add Tag</span>
                          </button>

                          {showTagPopoverId === `problem-${problemDetail.id}` && (
                            <div className={styles.tagPopover}>
                              <div className={styles.popoverHeader}>
                                <input
                                  type="text"
                                  placeholder="Filter/create tag..."
                                  className={styles.tagSearchInput}
                                  value={tagSearchQuery}
                                  onChange={(e) => setTagSearchQuery(e.target.value)}
                                  autoFocus
                                />
                              </div>
                              <div className={styles.popoverList}>
                                {tagsData
                                  .filter(
                                    (t) =>
                                      t.name.toLowerCase().includes(tagSearchQuery.toLowerCase()) &&
                                      (!problemDetail.tags ||
                                        !problemDetail.tags.some((st: TagSummaryResponse) => st.id === t.id))
                                  )
                                  .map((t) => (
                                    <button
                                      key={t.id}
                                      type="button"
                                      className={styles.popoverItem}
                                      onClick={() =>
                                        handleAddTag(problemDetail.id, "PROBLEM", t.id)
                                      }
                                    >
                                      <span
                                        className={styles.tagColorPreview}
                                        style={{ backgroundColor: t.color || "var(--color-primary)" }}
                                      />
                                      <span>{t.name}</span>
                                    </button>
                                  ))}

                                {tagSearchQuery.trim() &&
                                  !tagsData.some(
                                    (t) => t.name.toLowerCase() === tagSearchQuery.toLowerCase()
                                  ) && (
                                    <button
                                      type="button"
                                      className={styles.popoverItemCreate}
                                      onClick={() =>
                                        handleCreateAndAddTag(problemDetail.id, "PROBLEM")
                                      }
                                    >
                                      <Icons.Plus size={10} />
                                      <span>Create "{tagSearchQuery}"</span>
                                    </button>
                                  )}
                              </div>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Shell Panel: Error stack trace logs */}
                    <div className="px-5">
                      <div className={styles.shellContainer}>
                        <div className={styles.shellHeader}>
                          <div className={styles.shellTitleSection}>
                            <Icons.Terminal size={14} className="text-rose-500" />
                            <span className={styles.shellTitle}>Stack Trace Log Output</span>
                          </div>
                          {problemDetail.errorDescription && (
                            <button
                              className={styles.copyButton}
                              onClick={() =>
                                handleCopy(
                                  problemDetail.errorDescription || "",
                                  `err-${problemDetail.id}`
                                )
                              }
                            >
                              {copiedId === `err-${problemDetail.id}` ? (
                                <Icons.Check size={12} className="text-emerald-500" />
                              ) : (
                                <Icons.Copy size={12} />
                              )}
                            </button>
                          )}
                        </div>
                        <div className={styles.shellEditorWrapper}>
                          {problemDetail.errorDescription ? (
                            <div className="border border-border/50 rounded overflow-hidden">
                              <Editor
                                key={`err-${problemDetail.id}`}
                                height="220px"
                                language="plaintext"
                                theme="vs-dark"
                                value={problemDetail.errorDescription}
                                options={{
                                  readOnly: true,
                                  minimap: { enabled: false },
                                  fontSize: 12,
                                  lineNumbers: "on",
                                  scrollBeyondLastLine: false,
                                  automaticLayout: true,
                                  domReadOnly: true,
                                  padding: { top: 6, bottom: 6 },
                                }}
                              />
                            </div>
                          ) : (
                            <div className="text-xs text-muted-foreground p-8 text-center font-mono">
                              No stack trace logs documented for this error.
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    {/* Shell Panel: Solution Fix */}
                    <div className="px-5 pb-6">
                      <div className={styles.shellContainer}>
                        <div className={styles.shellHeader}>
                          <div className={styles.shellTitleSection}>
                            <Icons.CheckSquare size={14} className="text-emerald-500" />
                            <span className={styles.shellTitle}>Solution script / Code fix</span>
                          </div>
                          {problemDetail.solution && (
                            <button
                              className={styles.copyButton}
                              onClick={() =>
                                handleCopy(problemDetail.solution || "", `sol-${problemDetail.id}`)
                              }
                            >
                              {copiedId === `sol-${problemDetail.id}` ? (
                                <Icons.Check size={12} className="text-emerald-500" />
                              ) : (
                                <Icons.Copy size={12} />
                              )}
                            </button>
                          )}
                        </div>
                        <div className={styles.shellEditorWrapper}>
                          {problemDetail.solution ? (
                            <div className="border border-border/50 rounded overflow-hidden">
                              <Editor
                                key={`sol-${problemDetail.id}`}
                                height="220px"
                                language="plaintext"
                                theme="vs-dark"
                                value={problemDetail.solution}
                                options={{
                                  readOnly: true,
                                  minimap: { enabled: false },
                                  fontSize: 12,
                                  lineNumbers: "on",
                                  scrollBeyondLastLine: false,
                                  automaticLayout: true,
                                  domReadOnly: true,
                                  padding: { top: 6, bottom: 6 },
                                }}
                              />
                            </div>
                          ) : (
                            <div className="text-xs text-muted-foreground p-8 text-center font-mono">
                              No resolution script documented. Edit details to attach a code fix.
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className={styles.placeholder}>
                  <Icons.AlertTriangle size={48} className="text-muted-foreground animate-pulse" />
                  <div className={styles.placeholderText}>
                    No diagnostic node selected. Select an error from the navigator or click "Log
                    Problem Node" to record a new anomaly.
                  </div>
                </div>
              )}
            </div>
          </>
        )}
      </div>

      {/* Global Problem Form Modal */}
      <ProblemForm
        isOpen={isProblemFormOpen}
        onClose={() => {
          setIsProblemFormOpen(false);
          setEditingProblemId(undefined);
        }}
        projectId={projectId}
        problemId={editingProblemId}
      />

      {/* Global Confirm Delete Modal */}
      <ConfirmModal
        isOpen={confirmModal.isOpen}
        onClose={closeConfirmModal}
        onConfirm={confirmModal.onConfirm}
        title={confirmModal.title}
        message={confirmModal.message}
        itemName={confirmModal.itemName}
        warningText={confirmModal.warningText}
        isLoading={confirmModal.isLoading}
      />
    </div>
  );
};

const useParamsHelper = () => {
  const match = useMatch({ from: "/projects/$projectId" });
  return match.params;
};

export const ProjectDetailRouteComponent: React.FC = () => {
  return <ProjectDetailView />;
};
