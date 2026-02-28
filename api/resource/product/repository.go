package product

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	productModel "sl-api/api/model/product"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (productModel.Products, error) {
	products := make([]*productModel.Product, 0)

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) Create(product *productModel.Product) (*productModel.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) Read(id uuid.UUID) (*productModel.Product, error) {
	product := &productModel.Product{}
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *Repository) Update(product *productModel.Product) (int64, error) {
	result := r.db.Model(&productModel.Product{}).
		Select("Name", "Description", "UpdatedAt").
		Where("id = ?", product.ID).
		Updates(product)

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&productModel.Product{})

	return result.RowsAffected, result.Error
}
