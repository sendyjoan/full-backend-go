package http

import (
	"net/http"
	"strconv"

	"backend-service-internpro/internal/rbac"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Permission handlers

// CreatePermission godoc
// @Summary Create a new permission
// @Description Create a new permission in the system
// @Tags permissions
// @Accept json
// @Produce json
// @Param permission body rbac.CreatePermissionRequest true "Permission data"
// @Success 201 {object} rbac.CreatePermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions [post]
func (h *Handler) CreatePermission(c *gin.Context) {
	var req rbac.CreatePermissionRequest
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

	createdBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.CreatePermission(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create permission",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPermissions godoc
// @Summary Get permissions list
// @Description Get paginated list of permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} rbac.PermissionListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions [get]
func (h *Handler) GetPermissions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	response, err := h.rbacService.GetPermissions(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get permissions",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetPermissionByID godoc
// @Summary Get permission by ID
// @Description Get a specific permission by its ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} rbac.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions/{id} [get]
func (h *Handler) GetPermissionByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid permission ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetPermissionByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "permission not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Permission not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get permission",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdatePermission godoc
// @Summary Update a permission
// @Description Update an existing permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Param permission body rbac.UpdatePermissionRequest true "Permission data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions/{id} [put]
func (h *Handler) UpdatePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid permission ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.UpdatePermissionRequest
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

	if err := h.rbacService.UpdatePermission(c.Request.Context(), id, &req, updatedBy); err != nil {
		if err.Error() == "permission not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Permission not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update permission",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Permission updated successfully",
	})
}

// DeletePermission godoc
// @Summary Delete a permission
// @Description Delete an existing permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions/{id} [delete]
func (h *Handler) DeletePermission(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid permission ID",
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

	if err := h.rbacService.DeletePermission(c.Request.Context(), id, deletedBy); err != nil {
		if err.Error() == "permission not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Permission not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete permission",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Permission deleted successfully",
	})
}

// GetPermissionsByResource godoc
// @Summary Get permissions by resource
// @Description Get all permissions for a specific resource
// @Tags permissions
// @Accept json
// @Produce json
// @Param resource path string true "Resource name"
// @Success 200 {object} rbac.PermissionListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/permissions/resource/{resource} [get]
func (h *Handler) GetPermissionsByResource(c *gin.Context) {
	resource := c.Param("resource")

	response, err := h.rbacService.GetPermissionsByResource(c.Request.Context(), resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get permissions by resource",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Menu handlers

// CreateMenu godoc
// @Summary Create a new menu
// @Description Create a new menu in the system
// @Tags menus
// @Accept json
// @Produce json
// @Param menu body rbac.CreateMenuRequest true "Menu data"
// @Success 201 {object} rbac.CreateMenuResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus [post]
func (h *Handler) CreateMenu(c *gin.Context) {
	var req rbac.CreateMenuRequest
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

	createdBy, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.CreateMenu(c.Request.Context(), &req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create menu",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetMenus godoc
// @Summary Get menus list
// @Description Get paginated list of menus
// @Tags menus
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search term"
// @Success 200 {object} rbac.MenuListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus [get]
func (h *Handler) GetMenus(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	response, err := h.rbacService.GetMenus(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get menus",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMenuTree godoc
// @Summary Get menu tree
// @Description Get hierarchical menu structure
// @Tags menus
// @Accept json
// @Produce json
// @Success 200 {object} rbac.MenuTreeResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus/tree [get]
func (h *Handler) GetMenuTree(c *gin.Context) {
	response, err := h.rbacService.GetMenuTree(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get menu tree",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMenuByID godoc
// @Summary Get menu by ID
// @Description Get a specific menu by its ID
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} rbac.MenuResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus/{id} [get]
func (h *Handler) GetMenuByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid menu ID",
			Message: err.Error(),
		})
		return
	}

	response, err := h.rbacService.GetMenuByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "menu not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Menu not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get menu",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateMenu godoc
// @Summary Update a menu
// @Description Update an existing menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Param menu body rbac.UpdateMenuRequest true "Menu data"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus/{id} [put]
func (h *Handler) UpdateMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid menu ID",
			Message: err.Error(),
		})
		return
	}

	var req rbac.UpdateMenuRequest
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

	if err := h.rbacService.UpdateMenu(c.Request.Context(), id, &req, updatedBy); err != nil {
		if err.Error() == "menu not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Menu not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update menu",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Menu updated successfully",
	})
}

// DeleteMenu godoc
// @Summary Delete a menu
// @Description Delete an existing menu
// @Tags menus
// @Accept json
// @Produce json
// @Param id path string true "Menu ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/rbac/menus/{id} [delete]
func (h *Handler) DeleteMenu(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid menu ID",
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

	if err := h.rbacService.DeleteMenu(c.Request.Context(), id, deletedBy); err != nil {
		if err.Error() == "menu not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "Menu not found",
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete menu",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Menu deleted successfully",
	})
}
