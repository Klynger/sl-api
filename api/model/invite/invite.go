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

// PendingInvite is one live invite joined with the group it belongs to and the
// user who sent it, so a recipient has enough to recognize and act on it. It is
// the scan target for the listing queries.
type PendingInvite struct {
	ID             uuid.UUID
	GroupID        uuid.UUID
	GroupName      string
	SenderID       uuid.UUID
	SenderUsername string
}

func (p *PendingInvite) ToDto() *PendingInviteDTO {
	return &PendingInviteDTO{
		ID:             p.ID.String(),
		GroupID:        p.GroupID.String(),
		GroupName:      p.GroupName,
		SenderID:       p.SenderID.String(),
		SenderUsername: p.SenderUsername,
	}
}

type PendingInviteDTO struct {
	ID             string `json:"id"`
	GroupID        string `json:"groupId"`
	GroupName      string `json:"groupName"`
	SenderID       string `json:"senderId"`
	SenderUsername string `json:"senderUsername"`
}

type PendingInvites []*PendingInvite

func (ps PendingInvites) ToDto() []*PendingInviteDTO {
	dtos := make([]*PendingInviteDTO, 0, len(ps))
	for _, p := range ps {
		dtos = append(dtos, p.ToDto())
	}

	return dtos
}
