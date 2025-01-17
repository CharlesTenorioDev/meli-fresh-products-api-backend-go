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

	lastId, err := h.sv.Create(inbound)
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
			Id int64 `json:"id"`
		}{
			Id: lastId, //last id generated
		},
	})
}

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
