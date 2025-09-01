package http

import (
	"net/http"
	"strconv"

	"backend-service-internpro/internal/rbac"
	"backend-service-internpro/internal/rbac/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	rbacService service.Service
}

// NewHandler creates a new RBAC HTTP handler
func NewHandler(rbacService service.Service) *Handler {
	return &Handler{
		rbacService: rbacService,
	}
}

// RegisterRoutes registers all RBAC routes
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	rbac := router.Group("/rbac")
	{
		// Role routes
		roles := rbac.Group("/roles")
		{
			roles.POST("", h.CreateRole)
			roles.GET("", h.GetRoles)
			roles.GET("/:id", h.GetRoleByID)
			roles.PUT("/:id", h.UpdateRole)
			roles.DELETE("/:id", h.DeleteRole)
			roles.GET("/:id/permissions", h.GetRolePermissions)
			roles.POST("/:id/permissions", h.AssignPermissionsToRole)
			roles.GET("/:id/menus", h.GetRoleMenus)
			roles.POST("/:id/menus", h.AssignMenusToRole)
		}

		// Permission routes
		permissions := rbac.Group("/permissions")
		{
			permissions.POST("", h.CreatePermission)
			permissions.GET("", h.GetPermissions)
			permissions.GET("/:id", h.GetPermissionByID)
			permissions.PUT("/:id", h.UpdatePermission)
			permissions.DELETE("/:id", h.DeletePermission)
			permissions.GET("/resource/:resource", h.GetPermissionsByResource)
		}

		// Menu routes
		menus := rbac.Group("/menus")
		{
			menus.POST("", h.CreateMenu)
			menus.GET("", h.GetMenus)
			menus.GET("/tree", h.GetMenuTree)
			menus.GET("/:id", h.GetMenuByID)
			menus.PUT("/:id", h.UpdateMenu)
			menus.DELETE("/:id", h.DeleteMenu)
		}

		// User role routes
		users := rbac.Group("/users")
		{
			users.GET("/:user_id/roles", h.GetUserRoles)
			users.POST("/:user_id/roles", h.AssignRolesToUser)
			users.DELETE("/:user_id/roles", h.RemoveRolesFromUser)
			users.GET("/:user_id/permissions", h.GetUserPermissions)
			users.GET("/:user_id/menus", h.GetUserMenus)
			users.GET("/:user_id/accessible-menus", h.GetUserAccessibleMenus)
		}

		// Authorization check routes
		auth := rbac.Group("/auth")
		{
			auth.POST("/check-permission", h.CheckUserPermission)
			auth.POST("/check-role", h.CheckUserRole)
		}
	}
}

// Role handlers

// CreateRole godoc
// @Summary Create a new role
// @Description Create a new role in the system
// @Tags roles
// @Accept json
// @Produce json
// @Param role body rbac.CreateRoleRequest true "Role data"
// @Success 201 {object} rbac.CreateRoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles [post]
func (h *Handler) CreateRole(c *gin.Context) {
	var req rbac.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (should be set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	createdBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.CreateRole(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetRoles godoc
// @Summary Get roles list
// @Description Get paginated list of roles
// @Tags roles
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} rbac.RoleListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles [get]
func (h *Handler) GetRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	response, err := h.rbacService.GetRoles(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get roles",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRoleByID godoc
// @Summary Get role by ID
// @Description Get a specific role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} rbac.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id} [get]
func (h *Handler) GetRoleByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetRoleByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Role not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRole godoc
// @Summary Update a role
// @Description Update an existing role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body rbac.UpdateRoleRequest true "Role data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id} [put]
func (h *Handler) UpdateRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	updatedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.rbacService.UpdateRole(c.Request.Context(), id, &req, updatedBy); err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Role not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Role updated successfully",
	})
}

// DeleteRole godoc
// @Summary Delete a role
// @Description Delete an existing role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id} [delete]
func (h *Handler) DeleteRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	deletedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.rbacService.DeleteRole(c.Request.Context(), id, deletedBy); err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Role not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Role deleted successfully",
	})
}

// GetRolePermissions godoc
// @Summary Get role permissions
// @Description Get permissions assigned to a specific role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} rbac.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id}/permissions [get]
func (h *Handler) GetRolePermissions(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetRoleWithPermissions(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Role not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get role permissions",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AssignPermissionsToRole godoc
// @Summary Assign permissions to role
// @Description Assign permissions to a specific role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param permissions body rbac.AssignRolePermissionsRequest true "Permission IDs"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id}/permissions [post]
func (h *Handler) AssignPermissionsToRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.AssignRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	assignedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.rbacService.AssignPermissionsToRole(c.Request.Context(), id, &req, assignedBy); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to assign permissions to role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Permissions assigned to role successfully",
	})
}

// GetRoleMenus godoc
// @Summary Get role menus
// @Description Get menus assigned to a specific role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} rbac.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id}/menus [get]
func (h *Handler) GetRoleMenus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetRoleWithMenus(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "role not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Role not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get role menus",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// AssignMenusToRole godoc
// @Summary Assign menus to role
// @Description Assign menus with permissions to a specific role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param menus body rbac.AssignRoleMenusRequest true "Menu permissions"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/roles/{id}/menus [post]
func (h *Handler) AssignMenusToRole(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid role ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.AssignRoleMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	assignedBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.rbacService.AssignMenusToRole(c.Request.Context(), id, &req, assignedBy); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to assign menus to role",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Menus assigned to role successfully",
	})
}

// Response types
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
