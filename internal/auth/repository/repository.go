package repository

import (
	"time"

	"backend-service-internpro/internal/auth"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	FindUserByUsernameOrEmail(uore string) (*auth.User, error)
	FindUserByEmail(email string) (*auth.User, error)
	CreateRefreshToken(rt *auth.RefreshToken) error
	GetRefreshToken(hash string) (*auth.RefreshToken, error)
	RevokeRefreshToken(id uuid.UUID) error
	MarkOTPUsed(id uuid.UUID) error
	FindValidOTP(email, code, purpose string, now time.Time) (*auth.OTP, error)
	SaveOTP(o *auth.OTP) error
	UpdateUserPassword(userID uuid.UUID, passwordHash string) error
}

type repo struct{ db *gorm.DB }

func New(db *gorm.DB) Repository { return &repo{db} }

func (r *repo) FindUserByUsernameOrEmail(uore string) (*auth.User, error) {
	var u auth.User
	if err := r.db.
		Where("username = ? OR email = ?", uore, uore).
		First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) FindUserByEmail(email string) (*auth.User, error) {
	var u auth.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repo) CreateRefreshToken(rt *auth.RefreshToken) error { return r.db.Create(rt).Error }

func (r *repo) GetRefreshToken(hash string) (*auth.RefreshToken, error) {
	var rt auth.RefreshToken
	if err := r.db.Where("token_hash = ? AND revoked = 0 AND expires_at > NOW()", hash).
		First(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *repo) RevokeRefreshToken(id uuid.UUID) error {
	return r.db.Model(&auth.RefreshToken{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *repo) MarkOTPUsed(id uuid.UUID) error {
	return r.db.Model(&auth.OTP{}).Where("id = ?", id).Update("used", true).Error
}

func (r *repo) FindValidOTP(email, code, purpose string, now time.Time) (*auth.OTP, error) {
	var u auth.User
	if err := r.db.Select("id").Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	var o auth.OTP
	if err := r.db.Where("user_id = ? AND code = ? AND purpose = ? AND used = 0 AND expires_at > ?",
		u.ID, code, purpose, now).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repo) SaveOTP(o *auth.OTP) error { return r.db.Create(o).Error }

func (r *repo) UpdateUserPassword(userID uuid.UUID, passwordHash string) error {
	return r.db.Model(&auth.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}
