package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemRepository "sl-api/api/repositories/list_item"
)

type DecreaseQuantityInput struct {
	ListID uuid.UUID
	ItemID uuid.UUID
	By     int
}

type DecreaseQuantityOutput struct {
	ItemID   uuid.UUID `json:"itemId"`
	Quantity int       `json:"quantity"`
}

// DecreaseQuantity lowers a list item's quantity and returns the new value. The
// caller must be a member of the list's group; a non-member, an unknown list,
// and an unknown item are all reported as ErrItemNotFound so the endpoint never
// reveals what it cannot show. A decrease that would take the quantity below 1
// is rejected with ErrQuantityBelowMinimum (removing an item is a separate
// delete).
func (s *ListService) DecreaseQuantity(ctx context.Context, input DecreaseQuantityInput) (*DecreaseQuantityOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	listItemRepo := listItemRepository.New(s.db)
	newQuantity, visible, decremented, err := listItemRepo.DecrementQuantityForMember(input.ListID, input.ItemID, userID, input.By)
	if err != nil {
		return nil, fmt.Errorf("failed to decrease item quantity: %w", err)
	}

	if !visible {
		return nil, ErrItemNotFound
	}

	if !decremented {
		return nil, ErrQuantityBelowMinimum
	}

	return &DecreaseQuantityOutput{
		ItemID:   input.ItemID,
		Quantity: newQuantity,
	}, nil
}
