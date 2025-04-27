package seeders

import (
	"context"
	"fmt"
	"github.com/agpprastyo/career-link/internal/user/entity"

	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/google/uuid"
)

func SeedAdmins(db *database.PostgresDB) error {

	// Super Admin
	adminSuperUserUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminSuperUserUUID UUID")
	}
	adminRoleSuperUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminRoleSuperUUID UUID")
	}
	adminRoleSuperPassword, _ := utils.HashPassword("adminsuper123")

	// Admin
	adminUserUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminUserUUID UUID")
	}
	adminRoleUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminRoleUUID UUID")
	}
	adminRolePassword, _ := utils.HashPassword("admin123")

	// admin viewer
	adminViewerUserUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminViewerUserUUID UUID")
	}
	adminViewerRoleUUID, err := uuid.NewV7()
	if err != nil {
		panic("failed to generate adminViewerRoleUUID UUID")
	}
	adminViewerRolePassword, _ := utils.HashPassword("adminviewer123")

	users := []entity.User{
		{
			ID:       adminSuperUserUUID,
			Username: "adminSuper",
			Email:    "adminsuper@example.com",
			Password: adminRoleSuperPassword,
			Role:     entity.AdminRole,
			IsActive: true,
		},
		{
			ID:       adminUserUUID,
			Username: "adminAdmin",
			Email:    "adminAdmin@example.com",
			Password: adminRolePassword,
			Role:     entity.AdminRole,
			IsActive: true,
		},
		{
			ID:       adminViewerUserUUID,
			Username: "adminViewer",
			Email:    "adminviewer@example.com",
			Password: adminViewerRolePassword,
			Role:     entity.AdminRole,
			IsActive: true,
		},
	}

	admin := []entity.Admin{
		{ID: adminRoleSuperUUID, UserID: adminSuperUserUUID, Role: entity.AdminRoleSuper, IsActive: true},
		{ID: adminRoleUUID, UserID: adminUserUUID, Role: entity.AdminRoleAdmin, IsActive: true},
		{ID: adminViewerRoleUUID, UserID: adminViewerUserUUID, Role: entity.AdminRoleViewer, IsActive: true},
	}

	// Begin transaction
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Defer rollback in case of error
	defer func() {
		if r := recover(); r != nil {
			err := tx.Rollback()
			if err != nil {
				return
			}
			panic(r) // re-throw panic after rollback
		}
	}()

	// Insert users
	for _, u := range users {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (id, username, email, password, role, is_active) 
             VALUES ($1, $2, $3, $4, $5, $6)
             ON CONFLICT (email) DO NOTHING`,
			u.ID, u.Username, u.Email, u.Password, u.Role, u.IsActive)
		if err != nil {
			err := tx.Rollback()
			if err != nil {
				return err
			}
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	// Insert admins
	// Insert admins
	for _, a := range admin {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO admins (id, user_id, role, is_active)
             VALUES ($1, $2, $3, $4)
             ON CONFLICT (id) DO NOTHING`, // Changed from user_id to id
			a.ID, a.UserID, a.Role, a.IsActive)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				return fmt.Errorf("rollback failed: %v, original error: %w", rollbackErr, err)
			}
			return fmt.Errorf("failed to insert admin: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}
