package listService_test

import (
	"context"
	"errors"
	"testing"

	"sl-api/api/middleware"
	listItemModel "sl-api/api/model/list_item"
	listService "sl-api/api/services/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestUpdateItemStatus_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	_, err = svc.UpdateItemStatus(context.Background(), listService.UpdateItemStatusInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		Status: listItemModel.StatusInCart,
	})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestUpdateItemStatus_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// No row updated: not a member, unknown list, or unknown item.
	mock.ExpectQuery("UPDATE list_items").WillReturnRows(sqlmock.NewRows(updatedItemCols()))

	_, err = svc.UpdateItemStatus(ctxWithUser(uuid.New()), listService.UpdateItemStatusInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		Status: listItemModel.StatusCancelled,
	})
	if !errors.Is(err, listService.ErrItemNotFound) {
		t.Fatalf("want ErrItemNotFound, got %v", err)
	}
}

func TestUpdateItemStatus_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	itemID := uuid.New()

	mock.ExpectQuery("UPDATE list_items").WillReturnRows(
		sqlmock.NewRows(updatedItemCols()).
			AddRow(itemID, uuid.New(), uuid.New(), uuid.New(), 2, "L", "in_cart", ""))

	out, err := svc.UpdateItemStatus(ctxWithUser(uuid.New()), listService.UpdateItemStatusInput{
		ListID: uuid.New(),
		ItemID: itemID,
		Status: listItemModel.StatusInCart,
	})
	testUtil.NoError(t, err)
	testUtil.Equal(t, out.ItemID, itemID)
	testUtil.Equal(t, out.Status, "in_cart")
}
