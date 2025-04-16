package main

import (
	"context"
	"github.com/agpprastyo/career-link/internal/wire"
	"github.com/agpprastyo/career-link/pkg/monitoring"
	"github.com/joho/godotenv"
	"os"
	"os/signal"

	"syscall"
	"time"
)

// @title Career Link API
// @version 1.0
// @description API for Career Link application
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@careerlink.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	err := godotenv.Load()
	// Load configuration
	//cfg := config.Load()

	server, err := wire.InitializeAPI()
	if err != nil {
		panic(err)
	}

	monitoring.SetupMonitoring(server.App, "career-link-api")

	server.Logger.Info("Initializing server...")

	server.Logger.Info("Starting server on port " + server.Config.Server.Port)
	server.Logger.Info("OpenAPI doc available at http://localhost:" + server.Config.Server.Port + "/swagger/index.html")
	if err := server.App.Listen(":" + server.Config.Server.Port); err != nil {
		server.Logger.WithError(err).Fatal("Server failed to start")
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	server.Logger.Println("Shutting down...")

	// Give server up to 10 seconds to finish processing requests
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		server.Logger.Fatalf("Server shutdown failed: %v", err)
	}

	server.Logger.Println("Server gracefully stopped")
}
