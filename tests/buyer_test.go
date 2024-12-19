package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const Api = "/api/v1/buyers"

type BuyerRouterSuite struct {
	hd *handler.BuyerHandlerDefault
	rt chi.Router
	suite.Suite
}

func (s *BuyerRouterSuite) SetupTest() {
	repo := repository.NewBuyerMap("../db/buyer.json")
	svc := service.NewBuyerService(repo)
	s.hd = handler.NewBuyerHandlerDefault(svc)
	s.rt = chi.NewRouter()
	s.rt.Route("/api/v1/buyers", func(r chi.Router) {
		r.Get("/", s.hd.GetAll)
	})
}

func (s *BuyerRouterSuite) TestGetAllBuyers() {
	r, err := http.NewRequest(http.MethodGet, Api, nil)
	require.NoError(s.T(), err)
	w := httptest.NewRecorder()

	var buyers struct {
		Data map[string]internal.Buyer `json:"data"`
	}
	s.rt.ServeHTTP(w, r)
	json.NewDecoder(w.Body).Decode(&buyers)

	ok := s.Run("the router replied with status ok", func() {
		require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("there are 5 entries registered", func() {
		require.Equal(s.T(), 5, len(buyers.Data))
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the first entry makes sense", func() {
		require.Equal(s.T(), buyers.Data["0"].ID, 0)
		require.Equal(s.T(), buyers.Data["0"].FirstName, "John")
		require.Equal(s.T(), buyers.Data["0"].LastName, "Doe")
		require.Equal(s.T(), buyers.Data["0"].CardNumberId, "1234567812345678")
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the endpoint replied with content-type application/json", func() {
		require.Equal(s.T(), w.HeaderMap["Content-Type"], []string{"application/json"})
	})
	if !ok {
		s.T().FailNow()
	}
}

func TestBuyerRouterSuite(t *testing.T) {
	suite.Run(t, new(BuyerRouterSuite))
}
