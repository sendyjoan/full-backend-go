package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"backend-service-internpro/internal/rbac"
	"backend-service-internpro/internal/rbac/repository"

	"github.com/google/uuid"
)

type service struct {
	repo repository.Repository
}

// NewService creates a new RBAC service
func NewService(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

// Role services
func (s *service) CreateRole(ctx context.Context, req *rbac.CreateRoleRequest, createdBy uuid.UUID) (*rbac.CreateRoleResponse, error) {
	// Validate slug uniqueness
	if err := s.ValidateRoleSlug(ctx, req.Slug, nil); err != nil {
		return nil, err
	}

	role := &rbac.RoleEntity{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    req.IsActive != nil && *req.IsActive,
		CreatedBy:   &createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreateRole(ctx, role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &rbac.CreateRoleResponse{
		ID:      role.ID,
		Message: "Role created successfully",
	}, nil
}

func (s *service) GetRoleByID(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error) {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	return &rbac.RoleResponse{
		Role: role.ToRole(),
	}, nil
}

func (s *service) GetRoles(ctx context.Context, page, limit int, search string) (*rbac.RoleListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	roles, total, err := s.repo.GetRoles(ctx, page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}

	var roleList []rbac.Role
	for _, role := range roles {
		roleList = append(roleList, role.ToRole())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &rbac.RoleListResponse{
		Data: roleList,
		Meta: rbac.RBACMetadata{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: int(total),
		},
	}, nil
}

func (s *service) UpdateRole(ctx context.Context, id uuid.UUID, req *rbac.UpdateRoleRequest, updatedBy uuid.UUID) error {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Update fields if provided
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Slug != nil {
		// Validate slug uniqueness
		if err := s.ValidateRoleSlug(ctx, *req.Slug, &id); err != nil {
			return err
		}
		role.Slug = *req.Slug
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.IsActive != nil {
		role.IsActive = *req.IsActive
	}

	role.UpdatedBy = &updatedBy
	role.UpdatedAt = time.Now()

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

func (s *service) DeleteRole(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	role.DeletedBy = &deletedBy
	now := time.Now()
	role.DeletedAt = &now

	if err := s.repo.UpdateRole(ctx, role); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

func (s *service) GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error) {
	role, err := s.repo.GetRoleWithPermissions(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role with permissions: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	return &rbac.RoleResponse{
		Role: role.ToRole(),
	}, nil
}

func (s *service) GetRoleWithMenus(ctx context.Context, id uuid.UUID) (*rbac.RoleResponse, error) {
	role, err := s.repo.GetRoleWithMenus(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role with menus: %w", err)
	}
	if role == nil {
		return nil, errors.New("role not found")
	}

	return &rbac.RoleResponse{
		Role: role.ToRole(),
	}, nil
}

func (s *service) AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, req *rbac.AssignRolePermissionsRequest, assignedBy uuid.UUID) error {
	// Check if role exists
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Validate permissions exist
	permissions, err := s.repo.GetPermissionsByIDs(ctx, req.PermissionIDs)
	if err != nil {
		return fmt.Errorf("failed to get permissions: %w", err)
	}
	if len(permissions) != len(req.PermissionIDs) {
		return errors.New("some permissions not found")
	}

	if err := s.repo.AssignPermissionsToRole(ctx, roleID, req.PermissionIDs, assignedBy); err != nil {
		return fmt.Errorf("failed to assign permissions to role: %w", err)
	}

	return nil
}

func (s *service) AssignMenusToRole(ctx context.Context, roleID uuid.UUID, req *rbac.AssignRoleMenusRequest, assignedBy uuid.UUID) error {
	// Check if role exists
	role, err := s.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return errors.New("role not found")
	}

	// Validate menus exist
	var menuIDs []uuid.UUID
	for _, mp := range req.MenuPermissions {
		menuIDs = append(menuIDs, mp.MenuID)
	}

	menus, err := s.repo.GetMenusByIDs(ctx, menuIDs)
	if err != nil {
		return fmt.Errorf("failed to get menus: %w", err)
	}
	if len(menus) != len(menuIDs) {
		return errors.New("some menus not found")
	}

	// Convert to entities
	var roleMenus []rbac.RoleMenuEntity
	for _, mp := range req.MenuPermissions {
		roleMenus = append(roleMenus, rbac.RoleMenuEntity{
			MenuID:    mp.MenuID,
			CanView:   mp.CanView,
			CanCreate: mp.CanCreate,
			CanEdit:   mp.CanEdit,
			CanDelete: mp.CanDelete,
		})
	}

	if err := s.repo.AssignMenusToRole(ctx, roleID, roleMenus, assignedBy); err != nil {
		return fmt.Errorf("failed to assign menus to role: %w", err)
	}

	return nil
}

// Permission services
func (s *service) CreatePermission(ctx context.Context, req *rbac.CreatePermissionRequest, createdBy uuid.UUID) (*rbac.CreatePermissionResponse, error) {
	// Validate slug uniqueness
	if err := s.ValidatePermissionSlug(ctx, req.Slug, nil); err != nil {
		return nil, err
	}

	permission := &rbac.PermissionEntity{
		ID:          uuid.New(),
		Name:        req.Name,
		Slug:        req.Slug,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
		IsActive:    req.IsActive != nil && *req.IsActive,
		CreatedBy:   &createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.CreatePermission(ctx, permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &rbac.CreatePermissionResponse{
		ID:      permission.ID,
		Message: "Permission created successfully",
	}, nil
}

func (s *service) GetPermissionByID(ctx context.Context, id uuid.UUID) (*rbac.PermissionResponse, error) {
	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return nil, errors.New("permission not found")
	}

	return &rbac.PermissionResponse{
		Permission: permission.ToPermission(),
	}, nil
}

func (s *service) GetPermissions(ctx context.Context, page, limit int, search string) (*rbac.PermissionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	permissions, total, err := s.repo.GetPermissions(ctx, page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}

	var permissionList []rbac.Permission
	for _, permission := range permissions {
		permissionList = append(permissionList, permission.ToPermission())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &rbac.PermissionListResponse{
		Data: permissionList,
		Meta: rbac.RBACMetadata{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: int(total),
		},
	}, nil
}

func (s *service) UpdatePermission(ctx context.Context, id uuid.UUID, req *rbac.UpdatePermissionRequest, updatedBy uuid.UUID) error {
	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return errors.New("permission not found")
	}

	// Update fields if provided
	if req.Name != nil {
		permission.Name = *req.Name
	}
	if req.Slug != nil {
		// Validate slug uniqueness
		if err := s.ValidatePermissionSlug(ctx, *req.Slug, &id); err != nil {
			return err
		}
		permission.Slug = *req.Slug
	}
	if req.Resource != nil {
		permission.Resource = *req.Resource
	}
	if req.Action != nil {
		permission.Action = *req.Action
	}
	if req.Description != nil {
		permission.Description = *req.Description
	}
	if req.IsActive != nil {
		permission.IsActive = *req.IsActive
	}

	permission.UpdatedBy = &updatedBy
	permission.UpdatedAt = time.Now()

	if err := s.repo.UpdatePermission(ctx, permission); err != nil {
		return fmt.Errorf("failed to update permission: %w", err)
	}

	return nil
}

func (s *service) DeletePermission(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	permission, err := s.repo.GetPermissionByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return errors.New("permission not found")
	}

	permission.DeletedBy = &deletedBy
	now := time.Now()
	permission.DeletedAt = &now

	if err := s.repo.UpdatePermission(ctx, permission); err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	return nil
}

func (s *service) GetPermissionsByResource(ctx context.Context, resource string) (*rbac.PermissionListResponse, error) {
	permissions, err := s.repo.GetPermissionsByResource(ctx, resource)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions by resource: %w", err)
	}

	var permissionList []rbac.Permission
	for _, permission := range permissions {
		permissionList = append(permissionList, permission.ToPermission())
	}

	return &rbac.PermissionListResponse{
		Data: permissionList,
		Meta: rbac.RBACMetadata{
			Page:       1,
			Limit:      len(permissionList),
			TotalPages: 1,
			TotalItems: len(permissionList),
		},
	}, nil
}

// Menu services
func (s *service) CreateMenu(ctx context.Context, req *rbac.CreateMenuRequest, createdBy uuid.UUID) (*rbac.CreateMenuResponse, error) {
	// Validate slug uniqueness
	if err := s.ValidateMenuSlug(ctx, req.Slug, nil); err != nil {
		return nil, err
	}

	// Validate parent menu exists if provided
	if req.ParentID != nil {
		parent, err := s.repo.GetMenuByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent menu: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent menu not found")
		}
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	menu := &rbac.MenuEntity{
		ID:        uuid.New(),
		Name:      req.Name,
		Slug:      req.Slug,
		URL:       req.URL,
		Icon:      req.Icon,
		ParentID:  req.ParentID,
		SortOrder: sortOrder,
		IsActive:  req.IsActive != nil && *req.IsActive,
		CreatedBy: &createdBy,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.CreateMenu(ctx, menu); err != nil {
		return nil, fmt.Errorf("failed to create menu: %w", err)
	}

	return &rbac.CreateMenuResponse{
		ID:      menu.ID,
		Message: "Menu created successfully",
	}, nil
}

func (s *service) GetMenuByID(ctx context.Context, id uuid.UUID) (*rbac.MenuResponse, error) {
	menu, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu: %w", err)
	}
	if menu == nil {
		return nil, errors.New("menu not found")
	}

	return &rbac.MenuResponse{
		Menu: menu.ToMenu(),
	}, nil
}

func (s *service) GetMenus(ctx context.Context, page, limit int, search string) (*rbac.MenuListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	menus, total, err := s.repo.GetMenus(ctx, page, limit, search)
	if err != nil {
		return nil, fmt.Errorf("failed to get menus: %w", err)
	}

	var menuList []rbac.Menu
	for _, menu := range menus {
		menuList = append(menuList, menu.ToMenu())
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &rbac.MenuListResponse{
		Data: menuList,
		Meta: rbac.RBACMetadata{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: int(total),
		},
	}, nil
}

func (s *service) GetMenuTree(ctx context.Context) (*rbac.MenuTreeResponse, error) {
	menus, err := s.repo.GetMenuTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu tree: %w", err)
	}

	var menuList []rbac.Menu
	for _, menu := range menus {
		menuList = append(menuList, menu.ToMenu())
	}

	return &rbac.MenuTreeResponse{
		Data: menuList,
	}, nil
}

func (s *service) UpdateMenu(ctx context.Context, id uuid.UUID, req *rbac.UpdateMenuRequest, updatedBy uuid.UUID) error {
	menu, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get menu: %w", err)
	}
	if menu == nil {
		return errors.New("menu not found")
	}

	// Update fields if provided
	if req.Name != nil {
		menu.Name = *req.Name
	}
	if req.Slug != nil {
		// Validate slug uniqueness
		if err := s.ValidateMenuSlug(ctx, *req.Slug, &id); err != nil {
			return err
		}
		menu.Slug = *req.Slug
	}
	if req.URL != nil {
		menu.URL = *req.URL
	}
	if req.Icon != nil {
		menu.Icon = *req.Icon
	}
	if req.ParentID != nil {
		// Validate parent menu exists and prevent circular reference
		if *req.ParentID == id {
			return errors.New("menu cannot be parent of itself")
		}
		if req.ParentID != nil {
			parent, err := s.repo.GetMenuByID(ctx, *req.ParentID)
			if err != nil {
				return fmt.Errorf("failed to get parent menu: %w", err)
			}
			if parent == nil {
				return errors.New("parent menu not found")
			}
		}
		menu.ParentID = req.ParentID
	}
	if req.SortOrder != nil {
		menu.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		menu.IsActive = *req.IsActive
	}

	menu.UpdatedBy = &updatedBy
	menu.UpdatedAt = time.Now()

	if err := s.repo.UpdateMenu(ctx, menu); err != nil {
		return fmt.Errorf("failed to update menu: %w", err)
	}

	return nil
}

func (s *service) DeleteMenu(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	menu, err := s.repo.GetMenuByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get menu: %w", err)
	}
	if menu == nil {
		return errors.New("menu not found")
	}

	menu.DeletedBy = &deletedBy
	now := time.Now()
	menu.DeletedAt = &now

	if err := s.repo.UpdateMenu(ctx, menu); err != nil {
		return fmt.Errorf("failed to delete menu: %w", err)
	}

	return nil
}

// User-Role services
func (s *service) AssignRolesToUser(ctx context.Context, userID uuid.UUID, req *rbac.AssignUserRolesRequest, assignedBy uuid.UUID) (*rbac.UserRoleResponse, error) {
	// Validate roles exist
	for _, roleID := range req.RoleIDs {
		role, err := s.repo.GetRoleByID(ctx, roleID)
		if err != nil {
			return nil, fmt.Errorf("failed to get role: %w", err)
		}
		if role == nil {
			return nil, fmt.Errorf("role with ID %s not found", roleID)
		}
	}

	if err := s.repo.AssignRolesToUser(ctx, userID, req.RoleIDs, assignedBy); err != nil {
		return nil, fmt.Errorf("failed to assign roles to user: %w", err)
	}

	return &rbac.UserRoleResponse{
		ID:      uuid.New(),
		Message: "Roles assigned to user successfully",
	}, nil
}

func (s *service) GetUserRoles(ctx context.Context, userID uuid.UUID) (*rbac.UserRoleListResponse, error) {
	userRoles, err := s.repo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	var roleList []rbac.UserRole
	for _, userRole := range userRoles {
		roleList = append(roleList, rbac.UserRole{
			ID:         userRole.ID,
			UserID:     userRole.UserID,
			RoleID:     userRole.RoleID,
			Role:       userRole.Role.ToRole(),
			AssignedAt: userRole.AssignedAt,
		})
	}

	return &rbac.UserRoleListResponse{
		Data: roleList,
		Meta: rbac.RBACMetadata{
			Page:       1,
			Limit:      len(roleList),
			TotalPages: 1,
			TotalItems: len(roleList),
		},
	}, nil
}

func (s *service) RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if err := s.repo.RemoveRolesFromUser(ctx, userID, roleIDs); err != nil {
		return fmt.Errorf("failed to remove roles from user: %w", err)
	}

	return nil
}

// Authorization services
func (s *service) CheckUserPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	return s.repo.CheckUserHasPermission(ctx, userID, resource, action)
}

