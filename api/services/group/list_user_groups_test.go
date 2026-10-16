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

// memberCols are the columns the group-membership detail queries select.
func memberCols() []string {
	return []string{"id", "user_id", "group_id", "roles", "user_name", "user_last_name", "group_name"}
}

func TestListUserGroups_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	_, err = svc.ListUserGroups(context.Background())
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestListUserGroups_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	mock.ExpectQuery("group_name").WillReturnRows(sqlmock.NewRows(memberCols()))

	groups, err := svc.ListUserGroups(ctxWithUser(uuid.New()))
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(groups), 0)
}

func TestListUserGroups_WithGroups(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := groupService.New(db)

	mock.ExpectQuery("group_name").WillReturnRows(
		sqlmock.NewRows(memberCols()).
			AddRow(uuid.New(), uuid.New(), uuid.New(), "{owner}", "Dev", "User", "Home").
			AddRow(uuid.New(), uuid.New(), uuid.New(), "{member}", "Dev", "User", "Work"))

	groups, err := svc.ListUserGroups(ctxWithUser(uuid.New()))
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(groups), 2)
	testUtil.Equal(t, groups[0].GroupName, "Home")
	testUtil.Equal(t, string(groups[0].Roles[0]), "owner")
}
