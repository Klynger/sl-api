package invite

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
