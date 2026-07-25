import React, { useState, useRef, useEffect } from "react";
import * as Icons from "lucide-react";
import type { SnippetLanguage } from "~types/api";
import styles from "./LanguageSelect.module.css";

const LANGUAGE_MAP: Record<SnippetLanguage, { label: string; icon: React.ComponentType<{ size?: number; className?: string }> }> = {
  PLAIN_TEXT: { label: "Plain Text", icon: Icons.FileText },
  JAVASCRIPT: { label: "JavaScript (JS)", icon: Icons.FileCode2 },
  TYPESCRIPT: { label: "TypeScript (TS)", icon: Icons.FileCode },
  JSX: { label: "React JSX", icon: Icons.Code2 },
  TSX: { label: "React TSX", icon: Icons.Code2 },
  PYTHON: { label: "Python", icon: Icons.FileCode },
  JAVA: { label: "Java", icon: Icons.Coffee },
  KOTLIN: { label: "Kotlin", icon: Icons.FileCode },
  GO: { label: "Go (Golang)", icon: Icons.Cpu },
  RUST: { label: "Rust", icon: Icons.ShieldAlert },
  C: { label: "C", icon: Icons.Cpu },
  CPP: { label: "C++", icon: Icons.Cpu },
  CSHARP: { label: "C#", icon: Icons.FileCode },
  BASH: { label: "Bash Shell", icon: Icons.Terminal },
  ZSH: { label: "Zsh Shell", icon: Icons.Terminal },
  FISH: { label: "Fish Shell", icon: Icons.Terminal },
  SH: { label: "Shell Script (SH)", icon: Icons.Terminal },
  POWERSHELL: { label: "PowerShell", icon: Icons.Terminal },
  BATCH: { label: "Batch File (BAT)", icon: Icons.Terminal },
  PHP: { label: "PHP", icon: Icons.FileCode },
  RUBY: { label: "Ruby", icon: Icons.FileCode },
  SWIFT: { label: "Swift", icon: Icons.FileCode },
  DART: { label: "Dart", icon: Icons.FileCode },
  SCALA: { label: "Scala", icon: Icons.FileCode },
  LUA: { label: "Lua", icon: Icons.FileCode },
  PERL: { label: "Perl", icon: Icons.FileCode },
  R: { label: "R", icon: Icons.FileCode },
  ELIXIR: { label: "Elixir", icon: Icons.FileCode },
  HASKELL: { label: "Haskell", icon: Icons.FileCode },
  CLOJURE: { label: "Clojure", icon: Icons.FileCode },
  GROOVY: { label: "Groovy", icon: Icons.FileCode },
  HTML: { label: "HTML5", icon: Icons.Globe },
  CSS: { label: "CSS3", icon: Icons.Code2 },
  SCSS: { label: "SCSS", icon: Icons.Code2 },
  LESS: { label: "LESS", icon: Icons.Code2 },
  VUE: { label: "Vue.js", icon: Icons.Code2 },
  SVELTE: { label: "Svelte", icon: Icons.Code2 },
  JSON: { label: "JSON", icon: Icons.Braces },
  YAML: { label: "YAML", icon: Icons.FileJson },
  XML: { label: "XML", icon: Icons.FileJson },
  TOML: { label: "TOML", icon: Icons.FileJson },
  INI: { label: "INI Config", icon: Icons.Settings },
  ENV: { label: "Env (.env)", icon: Icons.Settings },
  CSV: { label: "CSV Data", icon: Icons.FileSpreadsheet },
  MARKDOWN: { label: "Markdown (MD)", icon: Icons.FileText },
  PROPERTIES: { label: "Properties", icon: Icons.Settings },
  DOCKERFILE: { label: "Dockerfile", icon: Icons.Boxes },
  DOCKER_COMPOSE: { label: "Docker Compose", icon: Icons.Boxes },
  NGINX: { label: "Nginx Config", icon: Icons.Server },
  APACHE: { label: "Apache Config", icon: Icons.Server },
  TERRAFORM: { label: "Terraform (HCL)", icon: Icons.Layers },
  ANSIBLE: { label: "Ansible", icon: Icons.Layers },
  KUBERNETES_YAML: { label: "Kubernetes YAML", icon: Icons.Boxes },
  HELM: { label: "Helm Chart", icon: Icons.Boxes },
  MAKEFILE: { label: "Makefile", icon: Icons.Wrench },
  CMAKE: { label: "CMake", icon: Icons.Wrench },
  GRADLE: { label: "Gradle", icon: Icons.Wrench },
  MAVEN_POM: { label: "Maven POM", icon: Icons.Wrench },
  SQL: { label: "SQL", icon: Icons.Database },
  PLSQL: { label: "PL/SQL", icon: Icons.Database },
  GRAPHQL: { label: "GraphQL", icon: Icons.Database },
  MONGODB: { label: "MongoDB Query", icon: Icons.Database },
  GITHUB_ACTIONS: { label: "GitHub Actions", icon: Icons.GitBranch },
  GITLAB_CI: { label: "GitLab CI", icon: Icons.GitBranch },
  JENKINSFILE: { label: "Jenkinsfile", icon: Icons.GitBranch },
  REGEX: { label: "Regular Expression", icon: Icons.Regex },
  DIFF: { label: "Git Diff", icon: Icons.GitCompare },
  LOG: { label: "Log File", icon: Icons.FileText },
};

