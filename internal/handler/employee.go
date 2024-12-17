package handler

import (
	"net/http"

	"github.com/bootcamp-go/web/response"
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
