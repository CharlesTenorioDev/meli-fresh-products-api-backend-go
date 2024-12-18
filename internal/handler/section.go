package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	errorss "github.com/meli-fresh-products-api-backend-t1/internal/errors"
)

func NewHandlerSection(svc internal.SectionService) *SectionHandler {
	return &SectionHandler{
		sv: svc,
	}
}

type SectionHandler struct {
	sv internal.SectionService
}

func (h *SectionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sections, err := h.sv.FindAll()
	if err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	RespondWithSections(w, sections, http.StatusOK)
}

func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var section internal.Section
	section, err = h.sv.FindByID(id)
	if err != nil {
		HandleError(w, err.Error(), http.StatusNotFound)
		return
	}

	RespondWithSection(w, section, http.StatusOK)
}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var reqBody RequestBodySection
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	section := internal.Section{
		SectionNumber:      reqBody.SectionNumber,
		CurrentTemperature: reqBody.CurrentTemperature,
		MinimumTemperature: reqBody.MinimumTemperature,
		CurrentCapacity:    reqBody.CurrentCapacity,
		MinimumCapacity:    reqBody.MinimumCapacity,
		MaximumCapacity:    reqBody.MaximumCapacity,
		WarehouseID:        reqBody.WarehouseID,
		ProductTypeID:      reqBody.ProductTypeID,
	}

	err := h.sv.Save(&section)
	if customErr, ok := err.(*errorss.CustomError); ok {
		HandleError(w, customErr.Message, customErr.StatusHttp)
		return
	} else if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	RespondWithSection(w, section, http.StatusCreated)
}

func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.sv.FindByID(id)
	if err != nil {
		HandleError(w, err.Error(), http.StatusNotFound)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	update, err := h.sv.Update(id, updates)
	if err != nil {
		HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	RespondWithSection(w, update, http.StatusOK)
}

func (h *SectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.sv.FindByID(id)
	if err != nil {
		HandleError(w, err.Error(), http.StatusNotFound)
		return
	}

	err = h.sv.Delete(id)
	if err != nil {
		HandleError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
