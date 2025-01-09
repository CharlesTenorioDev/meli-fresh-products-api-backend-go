package repository_test

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestLocalityMysql_Save(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	locality := &internal.Locality{
		ID:           1,
		LocalityName: "Test Locality",
		ProvinceName: "Test Province",
		CountryName:  "Test Country",
	}

	mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnResult(sqlmock.NewResult(1, 1))

	repo := repository.NewLocalityMysql(mockDB)

	err = repo.Save(locality)
	assert.NoError(t, err)
}

func TestLocalityMysql_Save_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	locality := &internal.Locality{
		ID:           1,
		LocalityName: "Test Locality",
		ProvinceName: "Test Province",
		CountryName:  "Test Country",
	}

	mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnError(fmt.Errorf("some error"))

	repo := repository.NewLocalityMysql(mockDB)

	err = repo.Save(locality)
	assert.Error(t, err)
}

func TestLocalityMysql_Save_Conflict(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	locality := &internal.Locality{
		ID:           1,
		LocalityName: "Test Locality",
		ProvinceName: "Test Province",
		CountryName:  "Test Country",
	}

	mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
		WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
		WillReturnError(&mysql.MySQLError{Number: 1062})

	repo := repository.NewLocalityMysql(mockDB)

	err = repo.Save(locality)
	assert.ErrorIs(t, err, internal.ErrLocalityConflict)
}

func TestLocalityMysql_ReportSellers(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name", "COUNT(s.id)"}).
		AddRow(1, "Test Locality", "Test Province", "Test Country", 2).
		AddRow(2, "Another Locality", "Another Province", "Another Country", 1)
	mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").WillReturnRows(rows)

	repo := repository.NewLocalityMysql(mockDB)

	localities, err := repo.ReportSellers()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(localities))
	assert.Equal(t, 1, localities[0].ID)
	assert.Equal(t, 2, localities[0].Sellers)
	assert.Equal(t, 2, localities[1].ID)
	assert.Equal(t, 1, localities[1].Sellers)
}

func TestLocalityMysql_ReportSellers_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").WillReturnError(sql.ErrNoRows)

	repo := repository.NewLocalityMysql(mockDB)

	localities, err := repo.ReportSellers()
	assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
	assert.Equal(t, 0, len(localities))
}

func TestLocalityMysql_ReportSellersByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name", "COUNT(s.id)"}).
		AddRow(1, "Test Locality", "Test Province", "Test Country", 2)
	mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").WithArgs(1).WillReturnRows(row)

	repo := repository.NewLocalityMysql(mockDB)

	localities, err := repo.ReportSellersByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(localities))
	assert.Equal(t, 1, localities[0].ID)
	assert.Equal(t, 2, localities[0].Sellers)
}

func TestLocalityMysql_ReportSellersByID_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

	repo := repository.NewLocalityMysql(mockDB)

	_, err = repo.ReportSellersByID(1)
	assert.Error(t, err)
}

func TestLocalityMysql_ReportSellersByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").WithArgs(1).WillReturnError(sql.ErrNoRows)

	repo := repository.NewLocalityMysql(mockDB)

	_, err = repo.ReportSellersByID(1)
	assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
}

func TestLocalityMysql_FindByID(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	row := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name"}).
		AddRow(1, "Test Locality", "Test Province", "Test Country")
	mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").WithArgs(1).WillReturnRows(row)

	repo := repository.NewLocalityMysql(mockDB)

	locality, err := repo.FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, 1, locality.ID)
}

func TestLocalityMysql_FindByID_Error(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").WithArgs(1).WillReturnError(fmt.Errorf("some error"))

	repo := repository.NewLocalityMysql(mockDB)

	_, err = repo.FindByID(1)
	assert.Error(t, err)
}

func TestLocalityMysql_FindByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").WithArgs(1).WillReturnError(sql.ErrNoRows)

	repo := repository.NewLocalityMysql(mockDB)

	_, err = repo.FindByID(1)
	assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
}
