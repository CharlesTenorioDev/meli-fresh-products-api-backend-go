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

type MysqlProductBatchTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.ProductBatchMysql
}

func (s *MysqlProductBatchTestSuite) Setup() {
	db, mock, e := sqlmock.New()

	require.NoError(s.T(), e)
	s.mock = mock
	s.rp = repository.NewProductBatchMysql(db)
}

//docker exec -it mysql /bin/bash
//mysql -u root -p

func (s *MysqlProductBatchTestSuite) TestRepository_CreateProductBatchUnitTest(t *testing.T) {
	s.T().Run("success", func(t *testing.T) {
		expectedId := 1
		prodBatch := internal.ProductBatch{
			BatchNumber:        1234,
			CurrentQuantity:    100,
			CurrentTemperature: 40.5,
			DueDate:            "2022-01-08",
			InitialQuantity:    120,
			ManufacturingDate:  "2022-01-01 ",
			ManufacturingHour:  15,
			MinumumTemperature: -8,
			ProductID:          1,
			SectionID:          3,
		}

		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(prodBatch.BatchNumber, prodBatch.CurrentQuantity, prodBatch.CurrentTemperature, prodBatch.DueDate,
				prodBatch.InitialQuantity, prodBatch.ManufacturingDate, prodBatch.ManufacturingHour, prodBatch.MinumumTemperature,
				prodBatch.ProductID, prodBatch.SectionID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.rp.Save(&prodBatch)

		require.NoError(t, err)
		require.EqualValues(t, expectedId, prodBatch.ID)
	})
	s.T().Run("fails, cid already exists", func(t *testing.T) {
		prodBatch := internal.ProductBatch{
			BatchNumber:        1234,
			CurrentQuantity:    100,
			CurrentTemperature: 40.5,
			DueDate:            "2022-01-08",
			InitialQuantity:    120,
			ManufacturingDate:  "2022-01-01 ",
			ManufacturingHour:  15,
			MinumumTemperature: -8,
			ProductID:          1,
			SectionID:          3,
		}

		s.Setup()
		s.mock.ExpectExec("INSERT").
			WithArgs(prodBatch.BatchNumber, prodBatch.CurrentQuantity, prodBatch.CurrentTemperature, prodBatch.DueDate,
				prodBatch.InitialQuantity, prodBatch.ManufacturingDate, prodBatch.ManufacturingHour, prodBatch.MinumumTemperature,
				prodBatch.ProductID, prodBatch.SectionID).
			WillReturnError(&mysql.MySQLError{
				Number: 1062,
			})

		err := s.rp.Save(&prodBatch)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrProductBatchAlreadyExists, err)
	})
}

func TestRepository_ReadProductBatchUnitTest(t *testing.T) {
}

func TestRepository_ProductBatchNumberProductBatchUnitTest(t *testing.T) {
}

func TestRepository_ReportProductsProductBatchUnitTest(t *testing.T) {
}

func TestRepository_ReportProductsByIDProductBatchUnitTest(t *testing.T) {
}
