package repository

import (
	"database/sql"

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
	SELECT COUNT(io.employee_id) inbound_orders_count, io.id, e.first_name, e.last_name, io.warehouse_id
	FROM inbound_orders io
	INNER JOIN employees e ON io.employee_id = e.id
	GROUP BY io.employee_id`
)

func (r *EmployeeMysql) GetAll() (db map[int]internal.Employee) {
	return
}

func (r *EmployeeMysql) GetById(id int) (emp internal.Employee, err error) {
	return
}
func (r *EmployeeMysql) Save(emp *internal.Employee) int {
	return 0
}

func (r *EmployeeMysql) Update(id int, employee internal.Employee) {
	return
}

func (r *EmployeeMysql) Delete(id int) {
	return
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

func (r *EmployeeMysql) ReportInboundOrdersById(employeeId int) (totalInboundOrders int, err error) {

	row := r.db.QueryRow(
		"SELECT COUNT(io.employee_id) report_inbound_orders FROM inbound_orders io WHERE employee_id = ?",
		employeeId,
	)

	row.Scan(&totalInboundOrders)
	if totalInboundOrders == 0 {
		err = sql.ErrNoRows
	}
	return
}
