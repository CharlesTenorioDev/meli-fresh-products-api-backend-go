package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewPurchaseOrderMysqlRepository(db *sql.DB) *PurchaseOrderRepository {
	return &PurchaseOrderRepository{db}
}

type PurchaseOrderRepository struct {
	db *sql.DB
}

func (r *PurchaseOrderRepository) FindByID(id int) (purchaseOrder internal.PurchaseOrder, err error) {
	query := `
		SELECT po.id, po.order_number, po.order_date, po.tracking_code, po.buyer_id, po.product_record_id
		FROM purchase_orders as po
		WHERE po.id = ?
	`
	row := r.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return internal.PurchaseOrder{}, err
	}

	// scanning the row
	err = row.Scan(&purchaseOrder.ID, &purchaseOrder.OrderNumber, &purchaseOrder.OrderDate, &purchaseOrder.TrackingCode, &purchaseOrder.BuyerID, &purchaseOrder.ProductRecordId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrPurchaseOrderNotFound
		}
		return internal.PurchaseOrder{}, err
	}

	return
}

// Save creates a new purchase order in the database
func (r *PurchaseOrderRepository) Save(purchaseOrder *internal.PurchaseOrder) error {
	// Checking if the purchase order already exists
	row := r.db.QueryRow("SELECT COUNT(*) FROM purchase_orders WHERE order_number = ?", purchaseOrder.OrderNumber)
	var count int
	row.Scan(&count)

	if count > 0 {
		return internal.ErrPurchaseOrderConflict
	}

	// Inserting the purchase order
	query := `
		INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, (*purchaseOrder).OrderNumber, (*purchaseOrder).OrderDate, (*purchaseOrder).TrackingCode, (*purchaseOrder).BuyerID, (*purchaseOrder).ProductRecordId)
	if err != nil {
		return err
	}

	// Get the ID of the last inserted purchase order
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Set the ID of the purchase order
	(*purchaseOrder).ID = int(id)

	return nil
}
