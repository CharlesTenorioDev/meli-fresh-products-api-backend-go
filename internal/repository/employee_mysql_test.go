package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MysqlEmployeeTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.EmployeeMysql
}

func (s *MysqlEmployeeTestSuite) Setup() {
	db, mock, e := sqlmock.New()

	require.NoError(s.T(), e)
	s.mock = mock
	s.rp = repository.NewEmployeeMysql(db)
}

func (s *MysqlEmployeeTestSuite) TestGetAll() {
	s.T().Run("success", func(t *testing.T) {
		expectedEmployees := []internal.Employee{
			{
				ID:           0,
				CardNumberID: "CID000",
				FirstName:    "Fabio",
				LastName:     "Nacarelli",
				WarehouseID:  1,
			},
			{
				ID:           1,
				CardNumberID: "CID001",
				FirstName:    "Matheus",
				LastName:     "Apostulo",
				WarehouseID:  2,
			},
		}
		s.Setup()
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
			AddRow(0, "CID000", "Fabio", "Nacarelli", 1).
			AddRow(1, "CID001", "Matheus", "Apostulo", 2)
		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)
		actualEmployees, e := s.rp.GetAll()

		require.NoError(t, e)
		require.Equal(t, expectedEmployees, actualEmployees)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("internal server error"))
		actualEmployees, e := s.rp.GetAll()

		require.Error(t, e)
		require.Zero(t, actualEmployees)
	})
}

func (s *MysqlEmployeeTestSuite) TestGetByID() {
	id := 0
	expectedEmployee := internal.Employee{
		ID:           id,
		CardNumberID: "CID000",
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		WarehouseID:  1,
	}
	s.Setup()
	row := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "warehouse_id"}).
		AddRow(id, "CID000", "Fabio", "Nacarelli", 1)
	s.mock.ExpectQuery("SELECT").WithArgs(id).WillReturnRows(row)
	actualEmployee, e := s.rp.GetByID(id)

	require.NoError(s.T(), e)
	require.Equal(s.T(), expectedEmployee, actualEmployee)
}

func (s *MysqlEmployeeTestSuite) TestSave() {
	s.T().Run("success", func(t *testing.T) {
		expectedId := 1
		expectedEmployee := internal.Employee{
			CardNumberID: "CID000",
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
		}
		s.Setup()
		s.mock.ExpectQuery("SELECT").
			WithArgs(expectedEmployee.CardNumberID).
			WillReturnError(sql.ErrNoRows)
		s.mock.ExpectExec("INSERT").WithArgs(
			expectedEmployee.CardNumberID,
			expectedEmployee.FirstName,
			expectedEmployee.LastName,
			expectedEmployee.WarehouseID,
		).WillReturnResult(sqlmock.NewResult(int64(1), 1))
		actualId, e := s.rp.Save(&expectedEmployee)

		require.NoError(t, e)
		require.EqualValues(t, expectedId, actualId)
	})
	s.T().Run("failure, user with card number id already exists", func(t *testing.T) {
		expectedEmployee := internal.Employee{
			CardNumberID: "CID000",
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
		}
		s.Setup()
		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		s.mock.ExpectQuery("SELECT").
			WithArgs(expectedEmployee.CardNumberID).WillReturnRows(rows)
		_, e := s.rp.Save(&expectedEmployee)

		require.Error(t, e)
	})
	s.T().Run("failure, generic error, != sql.ErrnoRows", func(t *testing.T) {
		expectedEmployee := internal.Employee{
			CardNumberID: "CID000",
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
		}
		s.Setup()
		s.mock.ExpectQuery("SELECT").
			WithArgs(expectedEmployee.CardNumberID).WillReturnError(errors.New("internal server error"))
		_, e := s.rp.Save(&expectedEmployee)

		require.Error(t, e)
	})
	s.T().Run("failure, insert fails", func(t *testing.T) {
		expectedEmployee := internal.Employee{
			CardNumberID: "CID000",
			FirstName:    "Fabio",
			LastName:     "Nacarelli",
			WarehouseID:  1,
		}
		s.mock.ExpectQuery("SELECT").
			WithArgs(expectedEmployee.CardNumberID).
			WillReturnError(sql.ErrNoRows)
		s.mock.ExpectExec("INSERT").WithArgs(
			expectedEmployee.CardNumberID,
			expectedEmployee.FirstName,
			expectedEmployee.LastName,
			expectedEmployee.WarehouseID,
		).WillReturnError(errors.New("insert failed"))
		_, e := s.rp.Save(&expectedEmployee)

		require.Error(t, e)
		require.Equal(t, "insert failed", e.Error())
	})
}

