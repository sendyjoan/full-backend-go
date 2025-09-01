package school

import (
	"time"

	"github.com/google/uuid"
)

// SchoolEntity represents the school entity for database operations
type SchoolEntity struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey"`
	Name      string     `gorm:"size:255;not null"`
	Address   *string    `gorm:"size:255"`
	Domain    *string    `gorm:"size:255;uniqueIndex"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt *time.Time `gorm:"index"`
	DeletedBy *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	Majorities []MajorityEntity `gorm:"foreignKey:SchoolID"`
	Classes    []ClassEntity    `gorm:"foreignKey:SchoolID"`
	Partners   []PartnerEntity  `gorm:"foreignKey:SchoolID"`
}

// TableName returns the table name for the SchoolEntity
func (SchoolEntity) TableName() string {
	return "schools"
}

// ToSchool converts SchoolEntity to School DTO
func (s *SchoolEntity) ToSchool() School {
	school := School{
		ID:        s.ID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}

	if s.Address != nil {
		school.Address = *s.Address
	}

	if s.Domain != nil {
		school.Domain = *s.Domain
	}

	return school
}

// MajorityEntity represents the majority entity for database operations
type MajorityEntity struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey"`
	SchoolID    uuid.UUID  `gorm:"type:char(36);not null;index"`
	Name        string     `gorm:"size:255;not null"`
	Description *string    `gorm:"type:text"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt   *time.Time `gorm:"index"`
	DeletedBy   *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	School  SchoolEntity  `gorm:"foreignKey:SchoolID;references:ID"`
	Classes []ClassEntity `gorm:"foreignKey:MajorityID"`
}

// TableName returns the table name for the MajorityEntity
func (MajorityEntity) TableName() string {
	return "majorities"
}

// ToMajority converts MajorityEntity to Majority DTO
func (m *MajorityEntity) ToMajority() Majority {
	majority := Majority{
		ID:        m.ID,
		SchoolID:  m.SchoolID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.Description != nil {
		majority.Description = *m.Description
	}

	return majority
}

// ClassEntity represents the class entity for database operations
type ClassEntity struct {
	ID          uuid.UUID  `gorm:"type:char(36);primaryKey"`
	SchoolID    uuid.UUID  `gorm:"type:char(36);not null;index"`
	MajorityID  uuid.UUID  `gorm:"type:char(36);not null;index"`
	Name        string     `gorm:"size:255;not null"`
	Description *string    `gorm:"type:text"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy   *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy   *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt   *time.Time `gorm:"index"`
	DeletedBy   *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	School   SchoolEntity   `gorm:"foreignKey:SchoolID;references:ID"`
	Majority MajorityEntity `gorm:"foreignKey:MajorityID;references:ID"`
}

// TableName returns the table name for the ClassEntity
func (ClassEntity) TableName() string {
	return "classes"
}

// ToClass converts ClassEntity to Class DTO
func (c *ClassEntity) ToClass() Class {
	class := Class{
		ID:         c.ID,
		SchoolID:   c.SchoolID,
		MajorityID: c.MajorityID,
		Name:       c.Name,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}

	if c.Description != nil {
		class.Description = *c.Description
	}

	return class
}

// PartnerEntity represents the partner entity for database operations
type PartnerEntity struct {
	ID            uuid.UUID  `gorm:"type:char(36);primaryKey"`
	SchoolID      uuid.UUID  `gorm:"type:char(36);not null;index"`
	Name          string     `gorm:"size:255;not null;index"`
	Website       *string    `gorm:"size:255"`
	Description   *string    `gorm:"type:text"`
	Address       *string    `gorm:"size:255"`
	ContactName   *string    `gorm:"size:255"`
	ContactPerson *string    `gorm:"size:255"`
	ContactEmail  *string    `gorm:"size:255"`
	CreatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy     *uuid.UUID `gorm:"type:char(36)"`
	UpdatedAt     time.Time  `gorm:"default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
	UpdatedBy     *uuid.UUID `gorm:"type:char(36)"`
	DeletedAt     *time.Time `gorm:"index"`
	DeletedBy     *uuid.UUID `gorm:"type:char(36)"`

	// Relationships
	School SchoolEntity `gorm:"foreignKey:SchoolID;references:ID"`
}

// TableName returns the table name for the PartnerEntity
func (PartnerEntity) TableName() string {
	return "partners"
}

// ToPartner converts PartnerEntity to Partner DTO
func (p *PartnerEntity) ToPartner() Partner {
	partner := Partner{
		ID:        p.ID,
		SchoolID:  p.SchoolID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}

	if p.Website != nil {
		partner.Website = *p.Website
	}

	if p.Description != nil {
		partner.Description = *p.Description
	}

	if p.Address != nil {
		partner.Address = *p.Address
	}

	if p.ContactName != nil {
		partner.ContactName = *p.ContactName
	}

	if p.ContactPerson != nil {
		partner.ContactPerson = *p.ContactPerson
	}

	if p.ContactEmail != nil {
		partner.ContactEmail = *p.ContactEmail
	}

	return partner
}

// OTPEntity represents the OTP entity for database operations
type OTPEntity struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index"`
	Code      string    `gorm:"size:6;not null;index"`
	Purpose   string    `gorm:"size:32;not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Used      bool      `gorm:"default:false;index"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the OTPEntity
func (OTPEntity) TableName() string {
	return "otps"
}

// RefreshTokenEntity represents the refresh token entity for database operations
type RefreshTokenEntity struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID    uuid.UUID `gorm:"type:char(36);not null;index"`
	TokenHash string    `gorm:"size:255;not null;index"`
	UserAgent *string   `gorm:"size:255"`
	IP        *string   `gorm:"size:64"`
	Revoked   bool      `gorm:"default:false;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the RefreshTokenEntity
func (RefreshTokenEntity) TableName() string {
	return "refresh_tokens"
}
