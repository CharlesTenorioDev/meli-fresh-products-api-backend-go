package repository_test

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
)

func TestLocalityMysql_ReportCarries(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"carries_registered"}).AddRow(5)
		mock.ExpectQuery("SELECT COUNT(c.locality_id) carries_registered FROM carries c WHERE locality_id = ?").WithArgs(1).WillReturnRows(row)

		r := repository.NewLocalityMysql(db)
		amountOfCarries, err := r.ReportCarries(1)

		assert.NoError(t, err)
		assert.Equal(t, 5, amountOfCarries)
	})

	t.Run("No carries found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT COUNT(c.locality_id) carries_registered FROM carries c WHERE locality_id = ?").WithArgs(1).WillReturnError(sql.ErrNoRows)

		r := repository.NewLocalityMysql(db)
		amountOfCarries, err := r.ReportCarries(1)

		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Equal(t, 0, amountOfCarries)
	})

	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT COUNT(c.locality_id) carries_registered FROM carries c WHERE locality_id = ?").WithArgs(1).WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		amountOfCarries, err := r.ReportCarries(1)

		assert.Error(t, err)
		assert.Equal(t, 0, amountOfCarries)
	})
}

func TestLocalityMysql_GetAmountOfCarriesForEveryLocality(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"carries_count", "locality_id", "locality_name"}).
			AddRow(5, 1, "Locality 1").
			AddRow(10, 2, "Locality 2")
		mock.ExpectQuery(repository.AmountOfCarriesForEveryLocalityQuery).WillReturnRows(rows)

		r := repository.NewLocalityMysql(db)
		carries, err := r.GetAmountOfCarriesForEveryLocality()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(carries))
		assert.Equal(t, 5, carries[0].CarriesCount)
		assert.Equal(t, 10, carries[1].CarriesCount)
	})

	t.Run("No carries found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(repository.AmountOfCarriesForEveryLocalityQuery).WillReturnError(sql.ErrNoRows)

		r := repository.NewLocalityMysql(db)
		carries, err := r.GetAmountOfCarriesForEveryLocality()

		assert.ErrorIs(t, err, sql.ErrNoRows)
		assert.Empty(t, carries)
	})

	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(repository.AmountOfCarriesForEveryLocalityQuery).WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		carries, err := r.GetAmountOfCarriesForEveryLocality()

		assert.Error(t, err)
		assert.Empty(t, carries)
	})

	t.Run("Row Scan Error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(repository.AmountOfCarriesForEveryLocalityQuery).WillReturnRows(sqlmock.NewRows([]string{"carries_count", "locality_id", "locality_name"}).AddRow(43, "f", 43))

		r := repository.NewLocalityMysql(db)
		carries, err := r.GetAmountOfCarriesForEveryLocality()

		assert.Error(t, err)
		assert.Empty(t, carries)
	})
}

func TestLocalityMysql_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		locality := &internal.Locality{
			ID:           1,
			LocalityName: "Locality 1",
			ProvinceName: "Province 1",
			CountryName:  "Country 1",
		}

		mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
			WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := repository.NewLocalityMysql(db)
		err = r.Save(locality)

		assert.NoError(t, err)
	})

	t.Run("Locality conflict", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		locality := &internal.Locality{
			ID:           1,
			LocalityName: "Locality 1",
			ProvinceName: "Province 1",
			CountryName:  "Country 1",
		}

		mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
			WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
			WillReturnError(&mysql.MySQLError{Number: 1062})

		r := repository.NewLocalityMysql(db)
		err = r.Save(locality)

		assert.ErrorIs(t, err, internal.ErrLocalityConflict)
	})

	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		locality := &internal.Locality{
			ID:           1,
			LocalityName: "Locality 1",
			ProvinceName: "Province 1",
			CountryName:  "Country 1",
		}

		mock.ExpectExec("INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)").
			WithArgs(locality.ID, locality.LocalityName, locality.ProvinceName, locality.CountryName).
			WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		err = r.Save(locality)

		assert.Error(t, err)
	})
}

func TestLocalityMysql_ReportSellers(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name", "COUNT(s.id)"}).
			AddRow(1, "Locality 1", "Province 1", "Country 1", 5).
			AddRow(2, "Locality 2", "Province 2", "Country 2", 10)
		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").
			WillReturnRows(rows)

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellers()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(localities))
		assert.Equal(t, 5, localities[0].Sellers)
		assert.Equal(t, 10, localities[1].Sellers)
	})

	t.Run("No localities found", func(t *testing.T) {
		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").
			WillReturnError(sql.ErrNoRows)

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellers()

		assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
		assert.Empty(t, localities)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").
			WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellers()

		assert.Error(t, err)
		assert.Empty(t, localities)
	})

	t.Run("Row Scan Error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name", "COUNT(s.id)"}).AddRow(1, "Locality 1", 5, "Province 1", "Country 1")
		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id").WillReturnRows(rows)

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellers()

		assert.Error(t, err)
		assert.Empty(t, localities)
	})
}

func TestLocalityMysql_ReportSellersByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		row := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name", "COUNT(s.id)"}).
			AddRow(1, "Locality 1", "Province 1", "Country 1", 5)
		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").
			WithArgs(1).
			WillReturnRows(row)

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellersByID(1)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(localities))
		assert.Equal(t, 5, localities[0].Sellers)
	})

	t.Run("No localities found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellersByID(1)

		assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
		assert.Empty(t, localities)
	})

	t.Run("Database error", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		localities, err := r.ReportSellersByID(1)

		assert.Error(t, err)
		assert.Empty(t, localities)
	})
}

func TestLocalityMysql_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "name", "province_name", "country_name"}).
			AddRow(1, "Locality 1", "Province 1", "Country 1")
		mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").
			WithArgs(1).
			WillReturnRows(row)

		r := repository.NewLocalityMysql(db)
		locality, err := r.FindByID(1)

		assert.NoError(t, err)
		assert.Equal(t, 1, locality.ID)
	})

	t.Run("Locality not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		r := repository.NewLocalityMysql(db)
		locality, err := r.FindByID(1)

		assert.ErrorIs(t, err, internal.ErrLocalityNotFound)
		assert.Equal(t, internal.Locality{}, locality)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		r := repository.NewLocalityMysql(db)
		locality, err := r.FindByID(1)

		assert.Error(t, err)
		assert.Equal(t, internal.Locality{}, locality)
	})
}
