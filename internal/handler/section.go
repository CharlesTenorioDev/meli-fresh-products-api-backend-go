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

type SectionJSON struct {
	SectionNumber      int     `json:"section_number"`
	CurrentTemperature float64 `json:"current_temperature"`
	MinimumTemperature float64 `json:"minimum_temperature"`
	CurrentCapacity    int     `json:"current_capacity"`
	MinimumCapacity    int     `json:"minimum_capacity"`
	MaximumCapacity    int     `json:"maximum_capacity"`
	WarehouseID        int     `json:"warehouse_id"`
	ProductTypeID      int     `json:"product_type_id"`
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
	/*var prodBatchs []internal.ProductBatch
	var err error

	idStr := r.URL.Query().Get("id")

	switch idStr {
	case "":
		prodBatchs, err = h.sv.ReportProducts()
	default:
		id, err := strconv.Atoi(idStr)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("id should be a number"))
			return
		}

		prodBatchs, err = h.sv.ReportProductsByID(id)
	}

	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ProductBatchNotFound) {
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError(err.Error()))
			return
		}
		response.JSON(w, http.StatusInternalServerError, nil)
		return
	}

	var localitiesJson []LocalityGetJson
	for _, locality := range prodBatchs {
		localitiesJson = append(localitiesJson, LocalityGetJson{
			ID:           locality.ID,
			LocalityName: locality.LocalityName,
			ProvinceName: locality.ProvinceName,
			CountryName:  locality.CountryName,
			SellersCount: locality.Sellers,
		})
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": localitiesJson,
	})*/
}

func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sectionJSON SectionJSON
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
		if errors.Is(err, service.SectionAlreadyExists) || errors.Is(err, service.SectionNumberAlreadyInUse) {
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
		if errors.Is(err, service.SectionNotFound) {
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
		response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
