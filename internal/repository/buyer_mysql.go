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

func (r *BuyerMysqlRepository) GetAll() (db map[int]internal.Buyer, err error) {
	db = make(map[int]internal.Buyer)
	query := `
		SELECT
			id, card_number_id, first_name, last_name
		FROM
			buyers;
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return
	}

	for rows.Next() {
		var buyer internal.Buyer

		rows.Scan(&buyer.ID, &buyer.CardNumberID, &buyer.FirstName, &buyer.LastName)
		db[buyer.ID] = buyer
	}

	err = rows.Err()
	return
}

func (r *BuyerMysqlRepository) Add(buyer *internal.Buyer) (id int64, err error) {
	query := `
		INSERT INTO buyers (card_number_id, first_name, last_name)
		VALUES (?, ?, ?)
	`

	result, err := r.db.Exec(query, (*buyer).CardNumberID, (*buyer).FirstName, (*buyer).LastName)
	if err != nil {
		return
	}

	id, err = result.LastInsertId()

	(*buyer).ID = int(id)
	return
}

func (r *BuyerMysqlRepository) Update(id int, buyer internal.BuyerPatch) (err error) {
	query :=
		`
		SELECT
			id, card_number_id, first_name, last_name
		FROM
			buyers
		WHERE
			id = ?;
	`
	row := r.db.QueryRow(query, id)

	var b internal.Buyer
	row.Scan(&b.ID, &b.CardNumberID, &b.FirstName, &b.LastName)

	buyer.Patch(&b)

	query = `
		UPDATE buyers
		SET
			card_number_id = ?, first_name = ?, last_name = ?
		WHERE
			id = ?;
	`

	_, err = r.db.Exec(query, buyer.CardNumberID, buyer.FirstName, buyer.LastName, id)
	return
}

func (r *BuyerMysqlRepository) Delete(id int) (rowsAffected int64, err error) {
	query := `
		DELETE FROM 
			buyers
		WHERE
			id = ?;
	`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return
	}

	rowsAffected, err = res.RowsAffected()
	return
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
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var purchaseOrder internal.PurchaseOrdersByBuyer

		rows.Scan(&purchaseOrder.BuyerID, &purchaseOrder.CardNumberID, &purchaseOrder.FirstName, &purchaseOrder.LastName, &purchaseOrder.PurchaseOrdersCount)

		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	err = rows.Err()
	return purchaseOrders, err
}

func (r *BuyerMysqlRepository) ReportPurchaseOrdersByID(id int) (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
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
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var purchaseOrder internal.PurchaseOrdersByBuyer

		err = rows.Scan(&purchaseOrder.BuyerID, &purchaseOrder.CardNumberID, &purchaseOrder.FirstName, &purchaseOrder.LastName, &purchaseOrder.PurchaseOrdersCount)

		purchaseOrders = append(purchaseOrders, purchaseOrder)
	}

	err = rows.Err()
	return purchaseOrders, err
}
