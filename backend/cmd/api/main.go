package main

import (
	"devaulty-backend/internal/adapter/in/mcp"
	"devaulty-backend/internal/adapter/in/scheduler"
	"devaulty-backend/internal/adapter/in/web"
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/usecase"
	"devaulty-backend/migrations"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"devaulty-backend/internal/adapter/out/persistence"
	"devaulty-backend/internal/adapter/out/security"

	"github.com/jmoiron/sqlx"
)

func main() {
	if os.Getenv("APP_ENV") != "dev" {
		debug.SetGCPercent(20)
	}
	log.Println("Starting Devaulty API")

	dataDir, err := resolveDataDir()
	if err != nil {
		log.Fatalf("Refusing to start Devaulty with an invalid or unstable data directory: %v", err)
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("Critical error while creating data directory %s: %v", dataDir, err)
	}

	dbPath := filepath.Join(dataDir, "devaulty.db")

	// Use embedded migrations in production for a self-contained binary.
	// In development (APP_ENV=dev), use file-based migrations for hot-reload convenience.
	var db *sqlx.DB
	if os.Getenv("APP_ENV") == "dev" {
		db, err = persistence.InitDB(dbPath, "migrations")
	} else {
		db, err = persistence.InitDBWithFS(dbPath, migrations.FS)
	}
	if err != nil {
		log.Fatalf("Critical error while initializing Devaulty database: %v", err)
	}
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Critical error while closing Devaulty database: %v", err)
		}
	}(db)
	log.Println("Devaulty DB initialized successfully")

	devaultyInternalToken := os.Getenv("DEVAULTY_INTERNAL_TOKEN")
	if devaultyInternalToken == "" {
		devaultyInternalToken = "dev-token" // Default token for development
	}

	appSettingRepo := persistence.NewAppSettingRepository(db)
	masterKeySession := security.NewMasterKeySessionHolder()
	keyDeriver := security.NewArgon2KeyDeriver()
	vaultUseCase := usecase.NewVaultUseCase(keyDeriver, masterKeySession, appSettingRepo)
	securityHandler := handler.NewSecurityHandler(vaultUseCase)

	itemTagRepo := persistence.NewItemTagRepository(db)

	projectRepo := persistence.NewProjectRepository(db)
	projectUseCase := usecase.NewProjectUseCase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUseCase)

	snippetRepo := persistence.NewSnippetRepository(db)
	snippetUseCase := usecase.NewSnippetUseCase(snippetRepo, projectRepo, itemTagRepo)
	snippetHandler := handler.NewSnippetHandler(snippetUseCase)

	linkRepo := persistence.NewLinkRepository(db)
	linkUseCase := usecase.NewLinkUseCase(linkRepo, projectRepo, itemTagRepo)
	linkHandler := handler.NewLinkHandler(linkUseCase)

	problemRepo := persistence.NewProblemRepository(db)
	problemUseCase := usecase.NewProblemUseCase(problemRepo, projectRepo, itemTagRepo)
	problemHandler := handler.NewProblemHandler(problemUseCase)

	noteRepo := persistence.NewNoteRepository(db)
	noteUseCase := usecase.NewNoteUseCase(noteRepo, projectRepo, itemTagRepo)
	noteHandler := handler.NewNoteHandler(noteUseCase)

	boardRepo := persistence.NewBoardRepository(db)
	boardColumnRepo := persistence.NewBoardColumnRepository(db)
	boardUseCase := usecase.NewBoardUseCase(boardRepo, boardColumnRepo, projectRepo, itemTagRepo)
	boardHandler := handler.NewBoardHandler(boardUseCase)

	boardColumnUseCase := usecase.NewBoardColumnUseCase(boardColumnRepo, boardRepo, projectRepo)
	boardColumnHandler := handler.NewBoardColumnHandler(boardColumnUseCase)

	cardRepo := persistence.NewCardRepository(db)
	cardUseCase := usecase.NewCardUseCase(cardRepo, boardRepo, boardColumnRepo, projectRepo, itemTagRepo)
	cardHandler := handler.NewCardHandler(cardUseCase)

	cryptoAdapter := security.NewAESGCMCrypto()
	credentialRepo := persistence.NewCredentialRepository(db)
	credentialUseCase := usecase.NewCredentialUseCase(credentialRepo, projectRepo, itemTagRepo, cryptoAdapter, masterKeySession, *vaultUseCase)
	credentialHandler := handler.NewCredentialHandler(credentialUseCase)

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo, projectRepo)
	tagHandler := handler.NewTagHandler(tagUseCase)

	itemTagUseCase := usecase.NewItemTagUseCase(itemTagRepo, tagRepo, projectRepo, snippetRepo, credentialRepo, linkRepo, problemRepo, noteRepo, boardRepo, cardRepo)
	itemTagHandler := handler.NewItemTagHandler(itemTagUseCase, tagUseCase, projectUseCase)

	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		runMCPServer(projectUseCase, snippetUseCase, problemUseCase, linkUseCase, noteUseCase, boardUseCase, boardColumnUseCase, cardUseCase, tagUseCase, itemTagUseCase)
		return
	}

	handlers := &web.Handlers{
		Project:     projectHandler,
		Snippet:     snippetHandler,
		Credential:  credentialHandler,
		Link:        linkHandler,
		Problem:     problemHandler,
		Note:        noteHandler,
		Board:       boardHandler,
		BoardColumn: boardColumnHandler,
		Card:        cardHandler,
		Tag:         tagHandler,
		ItemTag:     itemTagHandler,
		Security:    securityHandler,
	}
	r := web.SetupRouter(handlers, devaultyInternalToken)

	autoLock := scheduler.NewVaultAutoLock(masterKeySession)
	go autoLock.PurgeExpiredSession()

	addr := "127.0.0.1:0"
	if os.Getenv("APP_ENV") == "dev" {
		addr = "127.0.0.1:8080"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port

	// Structured session handshake line for Tauri IPC.
	// Tauri reads stdout via a private pipe and parses this exact format
	// to extract the dynamic port and security token.
	fmt.Printf("[DEVAULTY_SESSION] PORT=%d TOKEN=%s\n", port, devaultyInternalToken)

	log.Printf("Starting server on port %d...", port)
	if err := r.RunListener(listener); err != nil {
		log.Fatalf("Critical error while starting Devaulty API: %v", err)
	}
}

