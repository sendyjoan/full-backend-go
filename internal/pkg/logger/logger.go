package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

// Logger wraps slog with additional functionality
type Logger struct {
	*slog.Logger
}

// LogLevel represents logging levels
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
)

// New creates a new logger instance
func New(level LogLevel) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Format time as readable string
			if a.Key == slog.TimeKey {
				return slog.String(slog.TimeKey, a.Value.Time().Format(time.RFC3339))
			}
			return a
		},
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithContext adds context to logger
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.Logger.With(),
	}
}

// WithFields adds structured fields to logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
	}
}

// Auth logger with predefined fields
func (l *Logger) Auth() *Logger {
	return l.WithFields(map[string]interface{}{
		"service": "auth",
	})
}

// HTTP logger with predefined fields
func (l *Logger) HTTP() *Logger {
	return l.WithFields(map[string]interface{}{
		"layer": "http",
	})
}

// Repository logger with predefined fields
func (l *Logger) Repository() *Logger {
	return l.WithFields(map[string]interface{}{
		"layer": "repository",
	})
}

// Service logger with predefined fields
func (l *Logger) Service() *Logger {
	return l.WithFields(map[string]interface{}{
		"layer": "service",
	})
}

// Convenient methods for different log levels with context
func (l *Logger) InfoCtx(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.InfoContext(ctx, msg, args...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.ErrorContext(ctx, msg, args...)
}

func (l *Logger) WarnCtx(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.WarnContext(ctx, msg, args...)
}

func (l *Logger) DebugCtx(ctx context.Context, msg string, args ...interface{}) {
	l.Logger.DebugContext(ctx, msg, args...)
}

// Error with error object
func (l *Logger) ErrorWithErr(msg string, err error, args ...interface{}) {
	allArgs := append([]interface{}{"error", err.Error()}, args...)
	l.Logger.Error(msg, allArgs...)
}

// Request logging helpers
func (l *Logger) LogRequest(method, path, userAgent, ip string) {
	l.Info("incoming request",
		"method", method,
		"path", path,
		"user_agent", userAgent,
		"ip", ip,
	)
}

func (l *Logger) LogResponse(method, path string, statusCode int, duration time.Duration) {
	l.Info("request completed",
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration_ms", duration.Milliseconds(),
	)
}

// Login attempts logging
func (l *Logger) LogLoginAttempt(email string, success bool, ip string) {
	status := "failed"
	if success {
		status = "success"
	}

	l.Info("login attempt",
		"email", email,
		"status", status,
		"ip", ip,
	)
}

// Security event logging
func (l *Logger) LogSecurityEvent(event, email, ip, details string) {
	l.Warn("security event",
		"event", event,
		"email", email,
		"ip", ip,
		"details", details,
	)
}

// Database operation logging
func (l *Logger) LogDBOperation(operation, table string, duration time.Duration, err error) {
	if err != nil {
		l.Error("database operation failed",
			"operation", operation,
			"table", table,
			"duration_ms", duration.Milliseconds(),
			"error", err.Error(),
		)
	} else {
		l.Debug("database operation completed",
			"operation", operation,
			"table", table,
			"duration_ms", duration.Milliseconds(),
		)
	}
}

// Global logger instance
var globalLogger *Logger

// InitGlobalLogger initializes the global logger
func InitGlobalLogger(level LogLevel) {
	globalLogger = New(level)
}

// Global logger functions
func Info(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Info(msg, args...)
	} else {
		fmt.Printf("INFO: %s\n", msg)
	}
}

func Error(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Error(msg, args...)
	} else {
		fmt.Printf("ERROR: %s\n", msg)
	}
}

func Warn(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Warn(msg, args...)
	} else {
		fmt.Printf("WARN: %s\n", msg)
	}
}

func Debug(msg string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.Debug(msg, args...)
	} else {
		fmt.Printf("DEBUG: %s\n", msg)
	}
}

// Get global logger
func Global() *Logger {
	if globalLogger == nil {
		globalLogger = New(LevelInfo)
	}
	return globalLogger
}
