package listService

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	listItemModel "sl-api/api/model/list_item"
	productModel "sl-api/api/model/product"
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
	productRepo := productRepository.New(s.db)

	var wg sync.WaitGroup
	wg.Add(2)

	var list *listModel.List
	var listErr error
	var product *productModel.Product
	var productErr error

	go func() {
		defer wg.Done()
		list, listErr = listRepo.Read(input.ListID)
	}()

	go func() {
		defer wg.Done()
		product, productErr = productRepo.Read(input.ProductID)
	}()

	wg.Wait()

	if listErr != nil {
		if errors.Is(listErr, gorm.ErrRecordNotFound) {
			return nil, ErrListNotFound
		}
		return nil, fmt.Errorf("failed to read list: %w", listErr)
	}

	if err := s.ensureGroupMember(userID, list.GroupID); err != nil {
		return nil, err
	}

	if productErr != nil {
		if errors.Is(productErr, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to read product: %w", productErr)
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
