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

		// Get environment mode
		ginMode := gin.Mode()

		// In development mode, allow all origins
		// In production, be more restrictive
		if ginMode == gin.DebugMode || ginMode == gin.TestMode {
			// Development: Allow all origins
			if origin != "" {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		} else {
			// Production: Allow specific origins
			allowedOrigins := map[string]bool{
				"https://schooltechindonesia.com":             true,
				"https://www.schooltechindonesia.com":         true,
				"https://api.schooltechindonesia.com":         true,
				"https://staging-api.schooltechindonesia.com": true,
				"https://testing-api.schooltechindonesia.com": true,
				"http://localhost:3000":                       true,
				"http://localhost:3001":                       true,
				"http://localhost:8080":                       true,
				"http://localhost:8000":                       true,
				"https://localhost:3000":                      true,
				"https://localhost:3001":                      true,
				"https://localhost:8080":                      true,
				"https://localhost:8000":                      true,
				"http://127.0.0.1:3000":                       true,
				"http://127.0.0.1:3001":                       true,
				"http://127.0.0.1:8080":                       true,
				"http://127.0.0.1:8000":                       true,
				"https://127.0.0.1:3000":                      true,
				"https://127.0.0.1:3001":                      true,
				"https://127.0.0.1:8080":                      true,
				"https://127.0.0.1:8000":                      true,
				"https://localhost:5173":                      true,
			}

			if origin == "" {
				c.Header("Access-Control-Allow-Origin", "*")
			} else if allowedOrigins[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				// For production, allow any localhost origin for development
				if ginMode == gin.ReleaseMode {
					c.Header("Access-Control-Allow-Origin", "*")
				}
			}
		}

		// Set comprehensive CORS headers
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-HTTP-Method-Override, X-Forwarded-For, X-Real-IP")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Access-Control-Expose-Headers", "Authorization, Content-Length, X-CSRF-Token")

		// Handle preflight requests
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
