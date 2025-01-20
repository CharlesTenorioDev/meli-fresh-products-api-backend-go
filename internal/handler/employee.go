package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bootcamp-go/web/response"
	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/resterr"
)

type EmployeeHandlerDefault struct {
	sv internal.EmployeeService
}

func NewEmployeeDefault(sv internal.EmployeeService) *EmployeeHandlerDefault {
	return &EmployeeHandlerDefault{
		sv: sv,
	}
}

func (h *EmployeeHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	dataEmployee, err := h.sv.GetAll()

	if err != nil {
		response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response.JSON(w, http.StatusOK, map[string]any{
		"data": dataEmployee,
	})

}

func (h *EmployeeHandlerDefault) GetByID(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid id format", //status 400
		})
		return
	}

	emp, err := h.sv.GetByID(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]any{
			"data": "employee not found", //status 404
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": emp,
	})

}

func (h *EmployeeHandlerDefault) Create(w http.ResponseWriter, r *http.Request) {

	var employee internal.Employee

	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid body format", //status 400
		})
		return
	}

	err = h.sv.Save(&employee) // save employee in service

	// checks if card number Id field is already in use, because it's a unique field
	if err != nil {

		if errors.Is(err, service.ErrCardNumberIDInUse) {
			response.JSON(w, http.StatusConflict, map[string]any{
				"error": "card number id already in use", // Status 409
			})
		} else if errors.Is(err, service.ErrEmployeeInUse) {

			response.JSON(w, http.StatusConflict, map[string]any{
				"error": "employee already in use", // Status 409
			})
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
				"data": err.Error(), // status 422
			})
		}
		return
	}

	response.JSON(w, http.StatusCreated, map[string]any{
		"data": employee,
	})

}

func (h *EmployeeHandlerDefault) Update(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil || id <= 0 {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid id format", //status 400
		})
		return
	}

	var employee internal.Employee
	err = json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid body format",
		})
		return
	}

	employee.ID = id
	err = h.sv.Update(employee)

	if err != nil {
		if errors.Is(err, service.ErrEmployeeNotFound) {
			response.JSON(w, http.StatusNotFound, map[string]any{
				"data": err.Error(),
			})

		} else {
			response.JSON(w, http.StatusConflict, map[string]any{
				"data": err.Error(),
			})
		}
		return
	}

	updatedEmployee, err := h.sv.GetByID(employee.ID)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]any{
			"data": "error retrieving updated employee",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]any{
		"data": updatedEmployee,
	})
}

func (h *EmployeeHandlerDefault) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid id format",
		})
		return
	}
	err = h.sv.Delete(id)
	if err != nil {
		response.JSON(w, http.StatusNotFound, map[string]any{
			"data": "employee not found",
		})
		return
	}
	response.JSON(w, http.StatusNoContent, nil)
}

func (h *EmployeeHandlerDefault) ReportInboundOrders(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	idStr = strings.TrimSpace(idStr)

	switch {

	case idStr == "":
		inboundOrders, err := h.sv.CountInboundOrdersPerEmployee()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, resterr.NewInternalServerError("failed to fetch inbound orders"))
			return
		}
		response.JSON(w, http.StatusOK, map[string]any{
			"data": inboundOrders,
		})
		return

	default:
		id, err := strconv.Atoi(idStr)
		switch {
		case err != nil:
			response.JSON(w, http.StatusBadRequest, resterr.NewBadRequestError("id should be a number"))
			return
		}

		countInboundOrders, err := h.sv.ReportInboundOrdersByID(id)
		switch {
		case err != nil:
			response.JSON(w, http.StatusNotFound, resterr.NewNotFoundError("employee not found"))
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": countInboundOrders,
		})
	}
}
