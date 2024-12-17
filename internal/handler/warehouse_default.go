package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
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

type WarehouseCreateRequest struct {
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}

var (
	InternalError = "Internal Server Error"
)

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

// GetByID returns a warehouse by its ID
func (h *WarehouseDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		warehouse, err := h.sv.FindByID(idInt)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryNotFound):
				response.Error(w, http.StatusNotFound, "Warehouse not found")
			default:
				response.Error(w, http.StatusInternalServerError, InternalError)
			}
			return
		}

		warehouseJson := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseJson,
		})
	}
}

// Create creates a new warehouse
func (h *WarehouseDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestInput *WarehouseCreateRequest

		// decode the request
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid data")
			return
		}

		// create the warehouse
		warehouse := internal.Warehouse{
			ID:                 0,
			WarehouseCode:      requestInput.WarehouseCode,
			Address:            requestInput.Address,
			Telephone:          requestInput.Telephone,
			MinimumCapacity:    requestInput.MinimumCapacity,
			MinimumTemperature: requestInput.MinimumTemperature,
		}

		// save the warehouse
		err := h.sv.Save(&warehouse)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryDuplicated):
				response.Error(w, http.StatusConflict, "Warehouse already exists")
			default:
				response.Error(w, http.StatusInternalServerError, InternalError)
			}
			return
		}

		// return the warehouse
		warehouseJson := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": warehouseJson,
		})
	}
}

// Update updates a warehouse
func (h *WarehouseDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid ID format")
			return
		}

		// decode the request into a WarehousePatchUpdate
		var requestInput *internal.WarehousePatchUpdate
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid data")
			return
		}

		// Calling the service to update the warehouse
		warehouse, err := h.sv.Update(idInt, requestInput)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryNotFound):
				response.Error(w, http.StatusNotFound, "Warehouse not found")
			default:
				response.Error(w, http.StatusInternalServerError, InternalError)
			}
			return
		}

		// Returning the updated warehouse
		warehouseJson := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseJson,
		})

	}
}
