package handler

import (
	"net/http"
	"strconv"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
)

type EmployeeHandlerDefault struct {
	sv service.EmployeeService
}

func NewEmployeeDefault(sv service.EmployeeService) *EmployeeHandlerDefault {
	return &EmployeeHandlerDefault{
		sv: sv,
	}
}

func (h *EmployeeHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	dataEmployee := h.sv.GetAll()

	response.JSON(w, http.StatusOK, map[string]any{
		"data": dataEmployee,
	})
}

func (h *EmployeeHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{ //status 400
			"data": "invalid id format",
		})
		return
	}

	emp, err := h.sv.GetById(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]any{ //status 404
			"data": "employee not found",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": emp,
	})

}
