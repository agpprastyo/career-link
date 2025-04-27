package main

import (
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/database/seeders"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		// Using standard log initially just for this message
		// since our logger isn't created yet
		println("Warning: .env file not found, using environment variables")
	}

	// Load config
	cfg := config.Load()

	// Initialize logger early
	log := logger.New(cfg)

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Database, log)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func(db *database.PostgresDB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	// Run seeders
	if err := seeders.RunSeeders(db, *log); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
	}

	log.Info("Seeding completed successfully")
}
