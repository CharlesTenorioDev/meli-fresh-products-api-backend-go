package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type ProductRecordsHandlerDefault struct {
	pd internal.ProductRecordsService
}

// NewProductRecordsDefault creates a new instance of ProductRecordsHandlerDefault.
func NewProductRecordsDefault(pd internal.ProductRecordsService) *ProductRecordsHandlerDefault {
	return &ProductRecordsHandlerDefault{
		pd: pd,
	}
}

// Create handles the creation of a Product Record.
func (h *ProductRecordsHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var productRec internal.ProductRecords

	// Decodifica o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&productRec); err != nil {
		response.JSON(w, http.StatusUnprocessableEntity, "JSON inválido")

		return
	}

	// Chama o serviço para criar o registro
	createdProductRec, err := h.pd.Create(productRec)
	if err != nil {
		if errors.Is(err, service.ErrProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		}

		if errors.Is(err, service.ErrProductNotExists) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		}

		return
	}

	productRecJSON := internal.ProductRecordsJSON{
		LastUpdateDate: createdProductRec.LastUpdateDate,
		PurchasePrice:  createdProductRec.PurchasePrice,
		SalePrice:      createdProductRec.SalePrice,
		ProductID:      createdProductRec.ProductID,
	}
	// Retorna o registro criado com status 201
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": productRecJSON,
	})
}
