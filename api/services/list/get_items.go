package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	listItemRepository "sl-api/api/repositories/list_item"
)

// ListItems returns a list's items (with product names). The caller must be a
// member of the list's group; a non-member or a missing list is reported as
// ErrListNotFound so the endpoint never reveals a list the caller may not see.
func (s *ListService) ListItems(ctx context.Context, listID uuid.UUID) (listItemModel.ItemsWithProduct, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	listItemRepo := listItemRepository.New(s.db)
	items, authorized, err := listItemRepo.ListForMember(listID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list items: %w", err)
	}

	if !authorized {
		return nil, ErrListNotFound
	}

	return items, nil
}
