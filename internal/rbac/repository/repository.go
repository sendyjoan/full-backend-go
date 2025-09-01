package repository

import (
	"context"

	"backend-service-internpro/internal/rbac"

	"github.com/google/uuid"
)

type Repository interface {
	// Role methods
	CreateRole(ctx context.Context, role *rbac.RoleEntity) error
	GetRoleByID(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error)
	GetRoleBySlug(ctx context.Context, slug string) (*rbac.RoleEntity, error)
	GetRoles(ctx context.Context, page, limit int, search string) ([]rbac.RoleEntity, int64, error)
	UpdateRole(ctx context.Context, role *rbac.RoleEntity) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error)
	GetRoleWithMenus(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error)

	// Permission methods
	CreatePermission(ctx context.Context, permission *rbac.PermissionEntity) error
	GetPermissionByID(ctx context.Context, id uuid.UUID) (*rbac.PermissionEntity, error)
	GetPermissionBySlug(ctx context.Context, slug string) (*rbac.PermissionEntity, error)
	GetPermissions(ctx context.Context, page, limit int, search string) ([]rbac.PermissionEntity, int64, error)
	UpdatePermission(ctx context.Context, permission *rbac.PermissionEntity) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	GetPermissionsByResource(ctx context.Context, resource string) ([]rbac.PermissionEntity, error)
	GetPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]rbac.PermissionEntity, error)

	// Menu methods
	CreateMenu(ctx context.Context, menu *rbac.MenuEntity) error
	GetMenuByID(ctx context.Context, id uuid.UUID) (*rbac.MenuEntity, error)
	GetMenuBySlug(ctx context.Context, slug string) (*rbac.MenuEntity, error)
	GetMenus(ctx context.Context, page, limit int, search string) ([]rbac.MenuEntity, int64, error)
	GetMenuTree(ctx context.Context) ([]rbac.MenuEntity, error)
	UpdateMenu(ctx context.Context, menu *rbac.MenuEntity) error
	DeleteMenu(ctx context.Context, id uuid.UUID) error
	GetMenusByParentID(ctx context.Context, parentID *uuid.UUID) ([]rbac.MenuEntity, error)
	GetMenusByIDs(ctx context.Context, ids []uuid.UUID) ([]rbac.MenuEntity, error)

	// Role-Permission methods
	AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, assignedBy uuid.UUID) error
	RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]rbac.PermissionEntity, error)
	CheckRoleHasPermission(ctx context.Context, roleID uuid.UUID, permissionSlug string) (bool, error)

	// User-Role methods
	AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error
	RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]rbac.UserRoleEntity, error)
	GetUsersByRole(ctx context.Context, roleID uuid.UUID, page, limit int) ([]rbac.UserRoleEntity, int64, error)
	CheckUserHasRole(ctx context.Context, userID uuid.UUID, roleSlug string) (bool, error)

	// Role-Menu methods
	AssignMenusToRole(ctx context.Context, roleID uuid.UUID, menuPermissions []rbac.RoleMenuEntity, assignedBy uuid.UUID) error
	RemoveMenusFromRole(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error
	GetRoleMenus(ctx context.Context, roleID uuid.UUID) ([]rbac.RoleMenuEntity, error)
	GetUserMenus(ctx context.Context, userID uuid.UUID) ([]rbac.RoleMenuEntity, error)
	UpdateRoleMenuPermissions(ctx context.Context, roleMenuID uuid.UUID, canView, canCreate, canEdit, canDelete bool, updatedBy uuid.UUID) error

	// Complex queries
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]rbac.PermissionEntity, error)
	CheckUserHasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error)
	GetUserAccessibleMenus(ctx context.Context, userID uuid.UUID) ([]rbac.RoleMenuEntity, error)
}
