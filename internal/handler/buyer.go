package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
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
		"message": "success",
		"data":    all,
	})
}

func (h *BuyerHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse id",
		})
		return
	}

	buyer, err := h.s.FindByID(id - 1)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "success",
		"data":    buyer,
	})
}

func (h *BuyerHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse id",
		})
		return
	}

	var buyer internal.Buyer
	err = json.NewDecoder(r.Body).Decode(&buyer)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse body",
		})
		return
	}

	ok := buyer.Parse()
	if !ok {
		response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
			"message": "failed",
			"data":    "failed to parse entity",
		})
		return
	}

	buyer.ID = id
	err = h.s.Save(id, buyer)
	if err != nil {
		response.JSON(w, http.StatusConflict, map[string]any{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"message": "success",
		"data":    buyer,
	})
}

func (h *BuyerHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse id",
		})
		return
	}

	var buyer internal.BuyerPatch
	err = json.NewDecoder(r.Body).Decode(&buyer)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse body",
		})
		return
	}

	err = h.s.Update(id - 1, buyer)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]any{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"message": "success",
		"data":    buyer,
	})
}

func (h *BuyerHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"message": "failed",
			"data":    "failed to parse id",
		})
		return
	}

	err = h.s.Delete(id - 1)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]any{
			"message": "failed",
			"data":    err.Error(),
		})
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
