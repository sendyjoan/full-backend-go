package repository

import (
	"context"

	"backend-service-internpro/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *user.UserEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*user.UserEntity, error)
	GetByEmail(ctx context.Context, email string) (*user.UserEntity, error)
	GetByUsername(ctx context.Context, username string) (*user.UserEntity, error)
	Update(ctx context.Context, user *user.UserEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, offset, limit int) ([]user.UserEntity, int64, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(ctx context.Context, user *user.UserEntity) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*user.UserEntity, error) {
	var user user.UserEntity
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*user.UserEntity, error) {
	var user user.UserEntity
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*user.UserEntity, error) {
	var user user.UserEntity
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *user.UserEntity) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&user.UserEntity{}, "id = ?", id).Error
}

func (r *repository) List(ctx context.Context, offset, limit int) ([]user.UserEntity, int64, error) {
	var users []user.UserEntity
	var total int64

	// Count total records
	if err := r.db.WithContext(ctx).Model(&user.UserEntity{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated records
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
