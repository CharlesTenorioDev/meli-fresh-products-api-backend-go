package service

import "github.com/meli-fresh-products-api-backend-t1/internal"

// NewWarehouseDefault creates a new instance of the warehouse service
func NewWarehouseDefault(rp internal.WarehouseRepository) *WarehouseDefault {
	return &WarehouseDefault{
		rp: rp,
	}
}

// WarehouseDefault is the default implementation of the warehouse service
type WarehouseDefault struct {
	// rp is the repository used by the service
	rp internal.WarehouseRepository
}

// FindAll returns all warehouses
func (s *WarehouseDefault) FindAll() (warehouses []internal.Warehouse, err error) {
	warehouses, err = s.rp.FindAll()
	return
}
