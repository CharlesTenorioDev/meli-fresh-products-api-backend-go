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

// FindByID returns a warehouse
func (s *WarehouseDefault) FindByID(id int) (warehouse internal.Warehouse, err error) {
	warehouse, err = s.rp.FindByID(id)
	return
}

// Save creates a new warehouse
func (s *WarehouseDefault) Save(warehouse *internal.Warehouse) (err error) {
	allWarehouses, err := s.rp.FindAll()
	if err != nil {
		return err
	}

	// We`re gonna check if there is a warehouse with the same code
	for _, w := range allWarehouses {
		if w.WarehouseCode == warehouse.WarehouseCode {
			return internal.ErrWarehouseRepositoryDuplicated
		}
	}

	err = s.rp.Save(warehouse)
	return
}
