package main

import (
	"log"

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
	log.Println("Devaulty Backend initialized successfully")
}
