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
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
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
	newProduct, err := h.s.Create(product)
	if err != nil {
		if errors.Is(err, service.SellerNotExists) || errors.Is(err, service.ProductTypeNotExists) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else if errors.Is(err, service.ProductCodeAlreadyExists) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else if errors.Is(err, service.ProductUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		} else {
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		}
		return
	}
	response.JSON(w, http.StatusCreated, map[string]any{
		"data": map[string]any{
			"product_id": newProduct.Id,
		},
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

	updatedProduct, err := h.s.Update(product)
	if err != nil {
		if errors.Is(err, service.SellerNotExists) || errors.Is(err, service.ProductTypeNotExists) || errors.Is(err, service.ProductNotExists) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else if errors.Is(err, service.ProductCodeAlreadyExists) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else if errors.Is(err, service.ProductUnprocessableEntity) {
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

	if err := h.s.Delete(id); err != nil {

		if err.Error() == "product not found" {
            response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError("product not found"))
            return
        }

		response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
