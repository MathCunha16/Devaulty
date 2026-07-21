export type SnippetLanguage =
  | "BASH"
  | "FISH"
  | "ZSH"
  | "SH"
  | "POWERSHELL"
  | "BATCH"
  | "JAVA"
  | "KOTLIN"
  | "JAVASCRIPT"
  | "TYPESCRIPT"
  | "PYTHON"
  | "GO"
  | "RUST"
  | "C"
  | "CPP"
  | "CSHARP"
  | "PHP"
  | "RUBY"
  | "SWIFT"
  | "DART"
  | "SCALA"
  | "LUA"
  | "PERL"
  | "R"
  | "ELIXIR"
  | "HASKELL"
  | "CLOJURE"
  | "GROOVY"
  | "HTML"
  | "CSS"
  | "SCSS"
  | "LESS"
  | "JSX"
  | "TSX"
  | "VUE"
  | "SVELTE"
  | "JSON"
  | "YAML"
  | "XML"
  | "TOML"
  | "INI"
  | "ENV"
  | "CSV"
  | "MARKDOWN"
  | "PROPERTIES"
  | "DOCKERFILE"
  | "DOCKER_COMPOSE"
  | "NGINX"
  | "APACHE"
  | "TERRAFORM"
  | "ANSIBLE"
  | "KUBERNETES_YAML"
  | "HELM"
  | "MAKEFILE"
  | "CMAKE"
  | "GRADLE"
  | "MAVEN_POM"
  | "SQL"
  | "PLSQL"
  | "GRAPHQL"
  | "MONGODB"
  | "GITHUB_ACTIONS"
  | "GITLAB_CI"
  | "JENKINSFILE"
  | "REGEX"
  | "DIFF"
  | "LOG"
  | "PLAIN_TEXT";

export type SnippetType = "COMMAND" | "CODE";

export interface CreateProjectRequest {
  name: string;
  description?: string;
  icon?: string;
  color?: string;
}

export interface ProjectViewResponse {
  id: string; // uuid
  name: string;
  description?: string;
  icon?: string;
  color?: string;
  archived: boolean;
  createdAt: string; // date-time
  updatedAt: string; // date-time
}

export interface CreateSnippetRequest {
  title: string;
  description?: string;
  content: string;
  language: SnippetLanguage;
  snippetType: SnippetType;
}

export interface SnippetViewResponse {
  id: string; // uuid
  projectId: string; // uuid
  title: string;
  description?: string;
  content: string;
  language: SnippetLanguage;
  snippetType: SnippetType;
  tags?: TagSummaryResponse[];
  createdAt: string; // date-time
  updatedAt: string; // date-time
}

export interface UpdateSnippetRequest {
  title?: string;
  description?: string;
  content?: string;
  language?: SnippetLanguage;
  snippetType?: SnippetType;
}

export interface UpdateProjectRequest {
  name?: string;
  description?: string;
  icon?: string;
  color?: string;
}

export interface PageMetadata {
  size: number;
  number: number;
  totalElements: number;
  totalPages: number;
}

export interface PagedModelProjectViewResponse {
  content: ProjectViewResponse[];
  page: PageMetadata;
}

export interface PagedModelSnippetViewResponse {
  content: SnippetViewResponse[];
  page: PageMetadata;
}

export interface ApiErrorResponse {
  status: number;
  message: string;
  timestamp: string;
  errors?: Array<{
    field: string;
    message: string;
  }>;
}

export type ProblemStatus = "OPEN" | "WORKING_ON" | "RESOLVED" | "WONT_FIX";
export type ProblemSeverity = "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";

export interface TagSummaryResponse {
  id: string; // uuid
  name: string;
  color?: string;
}

export interface TagViewResponse {
  id: string; // uuid
  projectId: string; // uuid
  name: string;
  color?: string;
  createdAt: string; // date-time
  updatedAt: string; // date-time
}

export interface CreateTagRequest {
  name: string;
  color?: string;
}

export interface CreateProblemRequest {
  title: string;
  errorDescription?: string;
  solution?: string;
  status: ProblemStatus;
  severity: ProblemSeverity;
}

export interface ProblemViewResponse {
  id: string; // uuid
  projectId: string; // uuid
  title: string;
  errorDescription?: string;
  solution?: string;
  status: ProblemStatus;
  severity: ProblemSeverity;
  tags: TagSummaryResponse[];
  createdAt: string; // date-time
  updatedAt: string; // date-time
}

export interface ProblemSummaryResponse {
  id: string; // uuid
  projectId: string; // uuid
  title: string;
  status: ProblemStatus;
  severity: ProblemSeverity;
  tags: TagSummaryResponse[];
  createdAt: string; // date-time
  updatedAt: string; // date-time
}

export interface UpdateProblemRequest {
  title?: string;
  errorDescription?: string;
  solution?: string;
  severity?: ProblemSeverity;
}

export interface UpdateProblemStatusRequest {
  status: ProblemStatus;
}

export interface PagedModelProblemSummaryResponse {
  content: ProblemSummaryResponse[];
  page: PageMetadata;
}

