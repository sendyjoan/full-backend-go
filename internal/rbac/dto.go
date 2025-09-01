package rbac

import (
	"time"

	"github.com/google/uuid"
)

// Role represents a role in the system
type Role struct {
	ID          uuid.UUID    `json:"id" doc:"Role ID"`
	Name        string       `json:"name" doc:"Role name"`
	Slug        string       `json:"slug" doc:"Role slug"`
	Description string       `json:"description" doc:"Role description"`
	IsActive    bool         `json:"is_active" doc:"Role active status"`
	CreatedAt   time.Time    `json:"created_at" doc:"Role creation date"`
	UpdatedAt   time.Time    `json:"updated_at" doc:"Role last update date"`
	Permissions []Permission `json:"permissions,omitempty" doc:"Role permissions"`
	Menus       []Menu       `json:"menus,omitempty" doc:"Role menus"`
}

// Permission represents a permission in the system
type Permission struct {
	ID          uuid.UUID `json:"id" doc:"Permission ID"`
	Name        string    `json:"name" doc:"Permission name"`
	Slug        string    `json:"slug" doc:"Permission slug"`
	Resource    string    `json:"resource" doc:"Permission resource"`
	Action      string    `json:"action" doc:"Permission action"`
	Description string    `json:"description" doc:"Permission description"`
	IsActive    bool      `json:"is_active" doc:"Permission active status"`
	CreatedAt   time.Time `json:"created_at" doc:"Permission creation date"`
	UpdatedAt   time.Time `json:"updated_at" doc:"Permission last update date"`
}

// Menu represents a menu in the system
type Menu struct {
	ID        uuid.UUID  `json:"id" doc:"Menu ID"`
	Name      string     `json:"name" doc:"Menu name"`
	Slug      string     `json:"slug" doc:"Menu slug"`
	URL       string     `json:"url" doc:"Menu URL"`
	Icon      string     `json:"icon" doc:"Menu icon"`
	ParentID  *uuid.UUID `json:"parent_id" doc:"Parent menu ID"`
	SortOrder int        `json:"sort_order" doc:"Menu sort order"`
	IsActive  bool       `json:"is_active" doc:"Menu active status"`
	CreatedAt time.Time  `json:"created_at" doc:"Menu creation date"`
	UpdatedAt time.Time  `json:"updated_at" doc:"Menu last update date"`
	Children  []Menu     `json:"children,omitempty" doc:"Child menus"`
}

// RoleMenu represents role-menu relationship with permissions
type RoleMenu struct {
	ID        uuid.UUID `json:"id" doc:"Role-Menu ID"`
	RoleID    uuid.UUID `json:"role_id" doc:"Role ID"`
	MenuID    uuid.UUID `json:"menu_id" doc:"Menu ID"`
	CanView   bool      `json:"can_view" doc:"Can view permission"`
	CanCreate bool      `json:"can_create" doc:"Can create permission"`
	CanEdit   bool      `json:"can_edit" doc:"Can edit permission"`
	CanDelete bool      `json:"can_delete" doc:"Can delete permission"`
	Menu      Menu      `json:"menu" doc:"Menu details"`
	CreatedAt time.Time `json:"created_at" doc:"Assignment creation date"`
	UpdatedAt time.Time `json:"updated_at" doc:"Assignment last update date"`
}

// UserRole represents user-role relationship
type UserRole struct {
	ID         uuid.UUID `json:"id" doc:"User-Role ID"`
	UserID     uuid.UUID `json:"user_id" doc:"User ID"`
	RoleID     uuid.UUID `json:"role_id" doc:"Role ID"`
	Role       Role      `json:"role" doc:"Role details"`
	AssignedAt time.Time `json:"assigned_at" doc:"Role assignment date"`
}

// RBACMetadata represents pagination metadata for RBAC responses
type RBACMetadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// Role Request/Response DTOs
type RoleListResponse struct {
	Data []Role       `json:"data"`
	Meta RBACMetadata `json:"meta"`
}

// PaginatedRolesResponse represents paginated roles response for Huma
type PaginatedRolesResponse struct {
	Data []Role       `json:"data" doc:"List of roles"`
	Meta RBACMetadata `json:"meta" doc:"Pagination metadata"`
}

// RoleQueryParams represents query parameters for role listing
type RoleQueryParams struct {
	Page   int    `json:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int    `json:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
	Search string `json:"search" doc:"Search by name or slug"`
}

type RoleResponse struct {
	Role
}

type CreateRoleRequest struct {
	Name        string `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Role name"`
	Slug        string `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Role slug"`
	Description string `json:"description" form:"description" maxLength:"1000" doc:"Role description"`
	IsActive    *bool  `json:"is_active" form:"is_active" doc:"Role active status"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Role name"`
	Slug        *string `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Role slug"`
	Description *string `json:"description" form:"description" maxLength:"1000" doc:"Role description"`
	IsActive    *bool   `json:"is_active" form:"is_active" doc:"Role active status"`
}

type CreateRoleResponse struct {
	ID      uuid.UUID `json:"id" doc:"Created role ID"`
	Message string    `json:"message" doc:"Success message"`
}

type AssignRolePermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" doc:"List of permission IDs to assign"`
}

// AssignPermissionsToRoleRequest represents request to assign permissions to role
type AssignPermissionsToRoleRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids" required:"true" doc:"List of permission IDs to assign"`
}

// AssignRoleToUserRequest represents request to assign role to user
type AssignRoleToUserRequest struct {
	RoleID uuid.UUID `json:"role_id" required:"true" doc:"Role ID to assign"`
}

