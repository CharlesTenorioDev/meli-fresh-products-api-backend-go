package repository_test

import (
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
)

func TestSellerMysql_FindAll(t *testing.T) {
	// Inicializa o mock do banco de dados
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	// Inicializa o txdb com o mock do banco de dados

	// Define o comportamento esperado do mock
	rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
		AddRow(1, 123, "Test Seller", "Rua 1", "1234567890").
		AddRow(2, 456, "Another Seller", "Rua 2", "9876543210")
	mock.ExpectQuery("SELECT `s.id`, `s.cid`, `s.company_name`, `s.address`, `s.telephone` FROM `sellers` AS `s`").WillReturnRows(rows)

	// Cria o repositório
	repo := repository.NewSellerMysql(mockDB)

	// Executa a função FindAll dentro de uma transação

	sellers, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(sellers))
	assert.Equal(t, 1, sellers[0].ID)
	assert.Equal(t, 2, sellers[1].ID)

}

func TestSellerMysql_FindByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
		AddRow(1, 123, "Test Seller", "Rua 1", "1234567890")
	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers`  WHERE `id` = ?").WithArgs(1).WillReturnRows(row)

	repo := repository.NewSellerMysql(mockDB)

	seller, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, seller.ID)

}

func TestSellerMysql_FindByCID(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
		AddRow(1, 123, "Test Seller", "Rua 1", "1234567890")
	mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers` WHERE `cid` = ?").WithArgs(123).WillReturnRows(row)

	repo := repository.NewSellerMysql(mockDB)

	seller, err := repo.FindByCID(123)
	assert.NoError(t, err)
	assert.Equal(t, 123, seller.CID)

}

func TestSellerMysql_Save(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	seller := &internal.Seller{
		CID:         123,
		CompanyName: "Test Seller",
		Address:     "Rua 1",
		Telephone:   "1234567890",
		Locality:    1,
	}

	mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
		WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewSellerMysql(mockDB)

	err = repo.Save(seller)
	assert.NoError(t, err)
	assert.Equal(t, 1, seller.ID)

}

func TestSellerMysql_Update(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	seller := &internal.Seller{
		ID:          1,
		CID:         123,
		CompanyName: "Test Seller",
		Address:     "Rua 1",
		Telephone:   "1234567890",
		Locality:    1,
	}

	mock.ExpectExec("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?").
		WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality, seller.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewSellerMysql(mockDB)

	err = repo.Update(seller)
	assert.NoError(t, err)

}

func TestSellerMysql_Delete(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectExec("DELETE FROM `sellers` WHERE `id` = ?").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewSellerMysql(mockDB)

	err = repo.Delete(1)

	assert.NoError(t, err)
}

func TestSellerMysql_Save_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	seller := &internal.Seller{
		CID:         123,
		CompanyName: "Test Seller",
		Address:     "Rua 1",
		Telephone:   "1234567890",
		Locality:    1,
	}

	mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
		WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
		WillReturnError(fmt.Errorf("some error"))

	repo := repository.NewSellerMysql(mockDB)

	err = repo.Save(seller)
	assert.Error(t, err)

}

func TestSellerMysql_Save_Conflict(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	seller := &internal.Seller{
		CID:         123,
		CompanyName: "Test Seller",
		Address:     "Rua 1",
		Telephone:   "1234567890",
		Locality:    1,
	}

	mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
		WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
		WillReturnError(&mysql.MySQLError{Number: 1062})

	repo := repository.NewSellerMysql(mockDB)

	err = repo.Save(seller)
	assert.ErrorIs(t, err, internal.ErrSellerConflict)

}
