package usecase

import (
	"context"
	"errors"
	"github.com/agpprastyo/career-link/internal/user/dto"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/internal/user/mapper"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/google/uuid"
)

// CreateAdmin creates a new admin users
func (uc *UserUseCase) CreateAdmin(ctx context.Context, req dto.CreateAdminRequest) (*dto.CreateAdminResponse, error) {
	// Check if users with the same email or username already exist
	_, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, repository.ErrUserAlreadyExists
	}

	_, err = uc.repo.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, repository.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		uc.log.WithError(err).Error("Failed to hash password")
		return nil, err
	}

	userUUID, err := uuid.NewV7()
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate UUID")
		return nil, err
	}

	// Create user
	usr := &entity.User{
		ID:       userUUID,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entity.AdminRole,
	}

	adminUUID, err := uuid.NewV7()
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate UUID")
		return nil, err
	}

	adm := &entity.Admin{
		ID:       adminUUID,
		UserID:   userUUID,
		Role:     entity.AdministratorRole(req.AdminRole),
		IsActive: true,
	}

	if err := uc.repo.CreateAdminUser(ctx, adm, usr); err != nil {
		uc.log.WithError(err).Error("Failed to create admin")
		return nil, err
	}

	return &dto.CreateAdminResponse{
		Data: dto.DataCreateAdminResponse{
			User:  mapper.ToUserDTO(usr),
			Admin: mapper.ToAdminDTO(adm),
		},
		Message: "Admin created successfully",
	}, nil
}

// ToggleAdminStatus toggles the admin status of a users
func (uc *UserUseCase) ToggleAdminStatus(ctx context.Context, userID uuid.UUID) (bool, error) {

	// Update user in the database

	adminStatus, err := uc.repo.UpdateAdminStatus(ctx, userID)
	if err != nil {
		uc.log.WithError(err).WithField("users_id", userID).Error("Failed to toggle admin status")
		return false, err
	}

	return adminStatus, nil
}

// SoftDeleteAdmin marks an admin users as deleted
func (uc *UserUseCase) SoftDeleteAdmin(ctx context.Context, userID uuid.UUID) error {
	// Get the admin to check their role
	admin, err := uc.repo.GetAdminByUserID(ctx, userID)
	if err != nil {
		uc.log.WithError(err).WithField("user_id", userID).Error("Failed to get admin data")
		return err
	}

	if admin == nil {
		return repository.ErrUserNotFound
	}

	// Prevent deletion of super admins
	if admin.Role == entity.AdminRoleSuper {
		return errors.New("cannot delete super admin accounts")
	}

	err = uc.repo.SoftDeleteAdmin(ctx, userID)
	if err != nil {
		uc.log.WithError(err).WithField("user_id", userID).Error("Failed to soft delete admin")
		return err
	}

	return nil
}
