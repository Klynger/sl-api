package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	listItemRepository "sl-api/api/repositories/list_item"
)

type UpdateItemStatusInput struct {
	ListID uuid.UUID
	ItemID uuid.UUID
	Status listItemModel.Status
}

type UpdateItemStatusOutput struct {
	ItemID uuid.UUID `json:"itemId"`
	Status string    `json:"status"`
}

// UpdateItemStatus sets an item's status to any value in the allowed set,
// regardless of its current status (transitions are not ordered). The caller
// must be a member of the list's group; a non-member, an unknown list, and an
// unknown item are all reported as ErrItemNotFound so the endpoint never reveals
// what it cannot show.
func (s *ListService) UpdateItemStatus(ctx context.Context, input UpdateItemStatusInput) (*UpdateItemStatusOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	status := input.Status
	listItemRepo := listItemRepository.New(s.db)
	item, updated, err := listItemRepo.UpdateForMember(input.ListID, input.ItemID, userID, listItemRepository.ItemPatch{
		Status: &status,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update item status: %w", err)
	}

	if !updated {
		return nil, ErrItemNotFound
	}

	return &UpdateItemStatusOutput{
		ItemID: item.ID,
		Status: string(item.Status),
	}, nil
}
