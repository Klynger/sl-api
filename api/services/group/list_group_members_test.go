package groupService_test

import (
	"context"
	"errors"
	"testing"

	"sl-api/api/middleware"
	groupService "sl-api/api/services/group"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestListGroupMembers_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	_, err = svc.ListGroupMembers(context.Background(), uuid.New())
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestListGroupMembers_NotMember(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	_, err = svc.ListGroupMembers(ctxWithUser(uuid.New()), uuid.New())
	// A non-member is reported as a missing group, not revealed as a non-member.
	if !errors.Is(err, groupService.ErrGroupNotFound) {
		t.Fatalf("want ErrGroupNotFound, got %v", err)
	}
}

func TestListGroupMembers_Member(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery("group_name").WillReturnRows(
		sqlmock.NewRows(memberCols()).
			AddRow(uuid.New(), uuid.New(), uuid.New(), "{owner}", "Dev", "User", "Home").
			AddRow(uuid.New(), uuid.New(), uuid.New(), "{member}", "Fri", "End", "Home"))

	members, err := svc.ListGroupMembers(ctxWithUser(uuid.New()), uuid.New())
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(members), 2)
	testUtil.Equal(t, members[1].UserName, "Fri")
	testUtil.Equal(t, string(members[1].Roles[0]), "member")
}
