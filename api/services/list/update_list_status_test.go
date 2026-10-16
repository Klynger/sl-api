package listService_test

import (
	"context"
	"errors"
	"testing"

	"sl-api/api/middleware"
	listModel "sl-api/api/model/list"
	listService "sl-api/api/services/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// listCols are what UpdateStatusForMember returns via RETURNING.
func listCols() []string {
	return []string{"id", "group_id", "created_by", "name", "status"}
}

func TestUpdateListStatus_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	_, err = svc.UpdateListStatus(context.Background(), listService.UpdateListStatusInput{
		ListID: uuid.New(),
		Status: listModel.StatusActive,
	})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestUpdateListStatus_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// No row updated: not a member or unknown list.
	mock.ExpectQuery("UPDATE lists").WillReturnRows(sqlmock.NewRows(listCols()))

	_, err = svc.UpdateListStatus(ctxWithUser(uuid.New()), listService.UpdateListStatusInput{
		ListID: uuid.New(),
		Status: listModel.StatusCompleted,
	})
	if !errors.Is(err, listService.ErrListNotFound) {
		t.Fatalf("want ErrListNotFound, got %v", err)
	}
}

func TestUpdateListStatus_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	listID := uuid.New()

	mock.ExpectQuery("UPDATE lists").WillReturnRows(
		sqlmock.NewRows(listCols()).
			AddRow(listID, uuid.New(), uuid.New(), "Weekly groceries", "active"))

	out, err := svc.UpdateListStatus(ctxWithUser(uuid.New()), listService.UpdateListStatusInput{
		ListID: listID,
		Status: listModel.StatusActive,
	})
	testUtil.NoError(t, err)
	testUtil.Equal(t, out.ListID, listID)
	testUtil.Equal(t, out.Status, "active")
}
