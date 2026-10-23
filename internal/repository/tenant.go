package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"condominium-manager/internal/model"
)

var (
	ErrAlreadyExists = errors.New("tenant already exists")
	ErrNotFound      = errors.New("tenant not found")
)

// TenantFilter holds optional search criteria; zero values are ignored.
type TenantFilter struct {
	ID       *int32
	Name     string
	Lastname string
}

// TenantRepository abstracts the database access for tenants.
type TenantRepository interface {
	List(ctx context.Context) ([]model.Tenant, error)
	Search(ctx context.Context, f TenantFilter) ([]model.Tenant, error)
	// Create stores the tenant, returning ErrAlreadyExists on duplicates.
	Create(ctx context.Context, t *model.Tenant) error
	// Delete removes the tenant, returning ErrNotFound if it does not exist.
	Delete(ctx context.Context, id int32) error
	// Update replaces the tenant's fields, returning ErrNotFound or ErrAlreadyExists.
	Update(ctx context.Context, t *model.Tenant) error
}

type gormTenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &gormTenantRepository{db: db}
}

func (r *gormTenantRepository) List(ctx context.Context) ([]model.Tenant, error) {
	tenants := []model.Tenant{}
	err := r.db.WithContext(ctx).Find(&tenants).Error
	return tenants, err
}

func (r *gormTenantRepository) Search(ctx context.Context, f TenantFilter) ([]model.Tenant, error) {
	q := r.db.WithContext(ctx).Model(&model.Tenant{})
	if f.ID != nil {
		q = q.Where("id = ?", *f.ID)
	}
	if f.Name != "" {
		q = q.Where("name = ?", f.Name)
	}
	if f.Lastname != "" {
		q = q.Where("lastname = ?", f.Lastname)
	}
	tenants := []model.Tenant{}
	err := q.Find(&tenants).Error
	return tenants, err
}

func (r *gormTenantRepository) Create(ctx context.Context, t *model.Tenant) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Tenant{}).
		Where("name = ? AND lastname = ?", t.Name, t.Lastname).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrAlreadyExists
	}
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (r *gormTenantRepository) Delete(ctx context.Context, id int32) error {
	res := r.db.WithContext(ctx).Delete(&model.Tenant{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormTenantRepository) Update(ctx context.Context, t *model.Tenant) error {
	res := r.db.WithContext(ctx).Model(&model.Tenant{}).
		Where("id = ?", t.ID).
		Select("name", "lastname", "age").
		Updates(t)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return ErrAlreadyExists
		}
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
