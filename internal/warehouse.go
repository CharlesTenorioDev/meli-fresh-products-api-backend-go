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

var (
	// ErrWarehouseRepositoryNotFound is returned when the warehouse is not found
	ErrWarehouseRepositoryNotFound = errors.New("repository: warehouse not found")
	// ErrWarehouseRepositoryDuplicated is returned when the warehouse already exists
	ErrWarehouseRepositoryDuplicated = errors.New("repository: warehouse already exists")
)

// WarehouseRepository is an interface that contains the methods that the warehouse repository should support
type WarehouseRepository interface {
	// FindAll returns all the warehouses
	FindAll() ([]Warehouse, error)
	// FindByID returns the warehouse with the given ID
	FindByID(id int) (Warehouse, error)
	// Save saves the given warehouse
	Save(warehouse *Warehouse) error
}

// WarehouseService is an interface that contains the methods that the warehouse service should support
type WarehouseService interface {
	// FindAll returns all the warehouses
	FindAll() ([]Warehouse, error)
	// FindByID returns the warehouse with the given ID
	FindByID(id int) (Warehouse, error)
	// Save saves the given warehouse
	Save(warehouse *Warehouse) error
}
