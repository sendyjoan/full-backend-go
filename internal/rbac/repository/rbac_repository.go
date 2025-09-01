package repository

import (
	"context"
	"errors"
	"strings"

	"backend-service-internpro/internal/rbac"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates a new RBAC repository
func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

// Role methods
func (r *repository) CreateRole(ctx context.Context, role *rbac.RoleEntity) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *repository) GetRoleByID(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error) {
	var role rbac.RoleEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleBySlug(ctx context.Context, slug string) (*rbac.RoleEntity, error) {
	var role rbac.RoleEntity
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoles(ctx context.Context, page, limit int, search string) ([]rbac.RoleEntity, int64, error) {
	var roles []rbac.RoleEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&rbac.RoleEntity{}).Where("deleted_at IS NULL")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&roles).Error

	return roles, total, err
}

func (r *repository) UpdateRole(ctx context.Context, role *rbac.RoleEntity) error {
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *repository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&rbac.RoleEntity{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error) {
	var role rbac.RoleEntity
	err := r.db.WithContext(ctx).
		Preload("Permissions", "deleted_at IS NULL AND is_active = ?", true).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleWithMenus(ctx context.Context, id uuid.UUID) (*rbac.RoleEntity, error) {
	var role rbac.RoleEntity
	err := r.db.WithContext(ctx).
		Preload("Menus", "deleted_at IS NULL AND is_active = ?", true).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

// Permission methods
func (r *repository) CreatePermission(ctx context.Context, permission *rbac.PermissionEntity) error {
	return r.db.WithContext(ctx).Create(permission).Error
}

func (r *repository) GetPermissionByID(ctx context.Context, id uuid.UUID) (*rbac.PermissionEntity, error) {
	var permission rbac.PermissionEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *repository) GetPermissionBySlug(ctx context.Context, slug string) (*rbac.PermissionEntity, error) {
	var permission rbac.PermissionEntity
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &permission, nil
}

func (r *repository) GetPermissions(ctx context.Context, page, limit int, search string) ([]rbac.PermissionEntity, int64, error) {
	var permissions []rbac.PermissionEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&rbac.PermissionEntity{}).Where("deleted_at IS NULL")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(resource) LIKE ? OR LOWER(action) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("resource ASC, action ASC").Offset(offset).Limit(limit).Find(&permissions).Error

	return permissions, total, err
}

func (r *repository) UpdatePermission(ctx context.Context, permission *rbac.PermissionEntity) error {
	return r.db.WithContext(ctx).Save(permission).Error
}

func (r *repository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&rbac.PermissionEntity{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) GetPermissionsByResource(ctx context.Context, resource string) ([]rbac.PermissionEntity, error) {
	var permissions []rbac.PermissionEntity
	err := r.db.WithContext(ctx).
		Where("resource = ? AND deleted_at IS NULL AND is_active = ?", resource, true).
		Order("action ASC").
		Find(&permissions).Error
	return permissions, err
}

func (r *repository) GetPermissionsByIDs(ctx context.Context, ids []uuid.UUID) ([]rbac.PermissionEntity, error) {
	var permissions []rbac.PermissionEntity
	err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Find(&permissions).Error
	return permissions, err
}

// Menu methods
func (r *repository) CreateMenu(ctx context.Context, menu *rbac.MenuEntity) error {
	return r.db.WithContext(ctx).Create(menu).Error
}

func (r *repository) GetMenuByID(ctx context.Context, id uuid.UUID) (*rbac.MenuEntity, error) {
	var menu rbac.MenuEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *repository) GetMenuBySlug(ctx context.Context, slug string) (*rbac.MenuEntity, error) {
	var menu rbac.MenuEntity
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &menu, nil
}

func (r *repository) GetMenus(ctx context.Context, page, limit int, search string) ([]rbac.MenuEntity, int64, error) {
	var menus []rbac.MenuEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&rbac.MenuEntity{}).Where("deleted_at IS NULL")

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ? OR LOWER(url) LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("sort_order ASC, name ASC").Offset(offset).Limit(limit).Find(&menus).Error

	return menus, total, err
}

func (r *repository) GetMenuTree(ctx context.Context) ([]rbac.MenuEntity, error) {
	var menus []rbac.MenuEntity
	err := r.db.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL AND is_active = ?", true).Order("sort_order ASC")
		}).
		Where("parent_id IS NULL AND deleted_at IS NULL AND is_active = ?", true).
		Order("sort_order ASC").
		Find(&menus).Error
	return menus, err
}

