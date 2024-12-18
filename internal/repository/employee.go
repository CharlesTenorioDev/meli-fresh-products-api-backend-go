package repository

import (
	"encoding/json"
	"log"
	"os"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type EmployeeRepositoryDefault struct {
	db map[int]*internal.Employee
}

func NewEmployeeRepository() *EmployeeRepositoryDefault {

	var employees []internal.Employee
	db := make(map[int]*internal.Employee)
	file, err := os.Open("db/employees.json") // open file employees

	if err != nil {
		log.Fatal(err)

	}

	// decode json and memory store in employees
	err = json.NewDecoder(file).Decode(&employees)

	if err != nil {
		log.Fatal(err)
	}

	// save employees in db
	for _, employee := range employees {
		if employee.Id > 0 {
			db[employee.Id] = &employee
		} else {
			log.Fatal(employee)
		}
	}
	return &EmployeeRepositoryDefault{ // return repository with db employees updated
		db: db,
	}
}

type EmployeeRepository interface {
	GetAll() (db map[int]internal.Employee)
	Save(emp internal.Employee)
	Update(id int, employee internal.EmployeePatch)
	Delete(id int)
}

func (r *EmployeeRepositoryDefault) GetAll() (db map[int]internal.Employee) {
	db = make(map[int]internal.Employee)

	for key, value := range r.db { // get all employees in db
		db[key] = *value
	}
	return
}

func (r *EmployeeRepositoryDefault) Save(emp internal.Employee) {

	if emp.Id == 0 {
		emp.Id = len(r.db) + 1 //increment id
	}

	r.db[emp.Id] = &emp // add new employee in db
}

func (r *EmployeeRepositoryDefault) Update(id int, employee internal.EmployeePatch) {
	empPatch := r.db[id]
	employee.EmployeePatch(empPatch)
}

func (r *EmployeeRepositoryDefault) Delete(id int) {
	delete(r.db, id)
}
