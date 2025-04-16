package monitoring

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// SetupMonitoring initializes Prometheus metrics collection for Fiber
func SetupMonitoring(app *fiber.App, serviceName string) {
	// Create a new Prometheus middleware instance
	fiberPrometheus := fiberprometheus.New(serviceName)
	fiberPrometheus.RegisterAt(app, "/metrics")

	// Use Prometheus middleware to collect metrics
	app.Use(fiberPrometheus.Middleware)
}

var (
	// LoginAttempts tracks successful/failed logins
	LoginAttempts = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "career_link_login_attempts_total",
			Help: "Total number of login attempts",
		},
		[]string{"status"}, // "success" or "failure"
	)

	// DatabaseOperationDuration measures database operation duration
	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "career_link_db_operation_duration_seconds",
			Help:    "Duration of database operations in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation", "table"},
	)

	// ActiveUsers tracks currently active users
	ActiveUsers = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "career_link_active_users",
			Help: "Number of currently active users",
		},
	)
)
