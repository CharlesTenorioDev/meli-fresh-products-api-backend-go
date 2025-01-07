package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const api = "/api/v1/carries"

type CarriesTestSuite struct {
	hd *handler.CarriesHandlerDefault
	rt chi.Router
	db *sql.DB
	suite.Suite
}

func (c *CarriesTestSuite) SetupTest() {
	var err error
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh",
	}
	if _, err := sql.Open("txdb_carries", ""); err != nil {
		txdb.Register("txdb_carries", "mysql", cfg.FormatDSN())
	}
	c.db, err = sql.Open("txdb_carries", "")
	require.NoError(c.T(), err)
	rp := repository.NewCarriesMysql(c.db)
	sv := service.NewCarriesService(rp)
	c.hd = handler.NewCarriesHandlerDefault(sv)
}

func (c *CarriesTestSuite) TestCarriesDefault_GetAll() {
	c.T().Run("check if last entry is as expected", func(t *testing.T) {
		defer func(db *sql.DB) {
			db.Exec("DELETE FROM carries")
			db.Exec("ALTER TABLE carries AUTO_INCREMENT = 1")
		}(c.db)
		var data struct {
			Data []internal.Carries `json:"data"`
		}
		res, err := c.db.Exec("INSERT INTO carries (`cid`, `company_name`, `address`, `phone_number`, `locality_id`) VALUES (14, 'Mercado Livre', 'Location', '11977021487', 3)")
		require.NoError(t, err)
		lastId, err := res.LastInsertId()
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodGet, api, nil)
		w := httptest.NewRecorder()
		c.hd.GetAll(w, r)
		json.NewDecoder(w.Result().Body).Decode(&data)
		require.Equal(t, "Location", data.Data[lastId-1].Address)
		require.Equal(t, "11977021487", data.Data[lastId-1].PhoneNumber)
		require.Equal(t, "Mercado Livre", data.Data[lastId-1].CompanyName)
		require.Equal(t, 14, data.Data[lastId-1].Cid)
	})
}

func (c *CarriesTestSuite) TestCarriesDefault_Create() {
	c.T().Run("attempt to create", func(t *testing.T) {
		defer func(db *sql.DB) {
			db.Exec("DELETE FROM carries")
			db.Exec("ALTER TABLE carries AUTO_INCREMENT = 1")
		}(c.db)
		carry := internal.Carries{
			Cid:         15,
			CompanyName: "Iris",
			Address:     "Paulista",
			PhoneNumber: "11977021487",
			LocalityId:  3,
		}
		b, err := json.Marshal(carry)
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodPost, api, bytes.NewReader(b))
		w := httptest.NewRecorder()
		c.hd.Create(w, r)
		var resCarry internal.Carries
		row := c.db.QueryRow("SELECT cid, company_name, address, phone_number, locality_id FROM carries")
		row.Scan(
			&resCarry.Cid,
			&resCarry.CompanyName,
			&resCarry.Address,
			&resCarry.PhoneNumber,
			&resCarry.LocalityId,
		)
		require.Equal(c.T(), carry, resCarry)
	})
}

func TestBuyerTestSuite(t *testing.T) {
	suite.Run(t, new(CarriesTestSuite))
}
