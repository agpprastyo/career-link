package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	responseError "github.com/agpprastyo/career-link/internal/common/errors"
	"github.com/agpprastyo/career-link/internal/user/dto"
	"github.com/agpprastyo/career-link/internal/user/entity"

	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/mail/templates"
	"github.com/agpprastyo/career-link/pkg/monitoring"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/agpprastyo/career-link/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

// GetUser handles fetching user data
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Locals("user_id").(string)
	userUUID := uuid.MustParse(userID)

	userData, err := h.userUseCase.GetUsesByID(ctx, userUUID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return responseError.RespondWithError(c, fiber.StatusNotFound, "User not found")
		default:
			h.log.WithError(err).Error("Get user failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	// Remove password from response
	userData.Password = ""
	filename := fmt.Sprintf("avatars/%s", *userData.Avatar)
	avatarURL := ""

	// GetAvatarURL
	if userData.Avatar != nil {
		avatarURL, err = h.userRepo.GetAvatarURL(ctx, filename)
		if err != nil {
			h.log.WithError(err).Error("Failed to get avatar URL")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	userData.Avatar = &avatarURL

	return c.Status(fiber.StatusOK).JSON(userData)
}

// UpdateAvatar handles updating user avatar through file upload
func (h *UserHandler) UpdateAvatar(c *fiber.Ctx) error {
	// Get user ID from context
	userID := c.Locals("user_id").(string)

	// Get file from form
	file, err := c.FormFile("avatar")
	if err != nil {
		h.log.WithError(err).Error("Failed to get avatar file")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Avatar file is required")
	}

	// Validate file size (max 1MB)
	if file.Size > 1*1024*1024 {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "File too large (max 1MB)")
	}

	// Validate file type
	fileExt := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[strings.ToLower(fileExt)] {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid file type. Only jpg, jpeg, png, and gif allowed")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		h.log.WithError(err).Error("Failed to open uploaded file")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			h.log.WithError(err).Error("Failed to close file")
		}
	}(src)

	// Generate unique filename
	uniqueID, err := uuid.NewV7()
	if err != nil {
		h.log.WithError(err).Error("Failed to generate unique ID")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}
	filename := fmt.Sprintf("%s%s", uniqueID, fileExt)

	// Get file content type
	fileBytes, err := io.ReadAll(src)
	if err != nil {
		h.log.WithError(err).Error("Failed to read file content")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	// Reset file pointer
	_, err = src.Seek(0, 0)
	if err != nil {
		h.log.WithError(err).Error("Failed to reset file pointer")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	contentType := http.DetectContentType(fileBytes)
	if !strings.HasPrefix(contentType, "image/") {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "File must be an image")
	}

	// Validate image aspect ratio (must be 1:1)
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		h.log.WithError(err).Error("Failed to decode image")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid image format")
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width != height {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Image must be square (1:1 ratio)")
	}

	// Update user avatar in database
	updateReq := dto.UpdateAvatarRequest{
		UserID:      userID,
		FileName:    filename,
		FileContent: fileBytes,
		FileSize:    int64(len(fileBytes)),
		ContentType: contentType,
	}

	ctx := c.Context()
	avatarURL, err := h.userUseCase.UpdateAvatar(ctx, updateReq)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return responseError.RespondWithError(c, fiber.StatusNotFound, "User not found")
		default:
			h.log.WithError(err).Error("Update avatar failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Avatar updated successfully",
		"avatar_url": avatarURL,
	})
}

// Logout godoc
// @Summary User logout
// @Description Logout a user by deleting their session. Clients should also remove the JWT token from local storage.
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.LogoutResponse
// @Failure 401 {object} dto.ErrorUnauthorized
// @Failure 500 {object} dto.ErrorInternalServer
// @Router /logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	sessionID := fmt.Sprintf("session:%s", userID)

	// Delete session from Redis
	err := h.redisClient.Del(c.Context(), sessionID)
	if err != nil {
		h.log.Error("Failed to delete session from Redis")
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Logged out",
		})

	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logged out"})
}

