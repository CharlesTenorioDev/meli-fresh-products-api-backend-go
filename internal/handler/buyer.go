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

type BuyerHandlerDefault struct {
	s internal.BuyerService
}

func NewBuyerHandlerDefault(svc internal.BuyerService) *BuyerHandlerDefault {
	return &BuyerHandlerDefault{
		s: svc,
	}
}

// GetAll godoc
// @Summary Get all buyers
// @Description Retrieve all buyers from the database
// @Tags Buyers
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all buyers"
// @Router /api/v1/buyers [get]
func (h *BuyerHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	all := h.s.GetAll()

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
}

// GetByID godoc
// @Summary Get a buyer by Id
// @Description Retrieve a specific buyer from the database using their Id
// @Tags Buyers
// @Accept json
// @Produce json
// @Param id path int true "Buyer ID"
// @Success 200 {object} map[string]interface{} "Buyer data"
// @Failure 400 {object} resterr.RestErr "Failed to parse Id"
// @Failure 404 {object} resterr.RestErr "Buyer not found"
// @Router /api/v1/buyers/{id} [get]
func (h *BuyerHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse id"))
		return
	}

	buyer, err := h.s.FindByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": buyer,
	})
}

// Create godoc
// @Summary Create a new buyer
// @Description Add a new buyer to the database
// @Tags Buyers
// @Accept json
// @Produce json
// @Param buyer body internal.Buyer true "Buyer data"
// @Success 201 {object} map[string]interface{} "Created buyer"
// @Failure 409 {object} resterr.RestErr "buyer with given card number already registered"
// @Failure 400 {object} resterr.RestErr "Invalid input"
// @Failure 422 {object} resterr.RestErr "Failed to create buyer"
// @Router /api/v1/buyers [post]
func (h *BuyerHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var buyer internal.Buyer
	err := json.NewDecoder(r.Body).Decode(&buyer)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	err = h.s.Save(&buyer)
	if err != nil {
		if errors.Is(err, service.ErrBuyerAlreadyExists) || errors.Is(err, service.ErrCardNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": buyer,
	})
}

// Update godoc
// @Summary Update a buyer's information
// @Description Update the details of an existing buyer in the database
// @Tags Buyers
// @Accept json
// @Produce json
// @Param id path int true "Buyer ID"
// @Param buyer body internal.BuyerPatch true "Buyer patch data"
// @Success 200 {object} map[string]interface{} "Updated buyer"
// @Failure 400 {object} resterr.RestErr "Failed to parse id" or "Failed to parse body"
// @Failure 404 {object} resterr.RestErr "Buyer not found"
// @Failure 409 {object} resterr.RestErr "buyer with given card number already registered"
// @Router /api/v1/buyers/{id} [patch]
func (h *BuyerHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse id"))
		return
	}

	var buyer internal.BuyerPatch
	err = json.NewDecoder(r.Body).Decode(&buyer)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse body"))
		return
	}

	err = h.s.Update(id, buyer)
	if err != nil {
		if errors.Is(err, service.ErrBuyerNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": buyer,
	})
}

// Delete godoc
// @Summary Delete a buyer by Id
// @Description Remove a specific buyer from the database
// @Tags Buyers
// @Accept json
// @Produce json
// @Param id path int true "Buyer ID"
// @Success 204 {object} nil "No content"
// @Failure 400 {object} resterr.RestErr "Failed to parse Id"
// @Failure 404 {object} resterr.RestErr "Buyer not found"
// @Router /api/v1/buyers/{id} [delete]
func (h *BuyerHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse id"))
		return
	}

	err = h.s.Delete(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}

// ReportPurchaseOrders godoc
// @Summary Get purchase orders by buyer
// @Description Generate a report of purchase orders for a specific buyer or all buyers
// @Tags Buyers
// @Accept json
// @Produce json
// @Param id query int false "Buyer Id"
// @Success 200 {object} map[string]interface{} "Report data"
// @Failure 400 {object} resterr.RestErr "failed to parse id"
// @Failure 404 {object} resterr.RestErr "Buyer not found"
// @Failure 500 {object} resterr.RestErr "Internal server error"
// @Router /api/v1/buyers/report-purchase-orders [get]
func (h *BuyerHandlerDefault) ReportPurchaseOrders(w http.ResponseWriter, r *http.Request) {
	var purchaseOrdersByBuyer []internal.PurchaseOrdersByBuyer
	var err error

	// Check if there is an id query parameter and call the corresponding service method
	id := r.URL.Query().Get("id")
	if id != "" {
		idInt, parseErr := strconv.Atoi(id)
		if parseErr != nil {
			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse id"))
			return
		}
		purchaseOrdersByBuyer, err = h.s.ReportPurchaseOrdersById(idInt)
	} else {
		purchaseOrdersByBuyer, err = h.s.ReportPurchaseOrders()
	}

	if err != nil {
		switch {
		case errors.Is(err, service.ErrBuyerNotFound):
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		default:
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": purchaseOrdersByBuyer,
	})
}
