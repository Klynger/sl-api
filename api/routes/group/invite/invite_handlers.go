package inviteHandlers

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	"sl-api/api/services/group"
)

type API struct {
	validator           *validator.Validate
	inviteService       *groupService.CreateInviteService
	acceptInviteService *groupService.AcceptInviteService
}

func New(db *gorm.DB, validator *validator.Validate) *API {
	return &API{
		validator: validator,
		inviteService: groupService.NewCreateInviteService(
			db,
		),
		acceptInviteService: groupService.NewAcceptInviteService(
			db,
		),
	}
}
