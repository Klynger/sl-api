package listService_test

import (
	"errors"
	"testing"

	listService "sl-api/api/services/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

// addItemData mocks the single combined list+product lookup. Pass nil for a
// column to simulate that row being absent (NULL).
func addItemData(listGroupID, productID any) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"list_group_id", "product_id"}).AddRow(listGroupID, productID)
}

func memberCount(n int) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"count"}).AddRow(n)
}

func TestAddItem_ListNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("list_group_id").WillReturnRows(addItemData(nil, uuid.New()))

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrListNotFound) {
		t.Fatalf("want ErrListNotFound, got %v", err)
	}
}

func TestAddItem_ProductNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("list_group_id").WillReturnRows(addItemData(uuid.New(), nil))

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrProductNotFound) {
		t.Fatalf("want ErrProductNotFound, got %v", err)
	}
}

func TestAddItem_NotAMember(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("list_group_id").WillReturnRows(addItemData(uuid.New(), uuid.New()))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).WillReturnRows(memberCount(0))

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrNotAMember) {
		t.Fatalf("want ErrNotAMember, got %v", err)
	}
}

func TestAddItem_AlreadyOnList(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("list_group_id").WillReturnRows(addItemData(uuid.New(), uuid.New()))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).WillReturnRows(memberCount(1))
	mock.ExpectBegin()
	mock.ExpectExec(`^INSERT INTO "list_items"`).WillReturnError(&pgconn.PgError{Code: "23505"})
	mock.ExpectRollback()

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrItemAlreadyOnList) {
		t.Fatalf("want ErrItemAlreadyOnList, got %v", err)
	}
}

func TestAddItem_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery("list_group_id").WillReturnRows(addItemData(uuid.New(), uuid.New()))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).WillReturnRows(memberCount(1))
	mock.ExpectBegin()
	mock.ExpectExec(`^INSERT INTO "list_items"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	out, err := svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New(), Quantity: 2, Unit: "L"})
	testUtil.NoError(t, err)
	if out == nil || out.ItemID == uuid.Nil {
		t.Fatal("expected a non-nil item id")
	}
}
