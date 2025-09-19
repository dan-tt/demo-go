// Package logger provides structured logging functionality for the demo-go application
// using zap as the underlying logging framework with configurable levels and formats.
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.SugaredLogger to provide structured logging
type Logger struct {
	*zap.SugaredLogger
}

// Config holds logger configuration
type Config struct {
	Level       string `json:"level"`       // debug, info, warn, error
	Environment string `json:"environment"` // development, production
	Format      string `json:"format"`      // json, console
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:       getEnvOrDefault("LOG_LEVEL", "info"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
		Format:      getEnvOrDefault("LOG_FORMAT", "console"),
	}
}

// New creates a new logger instance with the given configuration
func New(config *Config) (*Logger, error) {
	var zapConfig zap.Config

	// Set environment-specific defaults
	if config.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.DisableStacktrace = true
	} else {
		zapConfig = zap.NewDevelopmentConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(level)

	// Set encoding format
	if config.Format == "json" {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Customize encoder config for better readability
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.CallerKey = "caller"
	zapConfig.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Build the logger
	zapLogger, err := zapConfig.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{
		SugaredLogger: zapLogger.Sugar(),
	}, nil
}

// NewDefault creates a logger with default configuration
func NewDefault() (*Logger, error) {
	return New(DefaultConfig())
}

// WithFields adds structured context fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(args...),
	}
}

// WithField adds a single structured context field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		SugaredLogger: l.SugaredLogger.With(key, value),
	}
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *Logger {
	return l.WithField("error", err.Error())
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return l.WithField("request_id", requestID)
}

// WithUserID adds a user ID field to the logger
func (l *Logger) WithUserID(userID string) *Logger {
	return l.WithField("user_id", userID)
}

// ForComponent creates a logger for a specific component
func (l *Logger) ForComponent(component string) *Logger {
	return l.WithField("component", component)
}

// ForRequest creates a logger for a specific HTTP request
func (l *Logger) ForRequest(method, path, requestID string) *Logger {
	return l.WithFields(map[string]interface{}{
		"method":     method,
		"path":       path,
		"request_id": requestID,
		"component":  "http",
	})
}

// ForService creates a logger for a service layer operation
func (l *Logger) ForService(service, operation string) *Logger {
	return l.WithFields(map[string]interface{}{
		"service":   service,
		"operation": operation,
		"layer":     "service",
	})
}

// ForRepository creates a logger for a repository layer operation
func (l *Logger) ForRepository(repository, operation string) *Logger {
	return l.WithFields(map[string]interface{}{
		"repository": repository,
		"operation":  operation,
		"layer":      "repository",
	})
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.SugaredLogger.Sync()
}

// getEnvOrDefault returns the environment variable value or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Global logger instance for convenience
var globalLogger *Logger

// InitGlobal initializes the global logger
func InitGlobal(config *Config) error {
	logger, err := New(config)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetGlobal returns the global logger instance
func GetGlobal() *Logger {
	if globalLogger == nil {
		// Fallback to default logger if not initialized
		logger, err := NewDefault()
		if err != nil {
			// Fallback to a basic logger if NewDefault fails
			panic("Failed to create default logger: " + err.Error())
		}
		globalLogger = logger
	}
	return globalLogger
}

// Debug logs a debug-level message using the global logger instance.
func Debug(args ...interface{}) {
	GetGlobal().Debug(args...)
}

// Info logs an info message using the global logger
func Info(args ...interface{}) {
	GetGlobal().Info(args...)
}

// Warn logs a warning message using the global logger
func Warn(args ...interface{}) {
	GetGlobal().Warn(args...)
}

// Error logs an error message using the global logger
func Error(args ...interface{}) {
	GetGlobal().Error(args...)
}

// Fatal logs a fatal message and exits using the global logger
func Fatal(args ...interface{}) {
	GetGlobal().Fatal(args...)
}

// Debugf logs a formatted debug message using the global logger
func Debugf(template string, args ...interface{}) {
	GetGlobal().Debugf(template, args...)
}

// Infof logs a formatted info message using the global logger
func Infof(template string, args ...interface{}) {
	GetGlobal().Infof(template, args...)
}

// Warnf logs a formatted warning message using the global logger
func Warnf(template string, args ...interface{}) {
	GetGlobal().Warnf(template, args...)
}

// Errorf logs a formatted error message using the global logger
func Errorf(template string, args ...interface{}) {
	GetGlobal().Errorf(template, args...)
}

// Fatalf logs a formatted fatal message and exits using the global logger
func Fatalf(template string, args ...interface{}) {
	GetGlobal().Fatalf(template, args...)
}
