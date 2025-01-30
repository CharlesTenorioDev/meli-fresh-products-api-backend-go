package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const (
	AmountOfCarriesForEveryLocalityQuery = `
	SELECT COUNT(c.locality_id) carries_count, c.locality_id, l.name locality_name
	FROM carries c
	INNER JOIN localities l
	ON l.id = c.locality_id
	GROUP BY c.locality_id;
	`
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

func (r *LocalityMysql) ReportCarries(localityID int) (amountOfCarries int, e error) {
	row := r.db.QueryRow(
		"SELECT COUNT(c.locality_id) carries_registered FROM carries c WHERE locality_id = ?",
		localityID,
	)

	e = row.Scan(&amountOfCarries)

	if amountOfCarries == 0 || e != nil {
		e = sql.ErrNoRows
	}

	return
}

func (r *LocalityMysql) GetAmountOfCarriesForEveryLocality() (c []internal.CarriesCountPerLocality, e error) {
	rows, e := r.db.Query(AmountOfCarriesForEveryLocalityQuery)
	if e != nil {
		return
	}

	for rows.Next() {
		var carryCountPerLocality internal.CarriesCountPerLocality

		e = rows.Scan(
			&carryCountPerLocality.CarriesCount,
			&carryCountPerLocality.LocalityID,
			&carryCountPerLocality.LocalityName,
		)
		if e != nil {
			return
		}

		c = append(c, carryCountPerLocality)
	}

	return
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
			}
		}
	}

	return
}

func (r *LocalityMysql) ReportSellers() (localities []internal.Locality, err error) {
	rows, err := r.db.Query("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id GROUP BY l.id")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrLocalityNotFound
		}

		return
	}

	for rows.Next() {
		var locality internal.Locality

		err = rows.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceName, &locality.CountryName, &locality.Sellers)
		if err != nil {
			return
		}

		localities = append(localities, locality)
	}

	return
}

// ReportSellersByID returns a seller from the database by its id
func (r *LocalityMysql) ReportSellersByID(id int) (localities []internal.Locality, err error) {
	// execute the query
	row := r.db.QueryRow("SELECT l.id, l.name, l.province_name, l.country_name, COUNT(s.id) FROM localities AS l LEFT JOIN sellers AS s ON l.id = s.locality_id WHERE l.id = ? GROUP BY l.id", id)

	var locality internal.Locality
	// scan the row into the seller
	err = row.Scan(&locality.ID, &locality.LocalityName, &locality.ProvinceName, &locality.CountryName, &locality.Sellers)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = internal.ErrLocalityNotFound
		}
		return
	}

	localities = append(localities, locality)

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
		}
	}

	return
}
