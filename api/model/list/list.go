package listModel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	StatusPlanning  Status = "planning"
	StatusActive    Status = "active"
	StatusCompleted Status = "completed"
)

type List struct {
	ID        uuid.UUID `gorm:"primarykey"`
	GroupID   uuid.UUID `gorm:"column:group_id"`
	CreatedBy uuid.UUID `gorm:"column:created_by"`
	Name      string
	Status    Status `gorm:"type:text;default:planning"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (l *List) ToDto() *DTO {
	return &DTO{
		ID:        l.ID.String(),
		GroupID:   l.GroupID.String(),
		CreatedBy: l.CreatedBy.String(),
		Name:      l.Name,
		Status:    string(l.Status),
	}
}

type DTO struct {
	ID        string `json:"id"`
	GroupID   string `json:"groupId"`
	CreatedBy string `json:"createdBy"`
	Name      string `json:"name"`
	Status    string `json:"status"`
}

type Lists []*List

func (ls Lists) ToDto() []*DTO {
	dtos := make([]*DTO, len(ls))
	for i, v := range ls {
		dtos[i] = v.ToDto()
	}

	return dtos
}
