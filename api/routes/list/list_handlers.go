package listHandlers

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	listService "sl-api/api/services/list"
)

type API struct {
	validator   *validator.Validate
	listService *listService.ListService
}

func New(db *gorm.DB, validator *validator.Validate) *API {
	return &API{
		validator:   validator,
		listService: listService.New(db),
	}
}
