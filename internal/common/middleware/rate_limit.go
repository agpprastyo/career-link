package middleware

import (
	"time"

	responseError "github.com/agpprastyo/career-link/internal/common/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimiterConfig defines the configuration for rate limiter middleware
type RateLimiterConfig struct {
	// Max number of requests allowed within the expiration duration
	Max int

	// Expiration time for the limiter
	Expiration time.Duration

	// KeyGenerator generates a key for rate limiting (default: IP address)
	KeyGenerator func(*fiber.Ctx) string

	// LimitReached is called when the rate limit is exceeded
	LimitReached fiber.Handler

	// SkipSuccessfulRequests skips successful requests (status < 400)
	SkipSuccessfulRequests bool
}

// DefaultRateLimiterConfig returns a default configuration for rate limiting
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return responseError.RespondWithError(c, fiber.StatusTooManyRequests,
				"Rate limit exceeded. Please try again later.")
		},
		SkipSuccessfulRequests: false,
	}
}

// RateLimiter creates a new rate limiting middleware with the given configuration
func RateLimiter(config RateLimiterConfig) fiber.Handler {
	// Apply default configuration if necessary
	if config.Max <= 0 {
		config.Max = DefaultRateLimiterConfig().Max
	}
	if config.Expiration <= 0 {
		config.Expiration = DefaultRateLimiterConfig().Expiration
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = DefaultRateLimiterConfig().KeyGenerator
	}
	if config.LimitReached == nil {
		config.LimitReached = DefaultRateLimiterConfig().LimitReached
	}

	// Create a fiber limiter configuration
	limiterConfig := limiter.Config{
		Max:                    config.Max,
		Expiration:             config.Expiration,
		KeyGenerator:           config.KeyGenerator,
		LimitReached:           config.LimitReached,
		SkipSuccessfulRequests: config.SkipSuccessfulRequests,
	}

	return limiter.New(limiterConfig)
}

// AuthRateLimiter returns a rate limiter for authentication endpoints
func AuthRateLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return responseError.RespondWithError(c, fiber.StatusTooManyRequests,
				"Too many authentication attempts. Please try again later.")
		},
	})
}

// RegistrationRateLimiter returns a rate limiter specifically for registration
func RegistrationRateLimiter() fiber.Handler {
	return RateLimiter(RateLimiterConfig{
		Max:        3,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return responseError.RespondWithError(c, fiber.StatusTooManyRequests,
				"Too many registration attempts. Please try again later.")
		},
	})
}
