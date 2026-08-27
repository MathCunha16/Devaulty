import React, { useState } from "react";
import { toast } from "sonner";
import * as Icons from "lucide-react";
import { CodeViewer } from "../../../components/CodeViewer";
import { ConfirmModal } from "../../../components/ConfirmModal";
import { TagManagerSection } from "../../../components/TagManagerSection";
import { formatUpdatedDate } from "../../../utils/dateUtils";
import { copyToClipboard } from "../../../utils/clipboardUtils";
import { ProblemForm } from "../components/ProblemForm";

import {
  useProblemsQuery,
  useProblemQuery,
  useUpdateProblemStatusMutation,
  useDeleteProblemMutation,
} from "../hooks/useProblems";
import type { ProblemStatus, ProblemSeverity } from "~types/api";
import styles from "../../../routes/projects.$projectId.module.css";

interface ProblemsWorkspaceProps {
  projectId: string;
  onOpenManageTagsModal: () => void;
  initialSelectedId?: string;
}

export const ProblemsWorkspace: React.FC<ProblemsWorkspaceProps> = ({
  projectId,
  onOpenManageTagsModal,
  initialSelectedId,
}) => {
  const { data: problemsData } = useProblemsQuery(projectId);
  const updateProblemStatusMutation = useUpdateProblemStatusMutation(projectId);
  const deleteProblemMutation = useDeleteProblemMutation(projectId);

  const [selectedProblemId, setSelectedProblemId] = useState<string | undefined>(initialSelectedId);
  const [problemSearchQuery, setProblemSearchQuery] = useState("");
  const [problemSeverityFilter, setProblemSeverityFilter] = useState<"ALL" | ProblemSeverity>("ALL");
  const [problemStatusFilter, setProblemStatusFilter] = useState<"ALL" | "UNRESOLVED" | ProblemStatus>("UNRESOLVED");

  React.useEffect(() => {
    if (initialSelectedId) {
      setSelectedProblemId(initialSelectedId);
    }
  }, [initialSelectedId]);

  const [isProblemFormOpen, setIsProblemFormOpen] = useState(false);
  const [editingProblemId, setEditingProblemId] = useState<string | undefined>(undefined);

  const [copiedId, setCopiedId] = useState<string | null>(null);

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

  const problems = problemsData?.content || [];

  const filteredProblems = problems.filter((p) => {
    const matchesSearch =
      p.title.toLowerCase().includes(problemSearchQuery.toLowerCase()) ||
      (p.tags && p.tags.some((t) => t.name.toLowerCase().includes(problemSearchQuery.toLowerCase())));
    const matchesSeverity = problemSeverityFilter === "ALL" || p.severity === problemSeverityFilter;

    let matchesStatus = true;
    if (problemStatusFilter === "UNRESOLVED") {
      matchesStatus = p.status === "OPEN" || p.status === "WORKING_ON";
    } else if (problemStatusFilter !== "ALL") {
      matchesStatus = p.status === problemStatusFilter;
    }

    return matchesSearch && matchesSeverity && matchesStatus;
  });

  const selectedProblem = problems.find((p) => p.id === selectedProblemId);
  const { data: problemDetail, isLoading: isLoadingProblemDetail } = useProblemQuery(
    projectId,
    selectedProblemId || ""
  );

  const handleCopy = async (content: string, id: string) => {

    const success = await copyToClipboard(content);
    if (success) {
      setCopiedId(id);
      toast.success("Copied to clipboard!");
      setTimeout(() => setCopiedId(null), 2000);
    } else {
      toast.error("Failed to copy content");
    }
  };


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

  return (
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
                    {formatUpdatedDate(problemDetail.updatedAt, problemDetail.createdAt) && (
                      <div className={styles.problemMetadataItem}>
                        <Icons.RefreshCw size={12} />
                        <span>
                          Updated:{" "}
                          {formatUpdatedDate(problemDetail.updatedAt, problemDetail.createdAt)}
                        </span>
                      </div>
                    )}
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

              {/* Tag Section */}
              <TagManagerSection
                itemId={problemDetail.id}
                itemType="PROBLEM"
                itemTags={problemDetail.tags}
                projectId={projectId}
                onOpenManageTagsModal={onOpenManageTagsModal}
                title="Diagnostic Labels & Tags"
              />

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
                      <CodeViewer
                        height="220px"
                        language="log"
                        code={problemDetail.errorDescription}
                      />
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
                      <CodeViewer
                        height="220px"
                        language="plaintext"
                        code={problemDetail.solution}
                      />
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
            <Icons.AlertCircle size={48} className="text-muted-foreground animate-pulse" />
            <div className={styles.placeholderText}>
              No diagnostic node selected. Select a problem from the list or click "Log Problem Node".
            </div>
          </div>
        )}
      </div>

      <ProblemForm
        isOpen={isProblemFormOpen}
        onClose={() => {
          setIsProblemFormOpen(false);
          setEditingProblemId(undefined);
        }}
        projectId={projectId}
        problemId={editingProblemId}
      />

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
    </>
  );
};
