package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
func (h *ProductHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	product, err := h.s.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {
	var product internal.Product

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	newProduct, err := h.s.Create(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(newProduct) // Retorna o produto criado
}

func (h *ProductHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do produto a partir da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var product internal.Product
	// Decodifica o corpo da requisição JSON para um objeto Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	product.Id = id // Define o ID do produto a ser atualizado

	// Chama o serviço para atualizar o produto
	updatedProduct, err := h.s.Update(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct) // Retorna o produto atualizado
}
func (h *ProductHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {
	// Obtém o ID do produto a partir da URL
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Chama o serviço para deletar o produto
	if err := h.s.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}