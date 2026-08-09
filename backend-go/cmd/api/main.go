package main

import (
	"devaulty-backend/internal/adapter/in/web"
	"devaulty-backend/internal/adapter/in/web/handler"
	"devaulty-backend/internal/usecase"
	"log"
	"net"
	"os"

	"devaulty-backend/internal/adapter/out/persistence"

	"github.com/jmoiron/sqlx"
)

func main() {
	log.Println("Starting Devaulty API")
	db, err := persistence.InitDB("data/devaulty.db", "migrations")
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

	tagRepo := persistence.NewTagRepository(db)
	tagUseCase := usecase.NewTagUseCase(tagRepo, projectRepo)
	tagHandler := handler.NewTagHandler(tagUseCase)

	noteRepo := persistence.NewNoteRepository(db)
	noteUseCase := usecase.NewNoteUseCase(noteRepo, projectRepo, itemTagRepo)
	noteHandler := handler.NewNoteHandler(noteUseCase)

	itemTagUseCase := usecase.NewItemTagUseCase(itemTagRepo, tagRepo, projectRepo, snippetRepo, linkRepo, problemRepo, noteRepo)
	itemTagHandler := handler.NewItemTagHandler(itemTagUseCase, tagUseCase, projectUseCase)

	handlers := &web.Handlers{
		Project: projectHandler,
		Snippet: snippetHandler,
		Link:    linkHandler,
		Problem: problemHandler,
		Note:    noteHandler,
		Tag:     tagHandler,
		ItemTag: itemTagHandler,
	}
	r := web.SetupRouter(handlers, devaultyInternalToken)

	addr := "127.0.0.1:0"
	if os.Getenv("APP_ENV") == "dev" {
		addr = "127.0.0.1:8080"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to bind listener: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	log.Printf("DEVAULTY_API_BASE_URL=%d", port) // Exposes for TAURI

	log.Printf("Starting server on port %d...", port)
	if err := r.RunListener(listener); err != nil {
		log.Fatalf("Critical error while starting Devaulty API: %v", err)
	}
}
