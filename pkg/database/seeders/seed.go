package seeders

import (
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/logger"
)

func RunSeeders(db *database.PostgresDB, log logger.Logger) error {
	log.Println("Running database seeders...")

	if err := SeedUsers(db); err != nil {
		log.Fatalf("Error seeding users: %v", err)
		return err
	}

	if err := SeedAdmins(db); err != nil {
		log.Fatalf("error seeding admins: %v", err)
		return err
	}

	log.Println("Database seeding completed successfully")
	return nil
}
