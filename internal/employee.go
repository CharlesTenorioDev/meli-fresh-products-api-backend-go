package internal

import "errors"

var ErrEmployeeNotFound = errors.New("employee not found")
var ErrEmployeeConflict = errors.New("employee already in use")

type Employee struct {
	ID           int    `json:"id"`
	CardNumberID string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseID  int    `json:"warehouse_id"`
}

type EmployeePatch struct {
	CardNumberID *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
}

type InboundOrdersPerEmployee struct {
	ID            int    `json:"id"`
	CardNumberID  string `json:"card_number_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	WarehouseID   int    `json:"warehouse_id"`
	CountInOrders int    `json:"inbound_orders_count"`
}

func (emp *Employee) RequirementsFields() (ok bool) {
	ok = true
	if emp.CardNumberID == "" || emp.FirstName == "" || emp.LastName == "" || emp.WarehouseID == 0 {
		ok = false
	}

	return
}

// EmployeePatch function to update employee data in repository
func (emp EmployeePatch) EmployeePatch(empUpdate *Employee) {
	if emp.CardNumberID != nil {
		empUpdate.CardNumberID = *emp.CardNumberID
	}

	if emp.FirstName != nil {
		empUpdate.FirstName = *emp.FirstName
	}

	if emp.LastName != nil {
		empUpdate.LastName = *emp.LastName
	}
}

type EmployeeRepository interface {
	GetAll() (db []Employee, err error)
	GetByID(id int) (emp Employee, err error)
	Save(emp *Employee) (int, error)
	Update(id int, employee Employee) error
	Delete(id int) error
	CountInboundOrdersPerEmployee() (io []InboundOrdersPerEmployee, err error)
	ReportInboundOrdersByID(employeeID int) (io InboundOrdersPerEmployee, err error)
}

type EmployeeService interface {
	GetAll() (db []Employee, err error)
	GetByID(id int) (emp Employee, err error)
	Save(emp *Employee) (err error)
	Update(employees Employee) (err error)
	Delete(id int) (err error)
	CountInboundOrdersPerEmployee() (io []InboundOrdersPerEmployee, err error)
	ReportInboundOrdersByID(employeeID int) (io InboundOrdersPerEmployee, err error)
}
