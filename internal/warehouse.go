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
	ErrWarehouseRepositoryNotFound = errors.New("Warehouse not found")
	// ErrWarehouseRepositoryDuplicated is returned when the warehouse already exists
	ErrWarehouseRepositoryDuplicated = errors.New("Warehouse already exists")
)

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
