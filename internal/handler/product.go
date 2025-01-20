package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type ProductHandlerDefault struct {
	s internal.ProductService
}

func NewProducHandlerDefault(phd internal.ProductService) *ProductHandlerDefault {
	return &ProductHandlerDefault{s: phd}
}

func (h *ProductHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.s.GetAll()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
}

func (h *ProductHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	product, err := h.s.GetByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (h *ProductHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	newProduct, err := h.s.Create(product)
	if err != nil {
		if errors.Is(err, service.ErrSellerNotExists) || errors.Is(err, service.ErrProductTypeNotExists) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		} else if errors.Is(err, service.ErrProductCodeAlreadyExists) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		} else if errors.Is(err, service.ErrProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		} else {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
		}

		return
	}

	var productJSON internal.ProductJSONPost
	productJSON.ProductCode = newProduct.ProductCode
	productJSON.Description = newProduct.Description
	productJSON.Height = newProduct.Height
	productJSON.Length = newProduct.Length
	productJSON.NetWeight = newProduct.NetWeight
	productJSON.ExpirationRate = newProduct.ExpirationRate
	productJSON.RecommendedFreezingTemperature = newProduct.RecommendedFreezingTemperature
	productJSON.Width = newProduct.Width
	productJSON.FreezingRate = newProduct.FreezingRate
	productJSON.ProductTypeID = newProduct.ProductTypeID
	productJSON.SellerID = newProduct.SellerID

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": productJSON,
	})
}

func (h *ProductHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	product.ID = id

	updatedProduct, err := h.s.Update(product)
	if err != nil {
		if errors.Is(err, service.ErrSellerNotExists) || errors.Is(err, service.ErrProductTypeNotExists) || errors.Is(err, service.ErrProductNotExists) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		} else if errors.Is(err, service.ErrProductCodeAlreadyExists) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		} else if errors.Is(err, service.ErrProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		} else {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
		}

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": updatedProduct,
	})
}

func (h *ProductHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	if err := h.s.Delete(id); err != nil {
		if err.Error() == "product not found" {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError("product not found"))

			return
		}

		response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))

		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

func (h *ProductHandlerDefault) ReportRecords(w http.ResponseWriter, r *http.Request) {
	// Extrair o parâmetro "id" da URL
	id := r.URL.Query().Get("id")

	if id != "" {
		productID, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "ID inválido")

			return
		}

		report, err := h.s.GetByIDRecord(productID)
		if err != nil {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError("product not found"))

			return
		}

		// Retorna o relatório
		response.JSON(w, http.StatusOK, map[string]any{
			"data": report,
		})

		return
	}

	report, err := h.s.GetAllRecord()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": report,
	})
}
