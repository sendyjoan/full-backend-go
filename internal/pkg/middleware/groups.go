package middleware

import (
	"backend-service-internpro/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthenticatedRoutes creates a router group with authentication middleware
func AuthenticatedRoutes(router *gin.Engine, jwtSecrets jwt.Secrets) *gin.RouterGroup {
	authGroup := router.Group("/")
	authGroup.Use(AuthMiddleware(jwtSecrets))
	return authGroup
}

// AdminRoutes creates a router group with admin authentication
func AdminRoutes(router *gin.Engine, jwtSecrets jwt.Secrets) *gin.RouterGroup {
	adminGroup := router.Group("/admin")
	adminGroup.Use(AuthMiddleware(jwtSecrets))
	// Add admin role checking here if needed
	return adminGroup
}

// APIKeyMiddleware validates API key for certain endpoints
func APIKeyMiddleware(validAPIKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			c.JSON(401, gin.H{"error": "API key required"})
			c.Abort()
			return
		}

		if apiKey != validAPIKey {
			c.JSON(401, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
