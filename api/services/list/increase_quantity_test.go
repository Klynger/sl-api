package listService_test

import (
	"context"
	"errors"
	"testing"

	"sl-api/api/middleware"
	listService "sl-api/api/services/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestIncreaseQuantity_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	_, err = svc.IncreaseQuantity(context.Background(), listService.IncreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		By:     1,
	})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestIncreaseQuantity_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// No row returned: the item does not exist, or the caller is not a member of
	// its list's group. Both collapse to ErrItemNotFound.
	mock.ExpectQuery("UPDATE list_items").WillReturnRows(sqlmock.NewRows([]string{"quantity"}))

	_, err = svc.IncreaseQuantity(ctxWithUser(uuid.New()), listService.IncreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		By:     1,
	})
	if !errors.Is(err, listService.ErrItemNotFound) {
		t.Fatalf("want ErrItemNotFound, got %v", err)
	}
}

func TestIncreaseQuantity_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	itemID := uuid.New()

	// RETURNING yields the new quantity after the atomic increment.
	mock.ExpectQuery("UPDATE list_items").WillReturnRows(
		sqlmock.NewRows([]string{"quantity"}).AddRow(5))

	out, err := svc.IncreaseQuantity(ctxWithUser(uuid.New()), listService.IncreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: itemID,
		By:     3,
	})
	testUtil.NoError(t, err)
	testUtil.Equal(t, out.Quantity, 5)
	testUtil.Equal(t, out.ItemID, itemID)
}
