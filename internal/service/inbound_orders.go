package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

type InboundOrderService struct {
	rp  internal.InboundOrdersRepository
	rpE internal.EmployeeRepository
	rpW internal.WarehouseRepository
}

func NewInboundOrderService(rp internal.InboundOrdersRepository, rpEmployee internal.EmployeeRepository, rpWarehouse internal.WarehouseRepository) *InboundOrderService {
	return &InboundOrderService{
		rp:  rp,
		rpE: rpEmployee,
		rpW: rpWarehouse,
	}
}

func (s *InboundOrderService) Create(internal.InboundOrders) (int64, error) {
	inboundOrder := internal.InboundOrders{}
	return s.rp.Create(inboundOrder)
}

func (s *InboundOrderService) FindAll() ([]internal.InboundOrders, error) {
	return s.rp.FindAll()
}
