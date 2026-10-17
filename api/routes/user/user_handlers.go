package userHandlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"

	groupMemberRepository "sl-api/api/repositories/group_member"
	"sl-api/api/repositories/user"
)

type API struct {
	validator        *validator.Validate
	repository       *userRepository.UserRepository
	memberRepository *groupMemberRepository.GroupMemberRepository
	authSessionStore *sessions.CookieStore
	authMaxAge       int
	isDebugging      bool
}

func New(db *gorm.DB, authSessionStore *sessions.CookieStore, authMaxAge int, isDebugging bool, validator *validator.Validate) *API {
	return &API{
		repository:       userRepository.New(db),
		memberRepository: groupMemberRepository.New(db),
		validator:        validator,
		authSessionStore: authSessionStore,
		authMaxAge:       authMaxAge,
		isDebugging:      isDebugging,
	}
}
