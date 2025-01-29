package repository_test

import (
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MysqlBuyerTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.BuyerMysqlRepository
}

func (s *MysqlBuyerTestSuite) Setup() {
	db, mock, e := sqlmock.New()

	require.NoError(s.T(), e)
	s.mock = mock
	s.rp = repository.NewBuyerMysqlRepository(db)
}

func (s *MysqlBuyerTestSuite) TestGetAll() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedBuyers := map[int]internal.Buyer{
			0: {
				ID:           0,
				CardNumberID: "CID001",
				FirstName:    "Fabio",
				LastName:     "Nacarelli",
			},
			1: {
				ID:           1,
				CardNumberID: "CID002",
				FirstName:    "Matheus",
				LastName:     "Apostulo",
			},
		}
		query := `
		SELECT
		id, card_number_id, first_name, last_name
		`
		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(0, "CID001", "Fabio", "Nacarelli").
			AddRow(1, "CID002", "Matheus", "Apostulo")
		s.mock.ExpectQuery(query).WillReturnRows(rows)

		buyers, err := s.rp.GetAll()

		require.NoError(t, err)
		require.Equal(t, buyers, expectedBuyers)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		query := `
		SELECT
		id, card_number_id, first_name, last_name
		`
		s.mock.ExpectQuery(query).WillReturnError(errors.New("internal server error"))

		buyers, err := s.rp.GetAll()

		require.Error(t, err)
		require.Zero(t, len(buyers))
	})
}

func (s *MysqlBuyerTestSuite) TestAdd() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		buyer := internal.Buyer{
			CardNumberID: "CID003",
			FirstName:    "Not",
			LastName:     "Found",
		}
		query := `
		INSERT INTO buyers
		`
		s.mock.ExpectExec(query).
			WithArgs(buyer.CardNumberID, buyer.FirstName, buyer.LastName).
			WillReturnResult(sqlmock.NewResult(3, 1))

		id, err := s.rp.Add(&buyer)

		require.NoError(t, err)
		require.Equal(t, int64(3), id)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		buyer := internal.Buyer{
			CardNumberID: "CID003",
			FirstName:    "Not",
			LastName:     "Found",
		}
		query := `
		INSERT INTO buyers
		`
		s.mock.ExpectExec(query).
			WithArgs(buyer.CardNumberID, buyer.FirstName, buyer.LastName).
			WillReturnError(errors.New("internal server error"))

		id, err := s.rp.Add(&buyer)

		require.Error(t, err)
		require.Zero(t, id)
	})
}

func (s *MysqlBuyerTestSuite) TestUpdate() {
	s.Setup()
	id := 1
	cardNumberId := "CID32131"
	firstName := "Apostulo"
	lastName := "Matheus"

	query := `
		SELECT
			id, card_number_id, first_name, last_name
	`
	rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
		AddRow(id, "CID002", "Matheus", "Apostulo")
	s.mock.ExpectQuery(query).WithArgs(id).
		WillReturnRows(rows)

	query = `
		UPDATE buyers
	`
	s.mock.ExpectExec(query).
		WithArgs(cardNumberId, firstName, lastName, id).
		WillReturnResult(sqlmock.NewResult(int64(id), 1))

	e := s.rp.Update(id, internal.BuyerPatch{
		CardNumberID: &cardNumberId,
		FirstName:    &firstName,
		LastName:     &lastName,
	})

	require.NoError(s.T(), e)
}

func (s *MysqlBuyerTestSuite) TestDelete() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		id := 1
		query := `
			DELETE FROM buyers
		`

		s.mock.ExpectExec(query).
			WithArgs(id).
			WillReturnResult(driver.RowsAffected(1))

		rowsAffected, err := s.rp.Delete(id)

		require.NoError(t, err)
		require.Equal(t, int64(1), rowsAffected)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		id := 10
		query := `
			DELETE FROM buyers
		`

		s.mock.ExpectExec(query).
			WithArgs(id).
			WillReturnError(errors.New("no such id"))

		_, err := s.rp.Delete(id)

		require.Error(t, err)
		require.Equal(t, "no such id", err.Error())
	})
}

func (s *MysqlBuyerTestSuite) TestReportPurchaseOrders() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		query := `SELECT*`
		expectedPurchaseOrders := []internal.PurchaseOrdersByBuyer{
			{
				BuyerID:             0,
				CardNumberID:        "CID001",
				FirstName:           "Fabio",
				LastName:            "Nacarelli",
				PurchaseOrdersCount: 10,
			},
			{
				BuyerID:             1,
				CardNumberID:        "CID002",
				FirstName:           "Matheus",
				LastName:            "Apostulo",
				PurchaseOrdersCount: 20,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(0, "CID001", "Fabio", "Nacarelli", 10).
			AddRow(1, "CID002", "Matheus", "Apostulo", 20)
		s.mock.ExpectQuery(query).WillReturnRows(rows)

		actualPurchaseOrders, err := s.rp.ReportPurchaseOrders()

		require.NoError(t, err)
		require.Equal(t, expectedPurchaseOrders, actualPurchaseOrders)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		query := `SELECT*`
		s.mock.ExpectQuery(query).WillReturnError(errors.New("internal server error"))

		_, err := s.rp.ReportPurchaseOrders()

		require.Error(t, err)
	})
}

func (s *MysqlBuyerTestSuite) TestReportPurchaseOrdersByID() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		id := 0
		query := `SELECT*`
		expectedPurchaseOrdersByBuyer := []internal.PurchaseOrdersByBuyer{
			{
				BuyerID:             0,
				CardNumberID:        "CID001",
				FirstName:           "Fabio",
				LastName:            "Nacarelli",
				PurchaseOrdersCount: 10,
			},
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(id, "CID001", "Fabio", "Nacarelli", 10)
		s.mock.ExpectQuery(query).WithArgs(id).WillReturnRows(rows)

		actualPurchaseOrdersByBuyer, err := s.rp.ReportPurchaseOrdersByID(id)

		require.NoError(t, err)
		require.Equal(t, expectedPurchaseOrdersByBuyer, actualPurchaseOrdersByBuyer)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.Setup()
		query := `SELECT*`
		s.mock.ExpectQuery(query).WillReturnError(errors.New("internal server error"))

		_, err := s.rp.ReportPurchaseOrdersByID(10)

		require.Error(t, err)
	})
}

func TestRepositoryMysqlBuyerUnit(t *testing.T) {
	suite.Run(t, new(MysqlBuyerTestSuite))
}
