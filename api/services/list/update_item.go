package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemRepository "sl-api/api/repositories/list_item"
)

type UpdateItemInput struct {
	ListID   uuid.UUID
	ItemID   uuid.UUID
	Quantity *int
	Unit     *string
	Note     *string
}

type UpdateItemOutput struct {
	ItemID   uuid.UUID `json:"itemId"`
	Quantity int       `json:"quantity"`
	Unit     string    `json:"unit"`
	Status   string    `json:"status"`
	Note     string    `json:"note"`
}

// UpdateItem changes an item's mutable fields (quantity, unit, note) and returns
// its new state. Only the fields present in the input are changed; omitting all
// of them is ErrNoFieldsToUpdate. The caller must be a member of the list's
// group; a non-member, an unknown list, and an unknown item are all reported as
// ErrItemNotFound so the endpoint never reveals what it cannot show.
func (s *ListService) UpdateItem(ctx context.Context, input UpdateItemInput) (*UpdateItemOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	fields := map[string]any{}
	if input.Quantity != nil {
		fields["quantity"] = *input.Quantity
	}
	if input.Unit != nil {
		fields["unit"] = *input.Unit
	}
	if input.Note != nil {
		fields["note"] = *input.Note
	}

	if len(fields) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	listItemRepo := listItemRepository.New(s.db)
	item, updated, err := listItemRepo.UpdateForMember(input.ListID, input.ItemID, userID, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	if !updated {
		return nil, ErrItemNotFound
	}

	return &UpdateItemOutput{
		ItemID:   item.ID,
		Quantity: item.Quantity,
		Unit:     item.Unit,
		Status:   string(item.Status),
		Note:     item.Note,
	}, nil
}
