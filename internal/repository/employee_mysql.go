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
	SELECT COUNT(i.employee_id) AS inbound_orders_count, i.id, e.first_name, e.last_name, i.warehouse_id
	FROM inbound_orders i
	INNER JOIN employees e ON i.employee_id = e.id
	GROUP BY i.employee_id, i.id;`

	InboundOrdersPerEmployeeByIdQuery = `
	SELECT COUNT(i.employee_id) AS inbound_orders_count, i.id, e.first_name, e.last_name, i.warehouse_id
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
		if err := rows.Scan(&emp.Id, &emp.CardNumberId, &emp.FirstName, &emp.LastName, &emp.WarehouseId); err != nil {
			return nil, err
		}

		db = append(db, emp)
	}
	err = rows.Err()
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrEmployeeNotFound
		}
		return
	}
	return
}

func (r *EmployeeMysql) GetById(id int) (emp internal.Employee, err error) {
	row := r.db.QueryRow("SELECT id, card_number_id, first_name, last_name, warehouse_id FROM employees WHERE id = ?", id)

	if err := row.Scan(&emp.Id, &emp.CardNumberId, &emp.FirstName, &emp.LastName, &emp.WarehouseId); err != nil {
		if err == sql.ErrNoRows {
			return emp, fmt.Errorf("employee not found")
		}
		return emp, err
	}
	return emp, nil
}

func (r *EmployeeMysql) Save(emp *internal.Employee) (int, error) {
	var err error

	var existingEmployee internal.Employee
	err = r.db.QueryRow(
		"SELECT id FROM employees WHERE card_number_id = ?",
		emp.CardNumberId).Scan(&existingEmployee.Id)

	if err == nil {
		return 0, fmt.Errorf("employee with card_number_id %s already exists", emp.CardNumberId)
	} else if err != sql.ErrNoRows {
		return 0, err
	}

	err = r.db.QueryRow(
		"INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES (?, ?, ?, ?) RETURNING id",
		emp.CardNumberId, emp.FirstName, emp.LastName, emp.WarehouseId).Scan(&emp.Id)

	if err != nil {
		return 0, err
	}

	return emp.Id, nil
}

func (r *EmployeeMysql) Update(id int, employee internal.Employee) error {
	_, err := r.db.Exec(
		"UPDATE employees SET card_number_id = ?, first_name = ?, last_name = ?, warehouse_id = ? WHERE id = ?",
		employee.CardNumberId, employee.FirstName, employee.LastName, employee.WarehouseId, id)
	if err != nil {
		return err
	}
	return nil
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
		row.Scan(&countInboundPerEmployee.CountInOrders, &countInboundPerEmployee.Id, &countInboundPerEmployee.FirstName, &countInboundPerEmployee.LastName, &countInboundPerEmployee.WarehouseId)

		io = append(io, countInboundPerEmployee)
	}
	return
}

func (r *EmployeeMysql) ReportInboundOrdersById(employeeId int) (io internal.InboundOrdersPerEmployee, err error) {

	row := r.db.QueryRow(
		InboundOrdersPerEmployeeByIdQuery,
		employeeId,
	)

	row.Scan(&io.CountInOrders, &io.Id, &io.FirstName, &io.LastName, &io.WarehouseId)

	return
}
