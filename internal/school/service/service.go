package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend-service-internpro/internal/school"
	"backend-service-internpro/internal/school/repository"
)

// SchoolService defines the interface for school service
type SchoolService interface {
	CreateSchool(ctx context.Context, req school.CreateSchoolRequest) (*school.School, error)
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*school.School, error)
	GetAllSchools(ctx context.Context, params school.QueryParams) (*school.PaginatedSchoolsResponse, error)
	UpdateSchool(ctx context.Context, id uuid.UUID, req school.UpdateSchoolRequest) (*school.School, error)
	DeleteSchool(ctx context.Context, id uuid.UUID) error

	// Majority methods
	CreateMajority(ctx context.Context, req school.CreateMajorityRequest) (*school.Majority, error)
	GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.Majority, error)
	GetAllMajorities(ctx context.Context, params school.QueryParams) (*school.PaginatedMajoritiesResponse, error)
	UpdateMajority(ctx context.Context, id uuid.UUID, req school.UpdateMajorityRequest) (*school.Majority, error)
	DeleteMajority(ctx context.Context, id uuid.UUID) error

	// Class methods
	CreateClass(ctx context.Context, req school.CreateClassRequest) (*school.Class, error)
	GetClassByID(ctx context.Context, id uuid.UUID) (*school.Class, error)
	GetAllClasses(ctx context.Context, params school.QueryParams) (*school.PaginatedClassesResponse, error)
	UpdateClass(ctx context.Context, id uuid.UUID, req school.UpdateClassRequest) (*school.Class, error)
	DeleteClass(ctx context.Context, id uuid.UUID) error

	// Partner methods
	CreatePartner(ctx context.Context, req school.CreatePartnerRequest) (*school.Partner, error)
	GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.Partner, error)
	GetAllPartners(ctx context.Context, params school.QueryParams) (*school.PaginatedPartnersResponse, error)
	UpdatePartner(ctx context.Context, id uuid.UUID, req school.UpdatePartnerRequest) (*school.Partner, error)
	DeletePartner(ctx context.Context, id uuid.UUID) error
}

// schoolService implements SchoolService
type schoolService struct {
	repo repository.SchoolRepository
}

// NewSchoolService creates a new school service
func NewSchoolService(repo repository.SchoolRepository) SchoolService {
	return &schoolService{
		repo: repo,
	}
}

