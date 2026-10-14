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

// updatedItemCols are what UpdateForMember returns via RETURNING.
func updatedItemCols() []string {
	return []string{"id", "list_id", "product_id", "added_by", "quantity", "unit", "status", "note"}
}

func TestUpdateItem_Unauthorized(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	quantity := 2
	_, err = svc.UpdateItem(context.Background(), listService.UpdateItemInput{
		ListID:   uuid.New(),
		ItemID:   uuid.New(),
		Quantity: &quantity,
	})
	if !errors.Is(err, middleware.ErrUnauthorized) {
		t.Fatalf("want ErrUnauthorized, got %v", err)
	}
}

func TestUpdateItem_NoFields(t *testing.T) {
	t.Parallel()

	db, _, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// No field set: rejected before any DB call.
	_, err = svc.UpdateItem(ctxWithUser(uuid.New()), listService.UpdateItemInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
	})
	if !errors.Is(err, listService.ErrNoFieldsToUpdate) {
		t.Fatalf("want ErrNoFieldsToUpdate, got %v", err)
	}
}

func TestUpdateItem_NotFound(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	// No row updated: not a member, unknown list, or unknown item.
	mock.ExpectQuery("UPDATE list_items").WillReturnRows(sqlmock.NewRows(updatedItemCols()))

	note := "milk is low"
	_, err = svc.UpdateItem(ctxWithUser(uuid.New()), listService.UpdateItemInput{
		ListID: uuid.New(),
		ItemID: uuid.New(),
		Note:   &note,
	})
	if !errors.Is(err, listService.ErrItemNotFound) {
		t.Fatalf("want ErrItemNotFound, got %v", err)
	}
}

func TestUpdateItem_Success(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)
	svc := listService.New(db)

	itemID := uuid.New()

	mock.ExpectQuery("UPDATE list_items").WillReturnRows(
		sqlmock.NewRows(updatedItemCols()).
			AddRow(itemID, uuid.New(), uuid.New(), uuid.New(), 3, "L", "pending", "chilled"))

	quantity, unit, note := 3, "L", "chilled"
	out, err := svc.UpdateItem(ctxWithUser(uuid.New()), listService.UpdateItemInput{
		ListID:   uuid.New(),
		ItemID:   itemID,
		Quantity: &quantity,
		Unit:     &unit,
		Note:     &note,
	})
	testUtil.NoError(t, err)
	testUtil.Equal(t, out.ItemID, itemID)
	testUtil.Equal(t, out.Quantity, 3)
	testUtil.Equal(t, out.Unit, "L")
	testUtil.Equal(t, out.Note, "chilled")
	testUtil.Equal(t, out.Status, "pending")
}
