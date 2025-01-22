package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type InboundOrdersHandler struct {
	sv internal.InboundOrderService
}

func NewInboundOrdersHandler(sv internal.InboundOrderService) *InboundOrdersHandler {
	return &InboundOrdersHandler{
		sv: sv,
	}
}

// Create godoc
// @Summary Create a new inbound order
// @Description Create a new inbound order with the provided details
// @Tags InboundOrders
// @Accept json
// @Produce json
// @Param inbound body internal.InboundOrders true "Inbound order data"
// @Success 201 {object} map[string]interface{} "Created inbound order with ID"
// @Failure 400 {object} map[string]interface{} "Invalid body format"
// @Failure 422 {object} map[string]interface{} "Required fields are missing"
// @Failure 409 {object} map[string]interface{} "Order number already exists" or "Employee not exists"
// @Router /api/v1/inbound-orders [post]
func (h *InboundOrdersHandler) Create(w http.ResponseWriter, r *http.Request) {
	var inbound internal.InboundOrders

	err := json.NewDecoder(r.Body).Decode(&inbound)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid body format", //status code 400
		})

		return
	}

	if okFields := inbound.ValidateFieldsOk(); !okFields {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
			"error": "required fields are missing", //status code 422
		})

		return
	}

	lastID, err := h.sv.Create(inbound)
	if err != nil {
		if err == internal.ErrOrderNumberAlreadyExists {
			response.JSON(w, http.StatusConflict, map[string]any{
				"error": "order number already exists", //status code 409
			})

			return
		}

		if err == internal.ErrEmployeeNotFound {
			response.JSON(w, http.StatusConflict, map[string]any{
				"error": "employee not exists", //status code 409
			})

			return
		}
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": struct {
			ID int64 `json:"id"`
		}{
			ID: lastID, //last id generated
		},
	})
}

// GetAll godoc
// @Summary Get all inbound orders
// @Description Retrieve a list of all inbound orders from the database
// @Tags InboundOrders
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of inbound orders"
// @Failure 500 {object} rest_err.RestErr "Failed to fetch inbounds orders"
// @Router /api/v1/inbound-orders [get]
func (h *InboundOrdersHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	allInbounds, err := h.sv.FindAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError("failed to fetch inbounds orders"))

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": allInbounds,
	})
}
