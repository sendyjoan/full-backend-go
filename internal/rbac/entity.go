package rbac

import (
	"time"

	"github.com/google/uuid"
)

// RoleEntity represents the role entity for database operations
type RoleEntity struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"size:100;not null;uniqueIndex"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	CreatedBy   *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt   time.Time
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt   *time.Time `gorm:"index"`
	DeletedBy   *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Permissions []PermissionEntity `gorm:"many2many:role_permissions;"`
	Menus       []MenuEntity       `gorm:"many2many:role_menus;"`
}

// TableName returns the table name for the RoleEntity
func (RoleEntity) TableName() string {
	return "roles"
}

// ToRole converts RoleEntity to Role DTO
func (r *RoleEntity) ToRole() Role {
	var permissions []Permission
	for _, p := range r.Permissions {
		permissions = append(permissions, p.ToPermission())
	}

	var menus []Menu
	for _, m := range r.Menus {
		menus = append(menus, m.ToMenu())
	}

	return Role{
		ID:          r.ID,
		Name:        r.Name,
		Slug:        r.Slug,
		Description: r.Description,
		IsActive:    r.IsActive,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		Permissions: permissions,
		Menus:       menus,
	}
}

// PermissionEntity represents the permission entity for database operations
type PermissionEntity struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name        string    `gorm:"size:100;not null;uniqueIndex"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex"`
	Resource    string    `gorm:"size:100;not null"`
	Action      string    `gorm:"size:50;not null"`
	Description string    `gorm:"type:text"`
	IsActive    bool      `gorm:"default:true"`
	CreatedAt   time.Time
	CreatedBy   *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt   time.Time
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt   *time.Time `gorm:"index"`
	DeletedBy   *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Roles []RoleEntity `gorm:"many2many:role_permissions;"`
}

// TableName returns the table name for the PermissionEntity
func (PermissionEntity) TableName() string {
	return "permissions"
}

// ToPermission converts PermissionEntity to Permission DTO
func (p *PermissionEntity) ToPermission() Permission {
	return Permission{
		ID:          p.ID,
		Name:        p.Name,
		Slug:        p.Slug,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// MenuEntity represents the menu entity for database operations
type MenuEntity struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Name      string     `gorm:"size:100;not null"`
	Slug      string     `gorm:"size:100;not null;uniqueIndex"`
	URL       string     `gorm:"size:255"`
	Icon      string     `gorm:"size:100"`
	ParentID  *uuid.UUID `gorm:"type:char(36);index"`
	SortOrder int        `gorm:"default:0"`
	IsActive  bool       `gorm:"default:true"`
	CreatedAt time.Time
	CreatedBy *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt time.Time
	UpdatedBy *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt *time.Time `gorm:"index"`
	DeletedBy *uuid.UUID `gorm:"type:char(36)"`

	// Self-referencing relationship
	Parent   *MenuEntity  `gorm:"foreignKey:ParentID"`
	Children []MenuEntity `gorm:"foreignKey:ParentID"`

	// Many-to-many with roles
	Roles []RoleEntity `gorm:"many2many:role_menus;"`
}

// TableName returns the table name for the MenuEntity
func (MenuEntity) TableName() string {
	return "menus"
}

// ToMenu converts MenuEntity to Menu DTO
func (m *MenuEntity) ToMenu() Menu {
	var children []Menu
	for _, child := range m.Children {
		children = append(children, child.ToMenu())
	}

	return Menu{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		URL:       m.URL,
		Icon:      m.Icon,
		ParentID:  m.ParentID,
		SortOrder: m.SortOrder,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Children:  children,
	}
}

// RolePermissionEntity represents the role_permissions junction table
type RolePermissionEntity struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	RoleID       uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:unique_role_permission"`
	PermissionID uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:unique_role_permission"`
	CreatedAt    time.Time
	CreatedBy    *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Role       RoleEntity       `gorm:"foreignKey:RoleID"`
	Permission PermissionEntity `gorm:"foreignKey:PermissionID"`
}

// TableName returns the table name for the RolePermissionEntity
func (RolePermissionEntity) TableName() string {
	return "role_permissions"
}

// UserRoleEntity represents the user_roles junction table
type UserRoleEntity struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey"`
	UserID     uuid.UUID  `gorm:"type:char(36);not null;index;uniqueIndex:unique_user_role"`
	RoleID     uuid.UUID  `gorm:"type:char(36);not null;index;uniqueIndex:unique_user_role"`
	AssignedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	AssignedBy *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Role RoleEntity `gorm:"foreignKey:RoleID"`
}

// TableName returns the table name for the UserRoleEntity
func (UserRoleEntity) TableName() string {
	return "user_roles"
}

// RoleMenuEntity represents the role_menus junction table
type RoleMenuEntity struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	RoleID    uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:unique_role_menu"`
	MenuID    uuid.UUID `gorm:"type:char(36);not null;index;uniqueIndex:unique_role_menu"`
	CanView   bool      `gorm:"default:true"`
	CanCreate bool      `gorm:"default:false"`
	CanEdit   bool      `gorm:"default:false"`
	CanDelete bool      `gorm:"default:false"`
	CreatedAt time.Time
	CreatedBy *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt time.Time
	UpdatedBy *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Role RoleEntity `gorm:"foreignKey:RoleID"`
	Menu MenuEntity `gorm:"foreignKey:MenuID"`
}

// TableName returns the table name for the RoleMenuEntity
func (RoleMenuEntity) TableName() string {
	return "role_menus"
}

// ToRoleMenu converts RoleMenuEntity to RoleMenu DTO
func (rm *RoleMenuEntity) ToRoleMenu() RoleMenu {
	return RoleMenu{
		ID:        rm.ID,
		RoleID:    rm.RoleID,
		MenuID:    rm.MenuID,
		CanView:   rm.CanView,
		CanCreate: rm.CanCreate,
		CanEdit:   rm.CanEdit,
		CanDelete: rm.CanDelete,
		Menu:      rm.Menu.ToMenu(),
		CreatedAt: rm.CreatedAt,
		UpdatedAt: rm.UpdatedAt,
	}
}
