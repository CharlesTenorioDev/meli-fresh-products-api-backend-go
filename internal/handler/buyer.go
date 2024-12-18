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

func (h *BuyerHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	all := h.s.GetAll()

	response.JSON(w, http.StatusOK, map[string]any{
		"data":    all,
	})
}

func (h *BuyerHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("failed to parse id"))
		return
	}

	buyer, err := h.s.FindByID(id)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data":    buyer,
	})
}

func (h *BuyerHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var buyer internal.Buyer
	err := json.NewDecoder(r.Body).Decode(&buyer)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	err = h.s.Save(&buyer)
	if err != nil {
		if (errors.Is(err, service.BuyerAlreadyExists) || errors.Is(err, service.CardNumberAlreadyInUse)){
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data":    buyer,
	})
}

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
		if errors.Is(err, service.BuyerNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data":    buyer,
	})
}

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