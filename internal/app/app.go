package app

import (
	"context"
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/agpprastyo/career-link/internal/common/health"
	"github.com/agpprastyo/career-link/internal/common/middleware"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/internal/user/delivery"
	"github.com/agpprastyo/career-link/internal/user/repository"
	"github.com/agpprastyo/career-link/internal/user/usecase"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/mail"
	"github.com/agpprastyo/career-link/pkg/minio"
	"github.com/agpprastyo/career-link/pkg/redis"

	"github.com/agpprastyo/career-link/pkg/token"
	"github.com/sirupsen/logrus"
)

// InitializeAPI manually initializes all dependencies
func InitializeAPI() (*Server, error) {
	// Initialize config
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg)

	// Initialize Redis client
	redisClient, err := redis.NewClient(*cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Redis: %w", err)
	}

	// Replace the existing Redis Ping check with this:
	if _, err := redisClient.Ping(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	// Initialize database
	db, err := database.NewPostgresDB(cfg.Database, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize token maker
	tokenMaker, err := token.NewJWTMaker(cfg.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create token maker: %w", err)
	}

	// Initialize mail client
	mailClient := mail.NewSendGridClient(cfg, log)

	// Initialize MinIO client
	minioClient, err := minio.NewClient(cfg, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO: %w", err)
	}

	// Initialize repositories
	userRepo := initUserRepository(db, log, mailClient, minioClient, cfg.Server.VerifyBaseURL, redisClient)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, log, tokenMaker)

	// Initialize handlers
	userHandler := delivery.NewUserHandler(userUseCase, log, cfg, tokenMaker, redisClient, userRepo)
	healthHandler := health.NewHandler(db, redisClient)

	// Create the server with all dependencies
	server := NewServer(cfg, log, userHandler, healthHandler)

	return server, nil
}

func initUserRepository(db *database.PostgresDB, log *logger.Logger, mail *mail.Client, minio *minio.Client, verifyBaseURL string, redis *redis.Client) *repository.UserRepository {
	if redis == nil {
		log.WithFields(logrus.Fields{
			"component": "UserRepository",
		}).Fatal("Redis client is nil")
	}
	return repository.NewUserRepository(db, log, mail, minio, verifyBaseURL, redis)
}

// Server represents the fully configured API server
type Server struct {
	App    *fiber.App
	Config *config.AppConfig
	Logger *logger.Logger
}

// NewServer creates and configures a new server instance with all routes
func NewServer(cfg *config.AppConfig, log *logger.Logger, userHandler *delivery.UserHandler, healthHandler *health.Handler) *Server {
	// Initialize router with global middleware
	app := fiber.New()
	app.Use(fiberLogger.New())
	app.Use(middleware.RecoveryMiddleware(log))

	// Register documentation routes
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

	// Register health check route
	app.Get("/health", healthHandler.Check)

	// Register user verification route at root level
	userHandler.RegisterUserVerifyRoute(app.Group("/"))

	// Register API routes
	api := app.Group("/api/v1")

	userHandler.RegisterUserRoutes(api.Group("/users"))
	userHandler.RegisterUserWithMiddlewareRoutes(api.Group("/profile"))
	userHandler.RegisterAdminRoutes(api.Group("/admin"))
	userHandler.RegisterSuperAdminRoutes(api.Group("/super-admin"))
	userHandler.RegisterCompanyRoutes(api.Group("/company"))

	// Add after all your route registrations
	api.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})

	return &Server{
		App:    app,
		Config: cfg,
		Logger: log,
	}
}

// Start begins listening for requests
func (s *Server) Start() error {
	go func() {
		s.Logger.Infof("Server is running on port %s", s.Config.Server.Port)
		if err := s.App.Listen(":" + s.Config.Server.Port); err != nil {
			s.Logger.Fatalf("Error starting server: %v", err)
		}
	}()

	s.Logger.Info("Server started successfully")
	return nil
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Server is shutting down...")
	return s.App.Shutdown()
}
