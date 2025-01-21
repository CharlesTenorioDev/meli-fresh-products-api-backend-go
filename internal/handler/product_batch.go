package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
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
	ProductId          int     `json:"product_id"`
	SectionId          int     `json:"section_id"`
}

// GetByID godoc
// @Summary Get product batch by Id
// @Description Fetch the details of a product batch using its unique Id
// @Tags ProductBatch
// @Accept json
// @Produce json
// @Param id path int true "Product Batch ID"
// @Success 200 {object} map[string]any "Product batch data"
// @Failure 400 {object} rest_err.RestErr "Invalid Id format"
// @Failure 404 {object} rest_err.RestErr "Product-batch not found"
// @Router /api/v1/product_batches/{id} [get]
func (h *ProductBatchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	var prodBatch internal.ProductBatch
	prodBatch, err = h.sv.FindByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
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
// @Failure 400 {object} rest_err.RestErr "Invalid input format"
// @Failure 409 {object} rest_err.RestErr "Product-batch with given product-batch number already registered" or "Product-batch already exists"
// @Failure 422 {object} rest_err.RestErr "Couldn't parse product-batch"
// @Router /api/v1/product_batches [post]
func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var prodBatchJSON RequestProductBatchJSON
	if err := json.NewDecoder(r.Body).Decode(&prodBatchJSON); err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	prodBatch := internal.ProductBatch{
		BatchNumber:        prodBatchJSON.BatchNumber,
		CurrentQuantity:    prodBatchJSON.CurrentQuantity,
		CurrentTemperature: prodBatchJSON.CurrentTemperature,
		DueDate:            prodBatchJSON.DueDate,
		InitialQuantity:    prodBatchJSON.InitialQuantity,
		ManufacturingDate:  prodBatchJSON.ManufacturingDate,
		ManufacturingHour:  prodBatchJSON.ManufacturingHour,
		MinumumTemperature: prodBatchJSON.MinumumTemperature,
		ProductId:          prodBatchJSON.ProductId,
		SectionId:          prodBatchJSON.SectionId,
	}

	err := h.sv.Save(&prodBatch)
	if err != nil {
		if errors.Is(err, internal.ErrProductBatchNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": prodBatch,
	})
}
