package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/handler"
	"github.com/meli-fresh-products-api-backend-t1/internal/repository"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/meli-fresh-products-api-backend-t1/utils/rest_err"
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
		r.Get("/{id}", s.hd.GetByID)
		r.Post("/", s.hd.Create)
		r.Patch("/{id}", s.hd.Update)
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
		require.Equal(s.T(), w.Header()["Content-Type"], []string{"application/json"})
	})
	if !ok {
		s.T().FailNow()
	}
}

func (s *BuyerRouterSuite) TestGetSpecificBuyers() {
	tests := []struct {
		ID                   int
		ExpectedHttpStatus   int
		ExpectedFirstName    string
		ExpectedLastName     string
		ExpectedCardNumberid string
	}{
		{
			ID:                   0,
			ExpectedHttpStatus:   http.StatusOK,
			ExpectedFirstName:    "John",
			ExpectedLastName:     "Doe",
			ExpectedCardNumberid: "1234567812345678",
		},
		{
			ID:                   1,
			ExpectedHttpStatus:   http.StatusOK,
			ExpectedFirstName:    "Jane",
			ExpectedLastName:     "Smith",
			ExpectedCardNumberid: "2345678923456789",
		},
		{
			ID:                   2,
			ExpectedHttpStatus:   http.StatusOK,
			ExpectedFirstName:    "Alice",
			ExpectedLastName:     "Johnson",
			ExpectedCardNumberid: "3456789034567890",
		},
	}

	for _, tt := range tests {
		path := Api + "/" + strconv.Itoa(tt.ID)
		idStr := strconv.Itoa(tt.ID)
		r, _ := http.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		s.rt.ServeHTTP(w, r)

		ok := s.Run("the router replied with status ok on api v1 buyers "+idStr, func() {
			require.Equal(s.T(), tt.ExpectedHttpStatus, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		var buyers struct {
			Data internal.Buyer `json:"data"`
		}
		err := json.NewDecoder(w.Body).Decode(&buyers)
		ok = s.Run("got no error while decoding the response on api v1 buyers "+idStr, func() {
			require.NoError(s.T(), err)
		})
		if !ok {
			s.T().FailNow()
		}

		ok = s.Run("the entries match on api v1 buyers "+idStr, func() {
			require.Equal(s.T(), tt.ID, buyers.Data.ID)
			require.Equal(s.T(), tt.ExpectedFirstName, buyers.Data.FirstName)
			require.Equal(s.T(), tt.ExpectedLastName, buyers.Data.LastName)
			require.Equal(s.T(), tt.ExpectedCardNumberid, buyers.Data.CardNumberId)
		})
		if !ok {
			s.T().FailNow()
		}
	}
}

func (s *BuyerRouterSuite) TestGetInvalidBuyer() {
	r, err := http.NewRequest(http.MethodGet, Api+"/10", nil)
	w := httptest.NewRecorder()
	require.NoError(s.T(), err)

	var restErr rest_err.RestErr
	s.rt.ServeHTTP(w, r)
	err = json.NewDecoder(w.Body).Decode(&restErr)
	ok := s.Run("had no problems while decoding the response", func() {
		require.NoError(s.T(), err)
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the endpoint replied with status not found", func() {
		require.Equal(s.T(), http.StatusNotFound, w.Result().StatusCode)
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the rest err returned is a 'not found'", func() {
		expectedError := *rest_err.NewNotFoundError("buyer not found")
		require.Equal(s.T(), expectedError, restErr)
	})
	if !ok {
		s.T().FailNow()
	}
}

func (s *BuyerRouterSuite) TestCreateBuyer() {
	buyer := internal.Buyer{
		FirstName:    "Fabio",
		LastName:     "Nacarelli",
		CardNumberId: "40028922",
	}
	ok := s.Run("the user is created successfully", func() {
		b, _ := json.Marshal(buyer)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
		require.NoError(s.T(), err)
		s.rt.ServeHTTP(w, r)
		ok := s.Run("the endpoint replied with status created", func() {
			require.Equal(s.T(), http.StatusCreated, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		var resBuyer struct {
			Data internal.Buyer `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&resBuyer)
		require.NoError(s.T(), err)
		ok = s.Run("the response fields are as expected", func() {
			require.Equal(s.T(), 5, resBuyer.Data.ID)
			require.Equal(s.T(), buyer.FirstName, resBuyer.Data.FirstName)
			require.Equal(s.T(), buyer.LastName, resBuyer.Data.LastName)
			require.Equal(s.T(), buyer.CardNumberId, resBuyer.Data.CardNumberId)
		})
		if !ok {
			s.T().FailNow()
		}
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the user was actually created", func() {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, Api+"/5", nil)
		require.NoError(s.T(), err)
		s.rt.ServeHTTP(w, r)
		ok := s.Run("the user exists", func() {
			require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		var resBuyer struct {
			Data internal.Buyer
		}
		err = json.NewDecoder(w.Body).Decode(&resBuyer)
		ok = s.Run("the entries match", func() {
			require.Equal(s.T(), buyer.FirstName, resBuyer.Data.FirstName)
			require.Equal(s.T(), buyer.LastName, resBuyer.Data.LastName)
			require.Equal(s.T(), buyer.CardNumberId, resBuyer.Data.CardNumberId)
		})
		if !ok {
			s.T().FailNow()
		}
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("create with the same card number id conflicts", func() {
		b, _ := json.Marshal(buyer)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
		require.NoError(s.T(), err)
		s.rt.ServeHTTP(w, r)
		ok := s.Run("the endpoint replied with status created", func() {
			require.Equal(s.T(), http.StatusConflict, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		ok = s.Run("the rest err is as expected", func() {
			expectedRestErr := *rest_err.NewConflictError("buyer with given card number already registered")
			var restErr rest_err.RestErr
			err = json.NewDecoder(w.Body).Decode(&restErr)
			require.NoError(s.T(), err)
			require.Equal(s.T(), expectedRestErr, restErr)
		})
		if !ok {
			s.T().FailNow()
		}
	})
	if !ok {
		s.T().FailNow()
	}
}

func (s *BuyerRouterSuite) TestCreateUnprocessableEntity() {
	b, _ := json.Marshal(internal.Buyer{
		FirstName:    "",
		LastName:     "",
		CardNumberId: "",
	})
	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodPost, Api, bytes.NewReader(b))
	require.NoError(s.T(), err)
	s.rt.ServeHTTP(w, r)
	ok := s.Run("the endpoint replies with unprocessable entity", func() {
		require.Equal(s.T(), http.StatusUnprocessableEntity, w.Result().StatusCode)
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the rest error is as expected", func() {
		expectedRestErr := *rest_err.NewUnprocessableEntityError("couldn't parse buyer")
		var resRestErr rest_err.RestErr
		json.NewDecoder(w.Body).Decode(&resRestErr)
		require.Equal(s.T(), expectedRestErr, resRestErr)
	})
	if !ok {
		s.T().FailNow()
	}
}

func (s *BuyerRouterSuite) TestPatchBuyer() {
	fname := "Not"
	lname := "Found"
	cardNumberId := "404"
	buyer := internal.BuyerPatch{
		FirstName:    &fname,
		LastName:     &lname,
		CardNumberId: &cardNumberId,
	}
	ok := s.Run("the user is modified successfully", func() {
		b, _ := json.Marshal(buyer)
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodPatch, Api+"/0", bytes.NewReader(b))
		require.NoError(s.T(), err)
		s.rt.ServeHTTP(w, r)
		ok := s.Run("the endpoint replied with status ok", func() {
			require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		var resBuyer struct {
			Data internal.BuyerPatch `json:"data"`
		}
		err = json.NewDecoder(w.Body).Decode(&resBuyer)
		require.NoError(s.T(), err)
		ok = s.Run("the response fields are as expected", func() {
			require.Equal(s.T(), *buyer.FirstName, *resBuyer.Data.FirstName)
			require.Equal(s.T(), *buyer.LastName, *resBuyer.Data.LastName)
			require.Equal(s.T(), *buyer.CardNumberId, *resBuyer.Data.CardNumberId)
		})
		if !ok {
			s.T().FailNow()
		}
	})
	if !ok {
		s.T().FailNow()
	}

	ok = s.Run("the user was actually modified", func() {
		w := httptest.NewRecorder()
		r, err := http.NewRequest(http.MethodGet, Api+"/0", nil)
		require.NoError(s.T(), err)
		s.rt.ServeHTTP(w, r)
		ok := s.Run("the user exists", func() {
			require.Equal(s.T(), http.StatusOK, w.Result().StatusCode)
		})
		if !ok {
			s.T().FailNow()
		}

		var resBuyer struct {
			Data internal.BuyerPatch
		}
		err = json.NewDecoder(w.Body).Decode(&resBuyer)
		ok = s.Run("the entries match", func() {
			require.Equal(s.T(), *buyer.FirstName, *resBuyer.Data.FirstName)
			require.Equal(s.T(), *buyer.LastName, *resBuyer.Data.LastName)
			require.Equal(s.T(), *buyer.CardNumberId, *resBuyer.Data.CardNumberId)
		})
		if !ok {
			s.T().FailNow()
		}
	})
	if !ok {
		s.T().FailNow()
	}

}

func TestBuyerRouterSuite(t *testing.T) {
	suite.Run(t, new(BuyerRouterSuite))
}
