package repository_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MysqlCarriesTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.CarriesMysql
}

func (s *MysqlCarriesTestSuite) Setup() {
	db, mock, e := sqlmock.New()

	require.NoError(s.T(), e)
	s.mock = mock
	s.rp = repository.NewCarriesMysql(db)
}

func (s *MysqlCarriesTestSuite) TestFindAll() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedCarries := []internal.Carries{
			{
				ID:          0,
				Cid:         "CID000",
				CompanyName: "Meli",
				Address:     "OneTwoThree",
				PhoneNumber: "11977021447",
				LocalityID:  0,
			},
			{
				ID:          1,
				Cid:         "CID001",
				CompanyName: "Go Meli Go",
				Address:     "FourFiveSix",
				PhoneNumber: "11977021447",
				LocalityID:  1,
			},
		}
		rows := sqlmock.NewRows(
			[]string{
				"id",
				"cid",
				"company_name",
				"address",
				"phone_number",
				"locality_id",
			},
		).
			AddRow(0, "CID000", "Meli", "OneTwoThree", "11977021447", 0).
			AddRow(1, "CID001", "Go Meli Go", "FourFiveSix", "11977021447", 1)
		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualCarries, e := s.rp.FindAll()

		require.NoError(t, e)
		require.Equal(t, expectedCarries, actualCarries)
	})
	s.T().Run("query fails", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

		actualCarries, e := s.rp.FindAll()

		require.Error(t, e)
		require.ErrorIs(t, sql.ErrNoRows, e)
		require.Zero(t, actualCarries)
	})
	s.T().Run("rows fails", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("row err"))

		actualCarries, e := s.rp.FindAll()

		require.Error(t, e)
		require.Equal(t, "row err", e.Error())
		require.Zero(t, actualCarries)
	})
}

func (s *MysqlCarriesTestSuite) TestCreate() {
	s.T().Run("success", func(t *testing.T) {
		expectedId := 0
		carry := internal.Carries{
			Cid:         "CID000",
			CompanyName: "Meli",
			Address:     "OneTwoThree",
			PhoneNumber: "11218921",
			LocalityID:  0,
		}
		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.PhoneNumber, carry.LocalityID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualID, e := s.rp.Create(carry)

		require.NoError(t, e)
		require.EqualValues(t, expectedId, actualID)
	})
	s.T().Run("fails, cid already exists", func(t *testing.T) {
		carry := internal.Carries{
			Cid:         "CID000",
			CompanyName: "Meli",
			Address:     "OneTwoThree",
			PhoneNumber: "11218921",
			LocalityID:  0,
		}
		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.PhoneNumber, carry.LocalityID).
			WillReturnError(&mysql.MySQLError{
				Number: 1062,
			})

		_, e := s.rp.Create(carry)

		require.Error(t, e)
		require.ErrorIs(t, repository.ErrCidAlreadyExists, e)
	})
	s.T().Run("fails, no such locality id", func(t *testing.T) {
		carry := internal.Carries{
			Cid:         "CID000",
			CompanyName: "Meli",
			Address:     "OneTwoThree",
			PhoneNumber: "11218921",
			LocalityID:  0,
		}
		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(carry.Cid, carry.CompanyName, carry.Address, carry.PhoneNumber, carry.LocalityID).
			WillReturnError(&mysql.MySQLError{
				Number: 1452,
			})

		_, e := s.rp.Create(carry)

		require.Error(t, e)
		require.ErrorIs(t, repository.ErrNoSuchLocalityID, e)
	})
}

func TestRepositoryMysqlCarriesUnit(t *testing.T) {
	suite.Run(t, new(MysqlCarriesTestSuite))
}
