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

func ctxWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
}

func TestCreateList_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	svc := listService.New(db)
	_, err = svc.CreateList(context.Background(), listService.CreateListInput{GroupID: uuid.New(), Name: "X"})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestCreateList_NotAMember(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	svc := listService.New(db)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	_, err = svc.CreateList(ctxWithUser(uuid.New()), listService.CreateListInput{GroupID: uuid.New(), Name: "X"})
	if !errors.Is(err, listService.ErrNotAMember) {
		t.Fatalf("want ErrNotAMember, got %v", err)
	}
}

func TestCreateList_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	svc := listService.New(db)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectBegin()
	mock.ExpectExec(`^INSERT INTO "lists"`).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	out, err := svc.CreateList(ctxWithUser(uuid.New()), listService.CreateListInput{GroupID: uuid.New(), Name: "Groceries"})
	testUtil.NoError(t, err)
	if out == nil || out.ListID == uuid.Nil {
		t.Fatal("expected a non-nil list id")
	}
}
