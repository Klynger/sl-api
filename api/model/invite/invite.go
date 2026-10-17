package inviteModel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invite struct {
	ID            uuid.UUID `gorm:"primarykey"`
	GroupID       uuid.UUID `gorm:"column:group_id"`
	SenderID      uuid.UUID `gorm:"column:sender_id"`       // The user who sent the invite
	InvitedUserID uuid.UUID `gorm:"column:invited_user_id"` // The user who is invited to the group
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

type Form struct {
	InvitedUsername string `json:"invitedUsername" validate:"required,max=50"`
}

// PendingInvite is one live invite joined with the group it belongs to, the
// user who sent it, and the user it is addressed to, so a recipient can
// recognize and act on it and a group can see who it is waiting on. It is the
// scan target for the listing queries.
type PendingInvite struct {
	ID              uuid.UUID
	GroupID         uuid.UUID
	GroupName       string
	SenderID        uuid.UUID
	SenderUsername  string
	InvitedUserID   uuid.UUID
	InvitedUsername string
}

func (p *PendingInvite) ToDto() *PendingInviteDTO {
	return &PendingInviteDTO{
		ID:              p.ID.String(),
		GroupID:         p.GroupID.String(),
		GroupName:       p.GroupName,
		SenderID:        p.SenderID.String(),
		SenderUsername:  p.SenderUsername,
		InvitedUserID:   p.InvitedUserID.String(),
		InvitedUsername: p.InvitedUsername,
	}
}

type PendingInviteDTO struct {
	ID              string `json:"id"`
	GroupID         string `json:"groupId"`
	GroupName       string `json:"groupName"`
	SenderID        string `json:"senderId"`
	SenderUsername  string `json:"senderUsername"`
	InvitedUserID   string `json:"invitedUserId"`
	InvitedUsername string `json:"invitedUsername"`
}

type PendingInvites []*PendingInvite

func (ps PendingInvites) ToDto() []*PendingInviteDTO {
	dtos := make([]*PendingInviteDTO, 0, len(ps))
	for _, p := range ps {
		dtos = append(dtos, p.ToDto())
	}

	return dtos
}
