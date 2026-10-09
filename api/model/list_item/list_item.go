package listItemModel

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusInCart    Status = "in_cart"
	StatusPurchased Status = "purchased"
)

type ListItem struct {
	ID        uuid.UUID `gorm:"primarykey"`
	ListID    uuid.UUID `gorm:"column:list_id"`
	ProductID uuid.UUID `gorm:"column:product_id"`
	AddedBy   uuid.UUID `gorm:"column:added_by"`
	Quantity  int       `gorm:"default:1"`
	Unit      string
	Status    Status `gorm:"type:text;default:pending"`
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (ListItem) TableName() string {
	return "list_items"
}

// EffectiveUnit resolves which unit to show for this item: its own unit when
// set, otherwise the unit predefined on the product it references.
func (i *ListItem) EffectiveUnit(productPredefinedUnit string) string {
	if i.Unit != "" {
		return i.Unit
	}

	return productPredefinedUnit
}

func (i *ListItem) ToDto() *DTO {
	return &DTO{
		ID:        i.ID.String(),
		ListID:    i.ListID.String(),
		ProductID: i.ProductID.String(),
		AddedBy:   i.AddedBy.String(),
		Quantity:  i.Quantity,
		Unit:      i.Unit,
		Status:    string(i.Status),
		Note:      i.Note,
	}
}

type DTO struct {
	ID        string `json:"id"`
	ListID    string `json:"listId"`
	ProductID string `json:"productId"`
	AddedBy   string `json:"addedBy"`
	Quantity  int    `json:"quantity"`
	Unit      string `json:"unit"`
	Status    string `json:"status"`
	Note      string `json:"note"`
}

type ListItems []*ListItem

func (is ListItems) ToDto() []*DTO {
	dtos := make([]*DTO, len(is))
	for i, v := range is {
		dtos[i] = v.ToDto()
	}

	return dtos
}

type Form struct {
	ProductID string `json:"productId" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"omitempty,min=1"`
	Unit      string `json:"unit" validate:"omitempty,max=50"`
	Note      string `json:"note"`
}
