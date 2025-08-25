package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend-service-internpro/internal/pkg/jwt"

	"github.com/danielgtaylor/huma/v2"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides JWT authentication for Gin
func AuthMiddleware(jwtSecrets jwt.Secrets) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must start with 'Bearer '",
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token is required",
			})
			return
		}

		claims, err := jwt.ParseAccess(tokenStr, jwtSecrets.Access)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		// Store user ID in context for use in handlers
		c.Set("userID", claims.UserID)
		c.Next()
	}
}

// ValidateToken validates JWT token for Huma handlers
func ValidateToken(authHeader string, jwtSecrets jwt.Secrets) (*jwt.Claims, error) {
	if authHeader == "" {
		return nil, huma.Error401Unauthorized("Authorization header is required")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, huma.Error401Unauthorized("Authorization header must start with 'Bearer '")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenStr == "" {
		return nil, huma.Error401Unauthorized("Token is required")
	}

	claims, err := jwt.ParseAccess(tokenStr, jwtSecrets.Access)
	if err != nil {
		return nil, huma.Error401Unauthorized("Invalid or expired token")
	}

	return claims, nil
}

// AuthContext holds authenticated user information
type AuthContext struct {
	UserID string
}

// GetAuthContext extracts authentication context from Gin context
func GetAuthContext(c *gin.Context) (*AuthContext, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return nil, false
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return nil, false
	}

	return &AuthContext{
		UserID: userIDStr,
	}, true
}

// RequireAuth creates a middleware that requires authentication for Huma operations
func RequireAuth(jwtSecrets jwt.Secrets) func(ctx context.Context, op *huma.Operation) *huma.ErrorDetail {
	return func(ctx context.Context, op *huma.Operation) *huma.ErrorDetail {
		// This would be used with Huma's middleware system if available
		// For now, we'll use manual validation in handlers
		return nil
	}
}

// WithAuth is a helper function to wrap handlers with authentication
type AuthenticatedHandler[I, O any] func(ctx context.Context, input *I, userID string) (*O, error)

func WithAuth[I any, O any](
	jwtSecrets jwt.Secrets,
	handler AuthenticatedHandler[I, O],
) func(ctx context.Context, input *I) (*O, error) {
	return func(ctx context.Context, input *I) (*O, error) {
		// Extract Authorization header from input using reflection or type assertion
		// This is a simplified version - in practice, you'd need to handle the header extraction
		// For now, we'll keep using the manual validation approach in handlers
		return nil, huma.Error500InternalServerError("WithAuth wrapper not fully implemented")
	}
}
