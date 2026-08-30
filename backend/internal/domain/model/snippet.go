package model

import "github.com/google/uuid"

type SnippetType string

const (
	SnippetTypeCommand SnippetType = "COMMAND"
	SnippetTypeCode    SnippetType = "CODE"
)

type SnippetLanguage string

const (
	// Shells / terminal
	SnippetLangBash       SnippetLanguage = "BASH"
	SnippetLangFish       SnippetLanguage = "FISH"
	SnippetLangZsh        SnippetLanguage = "ZSH"
	SnippetLangSh         SnippetLanguage = "SH"
	SnippetLangPowershell SnippetLanguage = "POWERSHELL"
	SnippetLangBatch      SnippetLanguage = "BATCH"

	// Code Languages
	SnippetLangJava       SnippetLanguage = "JAVA"
	SnippetLangKotlin     SnippetLanguage = "KOTLIN"
	SnippetLangJavascript SnippetLanguage = "JAVASCRIPT"
	SnippetLangTypescript SnippetLanguage = "TYPESCRIPT"
	SnippetLangPython     SnippetLanguage = "PYTHON"
	SnippetLangGo         SnippetLanguage = "GO"
	SnippetLangRust       SnippetLanguage = "RUST"
	SnippetLangC          SnippetLanguage = "C"
	SnippetLangCpp        SnippetLanguage = "CPP"
	SnippetLangCsharp     SnippetLanguage = "CSHARP"
	SnippetLangPhp        SnippetLanguage = "PHP"
	SnippetLangRuby       SnippetLanguage = "RUBY"
	SnippetLangSwift      SnippetLanguage = "SWIFT"
	SnippetLangDart       SnippetLanguage = "DART"
	SnippetLangScala      SnippetLanguage = "SCALA"
	SnippetLangLua        SnippetLanguage = "LUA"
	SnippetLangPerl       SnippetLanguage = "PERL"
	SnippetLangR          SnippetLanguage = "R"
	SnippetLangElixir     SnippetLanguage = "ELIXIR"
	SnippetLangHaskell    SnippetLanguage = "HASKELL"
	SnippetLangClojure    SnippetLanguage = "CLOJURE"
	SnippetLangGroovy     SnippetLanguage = "GROOVY"

	// Web / frontend
	SnippetLangHtml   SnippetLanguage = "HTML"
	SnippetLangCss    SnippetLanguage = "CSS"
	SnippetLangScss   SnippetLanguage = "SCSS"
	SnippetLangLess   SnippetLanguage = "LESS"
	SnippetLangJsx    SnippetLanguage = "JSX"
	SnippetLangTsx    SnippetLanguage = "TSX"
	SnippetLangVue    SnippetLanguage = "VUE"
	SnippetLangSvelte SnippetLanguage = "SVELTE"

	// Data / config / markup
	SnippetLangJson       SnippetLanguage = "JSON"
	SnippetLangYaml       SnippetLanguage = "YAML"
	SnippetLangXml        SnippetLanguage = "XML"
	SnippetLangToml       SnippetLanguage = "TOML"
	SnippetLangIni        SnippetLanguage = "INI"
	SnippetLangEnv        SnippetLanguage = "ENV"
	SnippetLangCsv        SnippetLanguage = "CSV"
	SnippetLangMarkdown   SnippetLanguage = "MARKDOWN"
	SnippetLangProperties SnippetLanguage = "PROPERTIES"

	// Infra / DevOps
	SnippetLangDockerfile     SnippetLanguage = "DOCKERFILE"
	SnippetLangDockerCompose  SnippetLanguage = "DOCKER_COMPOSE"
	SnippetLangNginx          SnippetLanguage = "NGINX"
	SnippetLangApache         SnippetLanguage = "APACHE"
	SnippetLangTerraform      SnippetLanguage = "TERRAFORM"
	SnippetLangAnsible        SnippetLanguage = "ANSIBLE"
	SnippetLangKubernetesYaml SnippetLanguage = "KUBERNETES_YAML"
	SnippetLangHelm           SnippetLanguage = "HELM"
	SnippetLangMakefile       SnippetLanguage = "MAKEFILE"
	SnippetLangCmake          SnippetLanguage = "CMAKE"
	SnippetLangGradle         SnippetLanguage = "GRADLE"
	SnippetLangMavenPom       SnippetLanguage = "MAVEN_POM"

	// Databases
	SnippetLangSql     SnippetLanguage = "SQL"
	SnippetLangPlsql   SnippetLanguage = "PLSQL"
	SnippetLangGraphql SnippetLanguage = "GRAPHQL"
	SnippetLangMongodb SnippetLanguage = "MONGODB"

	// CI/CD
	SnippetLangGithubActions SnippetLanguage = "GITHUB_ACTIONS"
	SnippetLangGitlabCi      SnippetLanguage = "GITLAB_CI"
	SnippetLangJenkinsfile   SnippetLanguage = "JENKINSFILE"

	// Others / Generics
	SnippetLangRegex     SnippetLanguage = "REGEX"
	SnippetLangDiff      SnippetLanguage = "DIFF"
	SnippetLangLog       SnippetLanguage = "LOG"
	SnippetLangPlainText SnippetLanguage = "PLAIN_TEXT"
)

