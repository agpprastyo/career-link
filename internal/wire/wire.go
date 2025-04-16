//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/agpprastyo/career-link/pkg/server"
	"github.com/google/wire"
)

// InitializeAPI sets up the complete API server with all dependencies
func InitializeAPI() (*server.Server, error) {
	wire.Build(
		server.NewServer,
		provideRouter,
		provideAppConfig,
		provideLogger,
		provideTokenMaker,
		provideRedisClient,
		provideDBConnection,
		provideMailClient,
		provideMinioClient,
		provideVerifyBaseURL,
		provideUserRepository,
		provideUserUseCase,
		provideUserHandler,
		provideHealthHandler,
	)
	return &server.Server{}, nil
}
