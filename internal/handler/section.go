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
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
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

type SectionsUpdateJSON struct {
	SectionNumber      *int     `json:"section_number"`
	CurrentTemperature *float64 `json:"current_temperature"`
	MinimumTemperature *float64 `json:"minimum_temperature"`
	CurrentCapacity    *int     `json:"current_capacity"`
	MinimumCapacity    *int     `json:"minimum_capacity"`
	MaximumCapacity    *int     `json:"maximum_capacity"`
	WarehouseID        *int     `json:"warehouse_id"`
	ProductTypeID      *int     `json:"product_type_id"`
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

// GetAll retrieves all sections
// @Summary Retrieve all sections
// @Description Fetches all sections available in the database
// @Tags Section
// @Produce json
// @Success 200 {object} []internal.Section "List of sections"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Router /api/v1/sections [get]
func (h *SectionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	sections, err := h.sv.FindAll()
	if err != nil {
		if errors.Is(err, internal.ErrSectionNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError("section not found"))

			return
		}

		response.JSON(w, http.StatusInternalServerError, nil)

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": sections,
	})
}

// GetByID retrieves a section by ID
// @Summary Retrieve a section by ID
// @Description Fetches the section based on the provided section ID
// @Tags Section
// @Produce json
// @Param id path int true "Section ID"
// @Success 200 {object} internal.Section "Section data"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Section not found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sections/{id} [get]
func (h *SectionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))

		return
	}

	var section internal.Section

	section, err = h.sv.FindByID(id)
	if err != nil {
		if errors.Is(err, internal.ErrSectionNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))

			return
		}

		response.JSON(w, http.StatusInternalServerError, nil)

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": section,
	})
}

// ReportProducts retrieves a report of products for all sections or a specific section
// @Summary Retrieve a report of products in a section
// @Description Fetches a report of products available in a section or across all sections
// @Tags Section
// @Produce json
// @Param id query int false "Section ID"
// @Success 200 {object} []ResponseReportProd "Report of products in sections"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Section not found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sections/report-products [get]
func (h *SectionHandler) ReportProducts(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		sections, err := h.sv.ReportProducts()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": sections,
		})

		return
	}

	idSection, err := strconv.Atoi(idStr)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	report, err := h.sv.ReportProductsByID(idSection)
	if err != nil {
		log.Println(err)

		if errors.Is(err, internal.ErrSectionNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		} else {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
		}

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": report,
	})
}

// Create creates a new section
// @Summary Create a new section
// @Description Creates a new section with the provided details on the request body
// @Tags Section
// @Accept json
// @Produce json
// @Param section body RequestSectionJSON true "Section Create Request"
// @Success 201 {object} internal.Section "Created Section"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 409 {object} rest_err.RestErr "Section with given section number already registered" or "Warehouse not found" or "Product-type not found"
// @Failure 422 {object} rest_err.RestErr "Couldn't parse section"
// @Router /api/v1/sections [post]
func (h *SectionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var sectionJSON map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&sectionJSON); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	requiredFields := []string{
		"section_number", "current_temperature", "minimum_temperature",
		"current_capacity", "minimum_capacity", "maximum_capacity",
		"warehouse_id", "product_type_id",
	}

	for _, field := range requiredFields {
		if sectionJSON[field] == nil {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(field+" is required"))
			return
		}
	}

	section := internal.Section{
		SectionNumber:      int(sectionJSON["section_number"].(float64)),
		CurrentTemperature: sectionJSON["current_temperature"].(float64),
		MinimumTemperature: sectionJSON["minimum_temperature"].(float64),
		CurrentCapacity:    int(sectionJSON["current_capacity"].(float64)),
		MinimumCapacity:    int(sectionJSON["minimum_capacity"].(float64)),
		MaximumCapacity:    int(sectionJSON["maximum_capacity"].(float64)),
		WarehouseID:        int(sectionJSON["warehouse_id"].(float64)),
		ProductTypeID:      int(sectionJSON["product_type_id"].(float64)),
	}

	err := h.sv.Save(&section)
	if err != nil {
		if errors.Is(err, internal.ErrSectionAlreadyExists) || errors.Is(err, internal.ErrSectionNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		}

		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": section,
	})
}

// Update updates an existing section
// @Summary Update an existing section
// @Description Updates a section with the provided Id and data on the request body
// @Tags Section
// @Accept json
// @Produce json
// @Param id path int true "Section ID"
// @Param updates body map[string]interface{} true "Updated section data"
// @Success 200 {object} internal.Section "Updated Section"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Section not found"
// @Failure 409 {object} rest_err.RestErr "Section with given section number already registered"
// @Router /api/v1/sections/{id} [patch]
func (h *SectionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	var body SectionsUpdateJSON
	if err = json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	stPatch := internal.SectionPatch{
		SectionNumber:      body.SectionNumber,
		CurrentTemperature: body.CurrentTemperature,
		MinimumTemperature: body.MinimumTemperature,
		CurrentCapacity:    body.CurrentCapacity,
		MinimumCapacity:    body.MinimumCapacity,
		MaximumCapacity:    body.MaximumCapacity,
		WarehouseID:        body.WarehouseID,
		ProductTypeID:      body.ProductTypeID,
	}

	section, err := h.sv.Update(id, stPatch)
	if err != nil {
		if errors.Is(err, internal.ErrSectionNotFound) {
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		} else if errors.Is(err, internal.ErrSectionUnprocessableEntity) {
			response.JSON(w, http.StatusUnprocessableEntity, resterr.NewUnprocessableEntityError(err.Error()))
		} else if errors.Is(err, internal.ErrSectionNumberAlreadyInUse) {
			response.JSON(w, http.StatusConflict, resterr.NewConflictError(err.Error()))
		}

		response.JSON(w, http.StatusInternalServerError, nil)

		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": section,
	})
}

// Delete deletes a section
// @Summary Delete a section
// @Description Deletes a section identified by its Id
// @Tags Section
// @Produce json
// @Param id path int true "Section ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} rest_err.RestErr "Bad Request"
// @Failure 404 {object} rest_err.RestErr "Section not found"
// @Failure 500 {object} rest_err.RestErr "Internal Server Error"
// @Router /api/v1/sections/{id} [delete]
func (h *SectionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError(err.Error()))
		return
	}

	err = h.sv.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrSectionNotFound):
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError(err.Error()))
		default:
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(ErrInternalServer))
		}

		return
	}

	response.JSON(w, http.StatusNoContent, nil)
}