// UpdatePassword handles password updates for authenticated users
func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	// Get user data from context using the User struct stored in Redis
	userID := c.Locals("user_id").(string)
	sessionID := fmt.Sprintf("session:%s", userID)

	// Get user data from Redis
	userJSON, err := h.redisClient.Get(c.Context(), sessionID)
	if err != nil {
		h.log.WithError(err).Error("Failed to get user data from session")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Session error")
	}

	// Deserialize JSON to User struct
	var userData entity.User
	if err := json.Unmarshal([]byte(userJSON), &userData); err != nil {
		h.log.WithError(err).Error("Failed to decode user data")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Session data error")
	}

	// Parse request body
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		Validator       validator.Validator
	}

	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to decode update password request")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate request
	req.Validator.CheckField(req.CurrentPassword != "", "CurrentPassword", "Current password is required")
	req.Validator.CheckField(req.NewPassword != "", "NewPassword", "New password is required")
	req.Validator.CheckField(len(req.NewPassword) >= 8, "NewPassword", "Password is too short")
	req.Validator.CheckField(len(req.NewPassword) <= 72, "NewPassword", "Password is too long")
	req.Validator.CheckField(validator.NotIn(req.NewPassword, utils.CommonPasswords...), "NewPassword", "Password is too common")
	req.Validator.CheckField(req.CurrentPassword != req.NewPassword, "NewPassword", "New password must be different from current password")

	if req.Validator.HasErrors() {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
	}

	// Create update password request for use case
	updateReq := dto.UpdatePasswordRequest{
		UserID:          userData.ID.String(),
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	// Call use case to update password
	ctx := c.Context()
	err = h.userUseCase.UpdatePassword(ctx, updateReq)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidCredentials):
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Current password is incorrect")
		case errors.Is(err, repository.ErrUserNotFound):
			return responseError.RespondWithError(c, fiber.StatusNotFound, "User not found")
		default:
			h.log.WithError(err).Error("Update password failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

// VerifyEmailGet handles email verification from GET requests with query parameters
func (h *UserHandler) VerifyEmailGet(c *fiber.Ctx) error {
	email := c.Query("email")
	tkn := c.Query("token")

	// Create the verification request
	req := dto.VerifyEmailRequest{
		Email: email,
		Token: tkn,
	}

	// Validate parameters
	req.Validator.CheckField(req.Email != "", "Email", "Email is required")
	req.Validator.CheckField(req.Token != "", "Token", "Token is required")

	if req.Validator.HasErrors() {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
	}

	ctx := c.Context()
	err := h.userUseCase.VerifyEmail(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidToken):
			return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid token")
		case errors.Is(err, repository.ErrTokenExpired):
			return responseError.RespondWithError(c, fiber.StatusBadRequest, "Token has expired, please request a new one")
		default:
			h.log.WithError(err).Error("Verify email failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	frontedURL := h.config.FrontendURL + "/login"
	fmt.Println("frontedURL : ", frontedURL)
	successHTML, err := templates.GetVerificationSuccessHTML(frontedURL)
	if err != nil {
		h.log.WithError(err).Error("Failed to load verification success template")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Error rendering success page")
	}

	return c.Status(fiber.StatusOK).Type("html").SendString(successHTML)
}

func (h *UserHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req dto.ResendVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to decode resend verification request")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// Validate email
	req.Validator.CheckField(req.Email != "", "Email", "Email is required")
	req.Validator.CheckField(validator.Matches(req.Email, validator.RgxEmail), "Email", "Must be a valid email address")

	if req.Validator.HasErrors() {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
	}

	ctx := c.Context()
	err := h.userUseCase.ResendVerificationEmail(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserNotFound):
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "If your email exists, a new verification link has been sent"})
		case errors.Is(err, repository.ErrUserAlreadyActive):
			return responseError.RespondWithError(c, fiber.StatusBadRequest, "User is already active")

		default:
			h.log.WithError(err).Error("Resend verification email failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "A new verification link has been sent to your email"})
}

// Register godoc
// @Summary User registration
// @Description Register a new user with email/username and password for job seeker and company
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.RegisterRequest true "User registration data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorBadRequest
// @Failure 401 {object} dto.ErrorUnauthorized
// @Failure 500 {object} dto.ErrorInternalServer
// @Router /register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to decode register request")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	// validate request
	req.Validator.CheckField(req.Role != string(entity.CompanyRole) && req.Role != string(entity.JobSeekerRole), "Role", "Role is required")
	req.Validator.CheckField(req.Email != "", "email ", "either email or username is required")
	req.Validator.CheckField(validator.Matches(req.Email, validator.RgxEmail), "Email", "Must be a valid email address")

	req.Validator.CheckField(req.Username != "", "username ", "either email or username is required")
	req.Validator.CheckField(len(req.Username) >= 3, "Username", "Username is too short")
	req.Validator.CheckField(len(req.Username) <= 50, "Username", "Username is too long")

	req.Validator.CheckField(req.Password != "", "Password", "Password is required")
	req.Validator.CheckField(len(req.Password) >= 8, "Password", "Password is too short")
	req.Validator.CheckField(len(req.Password) <= 72, "Password", "Password is too long")
	req.Validator.CheckField(validator.NotIn(req.Password, utils.CommonPasswords...), "Password", "Password is too common")

	if req.Validator.HasErrors() {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
	}

	if req.Role == string(entity.CompanyRole) {
		req.Validator.CheckField(req.CompanyProfile.Name != "", "CompanyName", "Company name is required")
		req.Validator.CheckField(len(req.CompanyProfile.Name) >= 3, "CompanyName", "Company name is too short")
		req.Validator.CheckField(len(req.CompanyProfile.Name) <= 50, "CompanyName", "Company name is too long")

		if req.Validator.HasErrors() {
			return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
		}
	} else if req.Role == string(entity.JobSeekerRole) {
		req.Validator.CheckField(req.JobSeekerProfile.FirstName != "", "FirstName", "First name is required")
		req.Validator.CheckField(len(req.JobSeekerProfile.FirstName) >= 3, "FirstName", "First name is too short")
		req.Validator.CheckField(len(req.JobSeekerProfile.FirstName) <= 50, "FirstName", "First name is too long")

		req.Validator.CheckField(req.JobSeekerProfile.LastName != "", "LastName", "Last name is required")
		req.Validator.CheckField(len(req.JobSeekerProfile.LastName) >= 3, "LastName", "Last name is too short")
		req.Validator.CheckField(len(req.JobSeekerProfile.LastName) <= 50, "LastName", "Last name is too long")

		if req.Validator.HasErrors() {
			return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
		}
	}

	ctx := c.Context()
	resp, err := h.userUseCase.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrUserAlreadyExists):
			return responseError.RespondWithError(c, fiber.StatusConflict, "User already exists")
		default:
			h.log.WithError(err).Error("Register failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	// Remove password from response
	resp.User.Password = ""

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login godoc
// @Summary User login
// @Description Authenticate a user with email/username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body user.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorBadRequest
// @Failure 401 {object} dto.ErrorUnauthorized
// @Failure 500 {object} dto.ErrorInternalServer
// @Router /login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	startTime := time.Now()
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Failed to decode login request")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid request payload")
	}

	ctx := c.Context()
	resp, err := h.userUseCase.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInvalidCredentials):
			monitoring.LoginAttempts.WithLabelValues("failure_invalid_credentials").Inc()
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "Invalid credentials")
		case errors.Is(err, repository.ErrUserNotActive):
			monitoring.LoginAttempts.WithLabelValues("failure_inactive_user").Inc()
			return responseError.RespondWithError(c, fiber.StatusUnauthorized, "User is not active, please verify your email")
		default:
			monitoring.LoginAttempts.WithLabelValues("failure_server_error").Inc()
			h.log.WithError(err).Error("Login failed")
			return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	// Remove password from response
	resp.User.Password = ""

	monitoring.LoginAttempts.WithLabelValues("success").Inc()
	monitoring.DatabaseOperationDuration.WithLabelValues("select", "users").Observe(time.Since(startTime).Seconds())

	return c.Status(fiber.StatusOK).JSON(resp)
}
