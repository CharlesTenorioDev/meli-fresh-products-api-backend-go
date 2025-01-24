package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type ProductRecordsHandlerDefault struct {
	s internal.ProductRecordsService
}

// NewProductRecordsDefault creates a new instance of ProductRecordsHandlerDefault.
func NewProductRecordsDefault(pd internal.ProductRecordsService) *ProductRecordsHandlerDefault {
	return &ProductRecordsHandlerDefault{
		s: pd,
	}
}

// Create handles the creation of a Product Record.
// Create godoc
// @Summary Create a product record
// @Description Creates a new product record with details on the database.
// @Tags ProductRecords
// @Accept json
// @Produce json
// @Param product_record body internal.ProductRecords true "Product Record Data"
// @Success 201 {object} map[string]interface{} "Created product record"
// @Failure 422 {object} resterr.RestErr "Invalid JSON"
// @Failure 409 {object} resterr.RestErr "Error ID doesn't exists"
// @Router /api/v1/productRecords [post]
func (h *ProductRecordsHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var productRec internal.ProductRecords

	// Decodifica o corpo da requisição JSON
	if err := json.NewDecoder(r.Body).Decode(&productRec); err != nil {
		response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		
		return
	}

	// Chama o serviço para criar o registro
	createdProductRec, err := h.s.Create(productRec)
	if err != nil {
		if errors.Is(err, internal.ErrProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		}

		if errors.Is(err, internal.ErrProductIdNotFound) {
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
