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

func TestListLists_NotMember(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).WillReturnRows(memberCount(0))

	_, err = svc.ListLists(ctxWithUser(uuid.New()), uuid.New())
	if !errors.Is(err, listService.ErrNotAMember) {
		t.Fatalf("want ErrNotAMember, got %v", err)
	}
}

func TestListLists_Member(t *testing.T) {
	t.Parallel()
	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)
	groupID := uuid.New()

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).WillReturnRows(memberCount(1))
	mock.ExpectQuery(`^SELECT (.+) FROM "lists"`).WillReturnRows(
		sqlmock.NewRows([]string{"id", "group_id", "name", "status"}).
			AddRow(uuid.New(), groupID, "Weekly groceries", "planning").
			AddRow(uuid.New(), groupID, "Party", "active"))

	lists, err := svc.ListLists(ctxWithUser(uuid.New()), groupID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(lists), 2)
	testUtil.Equal(t, lists[0].Name, "Weekly groceries")
}
