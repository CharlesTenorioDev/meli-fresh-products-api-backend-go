package repository

import (
	"encoding/json"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type EmployeeRepositoryDefault struct {
	db map[int]*internal.Employee
}

func NewEmployeeRepository(dbPath string) (mp *EmployeeRepositoryDefault, err error) {
	var employees []internal.Employee
	db := make(map[int]*internal.Employee)

	file, err := os.Open(dbPath) // open file employees
	if err != nil {
		return
	}

	err = json.NewDecoder(file).Decode(&employees)
	if err != nil {
		return
	}

	for _, employee := range employees {
		db[employee.ID] = &employee
	}

	mp = &EmployeeRepositoryDefault{
		db: db,
	}
	return
}

func (r *EmployeeRepositoryDefault) GetAll() (db map[int]internal.Employee) {
	db = make(map[int]internal.Employee)

	for key, value := range r.db { // get all employees in db
		db[key] = *value
	}

	return
}

func (r *EmployeeRepositoryDefault) Save(emp *internal.Employee) int {
	if emp.ID == 0 {
		emp.ID = len(r.db) + 1 //increment id
	}

	r.db[emp.ID] = emp // add new employee in db

	return emp.ID
}

func (r *EmployeeRepositoryDefault) Update(id int, employee internal.Employee) {
	if emp, ok := r.db[id]; ok {
		emp.CardNumberID = employee.CardNumberID
		emp.FirstName = employee.FirstName
		emp.LastName = employee.LastName
		emp.WarehouseID = employee.WarehouseID
	}
}

func (r *EmployeeRepositoryDefault) Delete(id int) {
	delete(r.db, id)
}
