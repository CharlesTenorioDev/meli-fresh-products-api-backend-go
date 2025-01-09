package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
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

type ProductBatchJSON struct {
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

func (h *ProductBatchHandler) Create(w http.ResponseWriter, r *http.Request) {
	var prodBatchJSON ProductBatchJSON
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
		if errors.Is(err, service.ProductBatchNumberAlreadyInUse) {
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