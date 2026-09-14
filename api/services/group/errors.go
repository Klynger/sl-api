package groupService

import "errors"

var (
	ErrInviteSelf              = errors.New("cannot invite yourself")
	ErrGroupNotFound           = errors.New("group not found")
	ErrAlreadyInvited          = errors.New("user already invited")
	ErrAlreadyMember           = errors.New("user is already a member")
	ErrNotAMember              = errors.New("user is not a member of this group")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
	ErrInviteNotFound          = errors.New("invite not found")
)
