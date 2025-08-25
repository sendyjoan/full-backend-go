package auth

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	Username     string    `gorm:"uniqueIndex;size:60;not null"`
	Email        string    `gorm:"uniqueIndex;size:120;not null"`
	Fullname     string    `gorm:"size:120;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OTP struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);index;not null"`
	Code      string    `gorm:"size:6;not null"`
	Purpose   string    `gorm:"size:32;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);index;not null"`
	TokenHash string    `gorm:"size:255;not null"`
	UserAgent string    `gorm:"size:255"`
	IP        string    `gorm:"size:64"`
	Revoked   bool      `gorm:"default:false"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