func (r *repository) UpdateMenu(ctx context.Context, menu *rbac.MenuEntity) error {
	return r.db.WithContext(ctx).Save(menu).Error
}

func (r *repository) DeleteMenu(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&rbac.MenuEntity{}).
		Where("id = ?", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *repository) GetMenusByParentID(ctx context.Context, parentID *uuid.UUID) ([]rbac.MenuEntity, error) {
	var menus []rbac.MenuEntity
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL AND is_active = ?", true)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Order("sort_order ASC").Find(&menus).Error
	return menus, err
}

func (r *repository) GetMenusByIDs(ctx context.Context, ids []uuid.UUID) ([]rbac.MenuEntity, error) {
	var menus []rbac.MenuEntity
	err := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Order("sort_order ASC").
		Find(&menus).Error
	return menus, err
}

// Role-Permission methods
func (r *repository) AssignPermissionsToRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID, assignedBy uuid.UUID) error {
	// First, remove existing permissions
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&rbac.RolePermissionEntity{}).Error; err != nil {
		return err
	}

	// Then add new permissions
	var rolePermissions []rbac.RolePermissionEntity
	for _, permissionID := range permissionIDs {
		rolePermissions = append(rolePermissions, rbac.RolePermissionEntity{
			ID:           uuid.New(),
			RoleID:       roleID,
			PermissionID: permissionID,
			CreatedBy:    &assignedBy,
		})
	}

	if len(rolePermissions) > 0 {
		return r.db.WithContext(ctx).Create(&rolePermissions).Error
	}

	return nil
}

func (r *repository) RemovePermissionsFromRole(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id IN ?", roleID, permissionIDs).
		Delete(&rbac.RolePermissionEntity{}).Error
}

func (r *repository) GetRolePermissions(ctx context.Context, roleID uuid.UUID) ([]rbac.PermissionEntity, error) {
	var permissions []rbac.PermissionEntity
	err := r.db.WithContext(ctx).
		Table("permissions").
		Select("permissions.*").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.deleted_at IS NULL AND permissions.is_active = ?", roleID, true).
		Find(&permissions).Error
	return permissions, err
}

func (r *repository) CheckRoleHasPermission(ctx context.Context, roleID uuid.UUID, permissionSlug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("role_permissions").
		Joins("INNER JOIN permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.slug = ? AND permissions.deleted_at IS NULL AND permissions.is_active = ?",
			roleID, permissionSlug, true).
		Count(&count).Error
	return count > 0, err
}

// User-Role methods
func (r *repository) AssignRolesToUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error {
	// First, remove existing roles
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&rbac.UserRoleEntity{}).Error; err != nil {
		return err
	}

	// Then add new roles
	var userRoles []rbac.UserRoleEntity
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, rbac.UserRoleEntity{
			ID:         uuid.New(),
			UserID:     userID,
			RoleID:     roleID,
			AssignedBy: &assignedBy,
		})
	}

	if len(userRoles) > 0 {
		return r.db.WithContext(ctx).Create(&userRoles).Error
	}

	return nil
}

func (r *repository) RemoveRolesFromUser(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id IN ?", userID, roleIDs).
		Delete(&rbac.UserRoleEntity{}).Error
}

func (r *repository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]rbac.UserRoleEntity, error) {
	var userRoles []rbac.UserRoleEntity
	err := r.db.WithContext(ctx).
		Preload("Role", "deleted_at IS NULL AND is_active = ?", true).
		Where("user_id = ?", userID).
		Find(&userRoles).Error
	return userRoles, err
}

func (r *repository) GetUsersByRole(ctx context.Context, roleID uuid.UUID, page, limit int) ([]rbac.UserRoleEntity, int64, error) {
	var userRoles []rbac.UserRoleEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&rbac.UserRoleEntity{}).Where("role_id = ?", roleID)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	err := query.Order("assigned_at DESC").Offset(offset).Limit(limit).Find(&userRoles).Error

	return userRoles, total, err
}

