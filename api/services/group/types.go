package group_service

import (
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
}

type InviteRepository interface {
	Create(invite *inviteModel.Invite) (*inviteModel.Invite, error)
}

type UserRepository interface {
	Read(userID uuid.UUID) (*userModel.User, error)
}
