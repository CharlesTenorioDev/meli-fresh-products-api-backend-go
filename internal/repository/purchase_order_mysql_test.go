package repository_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestPurchaseOrderMysql_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	query := `
		SELECT po.id, po.order_number, po.order_date, po.tracking_code, po.buyer_id, po.product_record_id
		FROM purchase_orders as po
		WHERE po.id = ?
	`

	t.Run("case 1: success - Purchase Order found", func(t *testing.T) {
		id := 1
		date := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedPO := internal.PurchaseOrder{
			ID:              1,
			OrderNumber:     "123ABC",
			OrderDate:       date,
			TrackingCode:    "tracking_code",
			BuyerID:         1,
			ProductRecordID: 1,
		}

		rows := sqlmock.NewRows([]string{"id", "order_number", "order_date", "tracking_code", "buyer_id", "product_record_id"}).
			AddRow(expectedPO.ID, expectedPO.OrderNumber, expectedPO.OrderDate, expectedPO.TrackingCode, expectedPO.BuyerID, expectedPO.ProductRecordID)

		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnRows(rows)

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		po, err := rp.FindByID(id)

		require.NoError(t, err)
		require.Equal(t, expectedPO, po)
	})

	t.Run("case 2: failure - Purchase Order not found", func(t *testing.T) {
		id := 1
		mock.ExpectQuery(query).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		_, err = rp.FindByID(id)

		require.Error(t, err)
		require.Equal(t, internal.ErrPurchaseOrderNotFound, err)
	})
}

func TestPurchaseOrderMysql_Save(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	query := `
		INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
		VALUES (?, ?, ?, ?, ?)
	`

	po := internal.PurchaseOrder{
		OrderNumber:     "123ABC",
		OrderDate:       time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
		TrackingCode:    "tracking_code",
		BuyerID:         1,
		ProductRecordID: 1,
	}

	t.Run("case 1: success - Purchase Order saved", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT(*) FROM purchase_orders WHERE order_number = ?").
			WithArgs(po.OrderNumber).
			WillReturnRows(rows)

		mock.ExpectExec(query).
			WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		err := rp.Save(&po)

		require.NoError(t, err)
	})

	t.Run("case 2: error - Purchase Order already exists", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
		mock.ExpectQuery("SELECT COUNT(*) FROM purchase_orders WHERE order_number = ?").
			WithArgs(po.OrderNumber).
			WillReturnRows(rows)

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		err := rp.Save(&po)

		require.Error(t, err)
		require.Equal(t, internal.ErrPurchaseOrderConflict, err)
	})

	t.Run("case 3: error - Error executing the query", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT(*) FROM purchase_orders WHERE order_number = ?").
			WithArgs(po.OrderNumber).
			WillReturnRows(rows)

		mock.ExpectExec(query).
			WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
			WillReturnError(sql.ErrConnDone)

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		err := rp.Save(&po)

		require.Error(t, err)
	})

	t.Run("case 4: error - Error retrieving the last inserted ID", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)
		mock.ExpectQuery("SELECT COUNT(*) FROM purchase_orders WHERE order_number = ?").
			WithArgs(po.OrderNumber).
			WillReturnRows(rows)

		mock.ExpectExec(query).
			WithArgs(po.OrderNumber, po.OrderDate, po.TrackingCode, po.BuyerID, po.ProductRecordID).
			WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error")))

		rp := repository.NewPurchaseOrderMysqlRepository(db)
		err := rp.Save(&po)

		require.Error(t, err)
	})

}
