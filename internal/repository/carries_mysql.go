package repository

import (
	"database/sql"

	"github.com/meli-fresh-products-api-backend-t1/internal"
)

const (
	GetAllCarriesQuery = "SELECT * FROM carries"
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
		e = rows.Scan(
			&carry.Id,
			&carry.Cid,
			&carry.CompanyName,
			&carry.Address,
			&carry.PhoneNumber,
			&carry.LocalityId,
		)
		if e != nil {
			return
		}

		carries = append(carries, carry)
	}

	return
}

func (r *CarriesMysql) Create(carry internal.Carries) (lastId int64, e error) {
	res, e := r.db.Exec(
		"INSERT INTO carries (`cid`, `company_name`, `address`, `phone_number`, `locality_id`) VALUES (?, ?, ?, ?, ?)",
		carry.Cid, carry.CompanyName, carry.Address, carry.PhoneNumber, carry.LocalityId,
	)
	if e != nil {
		return
	}

	lastId, e = res.LastInsertId()
	return
}
