package service

import (
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

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

// Method to check if a warehouse code already exists
func (s *WarehouseDefault) checkWarehouseCodeExists(warehouseCode string) (err error) {
	allWarehouses, err := s.rp.FindAll()
	if err != nil {
		return
	}

	// We`re gonna check if there is a warehouse with the same code
	for _, w := range allWarehouses {
		if w.WarehouseCode == warehouseCode {
			return internal.ErrWarehouseRepositoryDuplicated
		}
	}

	return nil
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
	// Validating the warehouse
	if err := warehouse.Validate(); err != nil {
		return fmt.Errorf("%w: %v", internal.ErrWarehouseUnprocessableEntity, err)
	}

	// We`re gonna check if there is a warehouse with the same code
	err = s.checkWarehouseCodeExists(warehouse.WarehouseCode)
	if err != nil {
		switch err {
		case internal.ErrWarehouseRepositoryDuplicated:
			return internal.ErrWarehouseRepositoryDuplicated
		default:
			return
		}
	}

	err = s.rp.Save(warehouse)

	return
}

// Update updates a warehouse
func (s *WarehouseDefault) Update(id int, warehousePatch *internal.WarehousePatchUpdate) (warehouse internal.Warehouse, err error) {
	warehouse, err = s.rp.FindByID(id)
	if err != nil {
		return internal.Warehouse{}, internal.ErrWarehouseRepositoryNotFound
	}

	// Update the warehouse that we found
	if warehousePatch.WarehouseCode != nil {
		// We`re gonna check if there is a warehouse with the same code
		err = s.checkWarehouseCodeExists(*warehousePatch.WarehouseCode)
		if err != nil && warehouse.WarehouseCode != *warehousePatch.WarehouseCode {
			switch err {
			case internal.ErrWarehouseRepositoryDuplicated:
				return internal.Warehouse{}, internal.ErrWarehouseRepositoryDuplicated
			default:
				return internal.Warehouse{}, err
			}
		}

		warehouse.WarehouseCode = *warehousePatch.WarehouseCode
	}

	if warehousePatch.Address != nil {
		warehouse.Address = *warehousePatch.Address
	}

	if warehousePatch.Telephone != nil {
		warehouse.Telephone = *warehousePatch.Telephone
	}

	if warehousePatch.MinimumCapacity != nil {
		warehouse.MinimumCapacity = *warehousePatch.MinimumCapacity
	}

	if warehousePatch.MinimumTemperature != nil {
		warehouse.MinimumTemperature = *warehousePatch.MinimumTemperature
	}

	// Save the updated warehouse
	err = s.rp.Update(&warehouse)

	return warehouse, err
}

// Delete deletes a warehouse
func (s *WarehouseDefault) Delete(id int) (err error) {
	_, err = s.rp.FindByID(id)
	if err != nil {
		return
	}

	err = s.rp.Delete(id)

	return
}
