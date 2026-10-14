package listService_test

import (
	"errors"
	"testing"

	listService "sl-api/api/services/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// itemCols are what the ListForMember query selects, in order.
func itemCols() []string {
	return []string{"id", "product_id", "product_name", "added_by", "quantity", "unit", "status", "note"}
}

func TestListItems_NotAuthorized(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("product_name").WillReturnRows(sqlmock.NewRows(itemCols())) // zero rows

	_, err = svc.ListItems(ctxWithUser(uuid.New()), uuid.New())
	if !errors.Is(err, listService.ErrListNotFound) {
		t.Fatalf("want ErrListNotFound, got %v", err)
	}
}

func TestListItems_AuthorizedEmpty(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// One row with all item columns NULL = member, empty list.
	mock.ExpectQuery("product_name").WillReturnRows(
		sqlmock.NewRows(itemCols()).AddRow(nil, nil, nil, nil, nil, nil, nil, nil))

	items, err := svc.ListItems(ctxWithUser(uuid.New()), uuid.New())
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(items), 0)
}

func TestListItems_WithItems(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("product_name").WillReturnRows(
		sqlmock.NewRows(itemCols()).
			AddRow(uuid.New(), uuid.New(), "Milk", uuid.New(), 2, "L", "pending", "").
			AddRow(uuid.New(), uuid.New(), "Eggs", uuid.New(), 1, "", "in_cart", "dozen"))

	items, err := svc.ListItems(ctxWithUser(uuid.New()), uuid.New())
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(items), 2)
	testUtil.Equal(t, items[0].ProductName, "Milk")
	testUtil.Equal(t, items[1].Quantity, 1)
}
