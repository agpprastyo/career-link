package usecase

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/agpprastyo/career-link/internal/common/pagination"

	"github.com/agpprastyo/career-link/internal/user/dto"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/google/uuid"
	"sync"
	"time"
)

// GetUsers retrieves a paginated list of users
func (uc *UserUseCase) GetUsers(ctx context.Context, paging pagination.Pagination) (pagination.PageResponse, error) {
	// Get users with pagination from repository
	users, total, err := uc.repo.GetUsersPaginated(ctx, paging)
	if err != nil {
		return pagination.PageResponse{}, err
	}

	// Create paginated response
	response := pagination.NewResponse(users, paging, total)
	return response, nil
}

// UpdateAvatar updates a user's avatar in the database
func (uc *UserUseCase) UpdateAvatar(ctx context.Context, req dto.UpdateAvatarRequest) (avatarURL string, err error) {
	// Parse user ID from string to UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		uc.log.WithError(err).Error("Invalid user ID format")
		return "", repository.ErrInvalidInput
	}

	// Get the user from the database
	userData, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get user for avatar update")
		return "", repository.ErrUserNotFound
	}

	if userData.Avatar != nil {
		bgCtx := context.Background()
		go func() {
			if err := uc.repo.DeleteAvatarFile(bgCtx, *userData.Avatar); err != nil {
				uc.log.WithError(err).Error("Failed to delete old avatar file")
			}
		}()
	}

	fileContent := bytes.NewReader(req.FileContent)
	filename := fmt.Sprintf("avatars/%s", req.FileName)

	// Upload avatar to MinIO
	avatarURL, err = uc.repo.UploadAvatarFile(ctx, filename, fileContent, req.FileSize, req.ContentType)
	if err != nil {
		uc.log.WithError(err).Error("Failed to upload avatar file")
		return "", err
	}

	// Update avatar URL in user data
	userData.Avatar = &req.FileName
	// Update user in the database
	err = uc.repo.UpdateUser(ctx, userData)
	if err != nil {
		uc.log.WithError(err).Error("Failed to update avatar in database")
		return "", err
	}

	return avatarURL, nil
}

func (uc *UserUseCase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// Parse user ID from string to UUID
	email := req.Email

	// Get the user from the database
	userData, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get user for password reset")
		return repository.ErrUserNotFound
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		uc.log.WithError(err).Error("Failed to hash new password")
		return err
	}

	// Update password in the database
	userData.Password = hashedPassword
	err = uc.repo.UpdateUser(ctx, userData)
	if err != nil {
		uc.log.WithError(err).Error("Failed to update password in database")
		return err
	}

	return nil
}

// VerifyPasswordResetToken checks if the password reset token is valid
func (uc *UserUseCase) VerifyPasswordResetToken(ctx context.Context, token string) (bool, error) {
	t, err := uc.repo.GetToken(ctx, token)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get token")
		return false, err
	}

	if t.Type != string(entity.PasswordReset) {
		return false, repository.ErrInvalidToken
	}

	if t.ExpiredAt.Before(time.Now()) {
		return false, repository.ErrTokenExpired
	}

	return true, nil
}

// UpdatePassword changes a user's password after verifying the current password
func (uc *UserUseCase) UpdatePassword(ctx context.Context, req dto.UpdatePasswordRequest) error {
	// Parse user ID from string to UUID
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		uc.log.WithError(err).Error("Invalid user ID format")
		return repository.ErrInvalidInput
	}

	// Get the user from the database
	userData, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get user for password update")
		return repository.ErrUserNotFound
	}

	// Verify the current password
	err = utils.VerifyPassword(userData.Password, req.CurrentPassword)
	if err != nil {
		uc.log.WithError(err).Error("Current password verification failed")
		return repository.ErrInvalidCredentials
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		uc.log.WithError(err).Error("Failed to hash new password")
		return err
	}

	// Update password in the database
	userData.Password = hashedPassword
	err = uc.repo.UpdateUser(ctx, userData)
	if err != nil {
		uc.log.WithError(err).Error("Failed to update password in database")
		return err
	}

	return nil
}

