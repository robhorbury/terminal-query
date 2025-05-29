package logger

import (
	"log/slog"
	"os"
	"sync"
)

// LogFormat is an enum-like type for log output formats
type LogFormat int

const (
	FormatText LogFormat = iota
	FormatJSON
)

var (
	log  *slog.Logger
	once sync.Once
)

// LoggerConfig holds logger settings.
type LoggerConfig struct {
	Level  slog.Level
	Format LogFormat
}

// Init initializes the logger with config.
// Call this once in main() before using the logger.
func Init(config LoggerConfig) {
	once.Do(func() {
		var handler slog.Handler

		switch config.Format {
		case FormatJSON:
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: config.Level})
		case FormatText:
			handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: config.Level})
		default:
			panic("unsupported log format")
		}

		log = slog.New(handler)
		slog.SetDefault(log)
	})
}

// Get returns the initialized logger.
func Get() *slog.Logger {
	if log == nil {
		panic("logger not initialized. Call logger.Init first.")
	}
	return log
}

