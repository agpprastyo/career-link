package main

import (
	"context"
	"github.com/agpprastyo/career-link/internal/app"
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
	_ = godotenv.Load()

	svr, err := app.InitializeAPI()
	if err != nil {
		panic("Failed to initialize API: " + err.Error())
	}

	// Setup monitoring
	monitoring.SetupMonitoring(svr.App, "career-link-api")

	svr.Logger.Info("Server initialized successfully")
	svr.Logger.Info("OpenAPI doc available at http://localhost:" + svr.Config.Server.Port + "/swagger/index.html")

	// Start the server
	err = svr.Start()
	if err != nil {
		svr.Logger.Fatalf("Failed to start server: %v", err)
	}

	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	svr.Logger.Info("Shutting down...")

	// Give server up to 10 seconds to finish processing requests
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := svr.Shutdown(shutdownCtx); err != nil {
		svr.Logger.Fatalf("Server shutdown failed: %v", err)
	}

	svr.Logger.Info("Server gracefully stopped")
}
