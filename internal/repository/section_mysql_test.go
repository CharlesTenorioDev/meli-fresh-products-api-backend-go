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

func (s *MysqlSectionTestSuite) TestRepository_ReadSectionUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		sections := []internal.Section{
			{ID: 1, SectionNumber: 123, CurrentTemperature: 22, MinimumTemperature: 15, CurrentCapacity: 60, MinimumCapacity: 30, MaximumCapacity: 100, WarehouseID: 1, ProductTypeID: 1},
			{ID: 2, SectionNumber: 456, CurrentTemperature: 23, MinimumTemperature: 16, CurrentCapacity: 70, MinimumCapacity: 35, MaximumCapacity: 110, WarehouseID: 2, ProductTypeID: 2},
		}
		s.Setup()
		rows := sqlmock.NewRows([]string{"id", "section_number", "current_temperature", "minimum_temperature", "current_capacity", "minimum_capacity", "maximum_capacity", "warehouse_id", "product_type_id"}).
			AddRow(sections[0].ID, sections[0].SectionNumber, sections[0].CurrentTemperature, sections[0].MinimumTemperature, sections[0].CurrentCapacity, sections[0].MinimumCapacity, sections[0].MaximumCapacity, sections[0].WarehouseID, sections[0].ProductTypeID).
			AddRow(sections[1].ID, sections[1].SectionNumber, sections[1].CurrentTemperature, sections[1].MinimumTemperature, sections[1].CurrentCapacity, sections[1].MinimumCapacity, sections[1].MaximumCapacity, sections[1].WarehouseID, sections[1].ProductTypeID)

		s.mock.ExpectQuery("SELECT `id`, `section_number`, `current_temperature`, `minimum_temperature`, `current_capacity`, `minimum_capacity`, `maximum_capacity`, `warehouse_id`, `product_type_id` FROM sections").
			WillReturnRows(rows)

		result, err := s.rp.FindAll()

		require.NoError(t, err)
		require.EqualValues(t, sections, result)
	})
	s.T().Run("no rows found", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT `id`, `section_number`, `current_temperature`, `minimum_temperature`, `current_capacity`, `minimum_capacity`, `maximum_capacity`, `warehouse_id`, `product_type_id` FROM sections").
			WillReturnError(sql.ErrNoRows)

		_, err := s.rp.FindAll()

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrSectionNotFound, err)
	})
}

func (s *MysqlSectionTestSuite) TestRepository_FindByIDSectionUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		s.Setup()
		section := internal.Section{ID: 1, SectionNumber: 123, CurrentTemperature: 22, MinimumTemperature: 15, CurrentCapacity: 50, MinimumCapacity: 30, MaximumCapacity: 100, WarehouseID: 2, ProductTypeID: 2}

		rows := sqlmock.NewRows(
			[]string{
				"id",
				"section_number",
				"current_temperature",
				"minimum_temperature",
				"current_capacity",
				"minimum_capacity",
				"maximum_capacity",
				"warehouse_id",
				"product_type_id",
			},
		).
			AddRow(1, 123, 22, 15, 50, 30, 100, 2, 2)

		s.mock.ExpectQuery("SELECT").WillReturnRows(rows)

		actualSection, err := s.rp.FindByID(1)

		require.NoError(t, err)
		require.Equal(t, section, actualSection)
	})

	s.T().Run("rows fails", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)

		actualSection, err := s.rp.FindByID(2)

		require.Error(t, err)
		require.EqualError(t, internal.ErrSectionNotFound, err.Error())
		require.Zero(t, actualSection)
	})
}

func (s *MysqlSectionTestSuite) TestRepository_SaveSectionUnitTest() {
	s.T().Run("success", func(t *testing.T) {
		expectedID := 1
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
		s.mock.ExpectExec("INSERT INTO sections .*").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID).
			WillReturnResult(sqlmock.NewResult(int64(expectedID), 1))

		err := s.rp.Save(&section)

		require.NoError(t, err)
		require.EqualValues(t, expectedID, section.ID)
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
		s.mock.ExpectExec("INSERT INTO sections .*").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID).
			WillReturnError(&mysql.MySQLError{
				Number: 1062, // Duplicate entry error
			})

		err := s.rp.Save(&section)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrSectionUnprocessableEntity, err)
	})
}

