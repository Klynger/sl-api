package listRepository_test

import (
	"testing"

	listModel "sl-api/api/model/list"
	"sl-api/api/repositories/list"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepository_ListByGroup(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := listRepository.New(db)
	groupID := uuid.New()

	mockRows := sqlmock.NewRows([]string{"id", "group_id", "created_by", "name", "status"}).
		AddRow(uuid.New(), groupID, uuid.New(), "Weekly groceries", listModel.StatusPlanning).
		AddRow(uuid.New(), groupID, uuid.New(), "Party supplies", listModel.StatusActive)

	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"").WillReturnRows(mockRows)

	lists, err := repo.ListByGroup(groupID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(lists), 2)
	testUtil.Equal(t, lists[0].Name, "Weekly groceries")
	testUtil.Equal(t, lists[1].Status, listModel.StatusActive)
}

func TestRepository_Read(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := listRepository.New(db)
	id := uuid.New()

	mockRows := sqlmock.NewRows([]string{"id", "group_id", "created_by", "name", "status"}).
		AddRow(id, uuid.New(), uuid.New(), "Weekly groceries", listModel.StatusPlanning)

	mock.ExpectQuery("^SELECT (.+) FROM \"lists\"").WillReturnRows(mockRows)

	list, err := repo.Read(id)
	testUtil.NoError(t, err)
	testUtil.Equal(t, list.ID, id)
	testUtil.Equal(t, list.Status, listModel.StatusPlanning)
}
