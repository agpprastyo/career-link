package server

import (
	"context"
	"time"

	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/pkg/database"
	customlogger "github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
)

// Server represents the HTTP server
type Server struct {
	app    *fiber.App
	db     *database.PostgresDB
	redis  *redis.Client
	config *config.AppConfig
	log    *customlogger.Logger
}

// New creates a new server instance
func New(cfg *config.AppConfig, db *database.PostgresDB, redisClient *redis.Client, log *customlogger.Logger) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
	})

	// Create the server
	s := &Server{
		app:    app,
		db:     db,
		redis:  redisClient,
		config: cfg,
		log:    log,
	}

	// Register middleware
	s.registerMiddleware()

	// Register routes (now from routes.go)
	s.setupRoutes()

	return s
}

// Start begins listening for requests
func (s *Server) Start() {
	go func() {
		s.log.Infof("Server is running on port %s", s.config.Server.Port)
		if err := s.app.Listen(":" + s.config.Server.Port); err != nil {
			s.log.Fatalf("Error starting server: %v", err)
		}
	}()
}

// Shutdown gracefully stops the server
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Server is shutting down...")
	return s.app.Shutdown()
}

// registerMiddleware adds middleware to the app
func (s *Server) registerMiddleware() {
	// Log all requests
	s.app.Use(fiberlogger.New())

	// Add CORS middleware
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization",
		AllowCredentials: false,
	}))
}
