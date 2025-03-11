package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger extends logrus.Logger
type Logger struct {
	*logrus.Logger
}

// Config holds configuration for the logger
type Config struct {
	Level      string
	JSONFormat bool
	Output     io.Writer
}

// New creates a new configured logger
func New(cfg Config) *Logger {
	logger := &Logger{Logger: logrus.New()}

	// Set output
	if cfg.Output != nil {
		logger.SetOutput(cfg.Output)
	} else {
		logger.SetOutput(os.Stdout)
	}

	// Set formatter
	if cfg.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02T15:04:05.000Z07:00",
			DisableColors:    false,
			DisableTimestamp: false,
		})
	}

	// Set level
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return logger
}

// Fields type to define log fields
type Fields logrus.Fields
