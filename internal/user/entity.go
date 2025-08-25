package user

import (
	"time"

	"github.com/google/uuid"
)

// UserEntity represents the user entity for database operations
type UserEntity struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	Username     string    `gorm:"uniqueIndex;size:60;not null"`
	Email        string    `gorm:"uniqueIndex;size:120;not null"`
	Fullname     string    `gorm:"size:120;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName returns the table name for the UserEntity
func (UserEntity) TableName() string {
	return "users"
}

// ToUser converts UserEntity to User DTO
func (u *UserEntity) ToUser() User {
	return User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Fullname:  u.Fullname,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
