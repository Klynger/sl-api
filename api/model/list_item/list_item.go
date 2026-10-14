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

// IncrementForm is the body for increasing an item's quantity. By is a pointer
// so an omitted field (increment by the default of 1) is distinguishable from an
// explicit value, which must be at least 1.
type IncrementForm struct {
	By *int `json:"by" validate:"omitempty,min=1"`
}

// ItemWithProduct is a list item joined with its product's name, for read
// endpoints where the client wants the product label inline rather than just an id.
type ItemWithProduct struct {
	ID          uuid.UUID
	ProductID   uuid.UUID
	ProductName string
	AddedBy     uuid.UUID
	Quantity    int
	Unit        string
	Status      Status
	Note        string
}

type ItemWithProductDTO struct {
	ID          string `json:"id"`
	ProductID   string `json:"productId"`
	ProductName string `json:"productName"`
	AddedBy     string `json:"addedBy"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	Status      string `json:"status"`
	Note        string `json:"note"`
}

func (i *ItemWithProduct) ToDto() *ItemWithProductDTO {
	return &ItemWithProductDTO{
		ID:          i.ID.String(),
		ProductID:   i.ProductID.String(),
		ProductName: i.ProductName,
		AddedBy:     i.AddedBy.String(),
		Quantity:    i.Quantity,
		Unit:        i.Unit,
		Status:      string(i.Status),
		Note:        i.Note,
	}
}

type ItemsWithProduct []*ItemWithProduct

func (is ItemsWithProduct) ToDto() []*ItemWithProductDTO {
	dtos := make([]*ItemWithProductDTO, len(is))
	for i, v := range is {
		dtos[i] = v.ToDto()
	}

	return dtos
}
