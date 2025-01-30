package repository

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const (
	GetAllCarriesQuery = "SELECT * FROM carries"
)

var (
	ErrCidAlreadyExists = errors.New("carry with this cid already exists")
	ErrNoSuchLocalityID = errors.New("there's no such locality id")
)

type CarriesMysql struct {
	db *sql.DB
}

func NewCarriesMysql(db *sql.DB) *CarriesMysql {
	return &CarriesMysql{db}
}

func (r *CarriesMysql) FindAll() (carries []internal.Carries, e error) {
	rows, e := r.db.Query(GetAllCarriesQuery)
	if e != nil {
		return
	}

	for rows.Next() {
		var carry internal.Carries

		rows.Scan(
			&carry.ID,
			&carry.Cid,
			&carry.CompanyName,
			&carry.Address,
			&carry.PhoneNumber,
			&carry.LocalityID,
		)

		carries = append(carries, carry)
	}

	e = rows.Err()
	return
}

func (r *CarriesMysql) Create(carry internal.Carries) (lastID int64, e error) {
	res, e := r.db.Exec(
		"INSERT INTO carries (`cid`, `company_name`, `address`, `phone_number`, `locality_id`) VALUES (?, ?, ?, ?, ?)",
		carry.Cid, carry.CompanyName, carry.Address, carry.PhoneNumber, carry.LocalityID,
	)
	if e != nil {
		mysqlErr, ok := e.(*mysql.MySQLError)
		if ok {
			switch mysqlErr.Number {
			case 1062:
				e = ErrCidAlreadyExists
			case 1452:
				e = ErrNoSuchLocalityID
			}
		}

		return
	}

	lastID, e = res.LastInsertId()

	return
}
