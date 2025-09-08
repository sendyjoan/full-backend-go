package service

import (
	"context"
	"errors"
	"math"
	"time"

	"backend-service-internpro/internal/pkg/constants"
	"backend-service-internpro/internal/pkg/response"
	"backend-service-internpro/internal/user"
	"backend-service-internpro/internal/user/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	CreateUser(ctx context.Context, req user.CreateUserRequest) (*user.CreateUserResponse, error)
	GetUserByID(ctx context.Context, id string) (*user.UserResponse, error)
	UpdateUser(ctx context.Context, id string, req user.UpdateUserRequest) (*user.UserBasicResponse, error)
	DeleteUser(ctx context.Context, id string) (*user.UserBasicResponse, error)
	ListUsers(ctx context.Context, page, limit int) (*user.UserListResponse, error)
}

type service struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, req user.CreateUserRequest) (*user.CreateUserResponse, error) {
	// Check if user with email already exists
	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check if user with username already exists
	if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.New("user with this username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user entity
	userEntity := &user.UserEntity{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		Fullname:     req.Fullname,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := s.repo.Create(ctx, userEntity); err != nil {
		return nil, errors.New("failed to create user")
	}

	createData := user.CreateUserData{
		ID: userEntity.ID,
	}

	return response.Success(constants.UserCreateSuccess, createData), nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*user.UserResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	userEntity, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	return response.Success(constants.UserDetailSuccess, userEntity.ToUser()), nil
}

func (s *service) UpdateUser(ctx context.Context, id string, req user.UpdateUserRequest) (*user.UserBasicResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Get existing user
	userEntity, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	// Check if username is being changed and if it's already taken
	if req.Username != "" && req.Username != userEntity.Username {
		if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
			return nil, errors.New("username already taken")
		}
		userEntity.Username = req.Username
	}

	// Update fields
	if req.Fullname != "" {
		userEntity.Fullname = req.Fullname
	}
	userEntity.UpdatedAt = time.Now()

	// Save changes
	if err := s.repo.Update(ctx, userEntity); err != nil {
		return nil, errors.New("failed to update user")
	}

	return response.SuccessWithoutData(constants.UserUpdateSuccess), nil
}

func (s *service) DeleteUser(ctx context.Context, id string) (*user.UserBasicResponse, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Check if user exists
	_, err = s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to get user")
	}

	// Delete user
	if err := s.repo.Delete(ctx, userID); err != nil {
		return nil, errors.New("failed to delete user")
	}

	return response.SuccessWithoutData(constants.UserDeleteSuccess), nil
}

func (s *service) ListUsers(ctx context.Context, page, limit int) (*user.UserListResponse, error) {
	offset := (page - 1) * limit

	userEntities, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, errors.New("failed to get users")
	}

	// Convert entities to DTOs
	users := make([]user.User, len(userEntities))
	for i, entity := range userEntities {
		users[i] = entity.ToUser()
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	listData := user.UserListData{
		Users: users,
		Meta: user.Metadata{
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
			TotalItems: int(total),
		},
	}

	return response.Success(constants.UserListSuccess, listData), nil
}
