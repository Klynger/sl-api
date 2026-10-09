package listService

import (
	"context"
	"fmt"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	groupMemberRepository "sl-api/api/repositories/group_member"
	listRepository "sl-api/api/repositories/list"

	"github.com/google/uuid"
)

type CreateListInput struct {
	GroupID uuid.UUID
	Name    string
	// CreatedBy is resolved from the authenticated user in ctx, not the request body.
}

type CreateListOutput struct {
	ListID uuid.UUID
}

// CreateList creates a new list in a group. The caller must be a member of that group.
func (s *ListService) CreateList(ctx context.Context, input CreateListInput) (*CreateListOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	memberRepo := groupMemberRepository.New(s.db)
	isMember, err := memberRepo.ExistsByUserAndGroup(userID, input.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}

	if !isMember {
		return nil, ErrNotAMember
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
