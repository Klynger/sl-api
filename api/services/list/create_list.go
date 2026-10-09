package listService

import (
	"context"
	"fmt"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	listRepository "sl-api/api/repositories/list"

	"github.com/google/uuid"
)

type CreateListInput struct {
	GroupID uuid.UUID
	Name    string
	// CreatedBy is resolved from the authenticated user in ctx, not the request body.
}

type CreateListOutput struct {
	ListID uuid.UUID `json:"listId"`
}

// CreateList creates a new list in a group. The caller must be a member of that group.
func (s *ListService) CreateList(ctx context.Context, input CreateListInput) (*CreateListOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.ensureGroupMember(userID, input.GroupID); err != nil {
		return nil, err
	}

	listRepo := listRepository.New(s.db)
	newList, err := listRepo.Create(&listModel.List{
		ID:        uuid.New(),
		GroupID:   input.GroupID,
		CreatedBy: userID,
		Name:      input.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return &CreateListOutput{ListID: newList.ID}, nil
}
