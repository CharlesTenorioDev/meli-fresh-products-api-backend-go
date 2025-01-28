package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductTypeMysql_FindByID_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"description",
	}).
		AddRow(1, "product type")
	mock.ExpectQuery(repository.FindByIDProductType).WillReturnRows(row)

	repo := repository.NewProductTypeMysql(mockDB)

	_, err = repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
}

func TestProductTypeMysql_FindByID_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"description",
	})

	mock.ExpectQuery(repository.FindByIDProductType).WithArgs(1).WillReturnRows(row)

	repo := repository.NewProductTypeMysql(mockDB)

	_, err = repo.FindByID(1)
	assert.Equal(t, internal.ErrProductTypeNotFound, err)
	assert.Equal(t, 1, product.ID)
}
