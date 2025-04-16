package server

import (
	"context"
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/agpprastyo/career-link/config"
	_ "github.com/agpprastyo/career-link/docs"
	"github.com/agpprastyo/career-link/internal/common/health"
	"github.com/agpprastyo/career-link/internal/common/middleware"
	"github.com/agpprastyo/career-link/internal/user/delivery"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// Server represents the fully configured API server
type Server struct {
	App    *fiber.App
	Config *config.AppConfig
	Logger *logger.Logger
}

func NewServer(app *fiber.App, cfg *config.AppConfig, log *logger.Logger, userHandler *delivery.Handler, healthHandler *health.Handler) *Server {

	app.Use(fiberLogger.New())
	app.Use(middleware.RecoveryMiddleware(log))

	app.Get("/swagger/*", middleware.NewSwaggerAuth(cfg), swagger.HandlerDefault)
	app.Get("/reference", middleware.NewSwaggerAuth(cfg), func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("%v", err))
		}

		return c.Type("html").SendString(htmlContent)
	})

	// Register routes
	root := app.Group("/")
	root.Get("/health", healthHandler.Check)

	api := app.Group("/api/v1")
	userHandler.RegisterUserRoutes(api)
	userHandler.RegisterUserVerifyRoute(api)
	userHandler.RegisterUserWithMiddlewareRoutes(api)
	userHandler.RegisterAdminRoutes(api)
	userHandler.RegisterSuperAdminRoutes(api)
	userHandler.RegisterCompanyRoutes(api)

	return &Server{
		App:    app,
		Config: cfg,
		Logger: log,
	}
}

// Start begins listening for requests
func (s *Server) Start() {
	go func() {
		s.Logger.Infof("Server is running on port %s", s.Config.Server.Port)
		if err := s.App.Listen(":" + s.Config.Server.Port); err != nil {
			s.Logger.Fatalf("Error starting server: %v", err)
		}
	}()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Server is shutting down...")
	return s.App.Shutdown()
}
