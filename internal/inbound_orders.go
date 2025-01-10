package internal

type InboundOrders struct {
	Id             int    `json:"id"`
	OrderDate      string `json:"order_date"`
	OrderNumber    string `json:"order_number"`
	EmployeeId     int    `json:"employee_id"`
	ProductBatchId int    `json:"product_batch_id"`
	WarehouseId    int    `json:"warehouse_id"`
}

type InboundOrderService interface {
	Create(InboundOrders) (int64, error)
	FindAll() ([]InboundOrders, error)
}

type InboundOrdersRepository interface {
	Create(InboundOrders) (int64, error)
	FindAll() ([]InboundOrders, error)
}

// validate required fields
func (io *InboundOrders) ValidateFieldsOk() bool {

	if io.OrderDate == "" || io.OrderNumber == "" || io.EmployeeId == 0 || io.ProductBatchId == 0 || io.WarehouseId == 0 {
		return false
	}
	return true
}
