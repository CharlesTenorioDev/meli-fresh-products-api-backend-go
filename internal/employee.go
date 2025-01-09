package internal

import "errors"

var ErrEmployeeNotFound = errors.New("employee not found")

type Employee struct {
	Id           int    `json:"id"`
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	WarehouseId  int    `json:"warehouse_id"`
}

type EmployeePatch struct {
	CardNumberId *string `json:"card_number_id,omitempty"`
	FirstName    *string `json:"first_name,omitempty"`
	LastName     *string `json:"last_name,omitempty"`
}

func (emp *Employee) RequirementsFields() (ok bool) {
	ok = true

	if emp.CardNumberId == "" || emp.FirstName == "" || emp.LastName == "" || emp.WarehouseId == 0 {
		ok = false
	}
	return
}

// function to update employee data in repository
func (emp EmployeePatch) EmployeePatch(empUpdate *Employee) {

	if emp.CardNumberId != nil {
		empUpdate.CardNumberId = *emp.CardNumberId
	}

	if emp.FirstName != nil {
		empUpdate.FirstName = *emp.FirstName
	}

	if emp.LastName != nil {
		empUpdate.LastName = *emp.LastName
	}
}

type EmployeeRepository interface {
	GetAll() (db map[int]Employee)
	Save(emp *Employee) int
	Update(id int, employee Employee)
	Delete(id int)
}

type EmployeeService interface {
	GetAll() map[int]Employee
	GetById(id int) (emp Employee, err error)
	Save(emp *Employee) (err error)
	Update(employees Employee) (err error)
	Delete(id int) (err error)
}