interface LanguageSelectProps {
  languages: SnippetLanguage[];
  value: SnippetLanguage;
  onChange: (value: SnippetLanguage) => void;
  disabled?: boolean;
  triggerId?: string;
  ariaLabelledBy?: string;
}

export const LanguageSelect: React.FC<LanguageSelectProps> = ({
  languages,
  value,
  onChange,
  disabled = false,
  triggerId,
  ariaLabelledBy,
}) => {
  const [isOpen, setIsOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [focusedIndex, setFocusedIndex] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const searchInputRef = useRef<HTMLInputElement>(null);

  const selectedOption = LANGUAGE_MAP[value] || { label: value, icon: Icons.Code2 };
  const SelectedIcon = selectedOption.icon;

  const filteredLanguages = languages.filter((lang) => {
    const info = LANGUAGE_MAP[lang];
    const label = info ? info.label : lang;
    const q = searchQuery.toLowerCase().trim();
    return label.toLowerCase().includes(q) || lang.toLowerCase().includes(q);
  });

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  useEffect(() => {
    if (isOpen) {
      const timer = setTimeout(() => searchInputRef.current?.focus(), 50);
      return () => clearTimeout(timer);
    }
  }, [isOpen]);

  const handleSelect = (lang: SnippetLanguage) => {
    onChange(lang);
    setIsOpen(false);
  };

  const handleToggleOpen = () => {
    if (disabled) return;
    if (!isOpen) {
      setSearchQuery("");
      setFocusedIndex(0);
    }
    setIsOpen(!isOpen);
  };


  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") {
      setIsOpen(false);
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      setFocusedIndex((prev) => (filteredLanguages.length > 0 ? (prev + 1) % filteredLanguages.length : 0));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setFocusedIndex((prev) =>
        filteredLanguages.length > 0 ? (prev - 1 + filteredLanguages.length) % filteredLanguages.length : 0
      );
    } else if (e.key === "Enter" && filteredLanguages.length > 0) {
      e.preventDefault();
      const targetLang = filteredLanguages[focusedIndex] || filteredLanguages[0];
      handleSelect(targetLang);
    }
  };

  return (
    <div className={styles.container} ref={containerRef}>
      <button
        id={triggerId}
        type="button"
        className={styles.trigger}
        onClick={handleToggleOpen}
        disabled={disabled}
        aria-labelledby={ariaLabelledBy}
        aria-expanded={isOpen}
        aria-haspopup="listbox"
      >
        <div className={styles.triggerContent}>
          <SelectedIcon size={14} className={styles.icon} />
          <span className={styles.label}>{selectedOption.label}</span>
        </div>
        <Icons.ChevronDown
          size={14}
          className={`${styles.arrow} ${isOpen ? styles.arrowOpen : ""}`}
        />
      </button>

      {isOpen && (
        <div className={styles.dropdown}>
          <div className={styles.searchBox}>
            <Icons.Search size={14} className={styles.searchIcon} />
            <input
              ref={searchInputRef}
              type="text"
              className={styles.searchInput}
              placeholder="Search language (e.g., Python, TS, Bash, SQL...)"
              value={searchQuery}
              onChange={(e) => {
                setSearchQuery(e.target.value);
                setFocusedIndex(0);
              }}
              onKeyDown={handleKeyDown}
            />
          </div>

          <div className={styles.list} role="listbox">
            {filteredLanguages.length === 0 ? (
              <div className={styles.empty}>No matching languages found</div>
            ) : (
              filteredLanguages.map((lang, idx) => {
                const info = LANGUAGE_MAP[lang] || { label: lang, icon: Icons.Code2 };
                const IconComp = info.icon;
                const isSelected = lang === value;
                const isFocused = idx === focusedIndex;

                return (
                  <button
                    key={lang}
                    type="button"
                    role="option"
                    aria-selected={isSelected}
                    className={`${styles.item} ${isSelected ? styles.itemSelected : ""} ${
                      isFocused ? styles.itemFocused : ""
                    }`}
                    onClick={() => handleSelect(lang)}
                  >
                    <div className={styles.itemLeft}>
                      <IconComp size={14} />
                      <span>{info.label}</span>
                    </div>
                    {isSelected && <Icons.Check size={12} />}
                  </button>
                );
              })
            )}
          </div>
        </div>
      )}
    </div>
  );
};
