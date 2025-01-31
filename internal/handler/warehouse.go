package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

// WarehouseJSON represents the warehouse in JSON format.
type WarehouseJSON struct {
	ID                 int     `json:"id"`
	WarehouseCode      string  `json:"warehouse_code"`
	Address            string  `json:"address"`
	Telephone          string  `json:"telephone"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MinimumTemperature float64 `json:"minimum_temperature"`
}

type WarehouseCreateRequest struct {
	WarehouseCode      *string  `json:"warehouse_code"`
	Address            *string  `json:"address"`
	Telephone          *string  `json:"telephone"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
}

var (
	ErrInternalServer = "Internal Server Error"
	ErrInvalidID      = "Invalid ID format"
	ErrInvalidData    = "Invalid data"
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
// @Summary Get all warehouses
// @Description Retrieve a list of all warehouses in the database
// @Tags Warehouse
// @Produce json
// @Success 200 {object} map[string]any "List of all warehouses"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/warehouses [get]
func (h *WarehouseDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Find all warehouses
		warehouses, err := h.sv.FindAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
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

// GetByID returns a warehouse by id
// @Summary Get warehouse by Id
// @Description Retrieve a warehouse's details by its Id
// @Tags Warehouse
// @Produce json
// @Param id path int true "Warehouse ID"
// @Success 200 {object} WarehouseJSON "Warehouse data"
// @Failure 400 {object} resterr.RestErr "Invalid ID format"
// @Failure 404 {object} resterr.RestErr "Warehouse not found"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/warehouses/{id} [get]
func (h *WarehouseDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(ErrInvalidID))
			return
		}

		warehouse, err := h.sv.FindByID(idInt)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryNotFound):
				response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
			}

			return
		}

		warehouseJSON := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseJSON,
		})
	}
}

// Create creates a new warehouse
// @Summary Create a new warehouse
// @Description Add a new warehouse to the database
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param warehouse body WarehouseCreateRequest true "Warehouse data"
// @Success 201 {object} WarehouseJSON "Created warehouse"
// @Failure 400 {object} resterr.RestErr "Invalid Data"
// @Failure 409 {object} resterr.RestErr "Warehouse already exists"
// @Failure 422 {object} resterr.RestErr "Unprocessable Entity"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/warehouses [post]
func (h *WarehouseDefault) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestInput *WarehouseCreateRequest

		// decode the request
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(ErrInvalidData))
			return
		}

		// validating the request input required fields
		causes := requestInput.ValidateRequiredFields()
		if len(causes) > 0 {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityWithCausesError(internal.ErrWarehouseUnprocessableEntity.Error(), causes))
			return
		}

		// create the warehouse
		warehouse := internal.Warehouse{
			ID:                 0,
			WarehouseCode:      *requestInput.WarehouseCode,
			Address:            *requestInput.Address,
			Telephone:          *requestInput.Telephone,
			MinimumCapacity:    *requestInput.MinimumCapacity,
			MinimumTemperature: *requestInput.MinimumTemperature,
		}

		// save the warehouse
		err := h.sv.Save(&warehouse)
		if err != nil {
			switch {
			case errors.As(err, &internal.DomainError{}):
				var domainError internal.DomainError
				errors.As(err, &domainError)
				var restCauses []resterr.Causes
				for _, cause := range domainError.Causes {
					restCauses = append(restCauses, resterr.Causes{
						Field:   cause.Field,
						Message: cause.Message,
					})
				}
				response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestValidationError(domainError.Message, restCauses))
			case errors.Is(err, internal.ErrWarehouseRepositoryDuplicated):
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			case errors.Is(err, internal.ErrWarehouseUnprocessableEntity):
				response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
			}

			return
		}

		// return the warehouse
		warehouseJSON := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusCreated, map[string]any{
			"data": warehouseJSON,
		})
	}
}

// Update updates a warehouse
// @Summary Update warehouse details
// @Description Modify an existing warehouse's data
// @Tags Warehouse
// @Accept json
// @Produce json
// @Param id path int true "Warehouse ID"
// @Param warehouse body internal.WarehousePatchUpdate true "Updated warehouse data"
// @Success 200 {object} WarehouseJSON "Updated warehouse"
// @Failure 400 {object} resterr.RestErr "Invalid ID format" or "Invalid Data"
// @Failure 404 {object} resterr.RestErr "Warehouse not found"
// @Failure 409 {object} resterr.RestErr "Warehouse already exists"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/warehouses/{id} [patch]
func (h *WarehouseDefault) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(ErrInvalidID))
			return
		}

		// decode the request into a WarehousePatchUpdate
		var requestInput *internal.WarehousePatchUpdate
		if err := json.NewDecoder(r.Body).Decode(&requestInput); err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(ErrInvalidData))
			return
		}

		// Calling the service to update the warehouse
		warehouse, err := h.sv.Update(idInt, requestInput)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryDuplicated):
				response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			case errors.Is(err, internal.ErrWarehouseRepositoryNotFound):
				response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
			}

			return
		}

		// Returning the updated warehouse
		warehouseJSON := WarehouseJSON{
			ID:                 warehouse.ID,
			WarehouseCode:      warehouse.WarehouseCode,
			Address:            warehouse.Address,
			Telephone:          warehouse.Telephone,
			MinimumCapacity:    warehouse.MinimumCapacity,
			MinimumTemperature: warehouse.MinimumTemperature,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": warehouseJSON,
		})
	}
}

// Delete deletes a warehouse
// @Summary Delete warehouse
// @Description Removes a warehouse from the database by its ID
// @Tags Warehouse
// @Param id path int true "Warehouse ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} resterr.RestErr "Invalid ID format"
// @Failure 404 {object} resterr.RestErr "Warehouse not found"
// @Failure 500 {object} resterr.RestErr "Internal Server Error"
// @Router /api/v1/warehouses/{id} [delete]
func (h *WarehouseDefault) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, err := strconv.Atoi(id)

		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(ErrInvalidID))
			return
		}

		err = h.sv.Delete(idInt)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrWarehouseRepositoryNotFound):
				response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
			}

			return
		}

		response.JSON(w, http.StatusNoContent, nil)
	}
}

// Validating the WarehouseCreateRequest required fields
func (r *WarehouseCreateRequest) ValidateRequiredFields() (causes []resterr.Causes) {
	if r.WarehouseCode == nil {
		causes = append(causes, resterr.Causes{
			Field:   "warehouse_code",
			Message: "warehouse code is required",
		})
	}
	if r.Address == nil {
		causes = append(causes, resterr.Causes{
			Field:   "address",
			Message: "address is required",
		})
	}
	if r.Telephone == nil {
		causes = append(causes, resterr.Causes{
			Field:   "telephone",
			Message: "telephone is required",
		})
	}
	if r.MinimumCapacity == nil {
		causes = append(causes, resterr.Causes{
			Field:   "minimum_capacity",
			Message: "minimum capacity is required",
		})
	}
	if r.MinimumTemperature == nil {
		causes = append(causes, resterr.Causes{
			Field:   "minimum_temperature",
			Message: "minimum temperature is required",
		})
	}
	return causes
}
