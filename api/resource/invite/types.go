package invite

import (
	"context"

	"github.com/google/uuid"

	groupMemberModel "sl-api/api/model/group_member"
	userModel "sl-api/api/model/user"
)

type InviteData struct {
	GroupExists      bool
	IsAlreadyInvited bool
	AuthedUser       *userModel.PublicUser
}
type GetInviteDataInput struct {
	ctx           context.Context
	GroupID       uuid.UUID
	AuthedUserID  uuid.UUID
	InvitedUserID uuid.UUID
}

type GetInviteDataOutput struct {
	GroupExists      bool
	IsAlreadyInvited bool
	IsAlreadyMember  bool
	AuthedUser       *userModel.PublicUser
	InvitedUser      *userModel.PublicUser
	AuthedUserMember *groupMemberModel.PublicGroupMember
}
