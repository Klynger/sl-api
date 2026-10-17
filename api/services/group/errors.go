package groupService

import "errors"

var (
	ErrInviteSelf              = errors.New("cannot invite yourself")
	ErrGroupNotFound           = errors.New("group not found")
	ErrAlreadyInvited          = errors.New("user already invited")
	ErrAlreadyMember           = errors.New("user is already a member")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInviteNotFound          = errors.New("invite not found")
	ErrUserNotFound            = errors.New("user not found")
)
