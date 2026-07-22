import React, { useState, useEffect } from "react";
import { Link, useMatch } from "@tanstack/react-router";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { marked } from "marked";
import DOMPurify from "dompurify";
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
import { TagsManagerModal } from "./TagsManagerModal";
import { NoteForm } from "~features/notes/components/NoteForm";
import { LinkForm } from "~features/links/components/LinkForm";
import {
  useNotesQuery,
  useNoteQuery,
  useDeleteNoteMutation,
  useArchiveNoteMutation,
  useUnarchiveNoteMutation,
} from "~features/notes/hooks/useNotes";
import {
  useLinksQuery,
  useLinkQuery,
  useDeleteLinkMutation,
} from "~features/links/hooks/useLinks";
import {
  useMasterPasswordSetupStatusQuery,
  useVaultStatusQuery,
} from "~features/security/hooks/useSecurity";
import { MasterPasswordSetupCard } from "~features/security/components/MasterPasswordSetupCard";
import { UnlockVaultCard } from "~features/security/components/UnlockVaultCard";
import { VaultSecurityBanner } from "~features/security/components/VaultSecurityBanner";
import {
  useCredentialsQuery,
  useDeleteCredentialMutation,
} from "~features/credentials/hooks/useCredentials";
import { CredentialForm } from "~features/credentials/components/CredentialForm";
import { CredentialDetailModal } from "~features/credentials/components/CredentialDetailModal";
import { useInactivityAutoLock } from "../hooks/useInactivityAutoLock";
import type {
  SnippetLanguage,
  SnippetType,
  ProblemStatus,
  ProblemSeverity,
  TagSummaryResponse,
  CredentialSecretType,
} from "~types/api";
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

  // Load notes details
  const { data: notesData } = useNotesQuery(projectId);
  const deleteNoteMutation = useDeleteNoteMutation(projectId);
  const archiveNoteMutation = useArchiveNoteMutation(projectId);
  const unarchiveNoteMutation = useUnarchiveNoteMutation(projectId);

  // Load links details
  const { data: linksData } = useLinksQuery(projectId);
  const deleteLinkMutation = useDeleteLinkMutation(projectId);

  // Workspace sub-navigation state
  const [activeTab, setActiveTab] = useState<
    "snippets" | "problems" | "credentials" | "notes" | "links"
  >("snippets");

  // Security queries
  const { data: isSetupRequired, isLoading: isSetupLoading } =
    useMasterPasswordSetupStatusQuery(activeTab === "credentials");
  const { data: vaultStatus, isLoading: isVaultStatusLoading } =
    useVaultStatusQuery(activeTab === "credentials");

  const isSecurityLoading =
    activeTab === "credentials" && (isSetupLoading || isVaultStatusLoading);
  const isVaultActive = vaultStatus?.active === true;
  const isVaultLocked = !isSecurityLoading && !isSetupRequired && !isVaultActive;

  // Credentials queries & mutations
  const { data: credentialsData, isLoading: isCredentialsLoading } =
    useCredentialsQuery(projectId, activeTab === "credentials" && isVaultActive);
  const deleteCredentialMutation = useDeleteCredentialMutation(projectId);

  // Auto-lock after 15 minutes of inactivity in credentials workspace
  useInactivityAutoLock(activeTab === "credentials" && isVaultActive, () => {
    setActiveTab("snippets");
  });

  // Toggle data-vault-active on root document when viewing credentials workspace
  useEffect(() => {
    if (activeTab === "credentials") {
      document.documentElement.dataset.vaultActive = "true";
    } else {
      delete document.documentElement.dataset.vaultActive;
    }
    return () => {
      delete document.documentElement.dataset.vaultActive;
    };
  }, [activeTab]);

  // Credentials UI states
  const [credentialSearchQuery, setCredentialSearchQuery] = useState("");
  const [credentialTypeFilter, setCredentialTypeFilter] = useState<
    "ALL" | CredentialSecretType
  >("ALL");
  const [isCredentialFormOpen, setIsCredentialFormOpen] = useState(false);
  const [editingCredentialId, setEditingCredentialId] = useState<
    string | undefined
  >(undefined);
  const [viewingCredentialId, setViewingCredentialId] = useState<
    string | undefined
  >(undefined);

  // Selected item states
  const [selectedSnippetId, setSelectedSnippetId] = useState<string | undefined>(undefined);
  const [selectedProblemId, setSelectedProblemId] = useState<string | undefined>(undefined);
  const [selectedNoteId, setSelectedNoteId] = useState<string | undefined>(undefined);
  const [selectedLinkId, setSelectedLinkId] = useState<string | undefined>(undefined);

  // Search/Filter states
  const [noteSearchQuery, setNoteSearchQuery] = useState("");
  const [linkSearchQuery, setLinkSearchQuery] = useState("");
  const [noteArchivedFilter, setNoteArchivedFilter] = useState<"ACTIVE" | "ARCHIVED" | "ALL">("ACTIVE");

  // Notes markdown preview toggle (persisted globally in localStorage)
  const [useMarkdown, setUseMarkdown] = useState<boolean>(() => {
    try {
      const saved = localStorage.getItem("devaulty_notes_markdown_preview");
      return saved !== "false"; // default to true
    } catch {
      return true;
    }
  });

  const handleToggleMarkdown = (enabled: boolean) => {
    setUseMarkdown(enabled);
    try {
      localStorage.setItem("devaulty_notes_markdown_preview", String(enabled));
    } catch {
      // ignore storage errors
    }
  };

  // Edit / Creation states
  const [isEditingSnippet, setIsEditingSnippet] = useState(false);
  const [isCreatingSnippet, setIsCreatingSnippet] = useState(false);
  const [isProblemFormOpen, setIsProblemFormOpen] = useState(false);
  const [editingProblemId, setEditingProblemId] = useState<string | undefined>(undefined);
  const [isNoteFormOpen, setIsNoteFormOpen] = useState(false);
  const [editingNoteId, setEditingNoteId] = useState<string | undefined>(undefined);
  const [isLinkFormOpen, setIsLinkFormOpen] = useState(false);
  const [editingLinkId, setEditingLinkId] = useState<string | undefined>(undefined);
  const [isTagsManagerOpen, setIsTagsManagerOpen] = useState(false);

  // Selected item detail queries
  const { data: noteDetail } = useNoteQuery(projectId, selectedNoteId || "");
  const { data: linkDetail } = useLinkQuery(projectId, selectedLinkId || "");

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

  // Close tag popover when clicking outside
  useEffect(() => {
    if (showTagPopoverId === null) return;

    const handleOutsideClick = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      if (!target.closest(`.${styles.addTagContainer}`)) {
        setShowTagPopoverId(null);
        setTagSearchQuery("");
      }
    };

    document.addEventListener("click", handleOutsideClick);
    return () => document.removeEventListener("click", handleOutsideClick);
  }, [showTagPopoverId]);

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
  const handleAddTag = async (
    itemId: string,
    itemType: "SNIPPET" | "PROBLEM" | "NOTE" | "LINK" | "CREDENTIAL",
    tagId: string
  ) => {
    try {
      await associateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag associated successfully");
      setShowTagPopoverId(null);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to associate tag");
    }
  };

  const handleRemoveTag = async (
    itemId: string,
    itemType: "SNIPPET" | "PROBLEM" | "NOTE" | "LINK" | "CREDENTIAL",
    tagId: string
  ) => {
    try {
      await disassociateTagMutation.mutateAsync({ itemType, itemId, tagId });
      toast.success("Tag removed successfully");
    } catch (err) {
      toast.error(err instanceof Error ? err.message : "Failed to remove tag");
    }
  };

  const handleCreateAndAddTag = async (
    itemId: string,
    itemType: "SNIPPET" | "PROBLEM" | "NOTE" | "LINK" | "CREDENTIAL"
  ) => {
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

  // Notes & Links handlers
  const handleDeleteNote = (noteId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete System Note",
      message: "Are you sure you want to permanently delete the note",
      itemName: title,
      warningText: "This action cannot be undone. The note contents will be permanently lost.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteNoteMutation.mutateAsync(noteId);
          toast.success("Note deleted successfully");
          if (selectedNoteId === noteId) setSelectedNoteId(undefined);
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete note");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const handleDeleteLink = (linkId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Web Link",
      message: "Are you sure you want to permanently delete the link",
      itemName: title,
      warningText: "This action cannot be undone. The link details will be permanently lost.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteLinkMutation.mutateAsync(linkId);
          toast.success("Link deleted successfully");
          if (selectedLinkId === linkId) setSelectedLinkId(undefined);
          closeConfirmModal();
        } catch {
          toast.error("Failed to delete link");
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const handleToggleArchiveNote = async (noteId: string, archived: boolean) => {
    try {
      if (archived) {
        await unarchiveNoteMutation.mutateAsync(noteId);
        toast.success("Note restored from archive");
      } else {
        await archiveNoteMutation.mutateAsync(noteId);
        toast.success("Note archived");
      }
    } catch {
      toast.error(`Failed to ${archived ? "restore" : "archive"} note`);
    }
  };

  const notes = notesData?.content || [];
  const filteredNotes = notes.filter((n) => {
    const matchesSearch = n.title.toLowerCase().includes(noteSearchQuery.toLowerCase());
    const matchesArchived =
      noteArchivedFilter === "ALL" ||
      (noteArchivedFilter === "ACTIVE" && !n.archived) ||
      (noteArchivedFilter === "ARCHIVED" && n.archived);
    return matchesSearch && matchesArchived;
  });

  const links = linksData?.content || [];
  const filteredLinks = links.filter((l) =>
    l.title.toLowerCase().includes(linkSearchQuery.toLowerCase()) ||
    (l.description && l.description.toLowerCase().includes(linkSearchQuery.toLowerCase()))
  );

  const credentials = credentialsData?.content || [];
  const filteredCredentials = credentials.filter((c) => {
    const matchesSearch =
      c.title.toLowerCase().includes(credentialSearchQuery.toLowerCase()) ||
      (c.relatedUrl &&
        c.relatedUrl.toLowerCase().includes(credentialSearchQuery.toLowerCase()));
    const matchesType =
      credentialTypeFilter === "ALL" || c.secretType === credentialTypeFilter;
    return matchesSearch && matchesType;
  });

  const handleDeleteCredential = (credId: string, title: string) => {
    setConfirmModal({
      isOpen: true,
      title: "Delete Vault Credential",
      message: "Are you sure you want to delete the credential",
      itemName: title,
      warningText:
        "This action cannot be undone. Encrypted secret data will be permanently deleted.",
      onConfirm: async () => {
        setConfirmModal((prev) => ({ ...prev, isLoading: true }));
        try {
          await deleteCredentialMutation.mutateAsync(credId);
          toast.success("Credential deleted successfully!");
          closeConfirmModal();
        } catch (err) {
          toast.error(
            err instanceof Error ? err.message : "Failed to delete credential."
          );
          setConfirmModal((prev) => ({ ...prev, isLoading: false }));
        }
      },
      isLoading: false,
    });
  };

  const handleTabChange = (
    tab: "snippets" | "problems" | "credentials" | "notes" | "links"
  ) => {
    setActiveTab(tab);
    setShowTagPopoverId(null);
    setTagSearchQuery("");
  };

  const renderMarkdown = (content: string | undefined) => {
    if (!content) {
      return '<span class="text-muted-foreground italic">No content documented. Click Edit to add details.</span>';
    }
    try {
      const rawHtml = marked.parse(content, { breaks: true, gfm: true }) as string;
      return DOMPurify.sanitize(rawHtml);
    } catch {
      return DOMPurify.sanitize(content);
    }
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

        <div className="flex items-center gap-2">
          <button
            type="button"
            onClick={() => setIsTagsManagerOpen(true)}
            className="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
          >
            <Icons.Tags size={12} />
            <span>Manage Tags</span>
          </button>

          <Link
            to="/"
            className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors"
          >
            <Icons.ArrowLeft size={12} />
            <span>Back to Dashboard</span>
          </Link>
        </div>
      </div>

      <div className={styles.pageLayout}>
        {/* Workspace Sidebar Tabs Selector */}
        <div className={styles.workspaceSidebar}>
          <button
            className={`${styles.workspaceTab} ${activeTab === "snippets" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("snippets")}
            title="Code Snippets"
          >
            <Icons.Code size={18} />
            <span className={styles.workspaceTabLabel}>Snippets</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "problems" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("problems")}
            title="Problems & Diagnostics"
          >
            <Icons.AlertCircle size={18} />
            <span className={styles.workspaceTabLabel}>Problems</span>
            {openProblemsCount > 0 && (
              <span className={styles.badgeCount}>{openProblemsCount}</span>
            )}
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "credentials" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("credentials")}
            title="Secure Credentials Vault"
          >
            <Icons.KeyRound size={18} />
            <span className={styles.workspaceTabLabel}>Credentials</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "notes" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("notes")}
            title="System Notes"
          >
            <Icons.FileText size={18} />
            <span className={styles.workspaceTabLabel}>Notes</span>
          </button>

          <button
            className={`${styles.workspaceTab} ${activeTab === "links" ? styles.workspaceTabActive : ""}`}
            onClick={() => handleTabChange("links")}
            title="Web Links"
          >
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
                                <div className="border-t border-border mt-2 pt-2 px-1">
                                  <button
                                    type="button"
                                    onClick={() => {
                                      setShowTagPopoverId(null);
                                      setIsTagsManagerOpen(true);
                                    }}
                                    className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                                  >
                                    <Icons.Settings size={10} />
                                    <span>Manage Tags</span>
                                  </button>
                                </div>
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
                                <div className="border-t border-border mt-2 pt-2 px-1">
                                  <button
                                    type="button"
                                    onClick={() => {
                                      setShowTagPopoverId(null);
                                      setIsTagsManagerOpen(true);
                                    }}
                                    className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                                  >
                                    <Icons.Settings size={10} />
                                    <span>Manage Tags</span>
                                  </button>
                                </div>
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

        {/* Tab 3: Notes Workspace */}
        {activeTab === "notes" && (
          <>
            {/* Left Side: Notes navigation list */}
            <div className={styles.leftPanel}>
              <button
                type="button"
                className={styles.newSnippetBtn}
                onClick={() => {
                  setEditingNoteId(undefined);
                  setIsNoteFormOpen(true);
                }}
              >
                <Icons.Plus size={14} />
                <span>Add Note</span>
              </button>

              <div className={styles.searchBar}>
                <Icons.Search className={styles.searchIcon} size={14} />
                <input
                  type="text"
                  placeholder="Search notes..."
                  className={styles.searchInput}
                  value={noteSearchQuery}
                  onChange={(e) => setNoteSearchQuery(e.target.value)}
                />
              </div>

              <div className={styles.filterTabs}>
                <button
                  type="button"
                  className={`${styles.filterTab} ${noteArchivedFilter === "ACTIVE" ? styles.filterTabActive : ""}`}
                  onClick={() => setNoteArchivedFilter("ACTIVE")}
                >
                  ACTIVE
                </button>
                <button
                  type="button"
                  className={`${styles.filterTab} ${noteArchivedFilter === "ARCHIVED" ? styles.filterTabActive : ""}`}
                  onClick={() => setNoteArchivedFilter("ARCHIVED")}
                >
                  ARCHIVED
                </button>
                <button
                  type="button"
                  className={`${styles.filterTab} ${noteArchivedFilter === "ALL" ? styles.filterTabActive : ""}`}
                  onClick={() => setNoteArchivedFilter("ALL")}
                >
                  ALL
                </button>
              </div>

              <div className={styles.snippetList}>
                {filteredNotes.length === 0 ? (
                  <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border font-mono">
                    No notes found
                  </div>
                ) : (
                  filteredNotes.map((n) => (
                    <button
                      key={n.id}
                      className={`${styles.snippetItem} ${selectedNoteId === n.id ? styles.snippetItemActive : ""}`}
                      onClick={() => setSelectedNoteId(n.id)}
                    >
                      <div className={styles.snippetItemHeader}>
                        <span className={styles.snippetItemTitle}>{n.title}</span>
                        {n.archived && <span className="text-[10px] text-amber-500 font-mono">ARCHIVED</span>}
                      </div>
                      {n.tags && n.tags.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-1">
                          {n.tags.map((tag) => (
                            <span
                              key={tag.id}
                              className={styles.badge}
                              style={{
                                backgroundColor: `${tag.color || "#8b5cf6"}15`,
                                color: tag.color || "#8b5cf6",
                                border: `1px solid ${tag.color || "#8b5cf6"}30`,
                                fontSize: "10px",
                              }}
                            >
                              {tag.name}
                            </span>
                          ))}
                        </div>
                      )}
                    </button>
                  ))
                )}
              </div>
            </div>

            {/* Right Side: Note Workspace details panel */}
            <div className={styles.rightPanel}>
              {noteDetail ? (
                <div className={styles.problemDetailScroll}>
                  <div className={styles.problemDetailContainer}>
                    <div className={styles.detailHeader}>
                      <div className={styles.detailTitleSection}>
                        <h2 className={styles.detailTitle}>{noteDetail.title}</h2>
                        <div className={styles.problemMetadataRow}>
                          <span>Created: {new Date(noteDetail.createdAt).toLocaleDateString()}</span>
                          <span>Updated: {new Date(noteDetail.updatedAt).toLocaleDateString()}</span>
                        </div>
                      </div>

                      <div className="flex gap-2">
                        <button
                          type="button"
                          onClick={() => {
                            setEditingNoteId(noteDetail.id);
                            setIsNoteFormOpen(true);
                          }}
                          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                        >
                          <Icons.Edit3 size={12} />
                          <span>Edit</span>
                        </button>
                        <button
                          type="button"
                          onClick={() => handleToggleArchiveNote(noteDetail.id, noteDetail.archived)}
                          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                        >
                          <Icons.Archive size={12} />
                          <span>{noteDetail.archived ? "Restore" : "Archive"}</span>
                        </button>
                        <button
                          type="button"
                          onClick={() => handleDeleteNote(noteDetail.id, noteDetail.title)}
                          className="flex items-center gap-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 border border-red-500/20 px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                        >
                          <Icons.Trash2 size={12} />
                          <span>Delete</span>
                        </button>
                      </div>
                    </div>

                    <div className="p-6 flex flex-col gap-6">
                      {/* Tags section */}
                      <div className="flex flex-col gap-2">
                        <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Tags</span>
                        <div className="flex flex-wrap items-center gap-1.5">
                          {noteDetail.tags &&
                            noteDetail.tags.map((tag: TagSummaryResponse) => (
                              <span
                                key={tag.id}
                                className={styles.badge}
                                style={{
                                  backgroundColor: `${tag.color || "#8b5cf6"}15`,
                                  color: tag.color || "#8b5cf6",
                                  border: `1px solid ${tag.color || "#8b5cf6"}30`,
                                  fontSize: "11px",
                                  padding: "0.2rem 0.5rem",
                                }}
                              >
                                <span>{tag.name}</span>
                                <button
                                  type="button"
                                  onClick={() => handleRemoveTag(noteDetail.id, "NOTE", tag.id)}
                                  className="ml-1 text-muted-foreground hover:text-foreground bg-transparent border-0 cursor-pointer"
                                >
                                  &times;
                                </button>
                              </span>
                            ))}

                          <div className={styles.addTagContainer}>
                            <button
                              type="button"
                              className={styles.addTagBtn}
                              onClick={() =>
                                setShowTagPopoverId(
                                  showTagPopoverId === `note-${noteDetail.id}`
                                    ? null
                                    : `note-${noteDetail.id}`
                                )
                              }
                            >
                              <Icons.Plus size={10} />
                              <span>Add Tag</span>
                            </button>

                            {showTagPopoverId === `note-${noteDetail.id}` && (
                              <div className={styles.tagPopover} style={{ bottom: "auto", top: "100%", marginTop: "0.375rem" }}>
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
                                        (!noteDetail.tags ||
                                          !noteDetail.tags.some((st: TagSummaryResponse) => st.id === t.id))
                                    )
                                    .map((t) => (
                                      <button
                                        key={t.id}
                                        type="button"
                                        className={styles.popoverItem}
                                        onClick={() => handleAddTag(noteDetail.id, "NOTE", t.id)}
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
                                        onClick={() => handleCreateAndAddTag(noteDetail.id, "NOTE")}
                                      >
                                        <Icons.Plus size={10} />
                                        <span>Create "{tagSearchQuery}"</span>
                                      </button>
                                    )}

                                  <div className="border-t border-border mt-2 pt-2 px-1">
                                    <button
                                      type="button"
                                      onClick={() => {
                                        setShowTagPopoverId(null);
                                        setIsTagsManagerOpen(true);
                                      }}
                                      className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                                    >
                                      <Icons.Settings size={10} />
                                      <span>Manage Tags</span>
                                    </button>
                                  </div>
                                </div>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>

                      {/* Content panel */}
                      <div className="flex-grow flex flex-col gap-2">
                        <div className="flex items-center justify-between">
                          <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Note Content</span>
                          <div className="flex items-center gap-1 bg-secondary/40 p-0.5 rounded border border-border">
                            <button
                              type="button"
                              onClick={() => handleToggleMarkdown(false)}
                              className={`px-2 py-1 text-[10px] font-mono rounded cursor-pointer transition-all ${
                                !useMarkdown
                                  ? "bg-primary text-primary-foreground shadow-sm font-bold"
                                  : "text-muted-foreground hover:text-foreground bg-transparent"
                              }`}
                              style={{ border: "none" }}
                            >
                              RAW
                            </button>
                            <button
                              type="button"
                              onClick={() => handleToggleMarkdown(true)}
                              className={`px-2 py-1 text-[10px] font-mono rounded cursor-pointer transition-all ${
                                useMarkdown
                                  ? "bg-primary text-primary-foreground shadow-sm font-bold"
                                  : "text-muted-foreground hover:text-foreground bg-transparent"
                              }`}
                              style={{ border: "none" }}
                            >
                              MARKDOWN
                            </button>
                          </div>
                        </div>

                        {useMarkdown ? (
                          <div
                            className={`bg-background/50 border border-border rounded p-6 text-sm leading-relaxed overflow-y-auto min-h-[300px] ${styles.markdownContainer}`}
                            dangerouslySetInnerHTML={{ __html: renderMarkdown(noteDetail.content) }}
                          />
                        ) : (
                          <div className="bg-background/50 border border-border rounded p-4 font-mono text-sm whitespace-pre-wrap leading-relaxed overflow-y-auto min-h-[300px]">
                            {noteDetail.content || (
                              <span className="text-muted-foreground italic">
                                No content documented. Click Edit to add details.
                              </span>
                            )}
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className={styles.placeholder}>
                  <Icons.FileText size={48} className="text-muted-foreground animate-pulse" />
                  <div className={styles.placeholderText}>
                    No note selected. Select a note from the navigator or click "Add Note" to write a new note.
                  </div>
                </div>
              )}
            </div>
          </>
        )}

        {/* Tab 3: Credentials Workspace */}
        {activeTab === "credentials" && (
          <div className="flex-1 flex flex-col h-full overflow-y-auto p-2">
            {isSecurityLoading ? (
              <div className="flex-1 flex items-center justify-center p-8">
                <Icons.Loader2 className="animate-spin text-muted-foreground" size={32} />
              </div>
            ) : isSetupRequired ? (
              <MasterPasswordSetupCard />
            ) : isVaultLocked ? (
              <UnlockVaultCard />
            ) : (
              <div className="flex-1 flex flex-col h-full overflow-hidden">
                <VaultSecurityBanner secondsLeft={vaultStatus?.secondsLeft} />

                <div className="flex-1 flex flex-col gap-4 overflow-hidden">
                  <div className="flex items-center justify-between gap-4 flex-wrap">
                    <div className="flex items-center gap-3">
                      <div className={styles.searchBar} style={{ minWidth: "260px" }}>
                        <Icons.Search className={styles.searchIcon} size={14} />
                        <input
                          type="text"
                          className={styles.searchInput}
                          placeholder="Search credentials in vault..."
                          value={credentialSearchQuery}
                          onChange={(e) => setCredentialSearchQuery(e.target.value)}
                        />
                      </div>

                      {/* Type filter tabs */}
                      <div className={styles.filterTabs}>
                        {(["ALL", "LOGIN", "API_KEY", "RAW_TEXT"] as const).map((type) => (
                          <button
                            key={type}
                            type="button"
                            className={`${styles.filterTab} ${credentialTypeFilter === type ? styles.filterTabActive : ""}`}
                            onClick={() => setCredentialTypeFilter(type)}
                          >
                            {type === "ALL" ? "ALL" : type === "RAW_TEXT" ? "TEXT" : type}
                          </button>
                        ))}
                      </div>
                    </div>

                    <button
                      type="button"
                      className="px-4 py-2 text-xs font-mono font-bold bg-primary text-primary-foreground rounded border border-primary hover:brightness-110 transition-all cursor-pointer flex items-center gap-1.5"
                      onClick={() => {
                        setEditingCredentialId(undefined);
                        setIsCredentialFormOpen(true);
                      }}
                    >
                      <Icons.Plus size={14} />
                      <span>New Credential</span>
                    </button>
                  </div>

                  {/* Grid / List of Credential Cards */}
                  {isCredentialsLoading ? (
                    <div className="flex flex-col items-center justify-center py-16 gap-3">
                      <Icons.Loader2 className="animate-spin text-primary" size={28} />
                      <span className="text-xs text-muted-foreground font-mono">LOADING CREDENTIALS...</span>
                    </div>
                  ) : filteredCredentials.length === 0 ? (
                    <div className={styles.placeholder}>
                      <Icons.KeyRound size={48} className="text-muted-foreground animate-pulse" />
                      <div className={styles.placeholderText}>
                        {credentialSearchQuery || credentialTypeFilter !== "ALL"
                          ? "No credentials match the active filters."
                          : "No credentials stored in this project yet. Click 'New Credential' to add one."}
                      </div>
                    </div>
                  ) : (
                    <div className={styles.credentialGrid}>
                      {filteredCredentials.map((cred) => (
                        <div
                          key={cred.id}
                          className={styles.credentialCard}
                          onClick={() => setViewingCredentialId(cred.id)}
                        >
                          <div className={styles.credentialCardHeader}>
                            <div className={styles.credentialTitleGroup}>
                              <span className={styles.credentialTypeBadge}>
                                {cred.secretType === "LOGIN" && <Icons.UserCheck size={12} />}
                                {cred.secretType === "API_KEY" && <Icons.KeyRound size={12} />}
                                {cred.secretType === "RAW_TEXT" && <Icons.Code2 size={12} />}
                                <span>{cred.secretType}</span>
                              </span>
                              <h3 className={styles.credentialCardTitle}>{cred.title}</h3>
                              {cred.relatedUrl && (
                                <a
                                  href={cred.relatedUrl}
                                  target="_blank"
                                  rel="noopener noreferrer"
                                  className={styles.credentialUrl}
                                  onClick={(e) => e.stopPropagation()}
                                >
                                  <Icons.ExternalLink size={12} />
                                  <span>{cred.relatedUrl.replace(/^https?:\/\//, "")}</span>
                                </a>
                              )}
                            </div>
                          </div>

                          <div className={styles.credentialCardFooter}>
                            {/* Tag list & tag popover */}
                            <div className={styles.tagList} onClick={(e) => e.stopPropagation()}>
                              {cred.tags && cred.tags.map((tag) => (
                                <span
                                  key={tag.id}
                                  className={styles.tagPill}
                                  style={{
                                    backgroundColor: `${tag.color || "#8b5cf6"}15`,
                                    color: tag.color || "#8b5cf6",
                                    borderColor: `${tag.color || "#8b5cf6"}30`,
                                  }}
                                >
                                  <span
                                    className={styles.tagDot}
                                    style={{ backgroundColor: tag.color || "#8b5cf6" }}
                                  />
                                  <span>{tag.name}</span>
                                  <button
                                    type="button"
                                    className={styles.tagRemoveBtn}
                                    onClick={(e) => {
                                      e.stopPropagation();
                                      handleRemoveTag(cred.id, "CREDENTIAL", tag.id);
                                    }}
                                    title="Remove tag"
                                  >
                                    <Icons.X size={10} />
                                  </button>
                                </span>
                              ))}

                              {/* Add Tag Popover Button */}
                              <div className={styles.addTagContainer}>
                                <button
                                  type="button"
                                  className={styles.addTagBtn}
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    setShowTagPopoverId(
                                      showTagPopoverId === `cred-${cred.id}` ? null : `cred-${cred.id}`
                                    );
                                  }}
                                >
                                  <Icons.Plus size={10} />
                                  <span>Tag</span>
                                </button>

                                {showTagPopoverId === `cred-${cred.id}` && (
                                  <div
                                    className={styles.tagPopover}
                                    style={{ bottom: "calc(100% + 6px)", top: "auto", zIndex: 100 }}
                                    onClick={(e) => e.stopPropagation()}
                                  >
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
                                            (!cred.tags || !cred.tags.some((st) => st.id === t.id))
                                        )
                                        .map((t) => (
                                          <button
                                            key={t.id}
                                            type="button"
                                            className={styles.popoverItem}
                                            onClick={(e) => {
                                              e.stopPropagation();
                                              handleAddTag(cred.id, "CREDENTIAL", t.id);
                                            }}
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
                                            onClick={(e) => {
                                              e.stopPropagation();
                                              handleCreateAndAddTag(cred.id, "CREDENTIAL");
                                            }}
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

                            {/* Actions */}
                            <div className={styles.credentialActions} onClick={(e) => e.stopPropagation()}>
                              <button
                                type="button"
                                className={styles.actionBtn}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  setEditingCredentialId(cred.id);
                                  setIsCredentialFormOpen(true);
                                }}
                                title="Edit"
                              >
                                <Icons.Edit2 size={12} />
                              </button>
                              <button
                                type="button"
                                className={`${styles.actionBtn} ${styles.actionBtnDanger}`}
                                onClick={(e) => {
                                  e.stopPropagation();
                                  handleDeleteCredential(cred.id, cred.title);
                                }}
                                title="Delete"
                              >
                                <Icons.Trash2 size={12} />
                              </button>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        )}

        {/* Tab 4: Links Workspace */}
        {activeTab === "links" && (
          <>
            {/* Left Side: Links navigation list */}
            <div className={styles.leftPanel}>
              <button
                type="button"
                className={styles.newSnippetBtn}
                onClick={() => {
                  setEditingLinkId(undefined);
                  setIsLinkFormOpen(true);
                }}
              >
                <Icons.Plus size={14} />
                <span>Add Link</span>
              </button>

              <div className={styles.searchBar}>
                <Icons.Search className={styles.searchIcon} size={14} />
                <input
                  type="text"
                  placeholder="Search links..."
                  className={styles.searchInput}
                  value={linkSearchQuery}
                  onChange={(e) => setLinkSearchQuery(e.target.value)}
                />
              </div>

              <div className={styles.snippetList}>
                {filteredLinks.length === 0 ? (
                  <div className="text-xs text-muted-foreground text-center py-8 border border-dashed rounded border-border font-mono">
                    No links found
                  </div>
                ) : (
                  filteredLinks.map((l) => (
                    <button
                      key={l.id}
                      className={`${styles.snippetItem} ${selectedLinkId === l.id ? styles.snippetItemActive : ""}`}
                      onClick={() => setSelectedLinkId(l.id)}
                    >
                      <div className={styles.snippetItemHeader}>
                        <span className={styles.snippetItemTitle}>{l.title}</span>
                        <Icons.ExternalLink size={12} className="text-muted-foreground shrink-0" />
                      </div>
                      <span className="text-[10px] text-muted-foreground font-mono truncate max-w-[220px] block">
                        {l.url}
                      </span>
                      {l.tags && l.tags.length > 0 && (
                        <div className="flex flex-wrap gap-1 mt-1">
                          {l.tags.map((tag) => (
                            <span
                              key={tag.id}
                              className={styles.badge}
                              style={{
                                backgroundColor: `${tag.color || "#8b5cf6"}15`,
                                color: tag.color || "#8b5cf6",
                                border: `1px solid ${tag.color || "#8b5cf6"}30`,
                                fontSize: "10px",
                              }}
                            >
                              {tag.name}
                            </span>
                          ))}
                        </div>
                      )}
                    </button>
                  ))
                )}
              </div>
            </div>

            {/* Right Side: Link Workspace details panel */}
            <div className={styles.rightPanel}>
              {linkDetail ? (
                <div className={styles.problemDetailScroll}>
                  <div className={styles.problemDetailContainer}>
                    <div className={styles.detailHeader}>
                      <div className={styles.detailTitleSection}>
                        <h2 className={styles.detailTitle}>{linkDetail.title}</h2>
                        <div className={styles.problemMetadataRow}>
                          <span>Created: {new Date(linkDetail.createdAt).toLocaleDateString()}</span>
                          <span>Updated: {new Date(linkDetail.updatedAt).toLocaleDateString()}</span>
                        </div>
                      </div>

                      <div className="flex gap-2">
                        <button
                          type="button"
                          onClick={() => {
                            setEditingLinkId(linkDetail.id);
                            setIsLinkFormOpen(true);
                          }}
                          className="flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground border border-border px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                        >
                          <Icons.Edit3 size={12} />
                          <span>Edit</span>
                        </button>
                        <button
                          type="button"
                          onClick={() => handleDeleteLink(linkDetail.id, linkDetail.title)}
                          className="flex items-center gap-1 text-xs text-red-400 hover:text-red-300 hover:bg-red-500/10 border border-red-500/20 px-2.5 py-1.5 rounded bg-card transition-colors cursor-pointer"
                        >
                          <Icons.Trash2 size={12} />
                          <span>Delete</span>
                        </button>
                      </div>
                    </div>

                    <div className="p-6 flex flex-col gap-6">
                      {/* URL card block */}
                      <div className="bg-background/50 border border-border rounded p-4 flex flex-col gap-3">
                        <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Web Address</span>
                        <div className="flex items-center justify-between gap-4">
                          <span className="font-mono text-sm text-primary truncate flex-grow">{linkDetail.url}</span>
                          <a
                            href={linkDetail.url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="flex items-center gap-1.5 text-xs text-primary-foreground bg-primary hover:bg-primary/95 border border-primary px-3.5 py-1.5 rounded transition-colors font-bold decoration-none cursor-pointer"
                          >
                            <span>Open Link</span>
                            <Icons.ExternalLink size={12} />
                          </a>
                        </div>
                      </div>

                      {/* Tags section */}
                      <div className="flex flex-col gap-2">
                        <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Tags</span>
                        <div className="flex flex-wrap items-center gap-1.5">
                          {linkDetail.tags &&
                            linkDetail.tags.map((tag: TagSummaryResponse) => (
                              <span
                                key={tag.id}
                                className={styles.badge}
                                style={{
                                  backgroundColor: `${tag.color || "#8b5cf6"}15`,
                                  color: tag.color || "#8b5cf6",
                                  border: `1px solid ${tag.color || "#8b5cf6"}30`,
                                  fontSize: "11px",
                                  padding: "0.2rem 0.5rem",
                                }}
                              >
                                <span>{tag.name}</span>
                                <button
                                  type="button"
                                  onClick={() => handleRemoveTag(linkDetail.id, "LINK", tag.id)}
                                  className="ml-1 text-muted-foreground hover:text-foreground bg-transparent border-0 cursor-pointer"
                                >
                                  &times;
                                </button>
                              </span>
                            ))}

                          <div className={styles.addTagContainer}>
                            <button
                              type="button"
                              className={styles.addTagBtn}
                              onClick={() =>
                                setShowTagPopoverId(
                                  showTagPopoverId === `link-${linkDetail.id}`
                                    ? null
                                    : `link-${linkDetail.id}`
                                )
                              }
                            >
                              <Icons.Plus size={10} />
                              <span>Add Tag</span>
                            </button>

                            {showTagPopoverId === `link-${linkDetail.id}` && (
                              <div className={styles.tagPopover} style={{ bottom: "auto", top: "100%", marginTop: "0.375rem" }}>
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
                                        (!linkDetail.tags ||
                                          !linkDetail.tags.some((st: TagSummaryResponse) => st.id === t.id))
                                    )
                                    .map((t) => (
                                      <button
                                        key={t.id}
                                        type="button"
                                        className={styles.popoverItem}
                                        onClick={() => handleAddTag(linkDetail.id, "LINK", t.id)}
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
                                        onClick={() => handleCreateAndAddTag(linkDetail.id, "LINK")}
                                      >
                                        <Icons.Plus size={10} />
                                        <span>Create "{tagSearchQuery}"</span>
                                      </button>
                                    )}

                                  <div className="border-t border-border mt-2 pt-2 px-1">
                                    <button
                                      type="button"
                                      onClick={() => {
                                        setShowTagPopoverId(null);
                                        setIsTagsManagerOpen(true);
                                      }}
                                      className="w-full flex items-center justify-center gap-1 text-[10px] text-muted-foreground hover:text-foreground py-1 font-mono uppercase tracking-wider bg-transparent border-0 cursor-pointer"
                                    >
                                      <Icons.Settings size={10} />
                                      <span>Manage Tags</span>
                                    </button>
                                  </div>
                                </div>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>

                      {/* Description Panel */}
                      <div className="flex-grow flex flex-col gap-2">
                        <span className="text-[10px] text-muted-foreground font-mono uppercase tracking-wider">Description</span>
                        <div className="bg-background/50 border border-border rounded p-4 font-mono text-sm whitespace-pre-wrap leading-relaxed overflow-y-auto min-h-[150px]">
                          {linkDetail.description || (
                            <span className="text-muted-foreground italic">
                              No description documented. Click Edit to add details.
                            </span>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ) : (
                <div className={styles.placeholder}>
                  <Icons.Link2 size={48} className="text-muted-foreground animate-pulse" />
                  <div className={styles.placeholderText}>
                    No link selected. Select a web link from the navigator or click "Add Link" to register a new destination.
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

      {/* Global Notes Form Modal */}
      <NoteForm
        isOpen={isNoteFormOpen}
        onClose={() => {
          setIsNoteFormOpen(false);
          setEditingNoteId(undefined);
        }}
        projectId={projectId}
        noteId={editingNoteId}
      />

      {/* Global Links Form Modal */}
      <LinkForm
        isOpen={isLinkFormOpen}
        onClose={() => {
          setIsLinkFormOpen(false);
          setEditingLinkId(undefined);
        }}
        projectId={projectId}
        linkId={editingLinkId}
      />

      {/* Global Credential Form Modal */}
      <CredentialForm
        isOpen={isCredentialFormOpen}
        onClose={() => {
          setIsCredentialFormOpen(false);
          setEditingCredentialId(undefined);
        }}
        projectId={projectId}
        credentialId={editingCredentialId}
      />

      {/* Global Credential Detail Modal */}
      {viewingCredentialId && (
        <CredentialDetailModal
          isOpen={!!viewingCredentialId}
          onClose={() => setViewingCredentialId(undefined)}
          projectId={projectId}
          credentialId={viewingCredentialId}
        />
      )}

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

      {/* Global Tags Manager Modal */}
      <TagsManagerModal
        isOpen={isTagsManagerOpen}
        onClose={() => setIsTagsManagerOpen(false)}
        projectId={projectId}
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
