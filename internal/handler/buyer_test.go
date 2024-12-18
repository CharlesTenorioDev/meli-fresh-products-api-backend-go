package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	DbPath = "../../db/buyer.json"
	Api    = "/api/v1/buyers"
)

type BuyerTestSuite struct {
	buyersFromFile map[int]internal.Buyer
	hd             *BuyerHandlerDefault
	suite.Suite
}

func (suite *BuyerTestSuite) SetupTest() {
	suite.buyersFromFile = make(map[int]internal.Buyer)
	file, err := os.Open(DbPath)
	if err != nil {
		log.Fatal("Failed to open db file: ", err)
	}

	var buyers []internal.Buyer
	json.NewDecoder(file).Decode(&buyers) // This should never fail, doesn't make sense to test with a invalid json, does it?

	for i, b := range buyers {
		suite.buyersFromFile[i] = b
	}
	rp := repository.NewBuyerMap(DbPath)
	sv := service.NewBuyerService(rp)
	suite.hd = NewBuyerHandlerDefault(sv)
}

func (suite *BuyerTestSuite) TestGetAllBuyers() {
	r := httptest.NewRequest(http.MethodGet, Api, nil)
	w := httptest.NewRecorder()
	suite.hd.GetAll(w, r)
	assert.Equal(suite.T(), 200, w.Result().StatusCode)

	var buyers struct {
		Data map[string]internal.Buyer `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&buyers)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), suite.buyersFromFile[0], buyers.Data["0"])
	require.Equal(suite.T(), suite.buyersFromFile[1], buyers.Data["1"])
	require.Equal(suite.T(), suite.buyersFromFile[2], buyers.Data["2"])
	require.Equal(suite.T(), suite.buyersFromFile[3], buyers.Data["3"])
	require.Equal(suite.T(), suite.buyersFromFile[4], buyers.Data["4"])
}

func (suite *BuyerTestSuite) TestGetBuyersById() {
	suite.Run("get several ids", func() {
		for i := range 5 {
			r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(i))
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()
			suite.hd.GetByID(w, r)
			var buyers struct {
				Data internal.Buyer `json:"data"`
			}
			err := json.NewDecoder(w.Body).Decode(&buyers)
			require.NoError(suite.T(), err)
			require.Equal(suite.T(), suite.buyersFromFile[i], buyers.Data)
		}
	})
}

func TestBuyerTestSuite(t *testing.T) {
	suite.Run(t, new(BuyerTestSuite))
}
