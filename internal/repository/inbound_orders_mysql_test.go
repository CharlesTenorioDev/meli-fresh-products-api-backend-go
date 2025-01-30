package repository_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestInboundMysqlCreate_Success(t *testing.T) {
	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	inbound := internal.InboundOrders{
		OrderDate:      "2025-01-01",
		OrderNumber:    "1111111",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    1,
	}

	mockRep.ExpectQuery("SELECT 1 FROM `inbound_orders` WHERE `order_number` = ?").
		WithArgs(inbound.OrderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	mockRep.ExpectQuery("SELECT 1 FROM `employees` WHERE `id` = ?").WithArgs(inbound.EmployeeID).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mockRep.ExpectExec("INSERT INTO `inbound_orders` (`order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id`) VALUES (?, ?, ?, ?, ?)").
		WithArgs(inbound.OrderDate, inbound.OrderNumber, inbound.EmployeeID, inbound.ProductBatchID, inbound.WarehouseID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	id, err := rep.Create(inbound)

	//assert
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)

}

func TestInboundMysqlCreate_OrderNumberExists(t *testing.T) {
	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	inbound := internal.InboundOrders{
		OrderDate:      "2025-01-01",
		OrderNumber:    "1111111",
		EmployeeID:     2,
		ProductBatchID: 5,
		WarehouseID:    7,
	}

	mockRep.ExpectQuery("SELECT 1 FROM `inbound_orders` WHERE `order_number` = ?").
		WithArgs(inbound.OrderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	_, err = rep.Create(inbound)

	//assert
	assert.ErrorIs(t, err, internal.ErrOrderNumberAlreadyExists)
}

func TestInboundMysqlCreate_EmployeeNotFound(t *testing.T) {
	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	inbound := internal.InboundOrders{
		OrderDate:      "2025-01-01",
		OrderNumber:    "45234252",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    580,
	}

	mockRep.ExpectQuery("SELECT 1 FROM `inbound_orders` WHERE `order_number` = ?").
		WithArgs(inbound.OrderNumber).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	mockRep.ExpectQuery("SELECT 1 FROM `employees` WHERE `id` = ?").WithArgs(inbound.EmployeeID).
		WillReturnError(sql.ErrNoRows)

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	_, err = rep.Create(inbound)

	//assert
	assert.ErrorIs(t, err, internal.ErrEmployeeNotFound)
}

func TestInboundMysqlCreate_InsertError(t *testing.T) {
	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	inboundInput := internal.InboundOrders{
		OrderDate:      "2025-01-01",
		OrderNumber:    "987654321",
		EmployeeID:     1,
		ProductBatchID: 1,
		WarehouseID:    580,
	}

	mockRep.ExpectQuery("SELECT 1 FROM `inbound_orders` WHERE `order_number` = ?").
		WithArgs(inboundInput.OrderNumber).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(0))

	mockRep.ExpectQuery("SELECT 1 FROM `employees` WHERE `id` = ?").
		WithArgs(inboundInput.EmployeeID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mockRep.ExpectExec("INSERT INTO `inbound_orders` (`order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id`) VALUES (?, ?, ?, ?, ?)").
		WithArgs(inboundInput.OrderDate, inboundInput.OrderNumber, inboundInput.EmployeeID, inboundInput.ProductBatchID, inboundInput.WarehouseID).
		WillReturnError(fmt.Errorf("failed to insert inbound order"))

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	id, err := rep.Create(inboundInput)

	//assert
	assert.Error(t, err)
	assert.Equal(t, int64(0), id)

}

func TestInboundMysqlGetAll_Success(t *testing.T) {

	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}).
		AddRow(1, "2025-01-01", "1111111", 1, 1, 1).
		AddRow(2, "2025-02-02", "2222222", 2, 2, 2).
		AddRow(3, "2025-03-03", "3333333", 3, 3, 3)
	mockRep.ExpectQuery("SELECT `id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id` FROM `inbound_orders`;").WillReturnRows(row)

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	inboundOrders, err := rep.FindAll()

	//assert
	assert.NoError(t, err)
	assert.Equal(t, 3, len(inboundOrders))
	assert.Equal(t, 1, inboundOrders[0].ID)
	assert.Equal(t, 2, inboundOrders[1].ID)
	assert.Equal(t, 3, inboundOrders[2].ID)

}

func TestInboundMysqlGetAll_QueryError(t *testing.T) {

	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT `id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id` FROM `inbound_orders`;").
		WillReturnError(fmt.Errorf("failed to execute query"))

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	inbound, err := rep.FindAll()

	//assert
	assert.Error(t, err)
	assert.Nil(t, inbound)

}

func TestInboundMysqlGetAll_ScanError(t *testing.T) {

	mockDB, mockRep, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}).
		AddRow(1, "2025-01-01", "1111111", "aaaaaaaaa", 1, "hello world")

	mockRep.ExpectQuery("SELECT `id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id` FROM `inbound_orders`;").
		WillReturnRows(row)

	mockRep.ExpectQuery("SELECT `id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id` FROM `inbound_orders`;").
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_date", "order_number", "employee_id", "product_batch_id", "warehouse_id"}).
			AddRow(1, "2025-01-01", "1111111", "aaaaaaaaa", 1, "hello world"))

	//create repository
	rep := repository.NewInboundOrderMysql(mockDB)

	//act
	inbound, err := rep.FindAll()

	//assert
	assert.Error(t, err)
	assert.Nil(t, inbound)

}
