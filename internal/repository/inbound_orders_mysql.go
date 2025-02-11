package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const (
	AllInboundsQuery = "SELECT `id`, `order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id` FROM `inbound_orders`;"
)

// InboundOrdersMysql create a new instance of the inbound orders repository
type InboundOrdersMysql struct {
	db *sql.DB
}

func NewInboundOrderMysql(db *sql.DB) *InboundOrdersMysql {
	return &InboundOrdersMysql{db}
}

func (rp *InboundOrdersMysql) Create(io internal.InboundOrders) (id int64, err error) {
	var exists bool

	rp.db.QueryRow("SELECT 1 FROM `inbound_orders` WHERE `order_number` = ?", io.OrderNumber).Scan(&exists) //check 1 line

	if exists {
		return 0, internal.ErrOrderNumberAlreadyExists
	}

	var empExists bool

	rp.db.QueryRow("SELECT 1 FROM `employees` WHERE `id` = ?", io.EmployeeID).Scan(&empExists) //check 1 line

	if !empExists {
		return 0, internal.ErrEmployeeNotFound
	}

	res, err := rp.db.Exec(
		"INSERT INTO `inbound_orders` (`order_date`, `order_number`, `employee_id`, `product_batch_id`, `warehouse_id`) VALUES (?, ?, ?, ?, ?)",
		io.OrderDate, io.OrderNumber, io.EmployeeID, io.ProductBatchID, io.WarehouseID,
	)
	if err != nil {
		return id, err
	}

	id, err = res.LastInsertId()

	return id, err
}

func (rp *InboundOrdersMysql) FindAll() (inbounds []internal.InboundOrders, err error) {
	row, err := rp.db.Query(AllInboundsQuery)

	if err != nil {
		return
	}

	for row.Next() {
		var inboundOrder internal.InboundOrders
		err = row.Scan(&inboundOrder.ID, &inboundOrder.OrderDate, &inboundOrder.OrderNumber, &inboundOrder.EmployeeID, &inboundOrder.ProductBatchID, &inboundOrder.WarehouseID)

		if err != nil {
			return
		}

		inbounds = append(inbounds, inboundOrder)
	}

	return
}
