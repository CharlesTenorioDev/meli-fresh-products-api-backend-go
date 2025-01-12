package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

type ProductRecordsHandlerDefault struct {
	pd internal.ProductRecordsService
}

// Construtor do handler
func NewProductRecordsDefault(pd internal.ProductRecordsService) *ProductRecordsHandlerDefault {
	return &ProductRecordsHandlerDefault{
		pd: pd,
	}
}

// Handler para criar um Product Record
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
		if errors.Is(err, service.ProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		}
		if errors.Is(err, service.ProductNotFound) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		}
		return
	}
	productRecJson := internal.ProductRecordsJson{
		LastUpdateDate: createdProductRec.LastUpdateDate,
		PurchasePrice:  createdProductRec.PurchasePrice,
		SalePrice:      createdProductRec.SalePrice,
		ProductID:      createdProductRec.ProductID,
	}
	// Retorna o registro criado com status 201
	response.JSON(w, http.StatusCreated, map[string]interface{}{
		"data": productRecJson,
	})
}

// Handler para gerar relatórios de registros por Product
func (h *ProductRecordsHandlerDefault) GetReportByProduct(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	productIDStr := queryParams.Get("id")

	// Se o parâmetro "id" for fornecido, converte para int
	if productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "ID inválido")
			return
		}

		// Chama o serviço para buscar o relatório por Product ID
		report, err := h.pd.GetByID(productID)
		if err != nil {
			response.JSON(w, http.StatusNotFound, "Produto não encontrado")
			return
		}

		// Retorna o relatório
		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data": report,
		})
		return
	}

	// Se nenhum ID for enviado, retorna o relatório de todos os Products
	allRecords, err := h.pd.GetAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "Erro ao buscar registros")
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": allRecords,
	})
}
