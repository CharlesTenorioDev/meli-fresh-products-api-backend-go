package internal

import "errors"

// Warehouse is a struct that represents a warehouse
type Warehouse struct {
	ID                 int
	WarehouseCode      string
	Address            string
	Telephone          string
	MinimumCapacity    int
	MinimumTemperature float64
}

// WarehousePatchUpdate is a struct to use in a patch request
type WarehousePatchUpdate struct {
	WarehouseCode      *string  `json:"warehouse_code"`
	Address            *string  `json:"address"`
	Telephone          *string  `json:"telephone"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
}

var (
	// ErrWarehouseRepositoryNotFound is returned when the warehouse is not found
	ErrWarehouseRepositoryNotFound = errors.New("warehouse not found")
	// ErrWarehouseRepositoryDuplicated is returned when the warehouse already exists
	ErrWarehouseRepositoryDuplicated = errors.New("warehouse already exists")
	// ErrWarehouseUnprocessableEntity is returned when the warehouse is unprocessable
	ErrWarehouseUnprocessableEntity = errors.New("unprocessable entity")
)

func (w *Warehouse) Validate() error {
	if w.WarehouseCode == "" {
		return errors.New("warehouse code is required")
	}

	if w.Address == "" {
		return errors.New("address is required")
	}

	if w.Telephone == "" {
		return errors.New("telephone is required")
	}

	if w.MinimumCapacity == 0 {
		return errors.New("minimum capacity is required")
	}

	if w.MinimumTemperature == 0 {
		return errors.New("minimum temperature is required")
	}

	return nil
}

// WarehouseRepository is an interface that contains the methods that the warehouse repository should support
type WarehouseRepository interface {
	// FindAll returns all the warehouses
	FindAll() ([]Warehouse, error)
	// FindByID returns the warehouse with the given ID
	FindByID(id int) (Warehouse, error)
	// Save saves the given warehouse
	Save(warehouse *Warehouse) error
	// Update updates the given warehouse
	Update(warehouse *Warehouse) error
	// Delete deletes the warehouse with the given ID
	Delete(id int) error
}

// WarehouseService is an interface that contains the methods that the warehouse service should support
type WarehouseService interface {
	// FindAll returns all the warehouses
	FindAll() ([]Warehouse, error)
	// FindByID returns the warehouse with the given ID
	FindByID(id int) (Warehouse, error)
	// Save saves the given warehouse
	Save(warehouse *Warehouse) error
	// Update updates the given warehouse
	Update(id int, warehousePatch *WarehousePatchUpdate) (Warehouse, error)
	// Delete deletes the warehouse with the given ID
	Delete(id int) error
}
