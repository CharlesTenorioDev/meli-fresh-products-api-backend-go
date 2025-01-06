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