// School methods
func (s *schoolService) CreateSchool(ctx context.Context, req school.CreateSchoolRequest) (*school.School, error) {
	entity := &school.SchoolEntity{
		ID:        uuid.New(),
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Address != "" {
		entity.Address = &req.Address
	}

	if req.Domain != "" {
		// Check if domain already exists
		existing, err := s.repo.GetByDomain(ctx, req.Domain)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil {
			return nil, errors.New("domain already exists")
		}
		entity.Domain = &req.Domain
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToSchool()
	return &result, nil
}

func (s *schoolService) GetSchoolByID(ctx context.Context, id uuid.UUID) (*school.School, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	result := entity.ToSchool()
	return &result, nil
}

func (s *schoolService) GetAllSchools(ctx context.Context, params school.QueryParams) (*school.PaginatedSchoolsResponse, error) {
	// Set default pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	entities, total, err := s.repo.GetAll(ctx, params)
	if err != nil {
		return nil, err
	}

	schools := make([]school.School, len(entities))
	for i, entity := range entities {
		schools[i] = entity.ToSchool()
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &school.PaginatedSchoolsResponse{
		Schools: schools,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *schoolService) UpdateSchool(ctx context.Context, id uuid.UUID, req school.UpdateSchoolRequest) (*school.School, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Address != "" {
		entity.Address = &req.Address
	}
	if req.Domain != "" {
		// Check if domain already exists (excluding current record)
		existing, err := s.repo.GetByDomain(ctx, req.Domain)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("domain already exists")
		}
		entity.Domain = &req.Domain
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToSchool()
	return &result, nil
}

func (s *schoolService) DeleteSchool(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("school not found")
		}
		return err
	}

	return s.repo.Delete(ctx, id)
}

// Majority methods
func (s *schoolService) CreateMajority(ctx context.Context, req school.CreateMajorityRequest) (*school.Majority, error) {
	// Verify school exists
	_, err := s.repo.GetByID(ctx, req.SchoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	entity := &school.MajorityEntity{
		ID:        uuid.New(),
		SchoolID:  req.SchoolID,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Description != "" {
		entity.Description = &req.Description
	}

	if err := s.repo.CreateMajority(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToMajority()
	return &result, nil
}

func (s *schoolService) GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.Majority, error) {
	entity, err := s.repo.GetMajorityByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("majority not found")
		}
		return nil, err
	}

	result := entity.ToMajority()
	if entity.School.ID != uuid.Nil {
		schoolDTO := entity.School.ToSchool()
		result.School = &schoolDTO
	}
	return &result, nil
}

func (s *schoolService) GetAllMajorities(ctx context.Context, params school.QueryParams) (*school.PaginatedMajoritiesResponse, error) {
	// Set default pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	entities, total, err := s.repo.GetAllMajorities(ctx, params)
	if err != nil {
		return nil, err
	}

	majorities := make([]school.Majority, len(entities))
	for i, entity := range entities {
		majorities[i] = entity.ToMajority()
		if entity.School.ID != uuid.Nil {
			schoolDTO := entity.School.ToSchool()
			majorities[i].School = &schoolDTO
		}
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &school.PaginatedMajoritiesResponse{
		Majorities: majorities,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *schoolService) UpdateMajority(ctx context.Context, id uuid.UUID, req school.UpdateMajorityRequest) (*school.Majority, error) {
	entity, err := s.repo.GetMajorityByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("majority not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.SchoolID != uuid.Nil {
		// Verify school exists
		_, err := s.repo.GetByID(ctx, req.SchoolID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("school not found")
			}
			return nil, err
		}
		entity.SchoolID = req.SchoolID
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Description != "" {
		entity.Description = &req.Description
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.UpdateMajority(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToMajority()
	return &result, nil
}

func (s *schoolService) DeleteMajority(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetMajorityByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("majority not found")
		}
		return err
	}

	return s.repo.DeleteMajority(ctx, id)
}

// Class methods
func (s *schoolService) CreateClass(ctx context.Context, req school.CreateClassRequest) (*school.Class, error) {
	// Verify school exists
	_, err := s.repo.GetByID(ctx, req.SchoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	// Verify majority exists
	_, err = s.repo.GetMajorityByID(ctx, req.MajorityID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("majority not found")
		}
		return nil, err
	}

	entity := &school.ClassEntity{
		ID:         uuid.New(),
		SchoolID:   req.SchoolID,
		MajorityID: req.MajorityID,
		Name:       req.Name,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if req.Description != "" {
		entity.Description = &req.Description
	}

	if err := s.repo.CreateClass(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToClass()
	return &result, nil
}

func (s *schoolService) GetClassByID(ctx context.Context, id uuid.UUID) (*school.Class, error) {
	entity, err := s.repo.GetClassByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, err
	}

	result := entity.ToClass()
	if entity.School.ID != uuid.Nil {
		schoolDTO := entity.School.ToSchool()
		result.School = &schoolDTO
	}
	if entity.Majority.ID != uuid.Nil {
		majorityDTO := entity.Majority.ToMajority()
		result.Majority = &majorityDTO
	}
	return &result, nil
}

func (s *schoolService) GetAllClasses(ctx context.Context, params school.QueryParams) (*school.PaginatedClassesResponse, error) {
	// Set default pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	entities, total, err := s.repo.GetAllClasses(ctx, params)
	if err != nil {
		return nil, err
	}

	classes := make([]school.Class, len(entities))
	for i, entity := range entities {
		classes[i] = entity.ToClass()
		if entity.School.ID != uuid.Nil {
			schoolDTO := entity.School.ToSchool()
			classes[i].School = &schoolDTO
		}
		if entity.Majority.ID != uuid.Nil {
			majorityDTO := entity.Majority.ToMajority()
			classes[i].Majority = &majorityDTO
		}
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &school.PaginatedClassesResponse{
		Classes: classes,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *schoolService) UpdateClass(ctx context.Context, id uuid.UUID, req school.UpdateClassRequest) (*school.Class, error) {
	entity, err := s.repo.GetClassByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.SchoolID != uuid.Nil {
		// Verify school exists
		_, err := s.repo.GetByID(ctx, req.SchoolID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("school not found")
			}
			return nil, err
		}
		entity.SchoolID = req.SchoolID
	}
	if req.MajorityID != uuid.Nil {
		// Verify majority exists
		_, err := s.repo.GetMajorityByID(ctx, req.MajorityID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("majority not found")
			}
			return nil, err
		}
		entity.MajorityID = req.MajorityID
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Description != "" {
		entity.Description = &req.Description
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.UpdateClass(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToClass()
	return &result, nil
}

func (s *schoolService) DeleteClass(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetClassByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("class not found")
		}
		return err
	}

	return s.repo.DeleteClass(ctx, id)
}

// Partner methods
func (s *schoolService) CreatePartner(ctx context.Context, req school.CreatePartnerRequest) (*school.Partner, error) {
	// Verify school exists
	_, err := s.repo.GetByID(ctx, req.SchoolID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	entity := &school.PartnerEntity{
		ID:        uuid.New(),
		SchoolID:  req.SchoolID,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if req.Website != "" {
		entity.Website = &req.Website
	}
	if req.Description != "" {
		entity.Description = &req.Description
	}
	if req.Address != "" {
		entity.Address = &req.Address
	}
	if req.ContactName != "" {
		entity.ContactName = &req.ContactName
	}
	if req.ContactPerson != "" {
		entity.ContactPerson = &req.ContactPerson
	}
	if req.ContactEmail != "" {
		entity.ContactEmail = &req.ContactEmail
	}

	if err := s.repo.CreatePartner(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToPartner()
	return &result, nil
}

func (s *schoolService) GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.Partner, error) {
	entity, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}

	result := entity.ToPartner()
	if entity.School.ID != uuid.Nil {
		schoolDTO := entity.School.ToSchool()
		result.School = &schoolDTO
	}
	return &result, nil
}

func (s *schoolService) GetAllPartners(ctx context.Context, params school.QueryParams) (*school.PaginatedPartnersResponse, error) {
	// Set default pagination
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	entities, total, err := s.repo.GetAllPartners(ctx, params)
	if err != nil {
		return nil, err
	}

	partners := make([]school.Partner, len(entities))
	for i, entity := range entities {
		partners[i] = entity.ToPartner()
		if entity.School.ID != uuid.Nil {
			schoolDTO := entity.School.ToSchool()
			partners[i].School = &schoolDTO
		}
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &school.PaginatedPartnersResponse{
		Partners: partners,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *schoolService) UpdatePartner(ctx context.Context, id uuid.UUID, req school.UpdatePartnerRequest) (*school.Partner, error) {
	entity, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}

	// Update fields if provided
	if req.SchoolID != uuid.Nil {
		// Verify school exists
		_, err := s.repo.GetByID(ctx, req.SchoolID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("school not found")
			}
			return nil, err
		}
		entity.SchoolID = req.SchoolID
	}
	if req.Name != "" {
		entity.Name = req.Name
	}
	if req.Website != "" {
		entity.Website = &req.Website
	}
	if req.Description != "" {
		entity.Description = &req.Description
	}
	if req.Address != "" {
		entity.Address = &req.Address
	}
	if req.ContactName != "" {
		entity.ContactName = &req.ContactName
	}
	if req.ContactPerson != "" {
		entity.ContactPerson = &req.ContactPerson
	}
	if req.ContactEmail != "" {
		entity.ContactEmail = &req.ContactEmail
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.UpdatePartner(ctx, entity); err != nil {
		return nil, err
	}

	result := entity.ToPartner()
	return &result, nil
}

func (s *schoolService) DeletePartner(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("partner not found")
		}
		return err
	}

	return s.repo.DeletePartner(ctx, id)
}
