package groupMemberRepository_test

import (
	"testing"

	"sl-api/api/repositories/group_member"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepository_ExistsByUserAndGroup(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := groupMemberRepository.New(db)
	userID := uuid.New()
	groupID := uuid.New()

	// A live membership row makes the count 1, so the check reports true.
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WithArgs(userID, groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	ok, err := repo.ExistsByUserAndGroup(userID, groupID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, ok, true)

	// No matching row means count 0, so the check reports false (not an error).
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "group_members"`).
		WithArgs(userID, groupID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	ok, err = repo.ExistsByUserAndGroup(userID, groupID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, ok, false)

	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_SharesGroup(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := groupMemberRepository.New(db)
	userA := uuid.New()
	userB := uuid.New()

	// A common group makes the EXISTS subquery true.
	mock.ExpectQuery("EXISTS").
		WithArgs(userA, userB).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	ok, err := repo.SharesGroup(userA, userB)
	testUtil.NoError(t, err)
	testUtil.Equal(t, ok, true)

	// No common group means the pair shares nothing.
	mock.ExpectQuery("EXISTS").
		WithArgs(userA, userB).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	ok, err = repo.SharesGroup(userA, userB)
	testUtil.NoError(t, err)
	testUtil.Equal(t, ok, false)

	testUtil.NoError(t, mock.ExpectationsWereMet())
}
