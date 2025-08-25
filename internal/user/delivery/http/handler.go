package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/user"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	// svc service.Service // akan diimplementasikan nanti
}

// New registers user management routes into the Huma API.
func New(api huma.API) {
	_ = &Handler{} // Will be used when service is implemented

	// Group /v1/users
	g := huma.NewGroup(api, "/v1/users")

	// GET /users - List all users
	huma.Register(g, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of users with pagination",
		Tags:    []string{"User Management"},
	}, func(ctx context.Context, in *struct {
		Page  int `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit int `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
	}) (*struct {
		Body user.UserListResponse
	}, error) {
		// TODO: Implement user list logic
		return &struct {
			Body user.UserListResponse
		}{
			Body: user.UserListResponse{
				Data: []user.User{
					{ID: 1, Name: "John Doe", Email: "john@example.com", Role: "teacher"},
					{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Role: "student"},
				},
				Meta: user.Metadata{
					Page:       in.Page,
					Limit:      in.Limit,
					TotalPages: 1,
					TotalItems: 2,
				},
			},
		}, nil
	})

	// GET /users/{id} - Get user by ID
	huma.Register(g, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}",
		Summary: "Get user details by ID",
		Tags:    []string{"User Management"},
	}, func(ctx context.Context, in *struct {
		ID int `path:"id" minimum:"1" doc:"User ID"`
	}) (*struct {
		Body user.UserResponse
	}, error) {
		// TODO: Implement get user by ID logic
		return &struct {
			Body user.UserResponse
		}{
			Body: user.UserResponse{
				User: user.User{
					ID:    in.ID,
					Name:  "John Doe",
					Email: "john@example.com",
					Role:  "teacher",
				},
			},
		}, nil
	})

	// POST /users - Create new user
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "",
		Summary: "Create a new user",
		Tags:    []string{"User Management"},
	}, func(ctx context.Context, in *struct {
		Body user.CreateUserRequest
	}) (*struct {
		Body user.CreateUserResponse
	}, error) {
		// TODO: Implement create user logic
		return &struct {
			Body user.CreateUserResponse
		}{
			Body: user.CreateUserResponse{
				ID:      123,
				Message: "User created successfully",
			},
		}, nil
	})

	// PUT /users/{id} - Update user
	huma.Register(g, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/{id}",
		Summary: "Update user information",
		Tags:    []string{"User Management"},
	}, func(ctx context.Context, in *struct {
		ID   int `path:"id" minimum:"1" doc:"User ID"`
		Body user.UpdateUserRequest
	}) (*struct {
		Body user.UserBasicResponse
	}, error) {
		// TODO: Implement update user logic
		return &struct {
			Body user.UserBasicResponse
		}{
			Body: user.UserBasicResponse{
				Message: "User updated successfully",
			},
		}, nil
	})

	// DELETE /users/{id} - Delete user
	huma.Register(g, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete user by ID",
		Tags:    []string{"User Management"},
	}, func(ctx context.Context, in *struct {
		ID int `path:"id" minimum:"1" doc:"User ID"`
	}) (*struct {
		Body user.UserBasicResponse
	}, error) {
		// TODO: Implement delete user logic
		return &struct {
			Body user.UserBasicResponse
		}{
			Body: user.UserBasicResponse{
				Message: "User deleted successfully",
			},
		}, nil
	})
}
