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

func init() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "meli_pass",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "melifresh",
	}
	txdb.Register("txdb", "mysql", cfg.FormatDSN())
}

func (c *CarriesTestSuite) SetupTest() {
	var err error
	c.db, err = sql.Open("txdb", "identier")
	require.NoError(c.T(), err)
	rp := repository.NewCarriesMysql(c.db)
	sv := service.NewCarriesService(rp)
	c.hd = handler.NewCarriesHandlerDefault(sv)
}

func (c *CarriesTestSuite) TestCarriesDefault_GetAll() {
	defer c.db.Close()
	c.T().Run("check if last entry is as expected", func(t *testing.T) {
		var data struct {
			Data []internal.Carries `json:"data"`
		}
		_, err := c.db.Exec("INSERT INTO carries (`cid`, `company_name`, `address`, `phone_number`, `locality_id`) VALUES ('CID#14', 'Mercado Livre', 'Location', '11977021487', 3)")
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodGet, api, nil)
		w := httptest.NewRecorder()
		c.hd.GetAll(w, r)
		json.NewDecoder(w.Result().Body).Decode(&data)
		last := len(data.Data) - 1
		require.Equal(t, "Location", data.Data[last].Address)
		require.Equal(t, "11977021487", data.Data[last].PhoneNumber)
		require.Equal(t, "Mercado Livre", data.Data[last].CompanyName)
		require.Equal(t, "CID#14", data.Data[last].Cid)
	})
}

func (c *CarriesTestSuite) TestCarriesDefault_Create() {
	defer c.db.Close()
	c.T().Run("attempt to create with valid entries", func(t *testing.T) {
		carry := internal.Carries{
			Cid:         "CID#15",
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
		row := c.db.QueryRow("SELECT cid, company_name, address, phone_number, locality_id FROM carries ORDER BY id DESC LIMIT 1")
		row.Scan(
			&resCarry.Cid,
			&resCarry.CompanyName,
			&resCarry.Address,
			&resCarry.PhoneNumber,
			&resCarry.LocalityId,
		)
		var data struct {
			Data struct {
				Id int `json:"id"`
			} `json:"data"`
		}
		json.NewDecoder(w.Result().Body).Decode(&data)
		require.Equal(c.T(), http.StatusCreated, w.Result().StatusCode)
		require.Equal(c.T(), carry, resCarry)
	})
	c.T().Run("attempt to create with invalid entries", func(t *testing.T) {
		carry := internal.Carries{}
		b, err := json.Marshal(carry)
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodPost, api, bytes.NewReader(b))
		w := httptest.NewRecorder()
		c.hd.Create(w, r)
		require.Equal(c.T(), http.StatusUnprocessableEntity, w.Result().StatusCode)
	})
	c.T().Run("attempt to create with repeated entries", func(t *testing.T) {
		carry := internal.Carries{
			Cid:         "CID#15",
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
		require.Equal(c.T(), http.StatusConflict, w.Result().StatusCode)
	})
	c.T().Run("attempt to create with a invalid locality id", func(t *testing.T) {
		carry := internal.Carries{
			Cid:         "CID#100",
			CompanyName: "Iris",
			Address:     "Paulista",
			PhoneNumber: "11977021487",
			LocalityId:  100,
		}
		b, err := json.Marshal(carry)
		require.NoError(t, err)
		r := httptest.NewRequest(http.MethodPost, api, bytes.NewReader(b))
		w := httptest.NewRecorder()
		c.hd.Create(w, r)
		require.Equal(c.T(), http.StatusConflict, w.Result().StatusCode)
	})
}

func TestCarrierTestSuite(t *testing.T) {
	suite.Run(t, new(CarriesTestSuite))
}
