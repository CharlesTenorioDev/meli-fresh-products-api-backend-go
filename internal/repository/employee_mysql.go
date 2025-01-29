package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

type EmployeeMysql struct {
	db *sql.DB
}

func NewEmployeeMysql(db *sql.DB) *EmployeeMysql {
	return &EmployeeMysql{db}
}

const (
	InboundOrdersPerEmployeeQuery = `
	SELECT COUNT(i.employee_id) AS inbound_orders_count, e.id, e.card_number_id, e.first_name, e.last_name, i.warehouse_id
	FROM inbound_orders i
	INNER JOIN employees e ON i.employee_id = e.id
	GROUP BY i.employee_id, i.id;`

	InboundOrdersPerEmployeeByIDQuery = `
	SELECT COUNT(i.employee_id) AS inbound_orders_count, i.id, e.card_number_id, e.first_name, e.last_name, i.warehouse_id
	FROM inbound_orders i
	INNER JOIN employees e ON i.employee_id = e.id
	WHERE i.employee_id = ?
	GROUP BY i.employee_id, i.id;`
)

func (r *EmployeeMysql) GetAll() (db []internal.Employee, err error) {
	rows, err := r.db.Query("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var emp internal.Employee
		rows.Scan(&emp.ID, &emp.CardNumberID, &emp.FirstName, &emp.LastName, &emp.WarehouseID)

		db = append(db, emp)
	}

	err = rows.Err()
	return
}

func (r *EmployeeMysql) GetByID(id int) (emp internal.Employee, err error) {
	row := r.db.QueryRow("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?", id)
	err = row.Scan(&emp.ID, &emp.CardNumberID, &emp.FirstName, &emp.LastName, &emp.WarehouseID)
	return
}

func (r *EmployeeMysql) Save(emp *internal.Employee) (id int64, err error) {
	var existingEmployee internal.Employee

	err = r.db.QueryRow(
		"SELECT id FROM employees WHERE card_number_id = ?",
		emp.CardNumberID).Scan(&existingEmployee.ID)
	if err == nil {
		err = fmt.Errorf("employee with card_number_id %s already exists", emp.CardNumberID)
		return
	} else if err != sql.ErrNoRows {
		return
	}

	result, err := r.db.Exec(
		"INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?)",
		emp.CardNumberID, emp.FirstName, emp.LastName, emp.WarehouseID)
	if err != nil {
		return
	}

	id, err = result.LastInsertId()

	return
}

func (r *EmployeeMysql) Update(id int, employee internal.Employee) (err error) {
	_, err = r.db.Exec(
		"UPDATE employees SET card_number_id = ?, first_name = ?, last_name = ?, warehouse_id = ? WHERE id = ?",
		employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID, id,
	)

	return
}

func (r *EmployeeMysql) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id = ?", id)
	return err
}

func (r *EmployeeMysql) CountInboundOrdersPerEmployee() (io []internal.InboundOrdersPerEmployee, err error) {
	row, err := r.db.Query(InboundOrdersPerEmployeeQuery)
	if err != nil {
		return
	}

	for row.Next() {
		var countInboundPerEmployee internal.InboundOrdersPerEmployee

		row.Scan(&countInboundPerEmployee.CountInOrders, &countInboundPerEmployee.ID, &countInboundPerEmployee.CardNumberID, &countInboundPerEmployee.FirstName, &countInboundPerEmployee.LastName, &countInboundPerEmployee.WarehouseID)

		io = append(io, countInboundPerEmployee)
	}

	return
}

func (r *EmployeeMysql) ReportInboundOrdersByID(employeeID int) (io internal.InboundOrdersPerEmployee, err error) {
	row := r.db.QueryRow(
		InboundOrdersPerEmployeeByIDQuery,
		employeeID,
	)

	err = row.Scan(&io.CountInOrders, &io.ID, &io.CardNumberID, &io.FirstName, &io.LastName, &io.WarehouseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrEmployeeNotFound
		}
	}

	return
}
