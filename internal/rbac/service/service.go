package service

import (
	"context"

	"backend-service-internpro/internal/rbac"

	"github.com/google/uuid"
)

type Service interface {
	// Role services
	CreateRole(ctx context.Context, req *rbac.CreateRoleRequest, createdBy uuid.UUID) (*rbac.CreateRoleResponse, error)
	GetRoleByID(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error)
	GetRoles(ctx context.Context, page, limit int, search string) (*rbac.RoleListResponse, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req *rbac.UpdateRoleRequest, updatedBy uuid.UUID) error
	DeleteRole(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error)
	GetRoleWithMenus(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error)
	AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *rbac.AssignRolePermissionsRequest, assignedBy uuid.UUID) error
	AssignMenusToRole(ctx context.Context, roleID uuid.UUID, req *rbac.AssignRoleMenusRequest, assignedBy uuid.UUID) error

	// Permission services
	CreatePermission(ctx context.Context, req *rbac.CreatePermissionRequest, createdBy uuid.UUID) (*rbac.CreatePermissionResponse, error)
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*rbac.PermissionResponse, error)
	GetPermissions(ctx context.Context, page, limit int, search string) (*rbac.PermissionListResponse, error)
	UpdatePermission(ctx context.Context, id uuid.UUID, req *rbac.UpdatePermissionRequest, updatedBy uuid.UUID) error
	DeletePermission(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error
	GetPermissionsByResource(ctx context.Context, resource string) (*rbac.PermissionListResponse, error)

	// Menu services
	CreateMenu(ctx context.Context, req *rbac.CreateMenuRequest, createdBy uuid.UUID) (*rbac.CreateMenuResponse, error)
	GetMenuByID(ctx context.Context, id uuid.UUID) (*rbac.MenuResponse, error)
	GetMenus(ctx context.Context, page, limit int, search string) (*rbac.MenuListResponse, error)
	GetMenuTree(ctx context.Context) (*rbac.MenuTreeResponse, error)
	UpdateMenu(ctx context.Context, id uuid.UUID, req *rbac.UpdateMenuRequest, updatedBy uuid.UUID) error
	DeleteMenu(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error

	// User-Role services
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, req *rbac.AssignUserRolesRequest, assignedBy uuid.UUID) (*rbac.UserRoleResponse, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) (*rbac.UserRoleListResponse, error)
	RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error

	// Authorization services
	CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)
	CheckUserRole(ctx context.Context, userID uuid.UUID, roleSlug string) (bool, error)
	GetUserPermissions(ctx context.Context, userID uuid.UUID) (*rbac.PermissionListResponse, error)
	GetUserMenus(ctx context.Context, userID uuid.UUID) (*rbac.UserMenuResponse, error)
	GetUserAccessibleMenus(ctx context.Context, userID uuid.UUID) (*rbac.UserMenuResponse, error)

	// Validation services
	ValidateRoleSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error
	ValidatePermissionSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error
	ValidateMenuSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error
}