func (s *MysqlSectionTestSuite) TestRepository_UpdateSectionUnitTest() {
	s.T().Run("update success", func(t *testing.T) {
		section := internal.Section{
			ID:                 1,
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
		s.mock.ExpectExec("UPDATE sections SET .* WHERE id = ?").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID, section.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := s.rp.Update(&section)

		require.NoError(t, err)
	})

	s.T().Run("update fails, section not found", func(t *testing.T) {
		section := internal.Section{
			ID:                 1,
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
		s.mock.ExpectExec("UPDATE sections SET .* WHERE id = ?").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID, section.ID).
			WillReturnError(&mysql.MySQLError{Number: 1064})

		err := s.rp.Update(&section)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrSectionNotFound, err)
	})

	s.T().Run("update fails, unprocessable entity", func(t *testing.T) {
		section := internal.Section{
			ID:                 1,
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
		s.mock.ExpectExec("UPDATE sections SET .* WHERE id = ?").
			WithArgs(section.SectionNumber, section.CurrentTemperature, section.MinimumTemperature, section.CurrentCapacity, section.MinimumCapacity, section.MaximumCapacity, section.WarehouseID, section.ProductTypeID, section.ID).
			WillReturnError(&mysql.MySQLError{Number: 1062})

		err := s.rp.Update(&section)

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrSectionUnprocessableEntity, err)
	})
}

func (s *MysqlSectionTestSuite) TestRepository_SectionNumberSectionUnitTest() {
	s.T().Run("exists", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM sections WHERE section_number = ?").
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		exists, err := s.rp.SectionNumberExists(123)

		require.NoError(t, err)
		require.True(t, exists)
	})

	s.T().Run("not exists", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM sections WHERE section_number = ?").
			WithArgs(123).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		exists, err := s.rp.SectionNumberExists(123)

		require.NoError(t, err)
		require.False(t, exists)
	})

	s.T().Run("error query", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM sections WHERE section_number = ?").
			WithArgs(123).
			WillReturnError(errors.New("error query"))

		_, err := s.rp.SectionNumberExists(123)

		require.Error(t, err)
		require.EqualError(t, errors.New("error query"), err.Error())
	})
}

func (s *MysqlSectionTestSuite) TestRepository_ReportProductsSectionUnitTest() {
	s.T().Run("returns report", func(t *testing.T) {
		report := []internal.ReportProduct{
			{SectionID: 1, SectionNumber: 123, ProductsCount: 100},
		}
		s.Setup()
		s.mock.ExpectQuery("SELECT .* FROM sections s LEFT JOIN product_batches pb ON s.id = pb.section_id").
			WillReturnRows(sqlmock.NewRows([]string{"section_id", "section_number", "products_count"}).
				AddRow(report[0].SectionID, report[0].SectionNumber, report[0].ProductsCount))

		result, err := s.rp.ReportProducts()
		require.NoError(t, err)
		require.EqualValues(t, report, result)
	})
	s.T().Run("no report available", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT .* FROM sections s LEFT JOIN product_batches pb ON s.id = pb.section_id").
			WillReturnError(sql.ErrNoRows)

		_, err := s.rp.ReportProducts()

		require.Error(t, err)
		require.ErrorIs(t, internal.ErrReportProductNotFound, err)
	})
}

func (s *MysqlSectionTestSuite) TestRepository_ReportProductsByIDSectionUnitTest() {
	s.T().Run("returns report for ID", func(t *testing.T) {
		report := internal.ReportProduct{SectionID: 1, SectionNumber: 123, ProductsCount: 100}
		s.Setup()
		s.mock.ExpectQuery("SELECT .* FROM sections s LEFT JOIN product_batches pb ON s.id = pb.section_id WHERE s.id = ?").
			WithArgs(report.SectionID).
			WillReturnRows(sqlmock.NewRows([]string{"section_id", "section_number", "products_count"}).
				AddRow(report.SectionID, report.SectionNumber, report.ProductsCount))

		result, err := s.rp.ReportProductsByID(report.SectionID)
		require.NoError(t, err)
		require.EqualValues(t, report, result)
	})

	s.T().Run("no report for ID", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT .* FROM sections s LEFT JOIN product_batches pb ON s.id = pb.section_id WHERE s.id = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		report, err := s.rp.ReportProductsByID(1)
		require.NoError(t, err)
		require.EqualValues(t, 0, report.ProductsCount)
	})

	s.T().Run("no report for ID error", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectQuery("SELECT .* FROM sections s LEFT JOIN product_batches pb ON s.id = pb.section_id WHERE s.id = ?").
			WithArgs(1).
			WillReturnError(errors.New("Error query"))

		_, err := s.rp.ReportProductsByID(1)
		require.Error(t, err)
		require.EqualError(t, errors.New("Error query"), err.Error())
	})
}

func (s *MysqlSectionTestSuite) TestRepository_DeleteSectionUnitTest() {
	s.T().Run("delete successfully", func(t *testing.T) {
		s.Setup()
		s.mock.ExpectExec("DELETE FROM sections WHERE id = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.rp.Delete(1)
		require.NoError(t, err)
	})
}

func TestRepositoryMysqlSectionTestSuite(t *testing.T) {
	suite.Run(t, new(MysqlSectionTestSuite))
}
