package group_service

import (
	"context"

	"github.com/google/uuid"

	groupModel "sl-api/api/model/group"
	groupMemberModel "sl-api/api/model/group_member"
	inviteModel "sl-api/api/model/invite"
	userModel "sl-api/api/model/user"
)

type GroupRepository interface {
	Create(group *groupModel.Group) (*groupModel.Group, error)
	Read(groupID uuid.UUID) (*groupModel.Group, error)
}

type GroupMemberRepository interface {
	Create(member *groupMemberModel.GroupMember) (*groupMemberModel.GroupMember, error)
	FindByUserAndGroup(userID, groupID uuid.UUID) (*groupMemberModel.GroupMemberDetails, error)
}

type InviteRepository interface {
	Create(invite *inviteModel.Invite) (*inviteModel.Invite, error)
	GetCreateInviteData(input *GetCreateInviteDataInput) (*GetCreateInviteDataOutput, error)
	GetAcceptInviteData(groupID, invitedUserID uuid.UUID) (*inviteModel.Invite, error)
	Delete(id uuid.UUID) error
}

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

type UserRepository interface {
	Read(userID uuid.UUID) (*userModel.User, error)
}
