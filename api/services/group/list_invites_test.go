package groupService_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"sl-api/api/middleware"
	groupService "sl-api/api/services/group"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestListUserInvites_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	svc := groupService.New(db)
	_, err = svc.ListUserInvites(context.Background())
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestListUserInvites_Empty(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mock.ExpectQuery("sender_username").WillReturnRows(sqlmock.NewRows(inviteCols()))

	svc := groupService.New(db)
	invites, err := svc.ListUserInvites(ctxWithUser(uuid.New()))
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(invites), 0)
}

func TestListUserInvites_WithInvites(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	userID := uuid.New()

	mock.ExpectQuery("sender_username").WillReturnRows(
		sqlmock.NewRows(inviteCols()).
			AddRow(uuid.New(), uuid.New(), "Home", uuid.New(), "alice", userID, "carol").
			AddRow(uuid.New(), uuid.New(), "Work", uuid.New(), "bob", userID, "carol"))

	svc := groupService.New(db)
	invites, err := svc.ListUserInvites(ctxWithUser(userID))
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(invites), 2)
	testUtil.Equal(t, invites[0].GroupName, "Home")
	testUtil.Equal(t, invites[0].SenderUsername, "alice")
	testUtil.Equal(t, invites[0].InvitedUserID, userID)
	testUtil.Equal(t, invites[0].InvitedUsername, "carol")
}

func TestListGroupInvites_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	svc := groupService.New(db)
	_, err = svc.ListGroupInvites(context.Background(), uuid.New())
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestListGroupInvites_NotMember(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	svc := groupService.New(db)
	_, err = svc.ListGroupInvites(ctxWithUser(uuid.New()), uuid.New())
	// A non-member is reported as a missing group, not revealed as a non-member.
	if !errors.Is(err, groupService.ErrGroupNotFound) {
		t.Fatalf("want ErrGroupNotFound, got %v", err)
	}
}

func TestListGroupInvites_Member(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	invitedUserID := uuid.New()
	mock.ExpectQuery(regexp.QuoteMeta(`INNER JOIN users AS invited ON invited.id = invites.invited_user_id`)).WillReturnRows(
		sqlmock.NewRows(inviteCols()).
			AddRow(uuid.New(), uuid.New(), "Home", uuid.New(), "alice", invitedUserID, "carol"))

	svc := groupService.New(db)
	invites, err := svc.ListGroupInvites(ctxWithUser(uuid.New()), uuid.New())
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(invites), 1)
	testUtil.Equal(t, invites[0].SenderUsername, "alice")
	// A group's listing names who each invite is waiting on, not just who sent it.
	testUtil.Equal(t, invites[0].InvitedUserID, invitedUserID)
	testUtil.Equal(t, invites[0].InvitedUsername, "carol")

	dto := invites.ToDto()[0]
	testUtil.Equal(t, dto.InvitedUserID, invitedUserID.String())
	testUtil.Equal(t, dto.InvitedUsername, "carol")
}

// inviteCols are the columns the pending-invite listing queries select, in order.
func inviteCols() []string {
	return []string{"id", "group_id", "group_name", "sender_id", "sender_username", "invited_user_id", "invited_username"}
}

func ctxWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID.String())
}