func (r *repository) CheckUserHasRole(ctx context.Context, userID uuid.UUID, roleSlug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("user_roles").
		Joins("INNER JOIN roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.slug = ? AND roles.deleted_at IS NULL AND roles.is_active = ?",
			userID, roleSlug, true).
		Count(&count).Error
	return count > 0, err
}

// Role-Menu methods
func (r *repository) AssignMenusToRole(ctx context.Context, roleID uuid.UUID, menuPermissions []rbac.RoleMenuEntity, assignedBy uuid.UUID) error {
	// First, remove existing menus
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&rbac.RoleMenuEntity{}).Error; err != nil {
		return err
	}

	// Then add new menus with permissions
	for i := range menuPermissions {
		menuPermissions[i].ID = uuid.New()
		menuPermissions[i].RoleID = roleID
		menuPermissions[i].CreatedBy = &assignedBy
		menuPermissions[i].UpdatedBy = &assignedBy
	}

	if len(menuPermissions) > 0 {
		return r.db.WithContext(ctx).Create(&menuPermissions).Error
	}

	return nil
}

func (r *repository) RemoveMenusFromRole(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND menu_id IN ?", roleID, menuIDs).
		Delete(&rbac.RoleMenuEntity{}).Error
}

func (r *repository) GetRoleMenus(ctx context.Context, roleID uuid.UUID) ([]rbac.RoleMenuEntity, error) {
	var roleMenus []rbac.RoleMenuEntity
	err := r.db.WithContext(ctx).
		Preload("Menu", "deleted_at IS NULL AND is_active = ?", true).
		Where("role_id = ?", roleID).
		Order("created_at ASC").
		Find(&roleMenus).Error
	return roleMenus, err
}

func (r *repository) GetUserMenus(ctx context.Context, userID uuid.UUID) ([]rbac.RoleMenuEntity, error) {
	var roleMenus []rbac.RoleMenuEntity
	err := r.db.WithContext(ctx).
		Select("DISTINCT role_menus.*").
		Table("role_menus").
		Joins("INNER JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Joins("INNER JOIN menus ON role_menus.menu_id = menus.id").
		Preload("Menu", "deleted_at IS NULL AND is_active = ?", true).
		Where("user_roles.user_id = ? AND menus.deleted_at IS NULL AND menus.is_active = ?", userID, true).
		Order("menus.sort_order ASC").
		Find(&roleMenus).Error
	return roleMenus, err
}

func (r *repository) UpdateRoleMenuPermissions(ctx context.Context, roleMenuID uuid.UUID, canView, canCreate, canEdit, canDelete bool, updatedBy uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&rbac.RoleMenuEntity{}).
		Where("id = ?", roleMenuID).
		Updates(map[string]interface{}{
			"can_view":   canView,
			"can_create": canCreate,
			"can_edit":   canEdit,
			"can_delete": canDelete,
			"updated_by": updatedBy,
		}).Error
}

// Complex queries
func (r *repository) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]rbac.PermissionEntity, error) {
	var permissions []rbac.PermissionEntity
	err := r.db.WithContext(ctx).
		Select("DISTINCT permissions.*").
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.deleted_at IS NULL AND permissions.is_active = ?", userID, true).
		Find(&permissions).Error
	return permissions, err
}

func (r *repository) CheckUserHasPermission(ctx context.Context, userID uuid.UUID, resource, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ? AND permissions.resource = ? AND permissions.action = ? AND permissions.deleted_at IS NULL AND permissions.is_active = ?",
			userID, resource, action, true).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) GetUserAccessibleMenus(ctx context.Context, userID uuid.UUID) ([]rbac.RoleMenuEntity, error) {
	var roleMenus []rbac.RoleMenuEntity
	err := r.db.WithContext(ctx).
		Select("role_menus.*").
		Table("role_menus").
		Joins("INNER JOIN user_roles ON role_menus.role_id = user_roles.role_id").
		Joins("INNER JOIN menus ON role_menus.menu_id = menus.id").
		Preload("Menu.Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL AND is_active = ?", true).Order("sort_order ASC")
		}).
		Preload("Menu", "deleted_at IS NULL AND is_active = ?", true).
		Where("user_roles.user_id = ? AND menus.deleted_at IS NULL AND menus.is_active = ? AND role_menus.can_view = ?",
			userID, true, true).
		Order("menus.sort_order ASC").
		Find(&roleMenus).Error
	return roleMenus, err
}
