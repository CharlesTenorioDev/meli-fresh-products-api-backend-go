package handler_test

import (
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

const (
	name = "txdb_localities"
)

type LocalityTestSuite struct {
	db *sql.DB
	rt chi.Router
	hd *handler.LocalityDefault
	suite.Suite
}

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh",
	}
	txdb.Register(name, "mysql", cfg.FormatDSN())
}

func (l *LocalityTestSuite) SetupTest() {
	var err error
	l.db, err = sql.Open(name, "")
	require.NoError(l.T(), err)
	rp := repository.NewLocalityMysql(l.db)
	sv := service.NewLocalityDefault(rp)
	l.hd = handler.NewLocalityDefault(sv)
}

func (l *LocalityTestSuite) TestLocalityDefault_ReportCarries() {
	defer l.db.Close()
	l.T().Run("carries with locality id 3", func(t *testing.T) {
		expectedAmountOfCarries := 1
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=3", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		var res struct {
			Data struct {
				AmountOfCarries int `json:"amount_of_carries"`
			} `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&res)
		require.Equal(t, expectedAmountOfCarries, res.Data.AmountOfCarries)
	})
	l.T().Run("carries without passing locality", func(t *testing.T) {
		expectedCarriesCountPerLocality := []internal.CarriesCountPerLocality{
			{
				CarriesCount: 1,
				LocalityId:   1,
				LocalityName: "New York City",
			},
			{
				CarriesCount: 1,
				LocalityId:   2,
				LocalityName: "Los Angeles",
			},
			{
				CarriesCount: 1,
				LocalityId:   3,
				LocalityName: "Chicago",
			},
			{
				CarriesCount: 1,
				LocalityId:   4,
				LocalityName: "Houston",
			},
			{
				CarriesCount: 1,
				LocalityId:   5,
				LocalityName: "Phoenix",
			},
		}
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		var res struct {
			Data []internal.CarriesCountPerLocality `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&res)
		require.Equal(t, expectedCarriesCountPerLocality, res.Data)
		require.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
	l.T().Run("carries with a location that does not exist", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api/v1/localities/reportCarries?id=100", nil)
		w := httptest.NewRecorder()
		l.hd.ReportCarries()(w, r)
		require.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestLocalityTestSuite(t *testing.T) {
	suite.Run(t, new(LocalityTestSuite))
}
