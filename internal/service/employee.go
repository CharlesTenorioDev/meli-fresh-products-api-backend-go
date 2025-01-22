package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	ErrEmployeeInUse       = errors.New("employee already in use")
	ErrCardNumberIDInUse   = errors.New("card number id already in use")
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

func (s *EmployeeDefault) GetByID(id int) (emp internal.Employee, err error) {
	return s.rp.GetByID(id)
}

func (s *EmployeeDefault) Save(emp *internal.Employee) (err error) {
	employees, err := s.rp.GetAll()
	if err != nil {
		return internal.ErrEmployeeNotFound
	}

	if emp.ID != 0 {
		return errors.New("employee already exists")
	}

	if cardNumberIDInUse(emp.CardNumberID, employees) {
		err = ErrCardNumberIDInUse
		return err
	}

	validate := emp.RequirementsFields()
	if !validate {
		err = errors.New("invalid entity data")
		return err
	}

	_, err = s.rpW.FindByID(emp.WarehouseID)
	if err != nil {
		return internal.ErrWarehouseRepositoryNotFound
	}

	_, err = s.rp.Save(emp)
	if err != nil {
		return err
	}

	return nil
}

func cardNumberIDInUse(cardID string, employees []internal.Employee) bool {
	for _, employee := range employees {
		if employee.CardNumberID == cardID {
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
		if employee.ID == emp.ID {
			existingEmployee = &employee
			break
		}
	}

	if existingEmployee == nil {
		return ErrEmployeeNotFound
	}

	if cardNumberIDInUse(emp.CardNumberID, data) && existingEmployee.CardNumberID != emp.CardNumberID {
		return ErrCardNumberIDInUse
	}

	if !emp.RequirementsFields() {
		return errors.New("required fields are missing")
	}

	_, err = s.rpW.FindByID(emp.WarehouseID)
	if err != nil {
		return ErrConflictInEmployee
	}

	err = s.rp.Update(emp.ID, emp)
	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeDefault) Delete(id int) (err error) {
	_, err = s.rp.GetByID(id)
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

func (s *EmployeeDefault) ReportInboundOrdersByID(employeeID int) (io internal.InboundOrdersPerEmployee, err error) {
	return s.rp.ReportInboundOrdersByID(employeeID)
}
