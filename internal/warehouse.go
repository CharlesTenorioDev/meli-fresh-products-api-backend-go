package internal

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
)

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
	ErrWarehouseUnprocessableEntity = errors.New("warehouse inputs are missing")
	// ErrWarehouseBadRequest is returned when the warehouse request is bad
	ErrWarehouseBadRequest = errors.New("warehouse inputs are invalid")
)

// Validate validates the business rules of the warehouse
func (w *Warehouse) Validate() (causes []Causes) {
	if validator.BlankString(w.WarehouseCode) {
		causes = append(causes, Causes{
			Field:   "warehouse_code",
			Message: "warehouse code is required",
		})
	}

	if !validator.String(w.WarehouseCode, 1, 255) && !validator.BlankString(w.WarehouseCode) {
		causes = append(causes, Causes{
			Field:   "warehouse_code",
			Message: "warehouse code is out of range",
		})
	}

	if validator.BlankString(w.Address) {
		causes = append(causes, Causes{
			Field:   "address",
			Message: "address is required",
		})
	}

	if !validator.String(w.Address, 1, 255) && !validator.BlankString(w.Address) {
		causes = append(causes, Causes{
			Field:   "address",
			Message: "address is out of range",
		})
	}

	if validator.BlankString(w.Telephone) {
		causes = append(causes, Causes{
			Field:   "telephone",
			Message: "telephone is required",
		})
	}

	if !validator.IsTelephone(w.Telephone) {
		causes = append(causes, Causes{
			Field:   "telephone",
			Message: `telephone number is invalid, should be formatted as XX XXXXX-XXXX`,
		})
	}

	if !validator.String(w.Telephone, 1, 255) && !validator.BlankString(w.Telephone) {
		causes = append(causes, Causes{
			Field:   "telephone",
			Message: "telephone is out of range",
		})
	}

	if validator.IntIsZero(w.MinimumCapacity) {
		causes = append(causes, Causes{
			Field:   "minimum_capacity",
			Message: "minimum capacity is required",
		})
	}

	if validator.IntIsNegative(w.MinimumCapacity) {
		causes = append(causes, Causes{
			Field:   "minimum_capacity",
			Message: "minimum capacity cannot be negative",
		})
	}

	if !validator.FloatBetween(w.MinimumTemperature, -273.15, 1000) {
		causes = append(causes, Causes{
			Field:   "minimum_temperature",
			Message: "minimum temperature is out of range",
		})
	}

	return causes
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
