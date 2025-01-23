package repository_test

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/assert"
)

var product = internal.Product{
	ID:                             1,
	ProductCode:                    "code 1",
	Description:                    "Example Product",
	Height:                         1,
	Length:                         1,
	NetWeight:                      1,
	ExpirationRate:                 1,
	RecommendedFreezingTemperature: 1,
	Width:                          1,
	FreezingRate:                   1,
	ProductTypeID:                  1,
	SellerID:                       1,
}

func TestProductMysql_FinAll(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{
		"id",
		"product_code",
		"description",
		"height",
		"length",
		"net_weight",
		"expiration_rate",
		"recommended_freezing_temperature",
		"width",
		"freezing_rate",
		"product_type_id",
		"seller_id",
	}).
		AddRow(1, "code 1", 1, 1, 1, 1, 1, "desc 1", 1, 1, 1, 1).
		AddRow(2, "code 2", 2, 2, 2, 2, 2, "desc 2", 2, 2, 2, 2)

	mock.ExpectQuery(repository.FindAllString).WillReturnRows(rows)

	repo := repository.NewProductSQL(mockDB)

	products, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(products))
	assert.Equal(t, 1, products[0].ID)
	assert.Equal(t, 2, products[1].ID)
}

func TestProductMysql_FinAByID_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"product_code",
		"description",
		"height",
		"length",
		"net_weight",
		"expiration_rate",
		"recommended_freezing_temperature",
		"width",
		"freezing_rate",
		"product_type_id",
		"seller_id",
	}).
		AddRow(1, "code 1", 1, 1, 1, 1, 1, "desc 1", 1, 1, 1, 1)

	mock.ExpectQuery(repository.FindByIDString).WithArgs(1).WillReturnRows(row)

	repo := repository.NewProductSQL(mockDB)

	product, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, product.ID)
}

func TestProductMysql_FinAByID_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"id",
		"product_code",
		"description",
		"height",
		"length",
		"net_weight",
		"expiration_rate",
		"recommended_freezing_temperature",
		"width",
		"freezing_rate",
		"product_type_id",
		"seller_id",
	})

	mock.ExpectQuery(repository.FindByIDString).WithArgs(1).WillReturnRows(row)

	repo := repository.NewProductSQL(mockDB)

	product, err := repo.FindByID(1)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrProductNotFound, err)
	assert.Empty(t, product)
}

func TestProductMysql_save_ok(t *testing.T) {
	// Cria um mock do banco de dados e configura a expectativa para a consulta SQL
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err) // Verifica se não teve erro ao criar o mock do banco
	defer mockDB.Close()   // Garante que o mock do banco será fechado no final do teste

	// Configura o mock para esperar que a query de "Save" seja executada com os argumentos que passamos
	// E simula a execução da query, retornando um resultado indicando que 1 linha foi afetada
	mock.ExpectExec(repository.SaveString).
		WithArgs(
			product.ID,
			product.Description,
			product.ExpirationRate,
			product.FreezingRate,
			product.Height,
			product.Length,
			product.NetWeight,
			product.ProductCode,
			product.RecommendedFreezingTemperature,
			product.Width,
			product.ProductTypeID,
			product.SellerID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	// Cria uma instância do repositório ProductSQL passando o mock do banco
	repo := repository.NewProductSQL(mockDB)

	// Chama o método Save para salvar o produto
	_, err = repo.Save(product)

	// Verifica se o método Save não gerou erro
	assert.NoError(t, err)
}

func TestProductMysql_save_error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.SaveString).
		WithArgs(
			product.ID,
			product.Description,
			product.ExpirationRate,
			product.FreezingRate,
			product.Height,
			product.Length,
			product.NetWeight,
			product.ProductCode,
			product.RecommendedFreezingTemperature,
			product.Width,
			product.ProductTypeID,
			product.SellerID,
		).WillReturnError(fmt.Errorf("some Error"))

	// Cria uma instância do repositório ProductSQL passando o mock do banco
	repo := repository.NewProductSQL(mockDB)

	// Chama o método Save para salvar o produto
	_, err = repo.Save(product)
	// Verifica se o método Save gerou erro
	assert.Error(t, err)
}

