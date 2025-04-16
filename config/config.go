package config

import (
	"io"
	"os"
	"strconv"
	"time"
)

// AppConfig holds all application configuration
type AppConfig struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Logger      LoggerConfig
	Mailgun     MailgunConfig
	SendGrid    SendGridConfig
	FrontendURL string
	JWT         JWT
	Minio       MinioConfig
	SwaggerAuth SwaggerAuthConfig
}

type DatabaseConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

type JWT struct {
	Secret string
}

type LoggerConfig struct {
	Level      string
	JSONFormat bool
	Output     io.Writer
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port          string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	VerifyBaseURL string
	BaseURL       string
}

type MailgunConfig struct {
	APIKey  string
	Domain  string
	BaseURL string // Optional, for EU domains
}

type SendGridConfig struct {
	APIKey    string
	FromName  string
	FromEmail string
}

type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	Location        string
}

// RedisConfig Config holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PoolSize int
}

// SwaggerAuthConfig holds configuration for Swagger authentication
type SwaggerAuthConfig struct {
	Username string
	Password string
}

// Load reads configuration from environment variables with defaults
func Load() *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Port:          getEnv("SERVER_PORT", "8080"),
			ReadTimeout:   getDuration("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout:  getDuration("SERVER_WRITE_TIMEOUT", 10*time.Second),
			VerifyBaseURL: getEnv("SERVER_VERIFY_BASE_URL", "http://localhost:8080"),
			BaseURL:       getEnv("SERVER_BASE_URL", "http://localhost:8080"),
		},
		Database: DatabaseConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "postgres"),
			DBName:      getEnv("DB_NAME", "career_link"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			MaxOpenConn: getInt("DB_MAX_OPEN_CONN", 10),
			MaxIdleConn: getInt("DB_MAX_IDLE_CONN", 5),
			MaxLifetime: getDuration("DB_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
			PoolSize: getInt("REDIS_POOL_SIZE", 10),
		},

		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			JSONFormat: getBool("LOG_JSON", false),
		},
		Mailgun: MailgunConfig{
			APIKey:  getEnv("MAILGUN_API_KEY", ""),
			Domain:  getEnv("MAILGUN_DOMAIN", ""),
			BaseURL: getEnv("MAILGUN_BASE_URL", ""),
		},
		SendGrid: SendGridConfig{
			APIKey:    getEnv("SENDGRID_API_KEY", ""),
			FromName:  getEnv("SENDGRID_FROM_NAME", "Career Link"),
			FromEmail: getEnv("SENDGRID_FROM_EMAIL", "prasetyo.agpr@gmail.com"),
		},
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3001"),
		JWT: JWT{
			Secret: getEnv("JWT", "secret"),
		},
		Minio: MinioConfig{
			Endpoint:        getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKeyID:     getEnv("MINIO_ACCESS_KEY_ID", "A04fCBKcIGUEHSmFATXw"),
			SecretAccessKey: getEnv("MINIO_SECRET_ACCESS_KEY", "HxBXkORZiaznivvEICKsfAeuvjv3SVSf2L7VOJfS"),
			UseSSL:          getBool("MINIO_USE_SSL", false),
			BucketName:      getEnv("MINIO_BUCKET_NAME", "career-link"),
			Location:        getEnv("MINIO_LOCATION", "us-east-1"),
		},
		SwaggerAuth: SwaggerAuthConfig{
			Username: getEnv("SWAGGER_AUTH_USERNAME", "admin"),
			Password: getEnv("SWAGGER_AUTH_PASSWORD", "admin123"),
		},
	}
}

// Add this helper function
func getBool(key string, fallback bool) bool {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return fallback
	}
	return value
}

// Helper functions for environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return fallback
	}
	return value
}

func getDuration(key string, fallback time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		return fallback
	}
	return value
}
