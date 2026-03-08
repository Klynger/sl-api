package productRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"sl-api/api/model/product"
)

type ProductRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) List() (productModel.Products, error) {
	products := make([]*productModel.Product, 0)

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Create(product *productModel.Product) (*productModel.Product, error) {
	if err := r.db.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) Read(id uuid.UUID) (*productModel.Product, error) {
	product := &productModel.Product{}
	if err := r.db.Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductRepository) Update(product *productModel.Product) (int64, error) {
	result := r.db.Model(&productModel.Product{}).
		Select("Name", "Description", "UpdatedAt").
		Where("id = ?", product.ID).
		Updates(product)

	return result.RowsAffected, result.Error
}

func (r *ProductRepository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Where("id = ?", id).Delete(&productModel.Product{})

	return result.RowsAffected, result.Error
}
