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

func TestDeleteItem_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	err = svc.DeleteItem(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestDeleteItem_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// Zero rows affected: the item does not exist, or the caller is not a member
	// of its list's group. Both collapse to ErrItemNotFound.
	mock.ExpectExec("DELETE FROM list_items").WillReturnResult(sqlmock.NewResult(0, 0))

	err = svc.DeleteItem(ctxWithUser(uuid.New()), uuid.New(), uuid.New())
	if !errors.Is(err, listService.ErrItemNotFound) {
		t.Fatalf("want ErrItemNotFound, got %v", err)
	}
}

func TestDeleteItem_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectExec("DELETE FROM list_items").WillReturnResult(sqlmock.NewResult(0, 1))

	err = svc.DeleteItem(ctxWithUser(uuid.New()), uuid.New(), uuid.New())
	testUtil.NoError(t, err)
}
