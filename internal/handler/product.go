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

type ProductHandlerDefault struct {
	s internal.ProductService
}

func NewProducHandlerDefault(phd internal.ProductService) *ProductHandlerDefault {
	return &ProductHandlerDefault{s: phd}
}

// GetAll godoc
// @Summary Get all products
// @Description Retrieves a list of all products in the database
// @Tags Product
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all products"
// @Failure 400 {object} resterr.RestErr "Bad request"
// @Router /api/v1/products [get]
func (h *ProductHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.s.GetAll()
	if err != nil {
		if errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			return
		}

		response.JSON(w, http.StatusInternalServerError, nil)

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": products,
	})
}

// GetByID godoc
// @Summary Get product by ID
// @Description Retrieves a single product by its Id
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]any "Product data"
// @Failure 400 {object} resterr.RestErr "Invalid Id format"
// @Failure 404 {object} resterr.RestErr "Product not found"
// @Router /api/v1/products/{id} [get]
func (h *ProductHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		err = internal.ErrProductBadRequest
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

// Create godoc
// @Summary Create a new product
// @Description Adds a new product to the system with the provided details in the request body
// @Tags Product
// @Accept json
// @Produce json
// @Param product body internal.Product true "Product data"
// @Success 201 {object} map[string]any "Product created"
// @Failure 400 {object} resterr.RestErr "Invalid request body"
// @Failure 404 {object} resterr.RestErr "Seller or Product Type not exists"
// @Failure 409 {object} resterr.RestErr "Product code already exists"
// @Failure 422 {object} resterr.RestErr "Unprocessable entity"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Router /api/v1/products [post]
func (h *ProductHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var product internal.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	product, err := h.s.Create(product)

	if err != nil {
		switch {
		case errors.Is(err, internal.ErrSellerIdNotFound),
			errors.Is(err, internal.ErrProductTypeIDNotFound),
			errors.Is(err, internal.ErrProductNotFound):
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))

		case errors.Is(err, internal.ErrProductCodeAlreadyExists),
			errors.Is(err, internal.ErrProductConflitEntity),
			errors.Is(err, internal.ErrProductConflit):
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))

		case errors.Is(err, internal.ErrProductUnprocessableEntity):
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))

		default:
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
		}

		return
	}

	productJSON := internal.ProductJSONPost{
		ProductCode:                    product.ProductCode,
		Description:                    product.Description,
		Height:                         product.Height,
		Length:                         product.Length,
		NetWeight:                      product.NetWeight,
		ExpirationRate:                 product.ExpirationRate,
		RecommendedFreezingTemperature: product.RecommendedFreezingTemperature,
		Width:                          product.Width,
		FreezingRate:                   product.FreezingRate,
		ProductTypeID:                  product.ProductTypeID,
		SellerID:                       product.SellerID,
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": productJSON,
	})
}

// Update godoc
// @Summary Update a product
// @Description Updates an existing product's data in the request body
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body internal.Product true "Updated product data"
// @Success 200 {object} map[string]any "Updated product"
// @Failure 400 {object} resterr.RestErr "Invalid request body"
// @Failure 404 {object} resterr.RestErr "Seller or Product Type not exists"
// @Failure 409 {object} resterr.RestErr "Product code already exists"
// @Failure 422 {object} resterr.RestErr "All fields must be valid and filled"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Failure 400 {object} resterr.RestErr "Invalid request body"
// @Failure 404 {object} resterr.RestErr "Seller or Product Type not exists"
// @Failure 409 {object} resterr.RestErr "Product code already exists"
// @Failure 422 {object} resterr.RestErr "All fields must be valid and filled"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Router /api/v1/products/{id} [patch]
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
		if errors.Is(err, internal.ErrSellerIdNotFound) || errors.Is(err, internal.ErrProductTypeIDNotFound) || errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		} else if errors.Is(err, internal.ErrProductCodeAlreadyExists) || errors.Is(err, internal.ErrProductConflit) || errors.Is(err, internal.ErrProductConflitEntity) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		} else if errors.Is(err, internal.ErrProductUnprocessableEntity) {
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

// Delete godoc
// @Summary Delete a product
// @Description Deletes a product by its Id
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204 {object} nil "No content"
// @Failure 400 {object} resterr.RestErr "Invalid Id format"
// @Failure 404 {object} resterr.RestErr "Product not found"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Failure 400 {object} resterr.RestErr "Invalid Id format"
// @Failure 404 {object} resterr.RestErr "Product not found"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Router /api/v1/products/{id} [delete]
func (h *ProductHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	err = h.s.Delete(id)

	if err != nil {
		if errors.Is(err, internal.ErrProductConflit) || errors.Is(err, internal.ErrProductConflitEntity) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
			return
		}

		if errors.Is(err, internal.ErrProductNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			return
		}

		response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))

		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// ReportRecords godoc
// @Summary Get product records
// @Description Retrieves records of products, or a specific record by product Id
// @Tags Product
// @Accept json
// @Produce json
// @Param id query int false "Product ID"
// @Success 200 {object} map[string]interface{} "Product records"
// @Failure 400 {object} resterr.RestErr "Invalid Id"
// @Failure 404 {object} resterr.RestErr "Product not found"
// @Failure 400 {object} resterr.RestErr "Invalid Id"
// @Failure 404 {object} resterr.RestErr "Product not found"
// @Router /api/v1/products/report-records [get]
func (h *ProductHandlerDefault) ReportRecords(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		productID, err := strconv.Atoi(id)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, resterr.NewNotFoundError(err.Error()))

			return
		}

		report, err := h.s.GetByIDRecord(productID)
		if err != nil {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
			return
		}

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