var SnippetTypes = []SnippetType{ // this is used on MCP
	SnippetTypeCommand,
	SnippetTypeCode,
}

var SnippetLanguages = []SnippetLanguage{ // this is used on MCP
	// Shells / terminal
	SnippetLangBash,
	SnippetLangFish,
	SnippetLangZsh,
	SnippetLangSh,
	SnippetLangPowershell,
	SnippetLangBatch,

	// Code Languages
	SnippetLangJava,
	SnippetLangKotlin,
	SnippetLangJavascript,
	SnippetLangTypescript,
	SnippetLangPython,
	SnippetLangGo,
	SnippetLangRust,
	SnippetLangC,
	SnippetLangCpp,
	SnippetLangCsharp,
	SnippetLangPhp,
	SnippetLangRuby,
	SnippetLangSwift,
	SnippetLangDart,
	SnippetLangScala,
	SnippetLangLua,
	SnippetLangPerl,
	SnippetLangR,
	SnippetLangElixir,
	SnippetLangHaskell,
	SnippetLangClojure,
	SnippetLangGroovy,

	// Web / frontend
	SnippetLangHtml,
	SnippetLangCss,
	SnippetLangScss,
	SnippetLangLess,
	SnippetLangJsx,
	SnippetLangTsx,
	SnippetLangVue,
	SnippetLangSvelte,

	// Data / config / markup
	SnippetLangJson,
	SnippetLangYaml,
	SnippetLangXml,
	SnippetLangToml,
	SnippetLangIni,
	SnippetLangEnv,
	SnippetLangCsv,
	SnippetLangMarkdown,
	SnippetLangProperties,

	// Infra / DevOps
	SnippetLangDockerfile,
	SnippetLangDockerCompose,
	SnippetLangNginx,
	SnippetLangApache,
	SnippetLangTerraform,
	SnippetLangAnsible,
	SnippetLangKubernetesYaml,
	SnippetLangHelm,
	SnippetLangMakefile,
	SnippetLangCmake,
	SnippetLangGradle,
	SnippetLangMavenPom,

	// Databases
	SnippetLangSql,
	SnippetLangPlsql,
	SnippetLangGraphql,
	SnippetLangMongodb,

	// CI/CD
	SnippetLangGithubActions,
	SnippetLangGitlabCi,
	SnippetLangJenkinsfile,

	// Others / Generics
	SnippetLangRegex,
	SnippetLangDiff,
	SnippetLangLog,
	SnippetLangPlainText,
}

type Snippet struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	ProjectID   uuid.UUID        `json:"projectId" db:"project_id"`
	Title       string           `json:"title" db:"title"`
	Description *string          `json:"description,omitempty" db:"description"`
	Content     string           `json:"content" db:"content"`
	Language    *SnippetLanguage `json:"language,omitempty" db:"language"`
	SnippetType SnippetType      `json:"snippetType" db:"snippet_type"`
	BaseEntity
}
