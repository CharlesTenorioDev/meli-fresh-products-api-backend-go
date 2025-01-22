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
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
)

type EmployeeHandlerDefault struct {
	sv internal.EmployeeService
}

func NewEmployeeDefault(sv internal.EmployeeService) *EmployeeHandlerDefault {
	return &EmployeeHandlerDefault{
		sv: sv,
	}
}

// GetAll godoc
// @Summary Get all employees
// @Description Retrieve a list of all employees from the database
// @Tags Employees
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "List of all employees"
// @Failure 500 {object} resterr.RestErr "internal server error"
// @Router /api/v1/employees [get]
func (h *EmployeeHandlerDefault) GetAll(w http.ResponseWriter, r *http.Request) {
	dataEmployee, err := h.sv.GetAll()

	if err != nil {
		response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response.JSON(w, http.StatusOK, map[string]any{
		"data": dataEmployee,
	})

}

// GetByID godoc
// @Summary Get an employee by Id
// @Description Retrieve a specific employee by Id
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 200 {object} map[string]interface{} "Employee data"
// @Failure 400 {object} map[string]interface{} "Invalid Id format"
// @Failure 404 {object} map[string]interface{} "Employee not found"
// @Router /api/v1/employees/{id} [get]
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

// Create godoc
// @Summary Create a new employee
// @Description Create a new employee record in the database
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body internal.Employee true "Employee data"
// @Success 201 {object} map[string]interface{} "Created employee"
// @Failure 400 {object} map[string]interface{} "Invalid body format"
// @Failure 409 {object} map[string]interface{} "Card number id already in use" or "Employee already in use"
// @Failure 422 {object} map[string]interface{} "Invalid entity data"
// @Router /api/v1/employees [post]
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

		if errors.Is(err, service.ErrCardNumberIdInUse) {
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

// Update godoc
// @Summary Update an existing employee
// @Description Update the details of an existing employee by Id
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body internal.Employee true "Employee data"
// @Success 200 {object} map[string]interface{} "Updated employee"
// @Failure 400 {object} map[string]interface{} "Invalid Id format" or "Invalid body format"
// @Failure 404 {object} map[string]interface{} "Employee not found"
// @Failure 409 {object} map[string]interface{} "Card number id already in use" or "Conflict in employee"
// @Router /api/v1/employees/{id} [patch]
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

// Delete godoc
// @Summary Delete an employee by Id
// @Description Remove an employee record from the database by Id
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Success 204 {object} nil "No Content"
// @Failure 400 {object} map[string]interface{} "Invalid Id format"
// @Failure 404 {object} map[string]interface{} "Employee not found"
// @Router /api/v1/employees/{id} [delete]
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

// ReportInboundOrders godoc
// @Summary Get inbound orders for employees
// @Description Retrieve the count all inbound orders per employee, or for a specific employee if Id is provided
// @Tags Employees
// @Accept json
// @Produce json
// @Param id query int false "Employee ID"
// @Success 200 {object} map[string]interface{} "Count of inbound orders per employee"
// @Failure 400 {object} resterr.RestErr "Id should be a number"
// @Failure 404 {object} resterr.RestErr "Employee not found"
// @Failure 500 {object} resterr.RestErr "Failed to fetch inbound orders"
// @Router /api/v1/employees/report-inbound-orders [get]
func (h *EmployeeHandlerDefault) ReportInboundOrders(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	idStr = strings.TrimSpace(idStr)

	switch {

	case idStr == "":
		inboundOrders, err := h.sv.CountInboundOrdersPerEmployee()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, rest_err.NewInternalServerError("failed to fetch inbound orders"))
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
			response.JSON(w, http.StatusBadRequest, rest_err.NewBadRequestError("id should be a number"))
			return
		}

		countInboundOrders, err := h.sv.ReportInboundOrdersById(id)
		switch {
		case err != nil:
			response.JSON(w, http.StatusNotFound, rest_err.NewNotFoundError("employee not found"))
			return
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"data": countInboundOrders,
		})
	}
}
