package migration

import (
	"log"
	"os"

	"backend-service-internpro/internal/auth"
	"backend-service-internpro/internal/rbac"
	"backend-service-internpro/internal/school"
	"backend-service-internpro/internal/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Migrate user tables first (independent tables)
	if err := db.AutoMigrate(&user.UserEntity{}); err != nil {
		return err
	}

	// Migrate school related tables
	if err := db.AutoMigrate(&school.SchoolEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&school.MajorityEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&school.ClassEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&school.PartnerEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&school.OTPEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&school.RefreshTokenEntity{}); err != nil {
		return err
	}

	// Migrate auth tables
	if err := db.AutoMigrate(&auth.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&auth.OTP{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&auth.RefreshToken{}); err != nil {
		return err
	}

	// Migrate RBAC tables
	if err := db.AutoMigrate(&rbac.RoleEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&rbac.PermissionEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&rbac.MenuEntity{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&rbac.UserRoleEntity{}); err != nil {
		return err
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// AutoMigrateIfDevelopment runs migration only in development environment
func AutoMigrateIfDevelopment(db *gorm.DB) error {
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "development" {
		log.Println("🔧 Development environment detected, running auto-migration...")
		if err := AutoMigrate(db); err != nil {
			return err
		}

		// Create initial RBAC data if needed
		if err := CreateInitialRBACData(db); err != nil {
			log.Printf("⚠️  Warning: Failed to create initial RBAC data: %v", err)
		}

		// Create dummy data for testing
		if err := CreateDummyData(db); err != nil {
			log.Printf("⚠️  Warning: Failed to create dummy data: %v", err)
		}
	} else {
		log.Printf("Production environment (%s) detected, skipping auto-migration", appEnv)
	}
	return nil
}

// CreateDummyData creates sample data for testing
func CreateDummyData(db *gorm.DB) error {
	log.Println("Creating dummy data...")

	// Check if user already exists
	var userCount int64
	db.Model(&user.UserEntity{}).Count(&userCount)

	if userCount == 0 {
		// Create a dummy user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		dummyUser := user.UserEntity{
			ID:           uuid.New(),
			Username:     "testuser",
			Email:        "test@example.com",
			Fullname:     "Test User",
			PasswordHash: string(hashedPassword),
			IsAdmin:      true, // Set as admin
		}

		if err := db.Create(&dummyUser).Error; err != nil {
			return err
		}

		// Assign super-admin role to the dummy user
		if err := assignSuperAdminRole(db, dummyUser.ID); err != nil {
			log.Printf("⚠️  Warning: Failed to assign super-admin role: %v", err)
		}

		log.Println("✅ Dummy data created successfully")
	} else {
		log.Println("✅ Dummy data already exists")
	}

	return nil
}

// CreateInitialRBACData creates initial RBAC data for testing
func CreateInitialRBACData(db *gorm.DB) error {
	log.Println("Creating initial RBAC data...")

	// Check if roles already exist
	var roleCount int64
	db.Model(&rbac.RoleEntity{}).Count(&roleCount)

	if roleCount == 0 {
		// Create default roles
		adminRole := rbac.RoleEntity{
			ID:          uuid.New(),
			Name:        "Administrator",
			Slug:        "admin",
			Description: "Full system access",
			IsActive:    true,
		}

		userRole := rbac.RoleEntity{
			ID:          uuid.New(),
			Name:        "User",
			Slug:        "user",
			Description: "Basic user access",
			IsActive:    true,
		}

		if err := db.Create(&adminRole).Error; err != nil {
			return err
		}

		if err := db.Create(&userRole).Error; err != nil {
			return err
		}

		// Create default permissions
		permissions := []rbac.PermissionEntity{
			{ID: uuid.New(), Name: "Create User", Slug: "create-user", Resource: "users", Action: "create", Description: "Create new users", IsActive: true},
			{ID: uuid.New(), Name: "Read User", Slug: "read-user", Resource: "users", Action: "read", Description: "View user details", IsActive: true},
			{ID: uuid.New(), Name: "Update User", Slug: "update-user", Resource: "users", Action: "update", Description: "Update user information", IsActive: true},
			{ID: uuid.New(), Name: "Delete User", Slug: "delete-user", Resource: "users", Action: "delete", Description: "Delete users", IsActive: true},
			{ID: uuid.New(), Name: "Manage Roles", Slug: "manage-roles", Resource: "roles", Action: "manage", Description: "Manage user roles", IsActive: true},
			{ID: uuid.New(), Name: "Manage Permissions", Slug: "manage-permissions", Resource: "permissions", Action: "manage", Description: "Manage permissions", IsActive: true},
		}

		for _, permission := range permissions {
			if err := db.Create(&permission).Error; err != nil {
				return err
			}
		}

		// Create default menus
		menus := []rbac.MenuEntity{
			{ID: uuid.New(), Name: "Dashboard", Slug: "dashboard", Icon: "dashboard", URL: "/dashboard", SortOrder: 1, IsActive: true},
			{ID: uuid.New(), Name: "Users", Slug: "users", Icon: "users", URL: "/users", SortOrder: 2, IsActive: true},
			{ID: uuid.New(), Name: "Roles", Slug: "roles", Icon: "security", URL: "/roles", SortOrder: 3, IsActive: true},
			{ID: uuid.New(), Name: "Permissions", Slug: "permissions", Icon: "key", URL: "/permissions", SortOrder: 4, IsActive: true},
		}

		for _, menu := range menus {
			if err := db.Create(&menu).Error; err != nil {
				return err
			}
		}

		log.Println("✅ Initial RBAC data created successfully")
	} else {
		log.Println("✅ Initial RBAC data already exists")
	}

	return nil
}

// assignSuperAdminRole assigns super-admin role to a user
func assignSuperAdminRole(db *gorm.DB, userID uuid.UUID) error {
	// Find super-admin role
	var superAdminRole rbac.RoleEntity
	if err := db.Where("slug = ?", "super-admin").First(&superAdminRole).Error; err != nil {
		// If super-admin role doesn't exist, create it
		superAdminRole = rbac.RoleEntity{
			ID:          uuid.New(),
			Name:        "Super Admin",
			Slug:        "super-admin",
			Description: "Super administrator with full access to all system features",
			IsActive:    true,
		}
		if err := db.Create(&superAdminRole).Error; err != nil {
			return err
		}
	}

	// Check if user-role assignment already exists
	var userRole rbac.UserRoleEntity
	err := db.Where("user_id = ? AND role_id = ?", userID, superAdminRole.ID).First(&userRole).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		// Create user-role assignment
		userRole = rbac.UserRoleEntity{
			ID:     uuid.New(),
			UserID: userID,
			RoleID: superAdminRole.ID,
		}
		if err := db.Create(&userRole).Error; err != nil {
			return err
		}
		log.Printf("✅ Assigned super-admin role to user %s", userID)
	}

	return nil
}
