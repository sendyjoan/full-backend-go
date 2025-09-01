package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/middleware"
	"backend-service-internpro/internal/rbac"
	"backend-service-internpro/internal/rbac/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type HumaHandler struct {
	rbacService service.Service
	jwtSecrets  jwt.Secrets
}

// NewHuma registers RBAC routes into the Huma API for Swagger documentation.
func NewHuma(api huma.API, rbacService service.Service, jwtSecrets jwt.Secrets) {
	h := &HumaHandler{
		rbacService: rbacService,
		jwtSecrets:  jwtSecrets,
	}

	// Role Management Routes
	roleGroup := huma.NewGroup(api, "/v1/roles")

	// GET /roles - List all roles
	huma.Register(roleGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of roles with pagination",
		Tags:    []string{"RBAC - Roles"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
		Search        string `query:"search" doc:"Search by name or slug"`
	}) (*struct {
		Body rbac.RoleListResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetRoles(ctx, in.Page, in.Limit, in.Search)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.RoleListResponse
		}{Body: *result}, nil
	})

	// GET /roles/{id} - Get role by ID
	huma.Register(roleGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}",
		Summary: "Get role by ID",
		Tags:    []string{"RBAC - Roles"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string    `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            uuid.UUID `path:"id" required:"true" doc:"Role ID"`
	}) (*struct {
		Body rbac.RoleResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetRoleByID(ctx, in.ID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		return &struct {
			Body rbac.RoleResponse
		}{Body: *result}, nil
	})

	// Permission Management Routes
	permissionGroup := huma.NewGroup(api, "/v1/permissions")

	// GET /permissions - List all permissions
	huma.Register(permissionGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of permissions with pagination",
		Tags:    []string{"RBAC - Permissions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
		Search        string `query:"search" doc:"Search by name, resource, or action"`
	}) (*struct {
		Body rbac.PermissionListResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetPermissions(ctx, in.Page, in.Limit, in.Search)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.PermissionListResponse
		}{Body: *result}, nil
	})

	// GET /permissions/{id} - Get permission by ID
	huma.Register(permissionGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}",
		Summary: "Get permission by ID",
		Tags:    []string{"RBAC - Permissions"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string    `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            uuid.UUID `path:"id" required:"true" doc:"Permission ID"`
	}) (*struct {
		Body rbac.PermissionResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetPermissionByID(ctx, in.ID)
		if err != nil {
			return nil, huma.Error404NotFound(err.Error())
		}

		return &struct {
			Body rbac.PermissionResponse
		}{Body: *result}, nil
	})

	// Menu Management Routes
	menuGroup := huma.NewGroup(api, "/v1/menus")

	// GET /menus - List all menus
	huma.Register(menuGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of menus with pagination",
		Tags:    []string{"RBAC - Menus"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
		Search        string `query:"search" doc:"Search by name or slug"`
	}) (*struct {
		Body rbac.MenuListResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetMenus(ctx, in.Page, in.Limit, in.Search)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.MenuListResponse
		}{Body: *result}, nil
	})

	// GET /menus/tree - Get menu tree
	huma.Register(menuGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/tree",
		Summary: "Get hierarchical menu tree",
		Tags:    []string{"RBAC - Menus"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
	}) (*struct {
		Body rbac.MenuTreeResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetMenuTree(ctx)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.MenuTreeResponse
		}{Body: *result}, nil
	})

	// User-Role Management Routes
	userRoleGroup := huma.NewGroup(api, "/v1/users")

	// GET /users/{id}/roles - Get user roles
	huma.Register(userRoleGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}/roles",
		Summary: "Get roles assigned to user",
		Tags:    []string{"RBAC - User Roles"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string    `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            uuid.UUID `path:"id" required:"true" doc:"User ID"`
	}) (*struct {
		Body rbac.UserRoleListResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetUserRoles(ctx, in.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.UserRoleListResponse
		}{Body: *result}, nil
	})

	// GET /users/{id}/permissions - Get user permissions
	huma.Register(userRoleGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}/permissions",
		Summary: "Get permissions available to user through roles",
		Tags:    []string{"RBAC - User Roles"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string    `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            uuid.UUID `path:"id" required:"true" doc:"User ID"`
	}) (*struct {
		Body rbac.PermissionListResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetUserPermissions(ctx, in.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.PermissionListResponse
		}{Body: *result}, nil
	})

	// GET /users/{id}/menus - Get user accessible menus
	huma.Register(userRoleGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}/menus",
		Summary: "Get menus accessible to user through roles",
		Tags:    []string{"RBAC - User Roles"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string    `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            uuid.UUID `path:"id" required:"true" doc:"User ID"`
	}) (*struct {
		Body rbac.UserMenuResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.rbacService.GetUserAccessibleMenus(ctx, in.ID)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body rbac.UserMenuResponse
		}{Body: *result}, nil
	})
}

// validateToken validates JWT token from Authorization header
func (h *HumaHandler) validateToken(authHeader string) error {
	_, err := middleware.ValidateToken(authHeader, h.jwtSecrets)
	return err
}
