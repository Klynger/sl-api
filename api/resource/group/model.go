package group

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Group struct {
	ID        uuid.UUID `gorm:"primarykey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (g *Group) ToDto() *DTO {
	return &DTO{
		ID:   g.ID.String(),
		Name: g.Name,
	}
}

type DTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Form struct {
	Name string `json:"name" validate:"required,max=255"`
}

func (f *Form) ToModel() *Group {
	return &Group{
		Name: f.Name,
	}
}
