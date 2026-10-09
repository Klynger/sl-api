package listService

import (
	"context"
	"errors"
	"fmt"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	groupMemberRepository "sl-api/api/repositories/group_member"
	listRepository "sl-api/api/repositories/list"
	listItemRepository "sl-api/api/repositories/list_item"
	productRepository "sl-api/api/repositories/product"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AddItemInput struct {
	ListID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Unit      string
	Note      string
}

type AddItemOutput struct {
	ItemID uuid.UUID
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
	list, err := listRepo.Read(input.ListID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrListNotFound
		}
		return nil, fmt.Errorf("failed to read list: %w", err)
	}

	memberRepo := groupMemberRepository.New(s.db)
	isMember, err := memberRepo.ExistsByUserAndGroup(userID, list.GroupID)
	if err != nil {
		return nil, fmt.Errorf("failed to check group membership: %w", err)
	}

	if !isMember {
		return nil, ErrNotAMember
	}

	productRepo := productRepository.New(s.db)
	product, err := productRepo.Read(input.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to read product: %w", err)
	}

	listItemRepo := listItemRepository.New(s.db)
	listItem, err := listItemRepo.Create(&listItemModel.ListItem{
		ID:        uuid.New(),
		ListID:    list.ID,
		ProductID: product.ID,
		AddedBy:   userID,
		Quantity:  input.Quantity,
		Note:      input.Note,
		Unit:      input.Unit,
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
