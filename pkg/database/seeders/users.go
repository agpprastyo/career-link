package seeders

import (
	"context"
	"fmt"

	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/google/uuid"
)

func SeedUsers(db *database.PostgresDB) error {
	// Admin user
	adminID, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	adminPassword, _ := utils.HashPassword("admin123")

	// Company user
	companyID, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	companyPassword, _ := utils.HashPassword("company123")

	// Job seeker
	seekerID, err := uuid.NewV7()
	if err != nil {
		fmt.Println(err)
	}
	seekerPassword, _ := utils.HashPassword("seeker123")

	users := []entity.User{
		{
			ID:       adminID,
			Username: "admin",
			Email:    "admin@example.com",
			Password: string(adminPassword),
			Role:     entity.AdminRole,
			IsActive: true,
		},
		{
			ID:       companyID,
			Username: "company",
			Email:    "company@example.com",
			Password: string(companyPassword),
			Role:     entity.CompanyRole,
			IsActive: true,
		},
		{
			ID:       seekerID,
			Username: "jobSeeker",
			Email:    "seeker@example.com",
			Password: string(seekerPassword),
			Role:     entity.JobSeekerRole,
			IsActive: true,
		},
		{
			ID:       seekerID,
			Username: "nogahe1533",
			Email:    "nogahe1533@clubemp.com",
			Password: string(seekerPassword),
			Role:     entity.JobSeekerRole,
			IsActive: false,
		},
	}

	for _, u := range users {
		_, err := db.ExecContext(context.Background(),
			`INSERT INTO users (id, username, email, password, role, is_active) 
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (email) DO NOTHING`,
			u.ID, u.Username, u.Email, u.Password, u.Role, u.IsActive)
		if err != nil {
			return err
		}
	}

	return nil
}
