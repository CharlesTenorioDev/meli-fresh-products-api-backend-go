package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MysqlSectionTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.SectionMysql
}

func (s *MysqlSectionTestSuite) Setup() {
	db, mock, e := sqlmock.New()

	require.NoError(s.T(), e)
	s.mock = mock
	s.rp = repository.NewSectionMysql(db)
}

//docker exec -it mysql /bin/bash
//mysql -u root -p

func (s *MysqlSectionTestSuite) TestRepository_CreateSectionUnitTest(t *testing.T) {
	s.T().Run("success", func(t *testing.T) {
		expectedId := 1
		section := internal.Section{
			SectionNumber:      123,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity,
				section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.rp.Save(&section)

		require.NoError(t, err)
		require.EqualValues(t, expectedId, section.ID)
	})
	s.T().Run("fails, cid already exists", func(t *testing.T) {
		section := internal.Section{
			SectionNumber:      123,
			CurrentTemperature: 22.5,
			MinimumTemperature: 15.0,
			CurrentCapacity:    50,
			MinimumCapacity:    30,
			MaximumCapacity:    100,
			WarehouseID:        2,
			ProductTypeID:      2,
		}

		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity,
				section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID).
			WillReturnError(&mysql.MySQLError{
				Number: 1062,
			})

		err := s.rp.Save(&section)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrSectionUnprocessableEntity, err)
	})
}

func TestRepository_ReadSectionUnitTest(t *testing.T) {
}

func TestRepository_SectionNumberSectionUnitTest(t *testing.T) {
}

func TestRepository_ReportProductsSectionUnitTest(t *testing.T) {
}

func TestRepository_ReportProductsByIDSectionUnitTest(t *testing.T) {
}
