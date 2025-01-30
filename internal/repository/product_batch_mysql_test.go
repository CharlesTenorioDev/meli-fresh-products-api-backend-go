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

type MysqlProductBatchTestSuite struct {
	suite.Suite
	mock sqlmock.Sqlmock
	rp   *repository.ProductBatchMysql
}

func (s *MysqlProductBatchTestSuite) Setup() {
	db, mock, err := sqlmock.New()

	require.NoError(s.T(), err)
	s.mock = mock
	s.rp = repository.NewProductBatchMysql(db)
}

//docker exec -it mysql /bin/bash
//mysql -u root -p

func (s *MysqlProductBatchTestSuite) TestRepository_CreateProductBatchUnitTest() {
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
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.rp.Save(&prodBatch)

		require.NoError(t, err)
		require.EqualValues(t, expectedId, prodBatch.ID)
	})

	s.T().Run("fails, unprocessable entity", func(t *testing.T) {
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
		require.ErrorIs(t, internal.ErrProductBatchUnprocessableEntity, err)
	})

	s.T().Run("fails, unprocessable entity", func(t *testing.T) {
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
			WillReturnError(errors.New("Error SQL Query"))

		err := s.rp.Save(&prodBatch)

		require.Error(t, err)
		require.EqualError(t, errors.New("Error SQL Query"), err.Error())
	})
}

func (s *MysqlProductBatchTestSuite) TestRepository_ReadProductBatchUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedProdBatch := internal.ProductBatch{
			ID:                 1,
			BatchNumber:        1234,
			CurrentQuantity:    100,
			CurrentTemperature: 40.5,
			DueDate:            "2022-01-08",
			InitialQuantity:    120,
			ManufacturingDate:  "2022-01-01",
			ManufacturingHour:  15,
			MinumumTemperature: -8,
			ProductID:          1,
			SectionID:          3,
		}

		rows := sqlmock.NewRows(
			[]string{
				"id",
				"batch_number",
				"current_quantity",
				"current_temperature",
				"due_date",
				"initial_quantity",
				"manufacturing_date",
				"manufacturing_hour",
				"minumum_temperature",
				"product_id",
				"section_id",
			},
		).
			AddRow(1, 1234, 100, 40.5, "2022-01-08", 120, "2022-01-01", 15, -8, 1, 3)

		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualCarries, err := s.rp.FindByID(expectedProdBatch.ID)

		require.NoError(t, err)
		require.Equal(t, expectedProdBatch, actualCarries)
	})

	s.T().Run("query fails", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

		actualCarries, err := s.rp.FindByID(1)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrProductBatchNotFound, err)
		require.Zero(t, actualCarries)
	})

	s.T().Run("rows fails", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("row err"))

		actualCarries, err := s.rp.FindByID(1)

		require.Error(t, err)
		require.Equal(t, "row err", err.Error())
		require.Zero(t, actualCarries)
	})
}

func (s *MysqlProductBatchTestSuite) TestRepository_ProductBatchNumberProductBatchUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedExistsNumber := true

		rows := sqlmock.NewRows(
			[]string{
				"count",
			},
		).
			AddRow(1)

		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualCarries, err := s.rp.ProductBatchNumberExists(1234)

		require.NoError(t, err)
		require.Equal(t, expectedExistsNumber, actualCarries)
	})

	s.T().Run("fail, error query", func(t *testing.T) {
		s.Setup()
		expectedExistsNumber := false
		expectedError := errors.New("Error query")

		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("Error query"))

		actualCarries, err := s.rp.ProductBatchNumberExists(1234)

		require.Error(t, err)
		require.Equal(t, expectedExistsNumber, actualCarries)
		require.EqualError(t, expectedError, err.Error())
	})
}

func (s *MysqlProductBatchTestSuite) TestRepository_ReportProductsProductBatchUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedProdBatches := []internal.ProductBatch{
			{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
		}

		rows := sqlmock.NewRows(
			[]string{
				"batch_number",
				"current_quantity",
				"current_temperature",
				"due_date",
				"initial_quantity",
				"manufacturing_date",
				"manufacturing_hour",
				"minumum_temperature",
				"product_code",
				"section_number",
			},
		).
			AddRow(1234, 100, 40.5, "2022-01-08", 120, "2022-01-01", 15, -8, 1, 3)

		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualProdBatches, err := s.rp.ReportProducts()

		require.NoError(t, err)
		require.Equal(t, expectedProdBatches, actualProdBatches)
	})

	s.T().Run("query error", func(t *testing.T) {
		s.Setup()

		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("query error"))

		actualProdBatches, err := s.rp.ReportProducts()

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrProductBatchNotFound, err)
		require.Nil(t, actualProdBatches)
	})
}

func (s *MysqlProductBatchTestSuite) TestRepository_ReportProductsByIDProductBatchUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		expectedProdBatches := []internal.ProductBatch{
			{
				BatchNumber:        1234,
				CurrentQuantity:    100,
				CurrentTemperature: 40.5,
				DueDate:            "2022-01-08",
				InitialQuantity:    120,
				ManufacturingDate:  "2022-01-01",
				ManufacturingHour:  15,
				MinumumTemperature: -8,
				ProductID:          1,
				SectionID:          3,
			},
		}

		row := sqlmock.NewRows(
			[]string{
				"batch_number",
				"current_quantity",
				"current_temperature",
				"due_date",
				"initial_quantity",
				"manufacturing_date",
				"manufacturing_hour",
				"minumum_temperature",
				"product_code",
				"section_number",
			},
		).
			AddRow(1234, 100, 40.5, "2022-01-08", 120, "2022-01-01", 15, -8, 1, 3)

		s.mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(row)

		actualProdBatches, err := s.rp.ReportProductsByID(1)

		require.NoError(t, err)
		require.Equal(t, expectedProdBatches, actualProdBatches)
	})

	s.T().Run("no rows found", func(t *testing.T) {
		s.Setup()

		s.mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(sql.ErrNoRows)

		actualProdBatches, err := s.rp.ReportProductsByID(1)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrProductBatchNotFound, err)
		require.Nil(t, actualProdBatches)
	})

	s.T().Run("query error", func(t *testing.T) {
		s.Setup()

		s.mock.ExpectQuery("SELECT").WillReturnError(errors.New("query error"))

		actualProdBatches, err := s.rp.ReportProductsByID(1)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrProductBatchNotFound, err)
		require.Nil(t, actualProdBatches)
	})
}

func TestRepositoryMysqlProductBatchTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlProductBatchTestSuite))
}
