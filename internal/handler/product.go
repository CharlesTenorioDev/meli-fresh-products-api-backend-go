package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
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
		if errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			return
		}
		response.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": products,
	})
}

func (h *ProductHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		err = internal.ErrProductBadRequest
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}
	product, err := h.s.GetByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		return
	}
	response.JSON(w, http.StatusOK, map[string]any{
		"data": product,
	})
}

func (h *ProductHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}
	product, err := h.s.Create(product)
	fmt.Print(product)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrSellerIdNotFound),
			errors.Is(err, internal.ErrProductTypeIdNotFound),
			errors.Is(err, internal.ErrProductNotFound):
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))

		case errors.Is(err, internal.ErrProductCodeAlreadyExists):
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))

		case errors.Is(err, internal.ErrProductUnprocessableEntity):
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))

		default:
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		}
		return
	}
	productJson := internal.ProductJsonPost{
		ProductCode:                    product.ProductCode,
		Description:                    product.Description,
		Height:                         product.Height,
		Length:                         product.Length,
		NetWeight:                      product.NetWeight,
		ExpirationRate:                 product.ExpirationRate,
		RecommendedFreezingTemperature: product.RecommendedFreezingTemperature,
		Width:                          product.Width,
		FreezingRate:                   product.FreezingRate,
		ProductTypeId:                  product.ProductTypeId,
		SellerId:                       product.SellerId,
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": productJson,
	})
}

func (h *ProductHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	product.Id = id
	fmt.Print(product)

	updatedProduct, err := h.s.Update(product)
	fmt.Print(updatedProduct)
	if err != nil {
		if errors.Is(err, internal.ErrSellerIdNotFound) || errors.Is(err, internal.ErrProductTypeIdNotFound) || errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else if errors.Is(err, internal.ErrProductCodeAlreadyExists) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else if errors.Is(err, internal.ErrProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		} else {
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
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
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	err = h.s.Delete(id)
	fmt.Printf("erro do handler %v", err)
	if err != nil {
		if errors.Is(err, internal.ErroProductConflit) || errors.Is(err, internal.ErroProductConflitEntity) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
			return
		}
		if errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			return
		}

		response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
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

		report, err := h.s.GetByIdRecord(productID)
		if err != nil {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError("product not found"))
			return
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data": report,
		})
		return
	}
	report, err := h.s.GetAllRecord()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"data": report,
	})
}
