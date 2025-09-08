package user

import (
	"backend-service-internpro/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID         uuid.UUID  `json:"id" doc:"User ID"`
	Username   string     `json:"username" doc:"User username"`
	Email      string     `json:"email" doc:"User email address"`
	Fullname   string     `json:"fullname" doc:"User full name"`
	IsAdmin    bool       `json:"is_admin" doc:"Whether user is admin"`
	SchoolID   *uuid.UUID `json:"school_id,omitempty" doc:"User school ID"`
	MajorityID *uuid.UUID `json:"majority_id,omitempty" doc:"User majority ID"`
	ClassID    *uuid.UUID `json:"class_id,omitempty" doc:"User class ID"`
	PartnerID  *uuid.UUID `json:"partner_id,omitempty" doc:"User partner ID"`
	CreatedAt  time.Time  `json:"created_at" doc:"User creation date"`
	UpdatedAt  time.Time  `json:"updated_at" doc:"User last update date"`
}

// Metadata represents pagination metadata
type Metadata struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// UserListData represents the data structure for user list
type UserListData struct {
	Users []User   `json:"users"`
	Meta  Metadata `json:"meta"`
}

// UserListResponse represents the response for user list
type UserListResponse = response.ApiResponse

// UserResponse represents a single user response
type UserResponse = response.ApiResponse

// CreateUserRequest represents request to create a new user
type CreateUserRequest struct {
	Username   string     `json:"username" form:"username" minLength:"1" maxLength:"60" doc:"User username"`
	Email      string     `json:"email" form:"email" format:"email" maxLength:"120" doc:"User email address"`
	Fullname   string     `json:"fullname" form:"fullname" minLength:"1" maxLength:"120" doc:"User full name"`
	Password   string     `json:"password" form:"password" minLength:"8" doc:"User password"`
	IsAdmin    bool       `json:"is_admin,omitempty" doc:"Whether user is admin"`
	SchoolID   *uuid.UUID `json:"school_id,omitempty" doc:"User school ID"`
	MajorityID *uuid.UUID `json:"majority_id,omitempty" doc:"User majority ID"`
	ClassID    *uuid.UUID `json:"class_id,omitempty" doc:"User class ID"`
	PartnerID  *uuid.UUID `json:"partner_id,omitempty" doc:"User partner ID"`
}

// CreateUserResponse represents response after creating a user
type CreateUserResponse = response.ApiResponse

// UpdateUserRequest represents request to update user
type UpdateUserRequest struct {
	Username string `json:"username" form:"username" minLength:"1" maxLength:"60" doc:"User username"`
	Fullname string `json:"fullname" form:"fullname" minLength:"1" maxLength:"120" doc:"User full name"`
}

// UserBasicResponse represents a basic response with message for user operations
type UserBasicResponse = response.ApiResponse

// CreateUserData represents the data structure for user creation response
type CreateUserData struct {
	ID uuid.UUID `json:"id" doc:"Created user ID"`
}
