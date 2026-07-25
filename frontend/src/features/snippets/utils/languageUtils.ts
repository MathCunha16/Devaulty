import type { SnippetLanguage } from "~types/api";

export const ALL_LANGUAGES: SnippetLanguage[] = [
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

export const mapLanguageToMonaco = (lang: SnippetLanguage): string => {
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
