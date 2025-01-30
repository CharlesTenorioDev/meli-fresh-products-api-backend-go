package service_test

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CarriesRepositoryMock struct {
	mock.Mock
}

type CarriesServiceTestSuite struct {
	rp *CarriesRepositoryMock
	sv *service.CarriesService
	suite.Suite
}

func NewCarriesRepositoryMock() *CarriesRepositoryMock {
	return &CarriesRepositoryMock{}
}

func (m *CarriesRepositoryMock) FindAll() ([]internal.Carries, error) {
	args := m.Called()
	return args.Get(0).([]internal.Carries), args.Error(1)
}

func (m *CarriesRepositoryMock) Create(carry internal.Carries) (lastID int64, e error) {
	args := m.Called(carry)
	return args.Get(0).(int64), args.Error(1)
}

func (s *CarriesServiceTestSuite) SetupTest() {
	s.rp = NewCarriesRepositoryMock()
	s.sv = service.NewCarriesService(s.rp)
}

func (s *CarriesServiceTestSuite) TestFindAll() {
	s.T().Run("success", func(t *testing.T) {
		expectedCarries := []internal.Carries{
			{
				ID:          0,
				Cid:         "CID000",
				CompanyName: "Meli",
				Address:     "OneTwoThree",
				PhoneNumber: "1122211",
				LocalityID:  0,
			},
			{
				ID:          1,
				Cid:         "CID001",
				CompanyName: "Go Meli Go",
				Address:     "FourFiveSix",
				PhoneNumber: "1122211",
				LocalityID:  1,
			},
		}
		s.SetupTest()
		s.rp.On("FindAll").Return(expectedCarries, nil)

		actualCarries, e := s.sv.FindAll()

		require.NoError(t, e)
		require.Equal(t, expectedCarries, actualCarries)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.SetupTest()
		s.rp.On("FindAll").Return([]internal.Carries{}, sql.ErrNoRows)

		actualCarries, e := s.sv.FindAll()

		require.Error(t, e)
		require.ErrorIs(t, sql.ErrNoRows, e)
		require.Zero(t, len(actualCarries))
	})
}

func (s *CarriesServiceTestSuite) TestCreate() {
	s.T().Run("success", func(t *testing.T) {
		expectedId := 0
		carry := internal.Carries{
			Cid:         "CID000",
			CompanyName: "Meli",
			Address:     "OneTwoThree",
			PhoneNumber: "119218912",
			LocalityID:  1,
		}
		s.SetupTest()
		s.rp.On("Create", carry).Return(int64(0), nil)

		lastID, e := s.sv.Create(carry)
		require.NoError(t, e)
		require.EqualValues(t, expectedId, lastID)
	})
	s.T().Run("failure", func(t *testing.T) {
		expectedId := -1
		carry := internal.Carries{
			Cid:         "CID000",
			CompanyName: "Meli",
			Address:     "OneTwoThree",
			PhoneNumber: "119218912",
			LocalityID:  1,
		}
		s.SetupTest()
		s.rp.On("Create", carry).Return(int64(-1), errors.New("internal server error"))

		lastID, e := s.sv.Create(carry)
		require.Error(t, e)
		require.Equal(t, "internal server error", e.Error())
		require.EqualValues(t, expectedId, lastID)
	})
}

func TestCarriesServiceTestSuite(t *testing.T) {
	suite.Run(t, new(CarriesServiceTestSuite))
}
