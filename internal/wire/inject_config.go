package wire

import (
	"github.com/agpprastyo/career-link/config"
	"github.com/agpprastyo/career-link/internal/common/health"
	"github.com/agpprastyo/career-link/pkg/database"
	"github.com/agpprastyo/career-link/pkg/logger"
	"github.com/agpprastyo/career-link/pkg/mail"
	"github.com/agpprastyo/career-link/pkg/minio"
	"github.com/agpprastyo/career-link/pkg/redis"
	"github.com/agpprastyo/career-link/pkg/token"
	"github.com/gofiber/fiber/v2"
	"time"
)

func provideAppConfig() *config.AppConfig {
	return config.Load()
}

func provideLogger(cfg *config.AppConfig) *logger.Logger {
	return logger.New(cfg)
}

func provideDBConnection(cfg *config.AppConfig) *database.PostgresDB {
	db, err := database.NewPostgresDB(cfg.Database, provideLogger(cfg))
	if err != nil {
		panic(err)
	}
	return db
}

func provideTokenMaker(cfg *config.AppConfig) token.Maker {
	maker, err := token.NewJWTMaker(cfg.JWT.Secret)
	if err != nil {
		panic(err)
	}
	return maker
}

func provideRedisClient(cfg *config.AppConfig) (*redis.Client, error) {
	return redis.NewClient(*cfg)
}

func provideMailClient(cfg *config.AppConfig, log *logger.Logger) *mail.Client {
	return mail.NewSendGridClient(cfg, log)
}

func provideMinioClient(cfg *config.AppConfig, log *logger.Logger) *minio.Client {
	client, err := minio.NewClient(cfg, log)
	if err != nil {
		panic(err)
	}
	return client
}

func provideVerifyBaseURL(cfg *config.AppConfig) string {
	return cfg.Server.BaseURL + "/api/verify"
}

func provideHealthHandler(db *database.PostgresDB, redisClient *redis.Client) *health.Handler {
	return health.NewHandler(db, redisClient)
}

func provideRouter(cfg *config.AppConfig, log *logger.Logger) *fiber.App {
	return fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  120 * time.Second,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Handle specific errors here
			log.Error("Error: ", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "11" + err.Error(),
			})
		},
	})
}
