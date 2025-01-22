package handler

import (
	"bytes"
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

func (suite *BuyerTestSuite) TestCreateBuyer() {
	// Insert the object
	bc := internal.Buyer{
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		CardNumberID: "80028922",
	}
	b, _ := json.Marshal(bc)
	r := httptest.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
	w := httptest.NewRecorder()
	suite.hd.Create(w, r)
	assert.Equal(suite.T(), http.StatusCreated, w.Result().StatusCode)

	// Test if it was actually inserted
	r = httptest.NewRequest(http.MethodGet, Api+"/5", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "5")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w = httptest.NewRecorder()
	suite.hd.GetByID(w, r)

	var buyerCreated struct {
		Data internal.Buyer `json:"data"`
	}
	json.NewDecoder(w.Body).Decode(&buyerCreated)
	require.Equal(suite.T(), 5, buyerCreated.Data.ID)
	require.Equal(suite.T(), bc.FirstName, buyerCreated.Data.FirstName)
	require.Equal(suite.T(), bc.LastName, buyerCreated.Data.LastName)
	require.Equal(suite.T(), bc.CardNumberID, buyerCreated.Data.CardNumberID)
}

func (suite *BuyerTestSuite) TestCreateBuyerWithMissingParameters() {
	bc := internal.Buyer{}
	b, _ := json.Marshal(bc)
	r := httptest.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
	w := httptest.NewRecorder()
	suite.hd.Create(w, r)
	assert.Equal(suite.T(), http.StatusUnprocessableEntity, w.Result().StatusCode)
}

func (suite *BuyerTestSuite) TestGetIdThatDoesntExist() {
	r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "200")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	suite.hd.GetByID(w, r)
	var buyers struct {
		Data internal.Buyer `json:"data"`
	}
	err := json.NewDecoder(w.Body).Decode(&buyers)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
}

func (suite *BuyerTestSuite) TestPatchBuyer() {
	fname := "Doe"
	lname := "John"
	cardNumberId := "404"
	suite.Run("apply changes", func() {
		bp := internal.BuyerPatch{
			FirstName:    &fname,
			LastName:     &lname,
			CardNumberID: &cardNumberId,
		}
		b, _ := json.Marshal(bp)
		r := httptest.NewRequest(http.MethodPatch, Api+"/{id}", bytes.NewReader(b))
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "0")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.Update(w, r)
		require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)

		var updatedBuyer struct {
			Data internal.Buyer `json:"data"`
		}
		json.NewDecoder(w.Body).Decode(&updatedBuyer)
		require.Equal(suite.T(), fname, updatedBuyer.Data.FirstName)
		require.Equal(suite.T(), lname, updatedBuyer.Data.LastName)
		require.Equal(suite.T(), cardNumberId, updatedBuyer.Data.CardNumberID)
	})
	suite.Run("check if changes were applied", func() {
		r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "0")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.GetByID(w, r)
		var buyers struct {
			Data internal.Buyer `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&buyers)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), http.StatusOK, w.Result().StatusCode)
		require.Equal(suite.T(), fname, buyers.Data.FirstName)
		require.Equal(suite.T(), lname, buyers.Data.LastName)
		require.Equal(suite.T(), cardNumberId, buyers.Data.CardNumberID)
	})
}

func (suite *BuyerTestSuite) TestPatchInvalidId() {
	fname := "Doe"
	lname := "John"
	cardNumberId := "404"
	bp := internal.BuyerPatch{
		FirstName:    &fname,
		LastName:     &lname,
		CardNumberID: &cardNumberId,
	}
	b, _ := json.Marshal(bp)
	r := httptest.NewRequest(http.MethodPatch, Api+"/{id}", bytes.NewReader(b))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "200")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	suite.hd.Update(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
}

func (suite *BuyerTestSuite) TestDeleteBuyer() {
	suite.Run("delete a buyer", func() {
		r := httptest.NewRequest(http.MethodDelete, Api+"/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "4")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.Delete(w, r)
		require.Equal(suite.T(), http.StatusNoContent, w.Result().StatusCode)
	})
	suite.Run("check if the buyer was actually deleted", func() {
		r := httptest.NewRequest(http.MethodGet, Api+"/{id}", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", "4")
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
		w := httptest.NewRecorder()
		suite.hd.GetByID(w, r)
		var buyers struct {
			Data internal.Buyer `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&buyers)
		require.NoError(suite.T(), err)
		require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
	})
}

func (suite *BuyerTestSuite) TestDeleteBuyerThatDoesntExist() {
	r := httptest.NewRequest(http.MethodDelete, Api+"/{id}", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "200")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()
	suite.hd.Delete(w, r)
	require.Equal(suite.T(), http.StatusNotFound, w.Result().StatusCode)
}

func TestBuyerTestSuite(t *testing.T) {
	suite.Run(t, new(BuyerTestSuite))
}
