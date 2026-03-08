package groupHandlers

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"sl-api/api/services/group"
)

type API struct {
	validator          *validator.Validate
	createGroupService *groupService.CreateService
}

func New(db *gorm.DB, validator *validator.Validate) *API {
	return &API{
		validator:          validator,
		createGroupService: groupService.NewCreateService(db),
	}
}
