package repository

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

// NewLocalityMysql creates a new instance of the seller repository
func NewLocalityMysql(db *sql.DB) *LocalityMysql {
	return &LocalityMysql{db}
}

// LocalityMysql is the mysql implementation of the seller repository
type LocalityMysql struct {
	// db is the database connection to mysql
	db *sql.DB
}

// Save saves a locality into the database
func (r *LocalityMysql) Save(locality *internal.Locality) (err error) {
	// execute the query
	_, err = r.db.Exec(
		"INSERT INTO `localities` (`id`, `name`, `province_name`, `country_name`) VALUES (?, ?, ?, ?)",
		(*locality).ID, (*locality).LocalityName, (*locality).ProvinceName, (*locality).CountryName,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrLocalityConflict
			default:
				// ...
			}
			return
		}

		return
	}

	return
}

// ReportSellers returns a seller from the database by its id
func (r *LocalityMysql) ReportSellers(id int) (locality internal.Locality, err error) {
	// execute the query
	row := r.db.QueryRow("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id", id)

	// scan the row into the seller
	err = row.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceName, &locality.CountryName, &locality.Sellers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrLocalityNotFound
			return
		}
		return
	}

	return
}

func (r *LocalityMysql) FindByID(id int) (locality internal.Locality, err error) {
	// execute the query
	row := r.db.QueryRow("SELECT `id`, `name`, `province_name`, `country_name` FROM `localities` WHERE `id` = ?", id)

	// scan the row into the seller
	err = row.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceName, &locality.CountryName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrLocalityNotFound
			return
		}
		return
	}

	return
}
