package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

type RequestSectionJSON struct {
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
}

type ResponseReportProd struct {
	SectionID     int `json:"section_id"`
	SectionNumber int `json:"section_number"`
	ProductsCount int `json:"products_count"`
}

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
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": sections,
	})
}

func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	var section internal.Section
	section, err = h.sv.FindByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": section,
	})
}

func (h *SectionHandler) ReportProducts(w http.ResponseWriter, r *http.Request) {
	var totalQuantity int
	var err error
	var section internal.Section

	idStr := r.URL.Query().Get("id")

	switch idStr {
	case "":
		totalQuantity, err = h.sv.ReportProducts()
		section = internal.Section{}
	default:
		idSection, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("id should be a number"))
			return
		}

		section, err = h.sv.FindByID(idSection)
		if err != nil {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			return
		}

		totalQuantity, err = h.sv.ReportProductsByID(idSection)
	}

	if err != nil {
		log.Println(err)
		if errors.Is(err, internal.ProductBatchNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			return
		}
		response.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	responseReport := ResponseReportProd{
		SectionID:     0,
		SectionNumber: 0,
		ProductsCount: totalQuantity,
	}

	if section.ID != 0 {
		responseReport.SectionID = section.ID
		responseReport.SectionNumber = section.SectionNumber
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": responseReport,
	})
}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sectionJSON RequestSectionJSON
	if err := json.NewDecoder(r.Body).Decode(&sectionJSON); err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	section := internal.Section{
		SectionNumber:      sectionJSON.SectionNumber,
		CurrentTemperature: sectionJSON.CurrentTemperature,
		MinimumTemperature: sectionJSON.MinimumTemperature,
		CurrentCapacity:    sectionJSON.CurrentCapacity,
		MinimumCapacity:    sectionJSON.MinimumCapacity,
		MaximumCapacity:    sectionJSON.MaximumCapacity,
		WarehouseID:        sectionJSON.WarehouseID,
		ProductTypeID:      sectionJSON.ProductTypeID,
	}

	err := h.sv.Save(&section)
	if err != nil {
		if errors.Is(err, internal.SectionAlreadyExists) || errors.Is(err, internal.SectionNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, rest_err.NewUnprocessableEntityError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": section,
	})
}

func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	section, err := h.sv.Update(id, updates)
	if err != nil {
		if errors.Is(err, internal.SectionNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		} else {
			response.JSON(w, http.StatusConflict, rest_err.NewConflictError(err.Error()))
		}
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": section,
	})
}

func (h *SectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError(err.Error()))
		return
	}

	err = h.sv.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.SectionNotFound):
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
		default:
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(ErrInternalServer))
		}
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
