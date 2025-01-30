package service_test

import (
	"errors"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InboundOrdersRepositoryMock struct {
	mock.Mock
}

type InboundOrderServiceTestSuite struct {
	rp *InboundOrdersRepositoryMock
	sv *service.InboundOrderService
	suite.Suite
}

func NewInboundOrdersRepositoryMock() *InboundOrdersRepositoryMock {
	return &InboundOrdersRepositoryMock{}
}

func (m *InboundOrdersRepositoryMock) FindAll() ([]internal.InboundOrders, error) {
	args := m.Called()
	return args.Get(0).([]internal.InboundOrders), args.Error(1)
}

func (m *InboundOrdersRepositoryMock) Create(inboundOrder internal.InboundOrders) (int64, error) {
	args := m.Called(inboundOrder)
	return args.Get(0).(int64), args.Error(1)
}

func (s *InboundOrderServiceTestSuite) SetupTest() {
	s.rp = NewInboundOrdersRepositoryMock()
	s.sv = service.NewInboundOrderService(
		s.rp,
		nil,
		nil,
		nil,
	)
}

func (s *InboundOrderServiceTestSuite) TestFindAll() {
	s.T().Run("success", func(t *testing.T) {
		expectedInboundOrders := []internal.InboundOrders{
			{
				ID:             0,
				OrderDate:      "17/12/2001",
				OrderNumber:    "ON0",
				EmployeeID:     0,
				ProductBatchID: 0,
				WarehouseID:    0,
			},
			{
				ID:             1,
				OrderDate:      "16/12/2001",
				OrderNumber:    "ON1",
				EmployeeID:     1,
				ProductBatchID: 1,
				WarehouseID:    1,
			},
		}
		s.SetupTest()
		s.rp.On("FindAll").Return(expectedInboundOrders, nil)

		actualInboundOrders, e := s.sv.FindAll()

		require.NoError(t, e)
		require.Equal(t, expectedInboundOrders, actualInboundOrders)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.SetupTest()
		s.rp.On("FindAll").Return([]internal.InboundOrders{}, errors.New("internal server error"))

		actualInboundOrders, e := s.sv.FindAll()

		require.Error(t, e)
		require.Equal(t, "internal server error", e.Error())
		require.Zero(t, len(actualInboundOrders))
	})
}

func (s *InboundOrderServiceTestSuite) TestCreate() {
	inboundOrder := internal.InboundOrders{
		ID:             0,
		OrderDate:      "17/12/2001",
		OrderNumber:    "ON00",
		EmployeeID:     0,
		ProductBatchID: 0,
		WarehouseID:    0,
	}
	s.T().Run("success", func(t *testing.T) {
		s.SetupTest()
		s.rp.On("Create", inboundOrder).Return(int64(0), nil)

		lastID, e := s.sv.Create(inboundOrder)

		require.NoError(t, e)
		require.EqualValues(t, 0, lastID)
	})
	s.T().Run("failure", func(t *testing.T) {
		s.SetupTest()
		s.rp.On("Create", inboundOrder).Return(int64(-1), errors.New("internal server error"))

		lastID, e := s.sv.Create(inboundOrder)

		require.Error(t, e)
		require.Equal(t, "internal server error", e.Error())
		require.EqualValues(t, -1, lastID)
	})
}

func TestInboundOrdersServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InboundOrderServiceTestSuite))
}
