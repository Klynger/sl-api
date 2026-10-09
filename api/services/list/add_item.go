package listService

import (
	"context"
	"errors"
	"fmt"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	listRepository "sl-api/api/repositories/list"
	listItemRepository "sl-api/api/repositories/list_item"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type AddItemInput struct {
	ListID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Unit      string
	Note      string
}

type AddItemOutput struct {
	ItemID uuid.UUID `json:"itemId"`
}

// AddItem adds a product to a list. The caller must be a member of the list's
// group. If the product is already on the list the add fails (changing the
// quantity of an existing item is a separate operation).
func (s *ListService) AddItem(ctx context.Context, input AddItemInput) (*AddItemOutput, error) {
	userID, err := middleware.GetAuthedUserIDFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	listRepo := listRepository.New(s.db)

	// One round-trip confirms the list and the product exist and yields the
	// list's group for the membership check below.
	data, err := listRepo.GetAddItemData(input.ListID, input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to read list and product: %w", err)
	}

	if !data.ListGroupID.Valid {
		return nil, ErrListNotFound
	}

	if !data.ProductID.Valid {
		return nil, ErrProductNotFound
	}

	if err := s.ensureGroupMember(userID, data.ListGroupID.UUID); err != nil {
		return nil, err
	}

	listItemRepo := listItemRepository.New(s.db)
	listItem, err := listItemRepo.Create(&listItemModel.ListItem{
		ID:        uuid.New(),
		ListID:    input.ListID,
		ProductID: input.ProductID,
		AddedBy:   userID,
		Quantity:  input.Quantity,
		Unit:      input.Unit,
		Note:      input.Note,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrItemAlreadyOnList
		}
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return &AddItemOutput{
		ItemID: listItem.ID,
	}, nil
}
