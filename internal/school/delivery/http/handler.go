package http

import (
	"context"
	"net/http"

	"backend-service-internpro/internal/pkg/jwt"
	"backend-service-internpro/internal/pkg/middleware"
	"backend-service-internpro/internal/school"
	"backend-service-internpro/internal/school/service"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type Handler struct {
	svc        service.SchoolService
	jwtSecrets jwt.Secrets
}

// New registers school management routes into the Huma API.
func New(api huma.API, svc service.SchoolService, jwtSecrets jwt.Secrets) {
	h := &Handler{
		svc:        svc,
		jwtSecrets: jwtSecrets,
	}

	// School routes
	schoolGroup := huma.NewGroup(api, "/v1/schools")

	// GET /schools - List all schools
	huma.Register(schoolGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of schools with pagination",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
		Search        string `query:"search" doc:"Search by name, domain, or address"`
	}) (*struct {
		Body school.PaginatedSchoolsResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		params := school.QueryParams{
			Page:   in.Page,
			Limit:  in.Limit,
			Search: in.Search,
		}

		result, err := h.svc.GetAllSchools(ctx, params)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.PaginatedSchoolsResponse
		}{Body: *result}, nil
	})

	// POST /schools - Create school
	huma.Register(schoolGroup, huma.Operation{
		Method:  http.MethodPost,
		Path:    "",
		Summary: "Create a new school",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string                     `header:"Authorization" required:"true" doc:"Bearer token"`
		Body          school.CreateSchoolRequest `json:"body"`
	}) (*struct {
		Body school.School
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.svc.CreateSchool(ctx, in.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.School
		}{Body: *result}, nil
	})

	// GET /schools/{id} - Get school by ID
	huma.Register(schoolGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "/{id}",
		Summary: "Get school by ID",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string `path:"id" doc:"School ID"`
	}) (*struct {
		Body school.School
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid school ID")
		}

		result, err := h.svc.GetSchoolByID(ctx, id)
		if err != nil {
			if err.Error() == "school not found" {
				return nil, huma.Error404NotFound("School not found")
			}
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.School
		}{Body: *result}, nil
	})

	// PUT /schools/{id} - Update school
	huma.Register(schoolGroup, huma.Operation{
		Method:  http.MethodPut,
		Path:    "/{id}",
		Summary: "Update school",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string                     `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string                     `path:"id" doc:"School ID"`
		Body          school.UpdateSchoolRequest `json:"body"`
	}) (*struct {
		Body school.School
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid school ID")
		}

		result, err := h.svc.UpdateSchool(ctx, id, in.Body)
		if err != nil {
			if err.Error() == "school not found" {
				return nil, huma.Error404NotFound("School not found")
			}
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.School
		}{Body: *result}, nil
	})

	// DELETE /schools/{id} - Delete school
	huma.Register(schoolGroup, huma.Operation{
		Method:  http.MethodDelete,
		Path:    "/{id}",
		Summary: "Delete school",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		ID            string `path:"id" doc:"School ID"`
	}) (*struct {
		Body map[string]string
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("Invalid school ID")
		}

		err = h.svc.DeleteSchool(ctx, id)
		if err != nil {
			if err.Error() == "school not found" {
				return nil, huma.Error404NotFound("School not found")
			}
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body map[string]string
		}{Body: map[string]string{"message": "School deleted successfully"}}, nil
	})

	// Majority routes
	majorityGroup := huma.NewGroup(api, "/v1/majorities")

	// GET /majorities - List all majorities
	huma.Register(majorityGroup, huma.Operation{
		Method:  http.MethodGet,
		Path:    "",
		Summary: "Get list of majorities with pagination",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string `header:"Authorization" required:"true" doc:"Bearer token"`
		Page          int    `query:"page" minimum:"1" default:"1" doc:"Page number"`
		Limit         int    `query:"limit" minimum:"1" maximum:"100" default:"10" doc:"Items per page"`
		Search        string `query:"search" doc:"Search by name or description"`
		SchoolID      string `query:"school_id" doc:"Filter by school ID"`
	}) (*struct {
		Body school.PaginatedMajoritiesResponse
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		params := school.QueryParams{
			Page:     in.Page,
			Limit:    in.Limit,
			Search:   in.Search,
			SchoolID: in.SchoolID,
		}

		result, err := h.svc.GetAllMajorities(ctx, params)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.PaginatedMajoritiesResponse
		}{Body: *result}, nil
	})

	// POST /majorities - Create majority
	huma.Register(majorityGroup, huma.Operation{
		Method:  http.MethodPost,
		Path:    "",
		Summary: "Create a new majority",
		Tags:    []string{"School Management"},
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}, func(ctx context.Context, in *struct {
		Authorization string                       `header:"Authorization" required:"true" doc:"Bearer token"`
		Body          school.CreateMajorityRequest `json:"body"`
	}) (*struct {
		Body school.Majority
	}, error) {
		if err := h.validateToken(in.Authorization); err != nil {
			return nil, err
		}

		result, err := h.svc.CreateMajority(ctx, in.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError(err.Error())
		}

		return &struct {
			Body school.Majority
		}{Body: *result}, nil
	})

	// Continue with other endpoints...
}

// validateToken validates JWT token from Authorization header
func (h *Handler) validateToken(authHeader string) error {
	_, err := middleware.ValidateToken(authHeader, h.jwtSecrets)
	return err
}