func (s *service) CheckUserRole(ctx context.Context, userID uuid.UUID, roleSlug string) (bool, error) {
	return s.repo.CheckUserHasRole(ctx, userID, roleSlug)
}

func (s *service) GetUserPermissions(ctx context.Context, userID uuid.UUID) (*rbac.PermissionListResponse, error) {
	permissions, err := s.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	var permissionList []rbac.Permission
	for _, permission := range permissions {
		permissionList = append(permissionList, permission.ToPermission())
	}

	return &rbac.PermissionListResponse{
		Data: permissionList,
		Meta: rbac.RBACMetadata{
			Page:       1,
			Limit:      len(permissionList),
			TotalPages: 1,
			TotalItems: len(permissionList),
		},
	}, nil
}

func (s *service) GetUserMenus(ctx context.Context, userID uuid.UUID) (*rbac.UserMenuResponse, error) {
	roleMenus, err := s.repo.GetUserMenus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user menus: %w", err)
	}

	var menuList []rbac.RoleMenu
	for _, roleMenu := range roleMenus {
		menuList = append(menuList, roleMenu.ToRoleMenu())
	}

	return &rbac.UserMenuResponse{
		Data: menuList,
	}, nil
}

func (s *service) GetUserAccessibleMenus(ctx context.Context, userID uuid.UUID) (*rbac.UserMenuResponse, error) {
	roleMenus, err := s.repo.GetUserAccessibleMenus(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accessible menus: %w", err)
	}

	var menuList []rbac.RoleMenu
	for _, roleMenu := range roleMenus {
		menuList = append(menuList, roleMenu.ToRoleMenu())
	}

	return &rbac.UserMenuResponse{
		Data: menuList,
	}, nil
}

// Validation services
func (s *service) ValidateRoleSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error {
	existingRole, err := s.repo.GetRoleBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to check role slug: %w", err)
	}
	if existingRole != nil && (excludeID == nil || existingRole.ID != *excludeID) {
		return errors.New("role slug already exists")
	}
	return nil
}

func (s *service) ValidatePermissionSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error {
	existingPermission, err := s.repo.GetPermissionBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to check permission slug: %w", err)
	}
	if existingPermission != nil && (excludeID == nil || existingPermission.ID != *excludeID) {
		return errors.New("permission slug already exists")
	}
	return nil
}

func (s *service) ValidateMenuSlug(ctx context.Context, slug string, excludeID *uuid.UUID) error {
	existingMenu, err := s.repo.GetMenuBySlug(ctx, slug)
	if err != nil {
		return fmt.Errorf("failed to check menu slug: %w", err)
	}
	if existingMenu != nil && (excludeID == nil || existingMenu.ID != *excludeID) {
		return errors.New("menu slug already exists")
	}
	return nil
}
