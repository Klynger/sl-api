package userHandlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	"sl-api/api/repositories/user"
)

type API struct {
	validator        *validator.Validate
	repository       *userRepository.UserRepository
	authSessionStore *sessions.CookieStore
	authMaxAge       int
	isDebugging      bool
}

func New(db *gorm.DB, authSessionStore *sessions.CookieStore, authMaxAge int, isDebugging bool, validator *validator.Validate) *API {
	return &API{
		repository:       userRepository.New(db),
		validator:        validator,
		authSessionStore: authSessionStore,
		authMaxAge:       authMaxAge,
		isDebugging:      isDebugging,
	}
}
