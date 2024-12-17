package handler

import (
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

func (h *BuyerHandlerDefault) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := h.s.GetAll()

		response.JSON(w, http.StatusOK, map[string]any{
			"message": "success",
			"data":    all,
		})
	}
}

func (h *BuyerHandlerDefault) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}
