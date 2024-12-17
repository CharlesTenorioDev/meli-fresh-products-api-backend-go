package service

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
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
	FindByID(id int) (b internal.Employee, err error)
	Save(id int, employee internal.Employee) (err error)
	Update(id int, employees internal.EmployeePatch) (err error)
	Delete(id int) (err error)
}

func (s *EmployeeDefault) GetAll() map[int]internal.Employee {
	return s.rp.GetAll()
}
