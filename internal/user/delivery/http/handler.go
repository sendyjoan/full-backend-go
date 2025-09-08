package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/pkg/constants"
	"backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/middleware"
	"backend-service-internpro/internal/pkg/response"
	"backend-service-internpro/internal/user"
	"backend-service-internpro/internal/user/service"

	"github.com/danielgtaylor/huma/v2"
)

type Handler struct {
	svc        service.Service
	jwtSecrets jwt.Secrets
}

// New registers user management routes into the Huma API.
func New(api huma.API, svc service.Service, jwtSecrets jwt.Secrets) {
	h := &Handler{
		svc:        svc,
		jwtSecrets: jwtSecrets,
	}

	// Group /v1/users
	g := huma.NewGroup(api, "/v1/users")

	// GET /users - List all users
	huma.Register(g, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of users with pagination",
		Tags:    []string{"User Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
	}) (*struct {
		Body user.UserListResponse
	}, error) {
		// Validate token
		if err := h.validateToken(in.Authorization); err != nil {
			return &struct {
				Body user.UserListResponse
			}{
				Body: *response.Error(constants.TokenInvalid),
			}, nil
		}

		resp, err := h.svc.ListUsers(ctx, in.Page, in.Limit)
		if err != nil {
			return &struct {
				Body user.UserListResponse
			}{
				Body: *response.Error(constants.InternalServerError),
			}, nil
		}

		return &struct {
			Body user.UserListResponse
		}{
			Body: *resp,
		}, nil
	})

	// GET /users/{id} - Get user by ID
	huma.Register(g, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}",
		Summary: "Get user details by ID",
		Tags:    []string{"User Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string `path:"id" format:"uuid" doc:"User ID"`
	}) (*struct {
		Body user.UserResponse
	}, error) {
		// Validate token
		if err := h.validateToken(in.Authorization); err != nil {
			return &struct {
				Body user.UserResponse
			}{
				Body: *response.Error(constants.TokenInvalid),
			}, nil
		}

		resp, err := h.svc.GetUserByID(ctx, in.ID)
		if err != nil {
			return &struct {
				Body user.UserResponse
			}{
				Body: *response.Error(constants.UserNotFound),
			}, nil
		}

		return &struct {
			Body user.UserResponse
		}{
			Body: *resp,
		}, nil
	})

	// POST /users - Create new user
	huma.Register(g, huma.Operation{
		Method:  http.MethodPost,
		Path:    "",
		Summary: "Create a new user",
		Tags:    []string{"User Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Body          user.CreateUserRequest
	}) (*struct {
		Body user.CreateUserResponse
	}, error) {
		// Validate token
		if err := h.validateToken(in.Authorization); err != nil {
			return &struct {
				Body user.CreateUserResponse
			}{
				Body: *response.Error(constants.TokenInvalid),
			}, nil
		}

		resp, err := h.svc.CreateUser(ctx, in.Body)
		if err != nil {
			return &struct {
				Body user.CreateUserResponse
			}{
				Body: *response.Error(constants.UserCreateFailed),
			}, nil
		}

		return &struct {
			Body user.CreateUserResponse
		}{
			Body: *resp,
		}, nil
	})

	// PUT /users/{id} - Update user
	huma.Register(g, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/{id}",
		Summary: "Update user information",
		Tags:    []string{"User Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string `path:"id" format:"uuid" doc:"User ID"`
		Body          user.UpdateUserRequest
	}) (*struct {
		Body user.UserBasicResponse
	}, error) {
		// Validate token
		if err := h.validateToken(in.Authorization); err != nil {
			return &struct {
				Body user.UserBasicResponse
			}{
				Body: *response.Error(constants.TokenInvalid),
			}, nil
		}

		resp, err := h.svc.UpdateUser(ctx, in.ID, in.Body)
		if err != nil {
			return &struct {
				Body user.UserBasicResponse
			}{
				Body: *response.Error(constants.UserUpdateFailed),
			}, nil
		}

		return &struct {
			Body user.UserBasicResponse
		}{
			Body: *resp,
		}, nil
	})

	// DELETE /users/{id} - Delete user
	huma.Register(g, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete user by ID",
		Tags:    []string{"User Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string `path:"id" format:"uuid" doc:"User ID"`
	}) (*struct {
		Body user.UserBasicResponse
	}, error) {
		// Validate token
		if err := h.validateToken(in.Authorization); err != nil {
			return &struct {
				Body user.UserBasicResponse
			}{
				Body: *response.Error(constants.TokenInvalid),
			}, nil
		}

		resp, err := h.svc.DeleteUser(ctx, in.ID)
		if err != nil {
			return &struct {
				Body user.UserBasicResponse
			}{
				Body: *response.Error(constants.UserDeleteFailed),
			}, nil
		}

		return &struct {
			Body user.UserBasicResponse
		}{
			Body: *resp,
		}, nil
	})
}

// validateToken validates the Authorization header and returns error if invalid
func (h *Handler) validateToken(authHeader string) error {
	_, err := middleware.ValidateToken(authHeader, h.jwtSecrets)
	return err
}
