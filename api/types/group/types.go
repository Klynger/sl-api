package group

import (
	"context"

	"github.com/google/uuid"

	groupMemberModel "sl-api/api/model/group_member"
	userModel "sl-api/api/model/user"
)

type GetCreateInviteDataInput struct {
	Ctx           context.Context
	GroupID       uuid.UUID
	AuthedUserID  uuid.UUID
	InvitedUserID uuid.UUID
}

type GetCreateInviteDataOutput struct {
	GroupExists      bool
	IsAlreadyInvited bool
	IsAlreadyMember  bool
	AuthedUser       *userModel.PublicUser
	InvitedUser      *userModel.PublicUser
	AuthedUserMember *groupMemberModel.PublicGroupMember
}
