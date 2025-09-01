package repository

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"backend-service-internpro/internal/school"
)

// SchoolRepository defines the interface for school repository
type SchoolRepository interface {
	Create(ctx context.Context, entity *school.SchoolEntity) error
	GetByID(ctx context.Context, id uuid.UUID) (*school.SchoolEntity, error)
	GetAll(ctx context.Context, params school.QueryParams) ([]school.SchoolEntity, int, error)
	Update(ctx context.Context, entity *school.SchoolEntity) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByDomain(ctx context.Context, domain string) (*school.SchoolEntity, error)

	// Majority methods
	CreateMajority(ctx context.Context, entity *school.MajorityEntity) error
	GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.MajorityEntity, error)
	GetAllMajorities(ctx context.Context, params school.QueryParams) ([]school.MajorityEntity, int, error)
	UpdateMajority(ctx context.Context, entity *school.MajorityEntity) error
	DeleteMajority(ctx context.Context, id uuid.UUID) error

	// Class methods
	CreateClass(ctx context.Context, entity *school.ClassEntity) error
	GetClassByID(ctx context.Context, id uuid.UUID) (*school.ClassEntity, error)
	GetAllClasses(ctx context.Context, params school.QueryParams) ([]school.ClassEntity, int, error)
	UpdateClass(ctx context.Context, entity *school.ClassEntity) error
	DeleteClass(ctx context.Context, id uuid.UUID) error

	// Partner methods
	CreatePartner(ctx context.Context, entity *school.PartnerEntity) error
	GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.PartnerEntity, error)
	GetAllPartners(ctx context.Context, params school.QueryParams) ([]school.PartnerEntity, int, error)
	UpdatePartner(ctx context.Context, entity *school.PartnerEntity) error
	DeletePartner(ctx context.Context, id uuid.UUID) error
}

// schoolRepository implements SchoolRepository
type schoolRepository struct {
	db *gorm.DB
}

// NewSchoolRepository creates a new school repository
func NewSchoolRepository(db *gorm.DB) SchoolRepository {
	return &schoolRepository{
		db: db,
	}
}

// School methods
func (r *schoolRepository) Create(ctx context.Context, entity *school.SchoolEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *schoolRepository) GetByID(ctx context.Context, id uuid.UUID) (*school.SchoolEntity, error) {
	var entity school.SchoolEntity
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *schoolRepository) GetAll(ctx context.Context, params school.QueryParams) ([]school.SchoolEntity, int, error) {
	var entities []school.SchoolEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&school.SchoolEntity{}).Where("deleted_at IS NULL")

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(domain) LIKE ? OR LOWER(address) LIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(total), nil
}

func (r *schoolRepository) Update(ctx context.Context, entity *school.SchoolEntity) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", entity.ID).Updates(entity).Error
}

func (r *schoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&school.SchoolEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

func (r *schoolRepository) GetByDomain(ctx context.Context, domain string) (*school.SchoolEntity, error) {
	var entity school.SchoolEntity
	err := r.db.WithContext(ctx).Where("domain = ? AND deleted_at IS NULL", domain).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Majority methods
func (r *schoolRepository) CreateMajority(ctx context.Context, entity *school.MajorityEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *schoolRepository) GetMajorityByID(ctx context.Context, id uuid.UUID) (*school.MajorityEntity, error) {
	var entity school.MajorityEntity
	err := r.db.WithContext(ctx).Preload("School").Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *schoolRepository) GetAllMajorities(ctx context.Context, params school.QueryParams) ([]school.MajorityEntity, int, error) {
	var entities []school.MajorityEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&school.MajorityEntity{}).Preload("School").Where("deleted_at IS NULL")

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern)
	}

	// Apply school filter
	if params.SchoolID != "" {
		schoolID, err := uuid.Parse(params.SchoolID)
		if err == nil {
			query = query.Where("school_id = ?", schoolID)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(total), nil
}

func (r *schoolRepository) UpdateMajority(ctx context.Context, entity *school.MajorityEntity) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", entity.ID).Updates(entity).Error
}

func (r *schoolRepository) DeleteMajority(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&school.MajorityEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// Class methods
func (r *schoolRepository) CreateClass(ctx context.Context, entity *school.ClassEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *schoolRepository) GetClassByID(ctx context.Context, id uuid.UUID) (*school.ClassEntity, error) {
	var entity school.ClassEntity
	err := r.db.WithContext(ctx).Preload("School").Preload("Majority").
		Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *schoolRepository) GetAllClasses(ctx context.Context, params school.QueryParams) ([]school.ClassEntity, int, error) {
	var entities []school.ClassEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&school.ClassEntity{}).
		Preload("School").Preload("Majority").Where("deleted_at IS NULL")

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?",
			searchPattern, searchPattern)
	}

	// Apply school filter
	if params.SchoolID != "" {
		schoolID, err := uuid.Parse(params.SchoolID)
		if err == nil {
			query = query.Where("school_id = ?", schoolID)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(total), nil
}

func (r *schoolRepository) UpdateClass(ctx context.Context, entity *school.ClassEntity) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", entity.ID).Updates(entity).Error
}

func (r *schoolRepository) DeleteClass(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&school.ClassEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}

// Partner methods
func (r *schoolRepository) CreatePartner(ctx context.Context, entity *school.PartnerEntity) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *schoolRepository) GetPartnerByID(ctx context.Context, id uuid.UUID) (*school.PartnerEntity, error) {
	var entity school.PartnerEntity
	err := r.db.WithContext(ctx).Preload("School").Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *schoolRepository) GetAllPartners(ctx context.Context, params school.QueryParams) ([]school.PartnerEntity, int, error) {
	var entities []school.PartnerEntity
	var total int64

	query := r.db.WithContext(ctx).Model(&school.PartnerEntity{}).Preload("School").Where("deleted_at IS NULL")

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(contact_name) LIKE ? OR LOWER(contact_email) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Apply school filter
	if params.SchoolID != "" {
		schoolID, err := uuid.Parse(params.SchoolID)
		if err == nil {
			query = query.Where("school_id = ?", schoolID)
		}
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	if err := query.Offset(offset).Limit(params.Limit).Find(&entities).Error; err != nil {
		return nil, 0, err
	}

	return entities, int(total), nil
}

func (r *schoolRepository) UpdatePartner(ctx context.Context, entity *school.PartnerEntity) error {
	return r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", entity.ID).Updates(entity).Error
}

func (r *schoolRepository) DeletePartner(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&school.PartnerEntity{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", gorm.Expr("NOW()")).Error
}
