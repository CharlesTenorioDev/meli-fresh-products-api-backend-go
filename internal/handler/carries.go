package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bootcamp-go/web/response"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type CarriesHandlerDefault struct {
	sv internal.CarriesService
}

func NewCarriesHandlerDefault(sv internal.CarriesService) *CarriesHandlerDefault {
	return &CarriesHandlerDefault{sv}
}

func (h *CarriesHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	all, err := h.sv.FindAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError("failed to fetch carries"))

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": all,
	})
}

func (h *CarriesHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var carry internal.Carries

	err := json.NewDecoder(r.Body).Decode(&carry)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError("failed to parse body"))

		return
	}

	if ok := carry.Ok(); !ok {
		response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError("missing fields"))
		return
	}

	lastID, err := h.sv.Create(carry)
	if err != nil {
		response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))

		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": struct {
			ID int64 `json:"id"`
		}{
			ID: lastID,
		},
	})
}
