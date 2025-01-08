package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

func NewBuyerMysqlRepository(db *sql.DB) *BuyerMysqlRepository {
	return &BuyerMysqlRepository{db}
}

type BuyerMysqlRepository struct {
	db *sql.DB
}

func (r *BuyerMysqlRepository) GetAll() (db map[int]internal.Buyer) {
	return nil
}

func (r *BuyerMysqlRepository) Add(buyer *internal.Buyer) {
	return
}

func (r *BuyerMysqlRepository) Update(id int, buyer internal.BuyerPatch) {
	return
}

func (r *BuyerMysqlRepository) Delete(id int) {
	return
}

func (r *BuyerMysqlRepository) ReportPurchaseOrders() (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	query := `
		SELECT
			b.id, b.card_number_id, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count
		FROM
			buyers as b
		INNER JOIN
			purchase_orders as po ON po.buyer_id = b.id
		GROUP BY
			b.id;
	`
	// executing the query
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	// iterating over the rows
	for rows.Next() {
		var purchaseOrder internal.PurchaseOrdersByBuyer
		err = rows.Scan(&purchaseOrder.BuyerID, &purchaseOrder.CardNumberId, &purchaseOrder.FirstName, &purchaseOrder.LastName, &purchaseOrder.PurchaseOrdersCount)
		if err != nil {
			return
		}
		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

func (r *BuyerMysqlRepository) ReportPurchaseOrdersById(id int) (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	query := `
		SELECT
			b.id, b.card_number_id, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count
		FROM
			buyers as b
		INNER JOIN
			purchase_orders as po ON po.buyer_id = b.id
		GROUP BY
			b.id
		HAVING
			b.id = ?;
	`
	// executing the query
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterating over the rows
	for rows.Next() {
		var purchaseOrder internal.PurchaseOrdersByBuyer
		err = rows.Scan(&purchaseOrder.BuyerID, &purchaseOrder.CardNumberId, &purchaseOrder.FirstName, &purchaseOrder.LastName, &purchaseOrder.PurchaseOrdersCount)
		if err != nil {
			return
		}
		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}
