package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ErrEmployeeInUse       = errors.New("employee already in use")
	ErrCardNumberIdInUse   = errors.New("card number id already in use")
	ErrEmployeeNotFound    = errors.New("employee not found")
	ErrUnprocessableEntity = errors.New("couldn't parse employee")
	ErrConflictInEmployee  = errors.New("conflict in employee")
)

func NewEmployeeServiceDefault(rp internal.EmployeeRepository, rpWarehouse internal.WarehouseRepository) *EmployeeDefault {
	return &EmployeeDefault{
		rp:  rp,
		rpW: rpWarehouse,
	}

}

type EmployeeDefault struct {
	rp  internal.EmployeeRepository
	rpW internal.WarehouseRepository
}

func (s *EmployeeDefault) GetAll() (emp []internal.Employee, err error) {
	emp, err = s.rp.GetAll()
	if err != nil {
		return nil, err
	}
	return emp, nil
}

func (s *EmployeeDefault) GetById(id int) (emp internal.Employee, err error) {
	return s.rp.GetById(id)
}

func (s *EmployeeDefault) Save(emp *internal.Employee) (err error) {
	employees, err := s.rp.GetAll()
	if err != nil {
		return internal.ErrEmployeeNotFound
	}

	if emp.Id != 0 {
		err = errors.New("employee already exists")
		return
	}

	if cardNumberIdInUse(emp.CardNumberId, employees) {
		err = ErrCardNumberIdInUse
		return err
	}

	validate := emp.RequirementsFields()
	if !validate {
		err = errors.New("invalid entity data")
		return err
	}

	_, err = s.rpW.FindByID(emp.WarehouseId)
	if err != nil {
		return internal.ErrWarehouseRepositoryNotFound
	}

	_, err = s.rp.Save(emp)
	if err != nil {
		return err
	}

	return nil

}

func cardNumberIdInUse(cardId string, employees []internal.Employee) bool {

	for _, employee := range employees {
		if employee.CardNumberId == cardId {
			return true
		}
	}
	return false
}

func (s *EmployeeDefault) Update(emp internal.Employee) (err error) {

	data, err := s.rp.GetAll()
	if err != nil {
		return err
	}

	var existingEmployee *internal.Employee
	for _, employee := range data {
		if employee.Id == emp.Id {
			existingEmployee = &employee
			break
		}
	}

	if existingEmployee == nil {
		err = ErrEmployeeNotFound
		return
	}

	if cardNumberIdInUse(emp.CardNumberId, data) && existingEmployee.CardNumberId != emp.CardNumberId {
		err = ErrCardNumberIdInUse
		return
	}

	if !emp.RequirementsFields() {
		err = errors.New("required fields are missing")
		return
	}

	_, err = s.rpW.FindByID(emp.WarehouseId)
	if err != nil {
		return ErrConflictInEmployee
	}

	err = s.rp.Update(emp.Id, emp)
	if err != nil {
		return err
	}

	return
}

func (s *EmployeeDefault) Delete(id int) (err error) {
	_, err = s.rp.GetById(id)
	if err != nil {
		if err == internal.ErrEmployeeNotFound {
			return ErrEmployeeNotFound
		}
		return err
	}
	return s.rp.Delete(id)
}

func (s *EmployeeDefault) CountInboundOrdersPerEmployee() (io []internal.InboundOrdersPerEmployee, err error) {
	return s.rp.CountInboundOrdersPerEmployee()
}

func (s *EmployeeDefault) ReportInboundOrdersById(employeeId int) (io internal.InboundOrdersPerEmployee, err error) {
	return s.rp.ReportInboundOrdersById(employeeId)
}
