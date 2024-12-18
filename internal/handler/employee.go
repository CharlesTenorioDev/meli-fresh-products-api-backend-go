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
		response.JSON(w, http.StatusBadRequest, map[string]any{
			"data": "invalid id format", //status 400
		})
		return
	}

	emp, err := h.sv.GetById(id)
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

func (h *EmployeeHandlerDefault) Save(w http.ResponseWriter, r *http.Request) {

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
		if errors.Is(err, service.EmployeeInUse) || errors.Is(err, service.CardNumberIdInUse) {
			response.JSON(w, http.StatusConflict, map[string]any{
				"data": err.Error(),
			})
		} else {
			response.JSON(w, http.StatusUnprocessableEntity, map[string]any{
				"data": err.Error(), //status 422
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

	employee.Id = id
	err = h.sv.Update(employee)

	if err != nil {
		if errors.Is(err, service.EmployeeNotFound) {
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

	updatedEmployee, err := h.sv.GetById(employee.Id)
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
			"data": err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusNoContent, nil) //status 204
}
