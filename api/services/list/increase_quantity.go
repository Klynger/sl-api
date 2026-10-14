package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemRepository "sl-api/api/repositories/list_item"
)

type IncreaseQuantityInput struct {
	ListID uuid.UUID
	ItemID uuid.UUID
	By     int
}

type IncreaseQuantityOutput struct {
	ItemID   uuid.UUID `json:"itemId"`
	Quantity int       `json:"quantity"`
}

// IncreaseQuantity adds to a list item's quantity and returns the new value. The
// caller must be a member of the list's group; a non-member, an unknown list,
// and an unknown item are all reported as ErrItemNotFound so the endpoint never
// reveals what it cannot show.
func (s *ListService) IncreaseQuantity(ctx context.Context, input IncreaseQuantityInput) (*IncreaseQuantityOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	listItemRepo := listItemRepository.New(s.db)
	newQuantity, updated, err := listItemRepo.IncrementQuantityForMember(input.ListID, input.ItemID, userID, input.By)
	if err != nil {
		return nil, fmt.Errorf("failed to increase item quantity: %w", err)
	}

	if !updated {
		return nil, ErrItemNotFound
	}

	return &IncreaseQuantityOutput{
		ItemID:   input.ItemID,
		Quantity: newQuantity,
	}, nil
}
