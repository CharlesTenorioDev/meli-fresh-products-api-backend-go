package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
)

var (
	EmployeeInUse     = errors.New("employee already in use")
	CardNumberIdInUse = errors.New("card number id already in use")
	EmployeeNotFound  = errors.New("employee not found")
)

func NewEmployeeServiceDefault(rp repository.EmployeeRepository) *EmployeeDefault {
	return &EmployeeDefault{
		rp: rp,
	}

}

type EmployeeDefault struct {
	rp repository.EmployeeRepository
}

type EmployeeService interface {
	GetAll() map[int]internal.Employee
	GetById(id int) (emp internal.Employee, err error)
	Save(emp internal.Employee) (err error)
	Update(id int, employees internal.EmployeePatch) (err error)
	Delete(id int) (err error)
}

func (s *EmployeeDefault) GetAll() map[int]internal.Employee {
	return s.rp.GetAll()
}

func (s *EmployeeDefault) GetById(id int) (emp internal.Employee, err error) {
	employee := s.rp.GetAll()
	emp, ok := employee[id]

	if !ok {
		err = EmployeeNotFound
	}
	return
}

func (s *EmployeeDefault) Save(emp internal.Employee) (err error) {
	employees := s.rp.GetAll()
	_, ok := employees[emp.Id]
	if ok {
		err = errors.New("employee already exists")
		return
	}

	validate := emp.RequirementsFields()
	if !validate {
		err = errors.New("invalid entity data")
		return
	}

	if cardNumberIdInUse(emp.CardNumberId, employees) {
		err = errors.New("card number id already in use, please try again")
		return
	}

	s.rp.Save(emp)
	return

}

func cardNumberIdInUse(cardId string, employees map[int]internal.Employee) bool {

	for _, employee := range employees {
		if employee.CardNumberId == cardId {
			return true
		}
	}
	return false
}

func (s *EmployeeDefault) Update(id int, emp internal.EmployeePatch) (err error) {

	data := s.rp.GetAll()
	_, ok := data[id] // search employee by id
	if !ok {
		err = EmployeeNotFound
		return
	}

	// check if card number id is already in use
	if cardNumberIdInUse(*emp.CardNumberId, data) {
		err = CardNumberIdInUse
		return
	}

	s.rp.Update(id, emp)
	return
}

func (s *EmployeeDefault) Delete(id int) (err error) {
	employee := s.rp.GetAll() // search for employee by id
	_, ok := employee[id]
	// if employee not found
	if !ok {
		err = EmployeeNotFound
		return
	}

	s.rp.Delete(id)
	return
}
