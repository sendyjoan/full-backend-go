package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend-service-internpro/internal/pkg/constants"
	"backend-service-internpro/internal/pkg/response"
	"backend-service-internpro/internal/school"
	"backend-service-internpro/internal/school/repository"
)

// SchoolService defines the interface for school service
type SchoolService interface {
	CreateSchool(ctx context.Context, req school.CreateSchoolRequest) (*school.SchoolResponse, error)
	GetSchoolByID(ctx context.Context, id uuid.UUID) (*school.SchoolResponse, error)
	GetAllSchools(ctx context.Context, params school.QueryParams) (*school.PaginatedSchoolsResponse, error)
	UpdateSchool(ctx context.Context, id uuid.UUID, req school.UpdateSchoolRequest) (*school.SchoolResponse, error)
	DeleteSchool(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error)

	// Majority methods
	CreateMajority(ctx context.Context, req school.CreateMajorityRequest) (*school.MajorityResponse, error)
	GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.MajorityResponse, error)
	GetAllMajorities(ctx context.Context, params school.QueryParams) (*school.PaginatedMajoritiesResponse, error)
	UpdateMajority(ctx context.Context, id uuid.UUID, req school.UpdateMajorityRequest) (*school.MajorityResponse, error)
	DeleteMajority(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error)

	// Class methods
	CreateClass(ctx context.Context, req school.CreateClassRequest) (*school.ClassResponse, error)
	GetClassByID(ctx context.Context, id uuid.UUID) (*school.ClassResponse, error)
	GetAllClasses(ctx context.Context, params school.QueryParams) (*school.PaginatedClassesResponse, error)
	UpdateClass(ctx context.Context, id uuid.UUID, req school.UpdateClassRequest) (*school.ClassResponse, error)
	DeleteClass(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error)

	// Partner methods
	CreatePartner(ctx context.Context, req school.CreatePartnerRequest) (*school.PartnerResponse, error)
	GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.PartnerResponse, error)
	GetAllPartners(ctx context.Context, params school.QueryParams) (*school.PaginatedPartnersResponse, error)
	UpdatePartner(ctx context.Context, id uuid.UUID, req school.UpdatePartnerRequest) (*school.PartnerResponse, error)
	DeletePartner(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error)
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
func (s *schoolService) CreateSchool(ctx context.Context, req school.CreateSchoolRequest) (*school.SchoolResponse, error) {
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
	return response.Success(constants.SchoolCreateSuccess, result), nil
}

func (s *schoolService) GetSchoolByID(ctx context.Context, id uuid.UUID) (*school.SchoolResponse, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	result := entity.ToSchool()
	return response.Success(constants.SchoolDetailSuccess, result), nil
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

	data := school.SchoolListData{
		Schools: schools,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(constants.SchoolListSuccess, data), nil
}

func (s *schoolService) UpdateSchool(ctx context.Context, id uuid.UUID, req school.UpdateSchoolRequest) (*school.SchoolResponse, error) {
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
		// Check if domain already exists and belongs to different school
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
	return response.Success(constants.SchoolUpdateSuccess, result), nil
}

func (s *schoolService) DeleteSchool(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error) {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("school not found")
		}
		return nil, err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}

	return response.SuccessWithoutData(constants.SchoolDeleteSuccess), nil
}

// Majority methods
func (s *schoolService) CreateMajority(ctx context.Context, req school.CreateMajorityRequest) (*school.MajorityResponse, error) {
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
	return response.Success(constants.MajorityCreateSuccess, result), nil
}

func (s *schoolService) GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.MajorityResponse, error) {
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
	return response.Success(constants.MajorityDetailSuccess, result), nil
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

	data := school.MajorityListData{
		Majorities: majorities,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(constants.MajorityListSuccess, data), nil
}

func (s *schoolService) UpdateMajority(ctx context.Context, id uuid.UUID, req school.UpdateMajorityRequest) (*school.MajorityResponse, error) {
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
	return response.Success(constants.MajorityUpdateSuccess, result), nil
}

func (s *schoolService) DeleteMajority(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error) {
	_, err := s.repo.GetMajorityByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("majority not found")
		}
		return nil, err
	}

	if err := s.repo.DeleteMajority(ctx, id); err != nil {
		return nil, err
	}

	return response.SuccessWithoutData(constants.MajorityDeleteSuccess), nil
}

// Class methods
func (s *schoolService) CreateClass(ctx context.Context, req school.CreateClassRequest) (*school.ClassResponse, error) {
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
	return response.Success(constants.ClassCreateSuccess, result), nil
}

func (s *schoolService) GetClassByID(ctx context.Context, id uuid.UUID) (*school.ClassResponse, error) {
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
	return response.Success(constants.ClassGetSuccess, result), nil
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

	data := school.ClassListData{
		Classes: classes,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(constants.ClassGetAllSuccess, data), nil
}

func (s *schoolService) UpdateClass(ctx context.Context, id uuid.UUID, req school.UpdateClassRequest) (*school.ClassResponse, error) {
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
	return response.Success(constants.ClassUpdateSuccess, result), nil
}

func (s *schoolService) DeleteClass(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error) {
	_, err := s.repo.GetClassByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class not found")
		}
		return nil, err
	}

	if err := s.repo.DeleteClass(ctx, id); err != nil {
		return nil, err
	}

	return response.SuccessWithoutData(constants.ClassDeleteSuccess), nil
}

// Partner methods
func (s *schoolService) CreatePartner(ctx context.Context, req school.CreatePartnerRequest) (*school.PartnerResponse, error) {
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
	return response.Success(constants.PartnerCreateSuccess, result), nil
}

func (s *schoolService) GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.PartnerResponse, error) {
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
	return response.Success(constants.PartnerGetSuccess, result), nil
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

	data := school.PartnerListData{
		Partners: partners,
		Pagination: school.PaginationResult{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response.Success(constants.PartnerGetAllSuccess, data), nil
}

func (s *schoolService) UpdatePartner(ctx context.Context, id uuid.UUID, req school.UpdatePartnerRequest) (*school.PartnerResponse, error) {
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
	return response.Success(constants.PartnerUpdateSuccess, result), nil
}

func (s *schoolService) DeletePartner(ctx context.Context, id uuid.UUID) (*school.BasicResponse, error) {
	_, err := s.repo.GetPartnerByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("partner not found")
		}
		return nil, err
	}

	if err := s.repo.DeletePartner(ctx, id); err != nil {
		return nil, err
	}

	return response.SuccessWithoutData(constants.PartnerDeleteSuccess), nil
}