// Permission Request/Response DTOs
type PermissionListResponse struct {
	Data []Permission `json:"data"`
	Meta RBACMetadata `json:"meta"`
}

// PaginatedPermissionsResponse represents paginated permissions response for Huma
type PaginatedPermissionsResponse struct {
	Data []Permission `json:"data" doc:"List of permissions"`
	Meta RBACMetadata `json:"meta" doc:"Pagination metadata"`
}

// PermissionQueryParams represents query parameters for permission listing
type PermissionQueryParams struct {
	Page   int    `json:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int    `json:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
	Search string `json:"search" doc:"Search by name, resource, or action"`
}

type PermissionResponse struct {
	Permission
}

type CreatePermissionRequest struct {
	Name        string `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Permission name"`
	Slug        string `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Permission slug"`
	Resource    string `json:"resource" form:"resource" minLength:"1" maxLength:"100" doc:"Permission resource"`
	Action      string `json:"action" form:"action" minLength:"1" maxLength:"50" doc:"Permission action"`
	Description string `json:"description" form:"description" maxLength:"1000" doc:"Permission description"`
	IsActive    *bool  `json:"is_active" form:"is_active" doc:"Permission active status"`
}

type UpdatePermissionRequest struct {
	Name        *string `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Permission name"`
	Slug        *string `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Permission slug"`
	Resource    *string `json:"resource" form:"resource" minLength:"1" maxLength:"100" doc:"Permission resource"`
	Action      *string `json:"action" form:"action" minLength:"1" maxLength:"50" doc:"Permission action"`
	Description *string `json:"description" form:"description" maxLength:"1000" doc:"Permission description"`
	IsActive    *bool   `json:"is_active" form:"is_active" doc:"Permission active status"`
}

type CreatePermissionResponse struct {
	ID      uuid.UUID `json:"id" doc:"Created permission ID"`
	Message string    `json:"message" doc:"Success message"`
}

// Menu Request/Response DTOs
type MenuListResponse struct {
	Data []Menu       `json:"data"`
	Meta RBACMetadata `json:"meta"`
}

// PaginatedMenusResponse represents paginated menus response for Huma
type PaginatedMenusResponse struct {
	Data []Menu       `json:"data" doc:"List of menus"`
	Meta RBACMetadata `json:"meta" doc:"Pagination metadata"`
}

// MenuQueryParams represents query parameters for menu listing
type MenuQueryParams struct {
	Page   int    `json:"page" minimum:"1" default:"1" doc:"Page number"`
	Limit  int    `json:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
	Search string `json:"search" doc:"Search by name or slug"`
}

type MenuResponse struct {
	Menu
}

type CreateMenuRequest struct {
	Name      string     `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Menu name"`
	Slug      string     `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Menu slug"`
	URL       string     `json:"url" form:"url" maxLength:"255" doc:"Menu URL"`
	Icon      string     `json:"icon" form:"icon" maxLength:"100" doc:"Menu icon"`
	ParentID  *uuid.UUID `json:"parent_id" form:"parent_id" doc:"Parent menu ID"`
	SortOrder *int       `json:"sort_order" form:"sort_order" doc:"Menu sort order"`
	IsActive  *bool      `json:"is_active" form:"is_active" doc:"Menu active status"`
}

type UpdateMenuRequest struct {
	Name      *string    `json:"name" form:"name" minLength:"1" maxLength:"100" doc:"Menu name"`
	Slug      *string    `json:"slug" form:"slug" minLength:"1" maxLength:"100" doc:"Menu slug"`
	URL       *string    `json:"url" form:"url" maxLength:"255" doc:"Menu URL"`
	Icon      *string    `json:"icon" form:"icon" maxLength:"100" doc:"Menu icon"`
	ParentID  *uuid.UUID `json:"parent_id" form:"parent_id" doc:"Parent menu ID"`
	SortOrder *int       `json:"sort_order" form:"sort_order" doc:"Menu sort order"`
	IsActive  *bool      `json:"is_active" form:"is_active" doc:"Menu active status"`
}

type CreateMenuResponse struct {
	ID      uuid.UUID `json:"id" doc:"Created menu ID"`
	Message string    `json:"message" doc:"Success message"`
}

type AssignRoleMenusRequest struct {
	MenuPermissions []MenuPermissionRequest `json:"menu_permissions" doc:"List of menu permissions to assign"`
}

type MenuPermissionRequest struct {
	MenuID    uuid.UUID `json:"menu_id" doc:"Menu ID"`
	CanView   bool      `json:"can_view" doc:"Can view permission"`
	CanCreate bool      `json:"can_create" doc:"Can create permission"`
	CanEdit   bool      `json:"can_edit" doc:"Can edit permission"`
	CanDelete bool      `json:"can_delete" doc:"Can delete permission"`
}

// User Role Request/Response DTOs
type UserRoleListResponse struct {
	Data []UserRole   `json:"data"`
	Meta RBACMetadata `json:"meta"`
}

type AssignUserRolesRequest struct {
	RoleIDs []uuid.UUID `json:"role_ids" doc:"List of role IDs to assign"`
}

type UserRoleResponse struct {
	ID      uuid.UUID `json:"id" doc:"Assignment ID"`
	Message string    `json:"message" doc:"Success message"`
}

// Menu Tree Response for hierarchical menu structure
type MenuTreeResponse struct {
	Data []Menu `json:"data"`
}

// User Menu Response with permissions
type UserMenuResponse struct {
	Data []RoleMenu `json:"data"`
}
