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
//
// A patch that re-sends an item's current values is still a valid update: it is
// applied and bumps updated_at. We deliberately do not diff against the stored
// values, treating updated_at as the time of the last accepted write rather than
// the last content change. Revisit if a client ever needs it to mean the latter.
func (s *ListService) UpdateItem(ctx context.Context, input UpdateItemInput) (*UpdateItemOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if input.Quantity == nil && input.Unit == nil && input.Note == nil {
		return nil, ErrNoFieldsToUpdate
	}

	listItemRepo := listItemRepository.New(s.db)
	item, updated, err := listItemRepo.UpdateForMember(input.ListID, input.ItemID, userID, listItemRepository.ItemPatch{
		Quantity: input.Quantity,
		Unit:     input.Unit,
		Note:     input.Note,
	})
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