// ResendVerificationEmail resends email verification link
func (uc *UserUseCase) ResendVerificationEmail(ctx context.Context, req dto.ResendVerificationRequest) error {
	usr, err := uc.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get user by email")
		return repository.ErrUserNotFound
	}

	if usr.IsActive {
		uc.log.WithField("user_id", usr.ID).Error("User is already active")
		return repository.ErrUserAlreadyActive
	}

	// Generate a new token ID for the verification link
	tokenID, err := uuid.NewV7()
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate UUID")
		return err
	}

	// Generate secure token for the verification link
	token, err := utils.GenerateUniqueToken(ctx, uc.repo.TokenExists)
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate unique token")
		return err
	}

	tokenExpiry := time.Now().Add(24 * time.Hour)
	tokenType := entity.EmailVerification

	// Store token in database for later verification via link
	if err := uc.repo.CreateToken(ctx, usr.ID, tokenID, token, tokenType, tokenExpiry); err != nil {
		uc.log.WithError(err).Error("Failed to create verification token")
		return err
	}

	// Send verification link email in background
	go func() {
		bgCtx := context.Background()
		if err := uc.repo.SendVerificationEmail(bgCtx, usr.Username, usr.Email, token); err != nil {
			uc.log.WithError(err).Error("Failed to send verification email with link")
		}
	}()

	return nil
}

// VerifyEmail Check registration token
func (uc *UserUseCase) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	token, err := uc.repo.GetToken(ctx, req.Token)
	if err != nil {
		uc.log.WithError(err).Error("Failed to get token")
		return err
	}

	if token.UsedAt != nil {
		uc.log.WithError(err).Error("Token already used")
		return repository.ErrTokenAlreadyUsed
	}

	if token.Type != "email_verification" {
		return repository.ErrInvalidToken
	}

	if token.ExpiredAt.Before(time.Now()) {
		return repository.ErrTokenExpired
	}

	if err := uc.repo.ActivateUser(ctx, token.UserID); err != nil {
		uc.log.WithError(err).Error("Failed to activate user")
		return err
	}
	// Mark token as used and send post-verification email in background
	go func() {
		// Create a new background context for goroutine operations
		bgCtx := context.Background()

		// Get user data for email
		usr, err := uc.repo.GetUserByID(bgCtx, token.UserID)
		if err != nil {
			uc.log.WithError(err).Error("Failed to get user by ID in background task")
			return
		}

		// Mark token as used
		if err := uc.repo.UpdateToken(bgCtx, token.ID, time.Now()); err != nil {
			uc.log.WithError(err).Error("Failed to update token")
		}

		// Send confirmation email
		if err := uc.repo.SendPostVerificationEmail(bgCtx, usr.Username, usr.Email); err != nil {
			uc.log.WithError(err).Error("Failed to send post-verification email")
		}
	}()

	return nil
}

// Register creates a new users
func (uc *UserUseCase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if users with the same email or username already exist in parallel
	var wg sync.WaitGroup
	var emailExists, usernameExists bool
	var mu sync.Mutex

	wg.Add(2)

	go func() {
		defer wg.Done()
		_, err := uc.repo.GetUserByEmail(ctx, req.Email)
		if err == nil {
			mu.Lock()
			emailExists = true
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		_, err := uc.repo.GetUserByUsername(ctx, req.Username)
		if err == nil {
			mu.Lock()
			usernameExists = true
			mu.Unlock()
		}
	}()

	wg.Wait()

	if emailExists {
		return nil, repository.ErrUserAlreadyExists
	}
	if usernameExists {
		return nil, repository.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		uc.log.WithError(err).Error("Failed to hash password")
		return nil, err
	}

	id, err := uuid.NewV7()
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate UUID")
		return nil, err
	}

	// Create user
	usr := &entity.User{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entity.Role(req.Role),
	}

	if err := uc.repo.CreateUser(ctx, usr); err != nil {
		uc.log.WithError(err).Error("Failed to create user")
		return nil, err
	}

	// Launch verification link generation and email sending in background
	go func() {
		// Create a new context for background processing
		bgCtx := context.Background()

		tokenID, err := uuid.NewV7()
		if err != nil {
			uc.log.WithError(err).Error("Failed to generate UUID in background")
			return
		}

		// Generate secure token for the verification link
		token, err := utils.GenerateUniqueToken(bgCtx, uc.repo.TokenExists)
		if err != nil {
			uc.log.WithError(err).Error("Failed to generate unique token in background")
			return
		}

		tokenExpiry := time.Now().Add(24 * time.Hour)
		tokenType := entity.EmailVerification

		// Store token in database for later verification
		if err := uc.repo.CreateToken(bgCtx, usr.ID, tokenID, token, tokenType, tokenExpiry); err != nil {
			uc.log.WithError(err).Error("Failed to create verification token in background")
			return
		}

		// Send email with verification link
		if err := uc.repo.SendVerificationEmail(bgCtx, usr.Username, usr.Email, token); err != nil {
			uc.log.WithError(err).Error("Failed to send verification email in background")
		}
	}()

	return &dto.RegisterResponse{
		User:    *usr,
		Message: "User created successfully. Please check your email to verify your account",
	}, nil
}