// runMCPServer starts the MCP server.
func runMCPServer(
	projectUseCase *usecase.ProjectUseCase,
	snippetUseCase *usecase.SnippetUseCase,
	problemUseCase *usecase.ProblemUseCase,
	linkUseCase *usecase.LinkUseCase,
	noteUseCase *usecase.NoteUseCase,
	boardUseCase *usecase.BoardUseCase,
	boardColumnUseCase *usecase.BoardColumnUseCase,
	cardUseCase *usecase.CardUseCase,
	tagUseCase *usecase.TagUseCase,
	itemTagUseCase *usecase.ItemTagUseCase) {
	fs := flag.NewFlagSet("mcp", flag.ExitOnError)
	readOnly := fs.Bool("readonly", false, "only register ready-only tools")
	disableDelete := fs.Bool("disable-delete", false, "disable delete tool")
	if err := fs.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Failed to parse command line arguments: %v", err)
	}

	opts := mcp.Options{
		ReadOnly:      *readOnly,
		DisableDelete: *disableDelete,
	}

	adapter := mcp.NewMCPServerAdapter(opts, projectUseCase, snippetUseCase, problemUseCase, linkUseCase, noteUseCase, boardUseCase, boardColumnUseCase, cardUseCase, tagUseCase, itemTagUseCase)
	if err := adapter.Serve(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

// resolveDataDir returns a stable, absolute app data directory and rejects
// transient paths derived from AppImage mounts or other ephemeral storage.
func resolveDataDir() (string, error) {
	candidate := os.Getenv("DEVAULTY_DATA_DIR")
	if candidate == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", errors.New("failed to resolve user config directory")
		}
		candidate = filepath.Join(configDir, "devaulty")
	}

	absPath, err := filepath.Abs(candidate)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path for %q: %w", candidate, err)
	}

	if resolved, err := filepath.EvalSymlinks(absPath); err == nil {
		absPath = resolved
	}

	absPath = filepath.Clean(absPath)

	normalizedAbsPath := normalizePathForComparison(absPath)
	tempDir, err := filepath.Abs(os.TempDir())
	if err == nil {
		tempDir = filepath.Clean(tempDir)
	}
	normalizedTempDir := normalizePathForComparison(tempDir)
	if normalizedTempDir != "" {
		if normalizedAbsPath == normalizedTempDir || strings.HasPrefix(normalizedAbsPath, normalizedTempDir+"/") {
			return "", fmt.Errorf("refusing to use unstable data directory %q (AppImage mounts and other temp directories are not supported for persistent storage)", absPath)
		}
	}

	lowerAbsPath := strings.ToLower(filepath.ToSlash(absPath))
	if strings.Contains(lowerAbsPath, "/tmp/") ||
		strings.Contains(lowerAbsPath, ".mount_") ||
		strings.Contains(lowerAbsPath, "/appimage") ||
		strings.Contains(lowerAbsPath, "\\appimage") {
		return "", fmt.Errorf("refusing to use unstable data directory %q (AppImage mounts and other temp directories are not supported for persistent storage)", absPath)
	}

	return absPath, nil
}

func normalizePathForComparison(p string) string {
	cleaned := filepath.Clean(p)
	cleaned = filepath.ToSlash(cleaned)
	cleaned = strings.TrimSuffix(cleaned, "/")
	if cleaned == "." {
		return ""
	}
	return strings.ToLower(cleaned)
}
