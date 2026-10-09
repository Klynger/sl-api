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

func listRow(listID, groupID uuid.UUID) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "group_id"}).AddRow(listID, groupID)
}

func TestAddItem_ListNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery(`^SELECT .* FROM "lists"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "group_id"})) // no rows

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: uuid.New(), ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrListNotFound) {
		t.Fatalf("want ErrListNotFound, got %v", err)
	}
}

func TestAddItem_NotAMember(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)
	listID, groupID := uuid.New(), uuid.New()

	mock.ExpectQuery(`^SELECT .* FROM "lists"`).WillReturnRows(listRow(listID, groupID))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: listID, ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrNotAMember) {
		t.Fatalf("want ErrNotAMember, got %v", err)
	}
}

func TestAddItem_ProductNotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)
	listID, groupID := uuid.New(), uuid.New()

	mock.ExpectQuery(`^SELECT .* FROM "lists"`).WillReturnRows(listRow(listID, groupID))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`^SELECT .* FROM "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"})) // no rows

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: listID, ProductID: uuid.New()})
	if !errors.Is(err, listService.ErrProductNotFound) {
		t.Fatalf("want ErrProductNotFound, got %v", err)
	}
}

func TestAddItem_AlreadyOnList(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)
	listID, groupID, productID := uuid.New(), uuid.New(), uuid.New()

	mock.ExpectQuery(`^SELECT .* FROM "lists"`).WillReturnRows(listRow(listID, groupID))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`^SELECT .* FROM "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(productID))
	mock.ExpectBegin()
	mock.ExpectExec(`^INSERT INTO "list_items"`).
		WillReturnError(&pgconn.PgError{Code: "23505"})
	mock.ExpectRollback()

	_, err = svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: listID, ProductID: productID})
	if !errors.Is(err, listService.ErrItemAlreadyOnList) {
		t.Fatalf("want ErrItemAlreadyOnList, got %v", err)
	}
}

func TestAddItem_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)
	listID, groupID, productID := uuid.New(), uuid.New(), uuid.New()

	mock.ExpectQuery(`^SELECT .* FROM "lists"`).WillReturnRows(listRow(listID, groupID))
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`^SELECT .* FROM "products"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(productID))
	mock.ExpectBegin()
	mock.ExpectExec(`^INSERT INTO "list_items"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	out, err := svc.AddItem(ctxWithUser(uuid.New()), listService.AddItemInput{ListID: listID, ProductID: productID, Quantity: 2, Unit: "L"})
	testUtil.NoError(t, err)
	if out == nil || out.ItemID == uuid.Nil {
		t.Fatal("expected a non-nil item id")
	}
}
