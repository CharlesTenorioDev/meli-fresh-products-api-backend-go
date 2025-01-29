package repository

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
)

func TestSellerMysql_FindAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
			AddRow(1, 123, "Company 1", "Address 1", "1234567890").
			AddRow(2, 456, "Company 2", "Address 2", "9876543210")
		mock.ExpectQuery("SELECT `s.id`, `s.cid`, `s.company_name`, `s.address`, `s.telephone` FROM `sellers` AS `s`").WillReturnRows(rows)

		r := NewSellerMysql(db)
		sellers, err := r.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, 2, len(sellers))
		assert.Equal(t, 1, sellers[0].ID)
		assert.Equal(t, 2, sellers[1].ID)
	})

	t.Run("No sellers found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `s.id`, `s.cid`, `s.company_name`, `s.address`, `s.telephone` FROM `sellers` AS `s`").WillReturnError(sql.ErrNoRows)

		r := NewSellerMysql(db)
		sellers, err := r.FindAll()

		assert.ErrorIs(t, err, internal.ErrSellerNotFound)
		assert.Empty(t, sellers)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT `s.id`, `s.cid`, `s.company_name`, `s.address`, `s.telephone` FROM `sellers` AS `s`").WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		sellers, err := r.FindAll()

		assert.Error(t, err)
		assert.Empty(t, sellers)
	})

	t.Run("Row Scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
			AddRow(1, "Company 1", "Address 1", 1342, "1234567890")
		mock.ExpectQuery("SELECT `s.id`, `s.cid`, `s.company_name`, `s.address`, `s.telephone` FROM `sellers` AS `s`").WillReturnRows(rows)

		r := NewSellerMysql(db)
		sellers, err := r.FindAll()

		assert.Error(t, err)
		assert.Empty(t, sellers)
	})

}

func TestSellerMysql_FindByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
			AddRow(1, 123, "Company 1", "Address 1", "1234567890")
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers`  WHERE `id` = ?").WithArgs(1).WillReturnRows(row)

		r := NewSellerMysql(db)
		seller, err := r.FindByID(1)

		assert.NoError(t, err)
		assert.Equal(t, 1, seller.ID)
	})

	t.Run("Seller not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers`  WHERE `id` = ?").WithArgs(1).WillReturnError(sql.ErrNoRows)

		r := NewSellerMysql(db)
		seller, err := r.FindByID(1)

		assert.ErrorIs(t, err, internal.ErrSellerNotFound)
		assert.Equal(t, internal.Seller{}, seller)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers`  WHERE `id` = ?").WithArgs(1).WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		seller, err := r.FindByID(1)

		assert.Error(t, err)
		assert.Equal(t, internal.Seller{}, seller)
	})
}

func TestSellerMysql_FindByCID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		row := sqlmock.NewRows([]string{"id", "cid", "company_name", "address", "telephone"}).
			AddRow(1, 123, "Company 1", "Address 1", "1234567890")
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers` WHERE `cid` = ?").WithArgs(123).WillReturnRows(row)

		r := NewSellerMysql(db)
		seller, err := r.FindByCID(123)

		assert.NoError(t, err)
		assert.Equal(t, 123, seller.CID)
	})

	t.Run("Seller not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers` WHERE `cid` = ?").WithArgs(123).WillReturnError(sql.ErrNoRows)

		r := NewSellerMysql(db)
		seller, err := r.FindByCID(123)

		assert.ErrorIs(t, err, internal.ErrSellerNotFound)
		assert.Equal(t, internal.Seller{}, seller)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT `id`, `cid`, `company_name`, `address`, `telephone` FROM `sellers` WHERE `cid` = ?").WithArgs(123).WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		seller, err := r.FindByCID(123)

		assert.Error(t, err)
		assert.Equal(t, internal.Seller{}, seller)
	})
}

func TestSellerMysql_Save(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		seller := &internal.Seller{
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := NewSellerMysql(db)
		err := r.Save(seller)

		assert.NoError(t, err)
		assert.Equal(t, 1, seller.ID)
	})

	t.Run("Seller conflict", func(t *testing.T) {
		seller := &internal.Seller{
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
			WillReturnError(&mysql.MySQLError{Number: 1062})

		r := NewSellerMysql(db)
		err := r.Save(seller)

		assert.ErrorIs(t, err, internal.ErrSellerConflict)
	})

	t.Run("Database error", func(t *testing.T) {
		seller := &internal.Seller{
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("INSERT INTO `sellers` (`cid`, `company_name`, `address`, `telephone`, `locality_id`) VALUES (?, ?, ?, ?, ?)").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality).
			WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		err := r.Save(seller)

		assert.Error(t, err)
	})
}

func TestSellerMysql_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality, seller.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := NewSellerMysql(db)
		err := r.Update(seller)

		assert.NoError(t, err)
	})

	t.Run("Seller not found", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality, seller.ID).
			WillReturnError(&mysql.MySQLError{Number: 1000})

		r := NewSellerMysql(db)
		err := r.Update(seller)

		assert.ErrorIs(t, err, internal.ErrSellerNotFound)
	})

	t.Run("Seller conflict", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality, seller.ID).
			WillReturnError(&mysql.MySQLError{Number: 1062})

		r := NewSellerMysql(db)
		err := r.Update(seller)

		assert.ErrorIs(t, err, internal.ErrSellerConflict)
	})

	t.Run("Database error", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         123,
			CompanyName: "Company 1",
			Address:     "Address 1",
			Telephone:   "1234567890",
			Locality:    1,
		}

		mock.ExpectExec("UPDATE `sellers` SET `cid` = ?, `company_name` = ?, `address` = ?, `telephone` = ?, `locality_id` = ? WHERE `id` = ?").
			WithArgs(seller.CID, seller.CompanyName, seller.Address, seller.Telephone, seller.Locality, seller.ID).
			WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		err := r.Update(seller)

		assert.Error(t, err)
	})
}

func TestSellerMysql_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM `sellers` WHERE `id` = ?").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		r := NewSellerMysql(db)
		err := r.Delete(1)

		assert.NoError(t, err)
	})

	t.Run("Database error", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM `sellers` WHERE `id` = ?").
			WithArgs(1).
			WillReturnError(errors.New("database error"))

		r := NewSellerMysql(db)
		err := r.Delete(1)

		assert.Error(t, err)
	})
}
