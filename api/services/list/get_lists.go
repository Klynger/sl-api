package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	listRepository "sl-api/api/repositories/list"
)

// ListLists returns the lists in a group. The caller must be a member.
func (s *ListService) ListLists(ctx context.Context, groupID uuid.UUID) (listModel.Lists, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if err := s.ensureGroupMember(userID, groupID); err != nil {
		return nil, err
	}

	listRepo := listRepository.New(s.db)
	lists, err := listRepo.ListByGroup(groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to list group lists: %w", err)
	}

	return lists, nil
}
