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

// decrementCols are what the DecrementQuantityForMember query selects.
func decrementCols() []string {
	return []string{"visible", "new_quantity"}
}

func TestDecreaseQuantity_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	_, err = svc.DecreaseQuantity(context.Background(), listService.DecreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		By:     1,
	})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestDecreaseQuantity_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// Not visible to the caller (not a member / unknown list / unknown item).
	mock.ExpectQuery("new_quantity").WillReturnRows(
		decrementRows().AddRow(false, nil))

	_, err = svc.DecreaseQuantity(ctxWithUser(uuid.New()), listService.DecreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		By:     1,
	})
	if !errors.Is(err, listService.ErrItemNotFound) {
		t.Fatalf("want ErrItemNotFound, got %v", err)
	}
}

func TestDecreaseQuantity_BelowMinimum(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// Visible, but the decrement would breach the floor of 1, so no row updated.
	mock.ExpectQuery("new_quantity").WillReturnRows(
		decrementRows().AddRow(true, nil))

	_, err = svc.DecreaseQuantity(ctxWithUser(uuid.New()), listService.DecreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		By:     5,
	})
	if !errors.Is(err, listService.ErrQuantityBelowMinimum) {
		t.Fatalf("want ErrQuantityBelowMinimum, got %v", err)
	}
}

func TestDecreaseQuantity_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	itemID := uuid.New()

	mock.ExpectQuery("new_quantity").WillReturnRows(
		decrementRows().AddRow(true, 4))

	out, err := svc.DecreaseQuantity(ctxWithUser(uuid.New()), listService.DecreaseQuantityInput{
		ListID: uuid.New(),
		ItemID: itemID,
		By:     2,
	})
	testUtil.NoError(t, err)
	testUtil.Equal(t, out.Quantity, 4)
	testUtil.Equal(t, out.ItemID, itemID)
}

func decrementRows() *sqlmock.Rows {
	return sqlmock.NewRows(decrementCols())
}
