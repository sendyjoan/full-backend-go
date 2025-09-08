package school

import (
	"backend-service-internpro/internal/pkg/response"
	"time"

	"github.com/google/uuid"
)

// School represents the school data transfer object
type School struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address,omitempty"`
	Domain    string    `json:"domain,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateSchoolRequest represents the request to create a school
type CreateSchoolRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=255"`
	Address string `json:"address,omitempty" validate:"max=255"`
	Domain  string `json:"domain,omitempty" validate:"max=255"`
}

// UpdateSchoolRequest represents the request to update a school
type UpdateSchoolRequest struct {
	Name    string `json:"name,omitempty" validate:"min=1,max=255"`
	Address string `json:"address,omitempty" validate:"max=255"`
	Domain  string `json:"domain,omitempty" validate:"max=255"`
}

// Majority represents the majority data transfer object
type Majority struct {
	ID          uuid.UUID `json:"id"`
	SchoolID    uuid.UUID `json:"school_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	School      *School   `json:"school,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMajorityRequest represents the request to create a majority
type CreateMajorityRequest struct {
	SchoolID    uuid.UUID `json:"school_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty"`
}

// UpdateMajorityRequest represents the request to update a majority
type UpdateMajorityRequest struct {
	SchoolID    uuid.UUID `json:"school_id,omitempty"`
	Name        string    `json:"name,omitempty" validate:"min=1,max=255"`
	Description string    `json:"description,omitempty"`
}

// Class represents the class data transfer object
type Class struct {
	ID          uuid.UUID `json:"id"`
	SchoolID    uuid.UUID `json:"school_id"`
	MajorityID  uuid.UUID `json:"majority_id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	School      *School   `json:"school,omitempty"`
	Majority    *Majority `json:"majority,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateClassRequest represents the request to create a class
type CreateClassRequest struct {
	SchoolID    uuid.UUID `json:"school_id" validate:"required"`
	MajorityID  uuid.UUID `json:"majority_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description,omitempty"`
}

// UpdateClassRequest represents the request to update a class
type UpdateClassRequest struct {
	SchoolID    uuid.UUID `json:"school_id,omitempty"`
	MajorityID  uuid.UUID `json:"majority_id,omitempty"`
	Name        string    `json:"name,omitempty" validate:"min=1,max=255"`
	Description string    `json:"description,omitempty"`
}

// Partner represents the partner data transfer object
type Partner struct {
	ID            uuid.UUID `json:"id"`
	SchoolID      uuid.UUID `json:"school_id"`
	Name          string    `json:"name"`
	Website       string    `json:"website,omitempty"`
	Description   string    `json:"description,omitempty"`
	Address       string    `json:"address,omitempty"`
	ContactName   string    `json:"contact_name,omitempty"`
	ContactPerson string    `json:"contact_person,omitempty"`
	ContactEmail  string    `json:"contact_email,omitempty"`
	School        *School   `json:"school,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreatePartnerRequest represents the request to create a partner
type CreatePartnerRequest struct {
	SchoolID      uuid.UUID `json:"school_id" validate:"required"`
	Name          string    `json:"name" validate:"required,min=1,max=255"`
	Website       string    `json:"website,omitempty" validate:"max=255"`
	Description   string    `json:"description,omitempty"`
	Address       string    `json:"address,omitempty" validate:"max=255"`
	ContactName   string    `json:"contact_name,omitempty" validate:"max=255"`
	ContactPerson string    `json:"contact_person,omitempty" validate:"max=255"`
	ContactEmail  string    `json:"contact_email,omitempty" validate:"email,max=255"`
}

// UpdatePartnerRequest represents the request to update a partner
type UpdatePartnerRequest struct {
	SchoolID      uuid.UUID `json:"school_id,omitempty"`
	Name          string    `json:"name,omitempty" validate:"min=1,max=255"`
	Website       string    `json:"website,omitempty" validate:"max=255"`
	Description   string    `json:"description,omitempty"`
	Address       string    `json:"address,omitempty" validate:"max=255"`
	ContactName   string    `json:"contact_name,omitempty" validate:"max=255"`
	ContactPerson string    `json:"contact_person,omitempty" validate:"max=255"`
	ContactEmail  string    `json:"contact_email,omitempty" validate:"email,max=255"`
}

// SchoolListData represents the data structure for school list
type SchoolListData struct {
	Schools    []School         `json:"schools"`
	Pagination PaginationResult `json:"pagination"`
}

// MajorityListData represents the data structure for majority list
type MajorityListData struct {
	Majorities []Majority       `json:"majorities"`
	Pagination PaginationResult `json:"pagination"`
}

// ClassListData represents the data structure for class list
type ClassListData struct {
	Classes    []Class          `json:"classes"`
	Pagination PaginationResult `json:"pagination"`
}

// PartnerListData represents the data structure for partner list
type PartnerListData struct {
	Partners   []Partner        `json:"partners"`
	Pagination PaginationResult `json:"pagination"`
}

// PaginatedSchoolsResponse represents the paginated response for schools
type PaginatedSchoolsResponse = response.ApiResponse

// PaginatedMajoritiesResponse represents the paginated response for majorities
type PaginatedMajoritiesResponse = response.ApiResponse

// PaginatedClassesResponse represents the paginated response for classes
type PaginatedClassesResponse = response.ApiResponse

// PaginatedPartnersResponse represents the paginated response for partners
type PaginatedPartnersResponse = response.ApiResponse

// SchoolResponse represents single school response
type SchoolResponse = response.ApiResponse

// MajorityResponse represents single majority response
type MajorityResponse = response.ApiResponse

// ClassResponse represents single class response
type ClassResponse = response.ApiResponse

// PartnerResponse represents single partner response
type PartnerResponse = response.ApiResponse

// BasicResponse represents basic response with message
type BasicResponse = response.ApiResponse

// PaginationResult represents pagination metadata
type PaginationResult struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// QueryParams represents common query parameters for listing
type QueryParams struct {
	Page     int    `json:"page" validate:"min=1"`
	Limit    int    `json:"limit" validate:"min=1,max=100"`
	Search   string `json:"search"`
	SchoolID string `json:"school_id"`
}
