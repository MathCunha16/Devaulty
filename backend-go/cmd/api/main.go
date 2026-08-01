package main

import (
	"devaulty-backend/internal/adapter/in/web"
	"devaulty-backend/internal/usecase"
	"log"
	"net"
	"os"

	"devaulty-backend/internal/adapter/out/persistence"

	"github.com/google/uuid"
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

	devaultyInternalToken := uuid.New().String()
	if os.Getenv("APP_ENV") == "dev" {
		devaultyInternalToken = "dev-token"
	}
	log.Printf("DEVAULTY_INTERNAL_TOKEN=%s", devaultyInternalToken) // Exposes for TAURI

	projectRepo := persistence.NewProjectRepository(db)
	projectUseCase := usecase.NewProjectUseCase(projectRepo)
	projectHandler := web.NewProjectHandler(projectUseCase)

	handlers := &web.Handlers{
		Project: projectHandler,
	}
	r := web.SetupRouter(handlers, devaultyInternalToken)

	addr := ":0"
	if os.Getenv("APP_ENV") == "dev" {
		addr = ":8080"
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
