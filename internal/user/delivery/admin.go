package delivery

import (
	"errors"
	"fmt"
	responseError "github.com/agpprastyo/career-link/internal/common/errors"
	"github.com/agpprastyo/career-link/internal/common/pagination"
	"github.com/agpprastyo/career-link/internal/user/dto"
	"github.com/agpprastyo/career-link/internal/user/entity"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/pkg/utils"
	"github.com/agpprastyo/career-link/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteAdmin godoc
// @Summary      Delete Admin & admin viewers
// @Description  Delete admin & admin viewers by super admin
// @Tags         Super Admin
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "id"
// @Success      200  {object}  entity.DeleteAdminResponse
// @Failure      400  {object}  entity.ErrorBadRequest
// @Failure      404  {object}  entity.ErrorNotFound
// @Failure      500  {object}  entity.ErrorInternalServer
// @Router       /admin/{id} [delete]
// DeleteAdmin handles deleting an admin user
func (h *Handler) DeleteAdmin(c *fiber.Ctx) error {
	ctx := c.Context()
	// Parse user ID from URL
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	// Delete admin user
	err = h.userUseCase.SoftDeleteAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return responseError.RespondWithError(c, fiber.StatusNotFound, "Admin not found")
		}
		h.log.WithError(err).Error("Failed to delete admin")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Failed to delete admin")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Admin successfully deleted",
	})
}

// UpdateAdminStatus handles updating the status of an admin user
func (h *Handler) UpdateAdminStatus(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse user ID from URL
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid user ID format")
	}

	// Toggle status with a single database operation
	newStatus, err := h.userUseCase.ToggleAdminStatus(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return responseError.RespondWithError(c, fiber.StatusNotFound, "Admin not found")
		}
		h.log.WithError(err).Error("Failed to toggle admin status")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Failed to update admin status")
	}

	statusMessage := "activated"
	if !newStatus {
		statusMessage = "deactivated"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   fmt.Sprintf("Admin successfully %s", statusMessage),
		"is_active": newStatus,
	})
}

// CreateAdmin handles creating a new admin user
func (h *Handler) CreateAdmin(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.CreateAdminRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.WithError(err).Error("Create admin - invalid input")
		return responseError.RespondWithError(c, fiber.StatusBadRequest, "Invalid input")
	}

	req.Validator.CheckField(req.Email != "", "email ", "either email or username is required")
	req.Validator.CheckField(validator.Matches(req.Email, validator.RgxEmail), "Email", "Must be a valid email address")

	req.Validator.CheckField(req.Username != "", "username ", "either email or username is required")
	req.Validator.CheckField(len(req.Username) >= 3, "Username", "Username is too short")
	req.Validator.CheckField(len(req.Username) <= 50, "Username", "Username is too long")

	req.Validator.CheckField(req.Password != "", "Password", "Password is required")
	req.Validator.CheckField(len(req.Password) >= 8, "Password", "Password is too short")
	req.Validator.CheckField(len(req.Password) <= 72, "Password", "Password is too long")
	req.Validator.CheckField(validator.NotIn(req.Password, utils.CommonPasswords...), "Password", "Password is too common")

	req.Validator.CheckField(entity.AdministratorRole(req.AdminRole) != entity.AdminRoleSuper, "Role", "Invalid role admin")
	req.Validator.CheckField(entity.AdministratorRole(req.AdminRole) == entity.AdminRoleViewer || entity.AdministratorRole(req.AdminRole) == entity.AdminRoleAdmin, "Role", "Role admin must be either admin or viewer")
	if req.Validator.HasErrors() {
		return responseError.RespondWithError(c, fiber.StatusBadRequest, req.Validator.FirstErrorMessage())
	}

	admin, err := h.userUseCase.CreateAdmin(ctx, req)
	if err != nil {
		h.log.WithError(err).Error("Create admin failed")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusCreated).JSON(admin)
}

// GetUsersByID handles fetching a user by ID
func (h *Handler) GetUsersByID(c *fiber.Ctx) error {
	ctx := c.Context()

	id := c.Params("id")
	userUUID := uuid.MustParse(id)
	usr, err := h.userUseCase.GetUsesByID(ctx, userUUID)
	if err != nil {
		h.log.WithError(err).Error("Get user by ID failed")
		return responseError.RespondWithError(c, fiber.StatusNotFound, "User not found")
	}

	usr.Password = ""

	return c.Status(fiber.StatusOK).JSON(usr)
}

// GetUsers handles fetching all users
func (h *Handler) GetUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	paging := pagination.ExtractFromRequest(c)

	users, err := h.userUseCase.GetUsers(ctx, paging)
	if err != nil {
		h.log.WithError(err).Error("Get users failed")
		return responseError.RespondWithError(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
