package middleware

import (
	"time"

	"backend-service-internpro/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware provides CORS support
func CORSMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Allow specific origins or localhost for development
		allowedOrigins := map[string]bool{
			"http://localhost:8080":    true,
			"https://localhost:8080":   true,
			"http://127.0.0.1:8080":    true,
			"https://127.0.0.1:8080":   true,
			"https://unpkg.com":        true,
			"https://cdn.jsdelivr.net": true,
		}

		if allowedOrigins[origin] || origin == "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// LoggingMiddleware provides request/response logging
func LoggingMiddleware() gin.HandlerFunc {
	appLogger := logger.Global()

	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()

		// Log incoming request
		appLogger.HTTP().LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.GetHeader("User-Agent"),
			c.ClientIP(),
		)

		// Process request
		c.Next()

		// Log response
		appLogger.HTTP().LogResponse(
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
		)
	})
}

// RecoveryMiddleware provides panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	appLogger := logger.Global()

	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		appLogger.Error("panic recovered",
			"error", recovered,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"ip", c.ClientIP(),
		)
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}
