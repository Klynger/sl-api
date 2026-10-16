package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	listRepository "sl-api/api/repositories/list"
)

type UpdateListStatusInput struct {
	ListID uuid.UUID
	Status listModel.Status
}

type UpdateListStatusOutput struct {
	ListID uuid.UUID `json:"listId"`
	Status string    `json:"status"`
}

// UpdateListStatus sets a list's status to any value in the allowed set,
// regardless of its current status (transitions are not ordered). The caller
// must be a member of the list's group; a non-member and an unknown list are
// both reported as ErrListNotFound so the endpoint never reveals a list it
// cannot show.
func (s *ListService) UpdateListStatus(ctx context.Context, input UpdateListStatusInput) (*UpdateListStatusOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	listRepo := listRepository.New(s.db)
	list, updated, err := listRepo.UpdateStatusForMember(input.ListID, userID, input.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update list status: %w", err)
	}

	if !updated {
		return nil, ErrListNotFound
	}

	return &UpdateListStatusOutput{
		ListID: list.ID,
		Status: string(list.Status),
	}, nil
}
