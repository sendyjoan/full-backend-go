package user

import (
	"time"

	"github.com/google/uuid"
)

// UserEntity represents the user entity for database operations
type UserEntity struct {
	ID           uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Username     string     `gorm:"uniqueIndex;size:60;not null"`
	Email        string     `gorm:"uniqueIndex;size:120;not null"`
	Fullname     string     `gorm:"size:120;not null"`
	PasswordHash string     `gorm:"size:255;not null"`
	IsAdmin      bool       `gorm:"default:false"`
	SchoolID     *uuid.UUID `gorm:"type:char(36);index"`
	MajorityID   *uuid.UUID `gorm:"type:char(36);index"`
	ClassID      *uuid.UUID `gorm:"type:char(36);index"`
	PartnerID    *uuid.UUID `gorm:"type:char(36);index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// TableName returns the table name for the UserEntity
func (UserEntity) TableName() string {
	return "users"
}

// ToUser converts UserEntity to User DTO
func (u *UserEntity) ToUser() User {
	user := User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Fullname:  u.Fullname,
		IsAdmin:   u.IsAdmin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.SchoolID != nil {
		user.SchoolID = u.SchoolID
	}

	if u.MajorityID != nil {
		user.MajorityID = u.MajorityID
	}

	if u.ClassID != nil {
		user.ClassID = u.ClassID
	}

	if u.PartnerID != nil {
		user.PartnerID = u.PartnerID
	}

	return user
}
