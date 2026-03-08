package productHandlers

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"sl-api/api/repositories/product"
)

type API struct {
	repository *productRepository.ProductRepository
	validator  *validator.Validate
}

func New(db *gorm.DB, v *validator.Validate) *API {
	return &API{
		repository: productRepository.New(db),
		validator:  v,
	}
}
