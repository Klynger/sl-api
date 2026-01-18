package product

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID `gorm:"primarykey"`
	Name        string    `gorm:"column:product_name"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (p *Product) ToDto() *DTO {
	return &DTO{
		ID:          p.ID.String(),
		Name:        p.Name,
		Description: p.Description,
	}
}

type DTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Form struct {
	Name        string `json:"name" validate:"required,max=255"`
	Description string `json:"description"`
}

func (f *Form) ToModel() *Product {

	return &Product{
		Name:        f.Name,
		Description: f.Description,
	}
}

type Products []*Product

func (ps Products) ToDto() []*DTO {
	dtos := make([]*DTO, len(ps))
	for i, v := range ps {
		dtos[i] = v.ToDto()
	}

	return dtos
}
