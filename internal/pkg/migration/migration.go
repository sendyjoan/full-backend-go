package migration

import (
	"log"

	"backend-service-internpro/internal/auth"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

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

	log.Println("✅ Database migrations completed successfully")
	return nil
}

// CreateDummyData creates sample data for testing
func CreateDummyData(db *gorm.DB) error {
	log.Println("Creating dummy data...")

	// Check if user already exists
	var userCount int64
	db.Model(&auth.User{}).Count(&userCount)

	if userCount == 0 {
		// Create a dummy user
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := auth.User{
			Username:     "testuser",
			Email:        "test@example.com",
			Fullname:     "Test User",
			PasswordHash: string(hashedPassword),
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		log.Println("✅ Dummy data created successfully")
	} else {
		log.Println("✅ Dummy data already exists")
	}

	return nil
}
