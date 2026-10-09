package listItemRepository_test

import (
	"regexp"
	"testing"

	listItemModel "sl-api/api/model/list_item"
	"sl-api/api/repositories/list_item"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepository_ListByList(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := listItemRepository.New(db)
	listID := uuid.New()

	mockRows := sqlmock.NewRows([]string{"id", "list_id", "product_id", "added_by", "quantity", "unit", "status"}).
		AddRow(uuid.New(), listID, uuid.New(), uuid.New(), 2, "L", listItemModel.StatusPending).
		AddRow(uuid.New(), listID, uuid.New(), uuid.New(), 1, "", listItemModel.StatusInCart)

	mock.ExpectQuery("^SELECT (.+) FROM \"list_items\"").WillReturnRows(mockRows)

	items, err := repo.ListByList(listID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(items), 2)
	testUtil.Equal(t, items[0].Quantity, 2)
	testUtil.Equal(t, items[1].Status, listItemModel.StatusInCart)
}

func TestRepository_FindByListAndProduct(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := listItemRepository.New(db)
	listID := uuid.New()
	productID := uuid.New()

	mockRows := sqlmock.NewRows([]string{"id", "list_id", "product_id", "quantity"}).
		AddRow(uuid.New(), listID, productID, 3)

	mock.ExpectQuery("^SELECT (.+) FROM \"list_items\" WHERE list_id = (.+) AND product_id = (.+)").
		WillReturnRows(mockRows)

	item, err := repo.FindByListAndProduct(listID, productID)
	testUtil.NoError(t, err)
	testUtil.Equal(t, item.ProductID, productID)
	testUtil.Equal(t, item.Quantity, 3)
}

func TestRepository_IncrementQuantity(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := listItemRepository.New(db)
	id := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "list_items" SET "quantity"=quantity + $1`)).
		WithArgs(2, mockDB.AnyTime{}, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	rows, err := repo.IncrementQuantity(id, 2)
	testUtil.NoError(t, err)
	testUtil.Equal(t, rows, int64(1))
	testUtil.NoError(t, mock.ExpectationsWereMet())
}

func TestListItem_EffectiveUnit(t *testing.T) {
	t.Parallel()

	withOwnUnit := &listItemModel.ListItem{Unit: "kg"}
	testUtil.Equal(t, withOwnUnit.EffectiveUnit("L"), "kg")

	inheriting := &listItemModel.ListItem{}
	testUtil.Equal(t, inheriting.EffectiveUnit("L"), "L")

	neither := &listItemModel.ListItem{}
	testUtil.Equal(t, neither.EffectiveUnit(""), "")
}