// Login authenticates a users and returns users data with auth token
func (uc *UserUseCase) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	var usr *entity.User
	var err error

	// Try to get users by email or username
	if req.Email != "" {
		usr, err = uc.repo.GetUserByEmail(ctx, req.Email)
	} else if req.Username != "" {
		usr, err = uc.repo.GetUserByUsername(ctx, req.Username)
	} else {
		return nil, errors.New("email or username required")
	}

	if err != nil {
		uc.log.WithError(err).Error("Failed to find users during login")
		return nil, repository.ErrInvalidCredentials
	}

	// Verify password
	if err = utils.VerifyPassword(usr.Password, req.Password); err != nil {
		uc.log.WithError(err).WithField("users_id", usr.ID).Error("Invalid password during login")
		return nil, repository.ErrInvalidCredentials
	}

	if !usr.IsActive {
		uc.log.WithField("users_id", usr.ID).Error("User is not active")
		return nil, repository.ErrUserNotActive
	}

	tokenString, err := uc.tokenMaker.CreateToken(usr.ID.String(), usr.Email, 24*time.Hour)
	if err != nil {
		uc.log.WithError(err).Error("Failed to generate token")
		return nil, errors.New("failed to generate token")
	}

	expiry := time.Now().Add(24 * time.Hour)

	sessionKey := fmt.Sprintf("session:%s", usr.ID)
	err = uc.repo.DeleteUserSession(ctx, sessionKey)
	if err != nil {
		uc.log.WithError(err).Error("Failed to delete old session")

	}

	// Store user data in Redis for session management
	if err := uc.repo.StoreUserSession(ctx, usr.ID.String(), tokenString, usr); err != nil {
		uc.log.WithError(err).Error("Failed to store user session in Redis")
		return nil, errors.New("failed to store session")
	}

	return &dto.LoginResponse{
		User:   *usr,
		Token:  tokenString,
		Expiry: expiry,
	}, nil
}

// GetUsesByID retrieves users by ID
func (uc *UserUseCase) GetUsesByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	usr, err := uc.repo.GetUserByID(ctx, id)
	if err != nil {
		uc.log.WithError(err).WithField("users_id", id).Error("Failed to get users by ID")
		return nil, repository.ErrUserNotFound
	}

	return usr, nil
}

// GetUserByEmail retrieves users by email
func (uc *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	usr, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		uc.log.WithError(err).WithField("email", email).Error("Failed to get users by email")
		return nil, repository.ErrUserNotFound
	}

	return usr, nil
}

// GetUsersByUsername  retrieves users by username
func (uc *UserUseCase) GetUsersByUsername(ctx context.Context, username string) (*entity.User, error) {
	usr, err := uc.repo.GetUserByUsername(ctx, username)
	if err != nil {
		uc.log.WithError(err).WithField("username", username).Error("Failed to get users by username")
		return nil, repository.ErrUserNotFound
	}

	return usr, nil
}
