package handler

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
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
