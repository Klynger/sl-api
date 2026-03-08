package productRepository_test

import (
	"testing"

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
