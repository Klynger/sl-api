package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Products, error) {
	products := make([]*Product, 0)

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Create(product *Product) (*Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) Read(id uuid.UUID) (*Product, error) {
	product := &Product{}
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) Update(product *Product) (int64, error) {
	result := r.db.Model(&Product{}).
		Select("Name", "Description", "UpdatedAt").
		Where("id = ?", product.ID).
		Updates(product)

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&Product{})

	return result.RowsAffected, result.Error
}
