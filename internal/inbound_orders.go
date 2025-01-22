package internal

import "errors"

var ErrOrderNumberAlreadyExists = errors.New("order number already exists")

type InboundOrders struct {
	ID             int    `json:"id"`
	OrderDate      string `json:"order_date"`
	OrderNumber    string `json:"order_number"`
	EmployeeID     int    `json:"employee_id"`
	ProductBatchID int    `json:"product_batch_id"`
	WarehouseID    int    `json:"warehouse_id"`
}

type InboundOrderService interface {
	Create(InboundOrders) (int64, error)
	FindAll() ([]InboundOrders, error)
}

type InboundOrdersRepository interface {
	Create(InboundOrders) (int64, error)
	FindAll() ([]InboundOrders, error)
}

// ValidateFieldsOk validates required fields
func (io *InboundOrders) ValidateFieldsOk() bool {
	if io.OrderDate == "" || io.OrderNumber == "" || io.EmployeeID == 0 || io.ProductBatchID == 0 || io.WarehouseID == 0 {
		return false
	}

	return true
}
