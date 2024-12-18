package handler

import (
	"encoding/json"
	"net/http"

	"github.com/meli-fresh-products-api-backend-t1/internal"
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
		http.Error(w, "Failed to get products", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)

}
