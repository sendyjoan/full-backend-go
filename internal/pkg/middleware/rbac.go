package middleware

import (
	"net/http"
	"strings"

	"backend-service-internpro/internal/rbac/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RBACMiddleware creates middleware for role-based access control
type RBACMiddleware struct {
	rbacService service.Service
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(rbacService service.Service) *RBACMiddleware {
	return &RBACMiddleware{
		rbacService: rbacService,
	}
}

// RequirePermission creates middleware that requires specific permission
func (m *RBACMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Check permission
		hasPermission, err := m.rbacService.CheckUserPermission(c.Request.Context(), uid, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check permission",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates middleware that requires specific role
func (m *RBACMiddleware) RequireRole(roleSlug string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Check role
		hasRole, err := m.rbacService.CheckUserRole(c.Request.Context(), uid, roleSlug)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check role",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient role privileges",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole creates middleware that requires any of the specified roles
func (m *RBACMiddleware) RequireAnyRole(roleSlugs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasAnyRole := false
		for _, roleSlug := range roleSlugs {
			hasRole, err := m.rbacService.CheckUserRole(c.Request.Context(), uid, roleSlug)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check role",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			if hasRole {
				hasAnyRole = true
				break
			}
		}

		if !hasAnyRole {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient role privileges",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission creates middleware that requires any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(permissions [][]string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasAnyPermission := false
		for _, permission := range permissions {
			if len(permission) < 2 {
				continue
			}
			resource := permission[0]
			action := permission[1]

			hasPermission, err := m.rbacService.CheckUserPermission(c.Request.Context(), uid, resource, action)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check permission",
					"message": err.Error(),
				})
				c.Abort()
				return
			}
			if hasPermission {
				hasAnyPermission = true
				break
			}
		}

		if !hasAnyPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin creates middleware that requires super admin role
func (m *RBACMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return m.RequireRole("super-admin")
}

// RequireAdmin creates middleware that requires admin or super admin role
func (m *RBACMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireAnyRole("super-admin", "admin")
}

// RequireResourceOwnership creates middleware that checks if user owns the resource
// This can be used for endpoints like /users/:id where user should only access their own data
func (m *RBACMiddleware) RequireResourceOwnership(resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceID := c.Param(resourceIDParam)
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": "Resource ID not provided",
			})
			c.Abort()
			return
		}

		// Check if user owns the resource (user ID matches resource ID)
		if userID.(string) != resourceID {
			// Check if user has admin privileges as fallback
			uid, err := uuid.Parse(userID.(string))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid user ID",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			hasAdminRole, err := m.rbacService.CheckUserRole(c.Request.Context(), uid, "super-admin")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to check admin role",
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			if !hasAdminRole {
				hasAdminRole, err = m.rbacService.CheckUserRole(c.Request.Context(), uid, "admin")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error":   "Failed to check admin role",
						"message": err.Error(),
					})
					c.Abort()
					return
				}
			}

			if !hasAdminRole {
				c.JSON(http.StatusForbidden, gin.H{
					"error":   "Forbidden",
					"message": "Access denied: you can only access your own resources",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// DynamicPermissionCheck creates middleware for dynamic permission checking based on HTTP method and path
func (m *RBACMiddleware) DynamicPermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (should be set by auth middleware)
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Unauthorized",
				"message": "User not authenticated",
			})
			c.Abort()
			return
		}

		// Parse user ID
		uid, err := uuid.Parse(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid user ID",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		// Determine resource and action based on path and method
		resource, action := m.extractResourceAndAction(c.Request.URL.Path, c.Request.Method)
		if resource == "" || action == "" {
			// If we can't determine resource/action, allow access
			// You might want to change this behavior based on your security requirements
			c.Next()
			return
		}

		// Check permission
		hasPermission, err := m.rbacService.CheckUserPermission(c.Request.Context(), uid, resource, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check permission",
				"message": err.Error(),
			})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// extractResourceAndAction extracts resource and action from URL path and HTTP method
func (m *RBACMiddleware) extractResourceAndAction(path, method string) (string, string) {
	// Remove leading slash and split path
	path = strings.TrimPrefix(path, "/")
	segments := strings.Split(path, "/")

	if len(segments) < 3 {
		return "", ""
	}

	// Expected format: /api/v1/resource or /api/v1/rbac/resource
	var resource string
	if segments[2] == "rbac" && len(segments) > 3 {
		resource = segments[3]
	} else {
		resource = segments[2]
	}

	// Map HTTP methods to actions
	var action string
	switch method {
	case "GET":
		action = "view"
	case "POST":
		action = "create"
	case "PUT", "PATCH":
		action = "edit"
	case "DELETE":
		action = "delete"
	default:
		return "", ""
	}

	return resource, action
}
