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
	db = make(map[int]internal.Buyer)
	query := `
		SELECT
			id, card_number_id, first_name, last_name
		FROM
			buyers;
	`

	/// executing the query
	rows, _ := r.db.Query(query)
	// iterating over the rows
	for rows.Next() {
		var buyer internal.Buyer
		rows.Scan(&buyer.ID, &buyer.CardNumberId, &buyer.FirstName, &buyer.LastName)
		db[buyer.ID] = buyer
	}

	return

}

func (r *BuyerMysqlRepository) Add(buyer *internal.Buyer) {
	// Inserting the buyer
	query := `
		INSERT INTO buyers (card_number_id, first_name, last_name)
		VALUES (?, ?, ?)
	`

	result, _ := r.db.Exec(query, (*buyer).CardNumberId, (*buyer).FirstName, (*buyer).LastName)

	// Get the ID of the last inserted purchase order
	id, _ := result.LastInsertId()

	// Set the ID of the purchase order
	(*buyer).ID = int(id)
}

func (r *BuyerMysqlRepository) Update(id int, buyer internal.BuyerPatch) {
	// Finding the buyer
	query :=
		`
		SELECT
			id, card_number_id, first_name, last_name
		FROM
			buyers
		WHERE
			id = ?;
	`
	// executing the query
	row := r.db.QueryRow(query, id)

	var b internal.Buyer
	// scanning the row
	row.Scan(&b.ID, &b.CardNumberId, &b.FirstName, &b.LastName)

	// applying the patch
	buyer.Patch(&b)

	query = `
		UPDATE buyers
		SET
			card_number_id = ?, first_name = ?, last_name = ?
		WHERE
			id = ?;
	`

	// applying the patch
	r.db.Exec(query, buyer.CardNumberId, buyer.FirstName, buyer.LastName, id)
}

func (r *BuyerMysqlRepository) Delete(id int) {
	query := `
		DELETE FROM 
			buyers
		WHERE
			id = ?;
	`
	r.db.Exec(query, id)
}

func (r *BuyerMysqlRepository) ReportPurchaseOrders() (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	query := `
		SELECT
			b.id, b.card_number_id, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count
		FROM
			buyers as b
		LEFT JOIN
			purchase_orders as po ON po.buyer_id = b.id
		GROUP BY
			b.id;
	`
	// executing the query
	rows, err := r.db.Query(query)
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

func (r *BuyerMysqlRepository) ReportPurchaseOrdersById(id int) (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	query := `
		SELECT
			b.id, b.card_number_id, b.first_name, b.last_name, COUNT(po.id) as purchase_orders_count
		FROM
			buyers as b
		LEFT JOIN
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
