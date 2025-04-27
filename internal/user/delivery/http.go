package delivery

import (
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/internal/common/middleware"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/internal/user/usecase"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/agpprastyo/career-link/pkg/token"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	_ "image/gif"  // Register GIF format
	_ "image/jpeg" // Register JPEG format
	_ "image/png"  // Register PNG format
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	config      *config.AppConfig
	userUseCase *usecase.UserUseCase
	userRepo    *repository.UserRepository
	log         *logger.Logger
	tokenMaker  token.Maker
	redisClient *redis.Client
}

// NewUserHandler creates a new user HTTP handler
func NewUserHandler(userUseCase *usecase.UserUseCase, log *logger.Logger, cfg *config.AppConfig, tokenMaker token.Maker, redisClient *redis.Client, repo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		userRepo:    repo,
		log:         log,
		config:      cfg,
		tokenMaker:  tokenMaker,
		redisClient: redisClient,
	}
}

func (h *UserHandler) RegisterUserRoutes(router fiber.Router) {
	router.Use(middleware.AuthRateLimiter())
	// TODO : implement delete user session in redis if relogin
	router.Post("/login", h.Login) // test with apidog
	// TODO: register user (company and jobSeeker) not completed yet
	router.Post("/register", h.Register)

	router.Post("/resend-verification-email", h.ResendVerificationEmail) // test with apidog

	// TODO : forgot password not completed yet
	//router.Post("/forgot-password", h.ForgotPassword)
}

func (h *UserHandler) RegisterUserVerifyRoute(router fiber.Router) {
	router.Get("/verify", h.VerifyEmailGet) // verify user email to mark as active

}

func (h *UserHandler) RegisterUserWithMiddlewareRoutes(router fiber.Router) {
	router.Use(middleware.RequireAuthMiddleware(h.tokenMaker, h.userRepo, h.log))
	router.Post("/update-password", h.UpdatePassword)
	router.Post("/logout", h.Logout)

	router.Post("/update-avatar", h.UpdateAvatar)
	router.Get("/user", h.GetUser)

}

// RegisterAdminRoutes admin routes
func (h *UserHandler) RegisterAdminRoutes(router fiber.Router) {
	router.Use(middleware.RequireAuthMiddleware(h.tokenMaker, h.userRepo, h.log))
	router.Use(middleware.RequireAdminMiddleware())
	router.Get("/admin/users", h.GetUsers)
	router.Get("/admin/users/:id", h.GetUsersByID)

	//router.Post("/admin/users", h.CreateUser)
	//router.Put("/admin/users/:id", h.UpdateUser)
	//router.Delete("/admin/users/:id", h.DeleteUser)
}

func (h *UserHandler) RegisterSuperAdminRoutes(router fiber.Router) {
	router.Use(middleware.RequireAuthMiddleware(h.tokenMaker, h.userRepo, h.log))
	router.Use(middleware.RequireSuperAdminMiddleware())

	router.Post("/admin", h.CreateAdmin)
	// active and deactivate admin
	router.Patch("/admin/:id/status", h.UpdateAdminStatus)
	router.Delete("/admin/:id", h.DeleteAdmin)

	router.Get("/test-admin", func(c *fiber.Ctx) error {
		return c.SendString("Hello Admin")
	})
}

func (h *UserHandler) RegisterCompanyRoutes(router fiber.Router) {
	router.Use(middleware.RequireAuthMiddleware(h.tokenMaker, h.userRepo, h.log))
	router.Use(middleware.RequireCompanyMiddleware())

	router.Get("/company", h.GetCompany)
	//router.Put("/company", h.UpdateCompany)
	//router.Delete("/company", h.DeleteCompany)
}

// GetCompany get company
func (h *UserHandler) GetCompany(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get the user ID from the context
	userID := c.Locals("user_id").(string)
	userUUID := uuid.MustParse(userID)
	company, err := h.userUseCase.GetCompany(ctx, userUUID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(company)
}
