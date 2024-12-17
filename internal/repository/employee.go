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
	for index, employee := range employees {
		db[index] = &employee
	}
	return &EmployeeRepositoryDefault{ // return repository with db employees updated
		db: db,
	}
}

type EmployeeRepository interface {
	GetAll() (db map[int]internal.Employee)
	GetById(id int) (db internal.Employee)
	Save(id int, employee internal.Employee)
	Update(id int, employee internal.EmployeePatch)
	Delete(id int)
}

func (r *EmployeeRepositoryDefault) GetAll() (db map[int]internal.Employee) {
	db = make(map[int]internal.Employee)

	for key, value := range r.db { // save employees in db
		db[key] = *value
	}
	return
}
