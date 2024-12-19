package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
)

var (
	EmployeeInUse       = errors.New("employee already in use")
	CardNumberIdInUse   = errors.New("card number id already in use")
	EmployeeNotFound    = errors.New("employee not found")
	UnprocessableEntity = errors.New("couldn't parse employee")
	ConflictInEmployee  = errors.New("conflict in employee")
)

func NewEmployeeServiceDefault(rp repository.EmployeeRepository, rpWarehouse internal.WarehouseRepository) *EmployeeDefault {
	return &EmployeeDefault{
		rp:  rp,
		rpW: rpWarehouse,
	}

}

type EmployeeDefault struct {
	rp  repository.EmployeeRepository
	rpW internal.WarehouseRepository
}

type EmployeeService interface {
	GetAll() map[int]internal.Employee
	GetById(id int) (emp internal.Employee, err error)
	Save(emp *internal.Employee) (err error)
	Update(employees internal.Employee) (err error)
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

func (s *EmployeeDefault) Save(emp *internal.Employee) (err error) {
	employees := s.rp.GetAll()

	if _, exists := employees[emp.Id]; exists && emp.Id != 0 {
		err = errors.New("employee already exists")
		return
	}

	validate := emp.RequirementsFields()
	if !validate {
		err = errors.New("invalid entity data")
		return
	}

	_, err = s.rpW.FindByID(emp.WarehouseId)
	if err != nil {
		return UnprocessableEntity
	}

	if cardNumberIdInUse(emp.CardNumberId, employees) {
		err = CardNumberIdInUse
		return
	}

	savedId := s.rp.Save(emp)
	emp.Id = savedId

	return nil

}

func cardNumberIdInUse(cardId string, employees map[int]internal.Employee) bool {

	for _, employee := range employees {
		if employee.CardNumberId == cardId {
			return true
		}
	}
	return false
}

func (s *EmployeeDefault) Update(emp internal.Employee) (err error) {

	data := s.rp.GetAll()
	existingEmployee, ok := data[emp.Id]
	if !ok {
		err = EmployeeNotFound
		return
	}

	// check if card number id is already in use
	if cardNumberIdInUse(emp.CardNumberId, data) && existingEmployee.CardNumberId != emp.CardNumberId {
		err = CardNumberIdInUse
		return
	}

	validate := emp.RequirementsFields()
	if !validate {
		err = errors.New("invalid entity data")
		return
	}

	_, err = s.rpW.FindByID(emp.WarehouseId)
	if err != nil {
		return ConflictInEmployee
	}

	s.rp.Update(emp.Id, emp)
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
