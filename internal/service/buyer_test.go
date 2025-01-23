package service_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func NewBuyerRepositoryMock() *BuyerRepositoryMock {
	return &BuyerRepositoryMock{}
}

type BuyerRepositoryMock struct {
	mock.Mock
}

func (rm *BuyerRepositoryMock) GetAll() map[int]internal.Buyer {
	args := rm.Called()
	return args.Get(0).(map[int]internal.Buyer)
}

func (rm *BuyerRepositoryMock) Add(buyer *internal.Buyer) {
	rm.Called(buyer)
}

func (rm *BuyerRepositoryMock) Update(id int, buyer internal.BuyerPatch) {
	rm.Called(id, buyer)
}

func (rm *BuyerRepositoryMock) Delete(id int) {
	rm.Called(id)
}

func (rm *BuyerRepositoryMock) ReportPurchaseOrders() (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	args := rm.Called()
	return args.Get(0).([]internal.PurchaseOrdersByBuyer), args.Error(1)
}

func (rm *BuyerRepositoryMock) ReportPurchaseOrdersByID(id int) (purchaseOrders []internal.PurchaseOrdersByBuyer, err error) {
	args := rm.Called(id)
	return args.Get(0).([]internal.PurchaseOrdersByBuyer), args.Error(1)
}

func TestBuyerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(BuyerServiceTestSuite))
}

type BuyerServiceTestSuite struct {
	rp *BuyerRepositoryMock
	sv *service.BuyerServiceDefault
	suite.Suite
}

func (s *BuyerServiceTestSuite) SetupTest() {
	rp := NewBuyerRepositoryMock()
	sv := service.NewBuyerService(rp)
	s.rp = rp
	s.sv = sv
}
