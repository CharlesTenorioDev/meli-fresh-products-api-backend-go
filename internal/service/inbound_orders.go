package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

type InboundOrderService struct {
	rp  internal.InboundOrdersRepository
	rpE internal.EmployeeRepository
	rpP internal.ProductBatchRepository
	rpW internal.WarehouseRepository
}

func NewInboundOrderService(rpInbound internal.InboundOrdersRepository, rpEmployee internal.EmployeeRepository, rpProductBatch internal.ProductBatchRepository, rpWarehouse internal.WarehouseRepository) *InboundOrderService {
	return &InboundOrderService{
		rp:  rpInbound,
		rpE: rpEmployee,
		rpP: rpProductBatch,
		rpW: rpWarehouse,
	}
}

func (s *InboundOrderService) Create(inboundOrder internal.InboundOrders) (int64, error) {
	return s.rp.Create(inboundOrder)
}

func (s *InboundOrderService) FindAll() ([]internal.InboundOrders, error) {
	return s.rp.FindAll()
}
