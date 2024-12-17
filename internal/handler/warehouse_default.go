package handler

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

// Structure to represent the warehouse in JSON format
type WarehouseJSON struct {
	ID                 int     `json:"id"`
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}

// NewWarehouseDefault Builder creates a new instance of the warehouse handler
func NewWarehouseDefault(sv internal.WarehouseService) *WarehouseDefault {
	return &WarehouseDefault{
		sv: sv,
	}
}

// WarehouseDefault is the default implementation of the warehouse handler
type WarehouseDefault struct {
	// sv is the service used by the handler
	sv internal.WarehouseService
}

// GetAll returns all warehouses
func (h *WarehouseDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find all warehouses
		warehouses, err := h.sv.FindAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		// Data to be returned
		var data []WarehouseJSON
		for _, warehouse := range warehouses {
			data = append(data, WarehouseJSON{
				ID:                 warehouse.ID,
				WarehouseCode:      warehouse.WarehouseCode,
				Address:            warehouse.Address,
				Telephone:          warehouse.Telephone,
				MinimumCapacity:    warehouse.MinimumCapacity,
				MinimumTemperature: warehouse.MinimumTemperature,
			})
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": data,
		})
	}
}
