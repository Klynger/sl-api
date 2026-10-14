package listService

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sl-api/api/middleware"
	listItemRepository "sl-api/api/repositories/list_item"
)

// DeleteItem removes a single item from a list. The caller must be a member of
// the list's group; a non-member, an unknown list, and an unknown item are all
// reported as ErrItemNotFound so the endpoint never reveals what it cannot show.
func (s *ListService) DeleteItem(ctx context.Context, listID, itemID uuid.UUID) error {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return err
	}

	listItemRepo := listItemRepository.New(s.db)
	deleted, err := listItemRepo.DeleteForMember(listID, itemID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	if !deleted {
		return ErrItemNotFound
	}

	return nil
}