func (s *MysqlEmployeeTestSuite) TestUpdate() {
	s.Setup()
	employee := internal.Employee{
		ID:           0,
		CardNumberID: "CID002",
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		WarehouseID:  2,
	}
	s.mock.ExpectExec("UPDATE").
		WithArgs(employee.CardNumberID, employee.FirstName, employee.LastName, employee.WarehouseID, employee.ID).
		WillReturnResult(sqlmock.NewResult(int64(employee.ID), 1))

	e := s.rp.Update(employee.ID, employee)

	require.NoError(s.T(), e)
}

func (s *MysqlEmployeeTestSuite) TestDelete() {
	s.Setup()
	id := 0
	s.mock.ExpectExec("DELETE").
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	e := s.rp.Delete(id)

	require.NoError(s.T(), e)
}

func (s *MysqlEmployeeTestSuite) TestCountInboundOrdersPerEmployee() {
	s.T().Run("success", func(t *testing.T) {
		expectedInbourdOrdersPerEmployee := []internal.InboundOrdersPerEmployee{
			{
				CountInOrders: 10,
				ID:            0,
				CardNumberID:  "CID000",
				FirstName:     "Fabio",
				LastName:      "Nacarelli",
				WarehouseID:   0,
			},
			{
				CountInOrders: 15,
				ID:            1,
				CardNumberID:  "CID001",
				FirstName:     "Matheus",
				LastName:      "Apostulo",
				WarehouseID:   1,
			},
		}
		s.Setup()
		rows := sqlmock.NewRows(
			[]string{
				"inbound_orders_count",
				"id",
				"card_number_id",
				"first_name",
				"last_name",
				"warehouse_id",
			},
		).
			AddRow(10, 0, "CID000", "Fabio", "Nacarelli", 0).
			AddRow(15, 1, "CID001", "Matheus", "Apostulo", 1)
		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualInboundOrdersPerEmployee, e := s.rp.CountInboundOrdersPerEmployee()

		require.NoError(t, e)
		require.Equal(t, expectedInbourdOrdersPerEmployee, actualInboundOrdersPerEmployee)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

		_, e := s.rp.CountInboundOrdersPerEmployee()

		require.ErrorIs(t, sql.ErrNoRows, e)
	})
}

func (s *MysqlEmployeeTestSuite) TestReportInboundOrdersByID() {
	s.T().Run("success", func(t *testing.T) {
		expectedIo := internal.InboundOrdersPerEmployee{
			CountInOrders: 10,
			ID:            0,
			CardNumberID:  "CID000",
			FirstName:     "Fabio",
			LastName:      "Nacarelli",
			WarehouseID:   0,
		}
		s.Setup()
		rows := sqlmock.NewRows(
			[]string{
				"inbound_orders_count",
				"id",
				"card_number_id",
				"first_name",
				"last_name",
				"warehouse_id",
			},
		).
			AddRow(10, 0, "CID000", "Fabio", "Nacarelli", 0)
		s.mock.ExpectQuery("SELECT").
			WithArgs(expectedIo.ID).
			WillReturnRows(rows)

		actualIo, e := s.rp.ReportInboundOrdersByID(expectedIo.ID)

		require.NoError(t, e)
		require.Equal(t, expectedIo, actualIo)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

		_, e := s.rp.ReportInboundOrdersByID(10)

		require.ErrorIs(t, internal.ErrEmployeeNotFound, e)
	})
}

func TestRepositoryMysqlEmployeeUnit(t *testing.T) {
	suite.Run(t, new(MysqlEmployeeTestSuite))
}
