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
	ProductId          int     `json:"product_id"`
	SectionId          int     `json:"section_id"`
}

func (h *ProductBatchHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	var prodBatch internal.ProductBatch
	prodBatch, err = h.sv.FindByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": prodBatch,
	})
}

func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var prodBatchJSON RequestProductBatchJSON
	if err := json.NewDecoder(r.Body).Decode(&prodBatchJSON); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
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
		ProductID:          prodBatchJSON.ProductId,
		SectionID:          prodBatchJSON.SectionId,
	}

	err := h.sv.Save(&prodBatch)
	if err != nil {
		if errors.Is(err, internal.ErrProductBatchNumberAlreadyInUse) {
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
