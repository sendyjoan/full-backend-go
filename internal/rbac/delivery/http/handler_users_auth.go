package http

import (
	"net/http"

	"backend-service-internpro/internal/rbac"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// User Role handlers

// GetUserRoles godoc
// @Summary Get user roles
// @Description Get roles assigned to a specific user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} rbac.UserRoleListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/roles [get]
func (h *Handler) GetUserRoles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetUserRoles(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user roles",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AssignRolesToUser godoc
// @Summary Assign roles to user
// @Description Assign roles to a specific user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param roles body rbac.AssignUserRolesRequest true "Role IDs"
// @Success 200 {object} rbac.UserRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/roles [post]
func (h *Handler) AssignRolesToUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.AssignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get assigner user ID from context
	assignerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	assignedBy, err := uuid.Parse(assignerID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid assigner user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.AssignRolesToUser(c.Request.Context(), userID, &req, assignedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to assign roles to user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RemoveRolesFromUser godoc
// @Summary Remove roles from user
// @Description Remove specific roles from a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param roles body rbac.AssignUserRolesRequest true "Role IDs to remove"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/roles [delete]
func (h *Handler) RemoveRolesFromUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.AssignUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	if err := h.rbacService.RemoveRolesFromUser(c.Request.Context(), userID, req.RoleIDs); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to remove roles from user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Roles removed from user successfully",
	})
}

// GetUserPermissions godoc
// @Summary Get user permissions
// @Description Get all permissions available to a user through their roles
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} rbac.PermissionListResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/permissions [get]
func (h *Handler) GetUserPermissions(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user permissions",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserMenus godoc
// @Summary Get user menus
// @Description Get all menus available to a user through their roles
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} rbac.UserMenuResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/menus [get]
func (h *Handler) GetUserMenus(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetUserMenus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user menus",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserAccessibleMenus godoc
// @Summary Get user accessible menus
// @Description Get hierarchical menu structure accessible to a user
// @Tags user-roles
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Success 200 {object} rbac.UserMenuResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/users/{user_id}/accessible-menus [get]
func (h *Handler) GetUserAccessibleMenus(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetUserAccessibleMenus(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user accessible menus",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Authorization check handlers

// CheckUserPermissionRequest represents the request for checking user permission
type CheckUserPermissionRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	Resource string    `json:"resource" binding:"required"`
	Action   string    `json:"action" binding:"required"`
}

// CheckUserRoleRequest represents the request for checking user role
type CheckUserRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleSlug string    `json:"role_slug" binding:"required"`
}

// CheckPermissionResponse represents the response for permission check
type CheckPermissionResponse struct {
	HasPermission bool `json:"has_permission"`
}

// CheckRoleResponse represents the response for role check
type CheckRoleResponse struct {
	HasRole bool `json:"has_role"`
}

// CheckUserPermission godoc
// @Summary Check user permission
// @Description Check if a user has a specific permission
// @Tags authorization
// @Accept json
// @Produce json
// @Param request body CheckUserPermissionRequest true "Permission check data"
// @Success 200 {object} CheckPermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/auth/check-permission [post]
func (h *Handler) CheckUserPermission(c *gin.Context) {
	var req CheckUserPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	hasPermission, err := h.rbacService.CheckUserPermission(c.Request.Context(), req.UserID, req.Resource, req.Action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to check user permission",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CheckPermissionResponse{
		HasPermission: hasPermission,
	})
}

// CheckUserRole godoc
// @Summary Check user role
// @Description Check if a user has a specific role
// @Tags authorization
// @Accept json
// @Produce json
// @Param request body CheckUserRoleRequest true "Role check data"
// @Success 200 {object} CheckRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/auth/check-role [post]
func (h *Handler) CheckUserRole(c *gin.Context) {
	var req CheckUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	hasRole, err := h.rbacService.CheckUserRole(c.Request.Context(), req.UserID, req.RoleSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to check user role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, CheckRoleResponse{
		HasRole: hasRole,
	})
}
