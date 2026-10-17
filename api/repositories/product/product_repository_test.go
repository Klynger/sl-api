package productRepository_test

import (
	"testing"

	productModel "sl-api/api/model/product"
	"sl-api/api/repositories/product"
	mockDB "sl-api/mock/db"
	testUtil "sl-api/util/test"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func TestRepository_List(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := productRepository.New(db)

	mockRows := sqlmock.NewRows([]string{"id", "product_name", "description"}).
		AddRow(uuid.New(), "Product 1", "Product description 2").
		AddRow(uuid.New(), "Product 2", "Product description 2")

	mock.ExpectQuery("^SELECT (.+) FROM \"products\"").WillReturnRows(mockRows)

	products, err := repo.List()
	testUtil.NoError(t, err)
	testUtil.Equal(t, len(products), 2)
}

func TestRepository_Update(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := productRepository.New(db)

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"products\"").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rows, err := repo.Update(&productModel.Product{ID: uuid.New(), Name: "Updated"})
	testUtil.NoError(t, err)
	testUtil.Equal(t, rows, int64(1))
}

func TestRepository_Delete(t *testing.T) {
	t.Parallel()

	db, mock, err := mockDB.NewMockDB()
	testUtil.NoError(t, err)

	repo := productRepository.New(db)

	// Products are soft deleted, so Delete issues an UPDATE of deleted_at.
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE \"products\" SET \"deleted_at\"").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	rows, err := repo.Delete(uuid.New())
	testUtil.NoError(t, err)
	testUtil.Equal(t, rows, int64(1))
}
