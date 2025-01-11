package service

import (
	"errors"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

var (
	EmployeeInUse       = errors.New("employee already in use")
	CardNumberIdInUse   = errors.New("card number id already in use")
	EmployeeNotFound    = errors.New("employee not found")
	UnprocessableEntity = errors.New("couldn't parse employee")
	ConflictInEmployee  = errors.New("conflict in employee")
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

	if emp.Id != 0 {
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
		return WarehouseNotFound
	}

	if cardNumberIdInUse(emp.CardNumberId, employees) {
		err = CardNumberIdInUse
		return
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
		err = EmployeeNotFound
		return
	}

	if cardNumberIdInUse(emp.CardNumberId, data) && existingEmployee.CardNumberId != emp.CardNumberId {
		err = CardNumberIdInUse
		return
	}

	if !emp.RequirementsFields() {
		err = errors.New("required fields are missing")
		return
	}

	_, err = s.rpW.FindByID(emp.WarehouseId)
	if err != nil {
		return ConflictInEmployee
	}

	err = s.rp.Update(emp.Id, emp)
	if err != nil {
		return err
	}

	return
}

func (s *EmployeeDefault) Delete(id int) (err error) {
	return s.rp.Delete(id)
}

func (s *EmployeeDefault) CountInboundOrdersPerEmployee() (io []internal.InboundOrdersPerEmployee, err error) {
	return s.rp.CountInboundOrdersPerEmployee()
}

func (s *EmployeeDefault) ReportInboundOrdersById(employeeId int) (io internal.InboundOrdersPerEmployee, err error) {
	return s.rp.ReportInboundOrdersById(employeeId)
}
