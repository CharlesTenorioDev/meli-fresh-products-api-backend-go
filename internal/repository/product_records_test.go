package repository_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/assert"
)

var productRecords = internal.ProductRecords{
	ID:             1,
	LastUpdateDate: time.Now().Truncate(24 * time.Hour),
	PurchasePrice:  1,
	SalePrice:      1,
	ProductID:      1,
}

func TestProductRecordsMysql_FinAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"last_update_date",
		"purchase_price",
		"sale_price",
		"product_id",
	}).
		AddRow(1, time.Now().Truncate(24*time.Hour), 1, 1, 1).
		AddRow(2, time.Now().Truncate(24*time.Hour), 2, 2, 2)

	mock.ExpectQuery(repository.FindAllProductRecords).WillReturnRows(rows)

	repo := repository.NewProductRecordsSQL(mockDB)

	products, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, 1, products[0].ID)
	assert.Equal(t, 2, products[1].ID)
}

func TestProductRecordsMysql_FindAll_query_error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery(repository.FindAllProductRecords).WillReturnError(errors.New("query execution error"))

	repo := repository.NewProductRecordsSQL(mockDB)

	productRecords, err := repo.FindAll()
	assert.Error(t, err)
	assert.Nil(t, productRecords)
	assert.Equal(t, "query execution error", err.Error())
}

func TestProductRecordsMysql_FindAll_scan_error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"last_update_date",
		"purchase_price",
		"sale_price",
		"product_id",
	}).
		AddRow(nil, nil, nil, nil, nil).
		RowError(1, errors.New("scan error"))

	mock.ExpectQuery(repository.FindAllProductRecords).WillReturnRows(rows)

	repo := repository.NewProductRecordsSQL(mockDB)

	productRecords, err := repo.FindAll()
	assert.Error(t, err) 
	assert.Nil(t, productRecords)
	assert.Equal(t, internal.ErrProductNotFound, err)
}


func TestProductRecordsMysql_FinAByID_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"last_update_date",
		"purchase_price",
		"sale_price",
		"product_id",
	}).
		AddRow(1, time.Now().Truncate(24*time.Hour), 1, 1, 1)

	mock.ExpectQuery(repository.FindByIDProductRecords).WillReturnRows(row)

	repo := repository.NewProductRecordsSQL(mockDB)

	product, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
}

func TestProductRecordsMysql_FinAByID_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"last_update_date",
		"purchase_price",
		"sale_price",
		"product_id",
	})

	mock.ExpectQuery(repository.FindByIDProductRecords).WillReturnRows(row)

	repo := repository.NewProductRecordsSQL(mockDB)

	product, err := repo.FindByID(1)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrProductRecordsNotFound, err)
	assert.Empty(t, product)
}

func TestProductRecordsMysql_Save_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.SaveProductRecords).
		WithArgs(
			productRecords.LastUpdateDate,
			productRecords.PurchasePrice,
			productRecords.SalePrice,
			productRecords.ProductID).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewProductRecordsSQL(mockDB)

	_, err = repo.Save(productRecords)
	assert.NoError(t, err)
}

func TestProductRecordsMysql_Save_error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.SaveProductRecords).
		WithArgs(
			productRecords.LastUpdateDate,
			productRecords.PurchasePrice,
			productRecords.SalePrice,
			productRecords.ProductID).WillReturnError(fmt.Errorf("some Error"))

	repo := repository.NewProductRecordsSQL(mockDB)

	_, err = repo.Save(productRecords)
	assert.Error(t, err)
}

func TestProductRecordsMysql_Save_conflict(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.SaveProductRecords).
		WithArgs(
			productRecords.LastUpdateDate,
			productRecords.PurchasePrice,
			productRecords.SalePrice,
			productRecords.ProductID).WillReturnError(&mysql.MySQLError{Number: 1062})

	repo := repository.NewProductRecordsSQL(mockDB)

	_, err = repo.Save(productRecords)
	assert.Error(t, err, internal.ErrProductRecordsConflict)
}