func TestProductMysql_save_conflict(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.SaveString).
		WithArgs(
			product.ID,
			product.Description,
			product.ExpirationRate,
			product.FreezingRate,
			product.Height,
			product.Length,
			product.NetWeight,
			product.ProductCode,
			product.RecommendedFreezingTemperature,
			product.Width,
			product.ProductTypeID,
			product.SellerID,
		).WillReturnError(&mysql.MySQLError{Number: 1062})

	// Cria uma instância do repositório ProductSQL passando o mock do banco
	repo := repository.NewProductSQL(mockDB)

	// Chama o método Save para salvar o produto
	_, err = repo.Save(product)
	// Verifica se o método Save gerou erro
	assert.Error(t, err, internal.ErrProductConflit)
}

func TestProductMysql_Update_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.UpdateString).
		WithArgs(
			product.Description,
			product.ExpirationRate,
			product.FreezingRate,
			product.Height,
			product.Length,
			product.NetWeight,
			product.ProductCode,
			product.RecommendedFreezingTemperature,
			product.Width,
			product.ProductTypeID,
			product.SellerID,
			product.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewProductSQL(mockDB)

	_, err = repo.Update(product)
	assert.NoError(t, err)
}

func TestProductMysql_Update_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.UpdateString).
		WithArgs(
			product.Description,
			product.ExpirationRate,
			product.FreezingRate,
			product.Height,
			product.Length,
			product.NetWeight,
			product.ProductCode,
			product.RecommendedFreezingTemperature,
			product.Width,
			product.ProductTypeID,
			product.SellerID,
			product.ID,
		).WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repository.NewProductSQL(mockDB)

	_, err = repo.Update(product)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrProductNotFound, err)
}

func TestProductMysql_Delete_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.DeleteString).
		WithArgs(
			product.ID,
		).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewProductSQL(mockDB)

	err = repo.Delete(1)
	assert.NoError(t, err)
}

func TestProductMysql_Delete_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec(repository.DeleteString).
		WithArgs(
			product.ID,
		).WillReturnResult(sqlmock.NewResult(0, 0))

	repo := repository.NewProductSQL(mockDB)

	err = repo.Delete(1)
	assert.Error(t, internal.ErrProductIdNotFound, err)
}

func TestProductMysql_FindAllRecord(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{
		"product_id",
		"description",
		"records_count",
	}).
		AddRow(1, "code 1", 1).
		AddRow(2, "code 2", 2)

	mock.ExpectQuery(repository.FindAllRecordString).WillReturnRows(rows)

	repo := repository.NewProductSQL(mockDB)

	productRecords, err := repo.FindAllRecord()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(productRecords))
	assert.Equal(t, 1, productRecords[0].ProductID)
	assert.Equal(t, 2, productRecords[1].ProductID)
}

func TestProductMysql_FinAByIDRecords_ok(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"product_id",
		"description",
		"records_count",
	}).
		AddRow(1, "code 1", 1)

	mock.ExpectQuery(repository.FindByIDRecordString).WithArgs(1).WillReturnRows(row)

	repo := repository.NewProductSQL(mockDB)

	productRecords, err := repo.FindByIDRecord(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, productRecords.ProductID)
}

func TestProductMysql_FinAByIDRecords_not_found(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{
		"product_id",
		"description",
		"records_count",
	})

	mock.ExpectQuery(repository.FindByIDRecordString).WithArgs(1).WillReturnRows(row)

	repo := repository.NewProductSQL(mockDB)

	productRecords, err := repo.FindByIDRecord(1)
	assert.Error(t, err)
	assert.Equal(t, internal.ErrProductIdNotFound, err)
	assert.Empty(t, productRecords)
}
