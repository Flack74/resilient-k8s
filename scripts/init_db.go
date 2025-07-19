package main

import (
	"log"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/config"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to the database
	db, err := storage.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize the schema
	log.Println("Initializing database schema...")
	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	log.Println("Database schema initialized successfully")
}