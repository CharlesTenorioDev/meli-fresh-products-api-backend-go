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

func NewHandlerProductBatch(svc internal.ProductBatchService) *ProductBatchHandler {
	return &ProductBatchHandler{
		sv: svc,
	}
}

type ProductBatchHandler struct {
	sv internal.ProductBatchService
}

type RequestProductBatchJSON struct {
	BatchNumber        int     `json:"batch_number"`
	CurrentQuantity    int     `json:"current_quantity"`
	CurrentTemperature float64 `json:"current_temperature"`
	DueDate            string  `json:"due_date"`
	InitialQuantity    int     `json:"initial_quantity"`
	ManufacturingDate  string  `json:"manufacturing_date"`
	ManufacturingHour  int     `json:"manufacturing_hour"`
	MinumumTemperature float64 `json:"minumum_temperature"`
	ProductID          int     `json:"product_id"`
	SectionID          int     `json:"section_id"`
}

// GetByID godoc
// @Summary Get product batch by Id
// @Description Fetch the details of a product batch using its unique Id
// @Tags ProductBatch
// @Accept json
// @Produce json
// @Param id path int true "Product Batch ID"
// @Success 200 {object} map[string]any "Product batch data"
// @Failure 400 {object} resterr.RestErr "Invalid Id format"
// @Failure 404 {object} resterr.RestErr "Product-batch not found"
// @Router /api/v1/product_batches/{id} [get]
func (h *ProductBatchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	var prodBatch internal.ProductBatch

	prodBatch, err = h.sv.FindByID(id)
	if err != nil {
		if errors.Is(err, internal.ErrProductBatchNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))

			return
		}

		response.JSON(w, http.StatusInternalServerError, nil)

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": prodBatch,
	})
}

// Create godoc
// @Summary Create a new product batch
// @Description Create a new product batch on the database
// @Tags ProductBatch
// @Accept json
// @Produce json
// @Param product_batch body RequestProductBatchJSON true "Product batch details"
// @Success 201 {object} map[string]any "Created product batch"
// @Failure 400 {object} resterr.RestErr "Invalid input format"
// @Failure 409 {object} resterr.RestErr "Product-batch with given product-batch number already registered" or "Product-batch already exists"
// @Failure 422 {object} resterr.RestErr "Couldn't parse product-batch"
// @Router /api/v1/product_batches [post]
func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var prodBatchJSON map[string]any
	if err := json.NewDecoder(r.Body).Decode(&prodBatchJSON); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	requiredFields := []string{
		"batch_number", "current_quantity", "current_temperature",
		"due_date", "initial_quantity", "manufacturing_date",
		"manufacturing_hour", "minumum_temperature", "product_id",
		"section_id",
	}

	for _, field := range requiredFields {
		if prodBatchJSON[field] == nil {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(field+" is required"))
			return
		}
	}

	prodBatch := internal.ProductBatch{
		BatchNumber:        int(prodBatchJSON["batch_number"].(float64)),
		CurrentQuantity:    int(prodBatchJSON["current_quantity"].(float64)),
		CurrentTemperature: prodBatchJSON["current_temperature"].(float64),
		DueDate:            prodBatchJSON["due_date"].(string),
		InitialQuantity:    int(prodBatchJSON["initial_quantity"].(float64)),
		ManufacturingDate:  prodBatchJSON["manufacturing_date"].(string),
		ManufacturingHour:  int(prodBatchJSON["manufacturing_hour"].(float64)),
		MinumumTemperature: prodBatchJSON["minumum_temperature"].(float64),
		ProductID:          int(prodBatchJSON["product_id"].(float64)),
		SectionID:          int(prodBatchJSON["section_id"].(float64)),
	}

	err := h.sv.Save(&prodBatch)
	if err != nil {
		if errors.Is(err, internal.ErrProductBatchAlreadyExists) || errors.Is(err, internal.ErrProductBatchNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		}

		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": prodBatch,
	})
}
