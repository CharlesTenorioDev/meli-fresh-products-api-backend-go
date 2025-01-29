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

func NewBuyerRepositoryMock() *BuyerRepositoryMock {
	return &BuyerRepositoryMock{}
}

type BuyerRepositoryMock struct {
	mock.Mock
}

func (rm *BuyerRepositoryMock) GetAll() (db map[int]internal.Buyer, err error) {
	args := rm.Called()
	return args.Get(0).(map[int]internal.Buyer), args.Error(1)
}

func (rm *BuyerRepositoryMock) Add(buyer *internal.Buyer) (id int64, err error) {
	args := rm.Called(buyer)
	return args.Get(0).(int64), args.Error(1)
}

func (rm *BuyerRepositoryMock) Update(id int, buyer internal.BuyerPatch) (err error) {
	args := rm.Called(id, buyer)
	return args.Error(0)
}

func (rm *BuyerRepositoryMock) Delete(id int) (rowsAffected int64, err error) {
	args := rm.Called(id)
	return args.Get(0).(int64), args.Error(1)
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

func (b *BuyerServiceTestSuite) TestBuyerService_Save() {

	b.T().Run("case 1 - Successfully saved a new buyer", func(t *testing.T) {
		b.SetupTest()

		buyer := internal.Buyer{
			ID:           1,
			CardNumberID: "3445342",
			FirstName:    "Paloma",
			LastName:     "Souza",
		}
		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		b.rp.On("Add", &buyer).Return(int64(1), nil)

		err := b.sv.Save(&buyer)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Add", 1)
		b.NoError(err)
	})

	b.T().Run("case 2 - Return an conflit error when trying to save a buyer with card number id already in use", func(t *testing.T) {
		b.SetupTest()

		buyer := internal.Buyer{
			ID:           1,
			CardNumberID: "3445342",
			FirstName:    "Brian",
			LastName:     "May",
		}
		b.rp.On("GetAll").Return(map[int]internal.Buyer{1: {ID: 1, CardNumberID: "3445342", FirstName: "Paloma", LastName: "Souza"}}, nil)
		err := b.sv.Save(&buyer)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Add", 0)
		b.Equal(service.ErrCardNumberAlreadyInUse, err)

	})

	b.T().Run("case 3 - Return an unprocessable entity error when trying to parse a buyer to save", func(t *testing.T) {
		b.SetupTest()

		buyer := internal.Buyer{
			ID:           1,
			CardNumberID: "",
			FirstName:    "Pah",
			LastName:     "Gabi",
		}
		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		err := b.sv.Save(&buyer)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Add", 0)
		b.Equal(err, service.ErrBuyerUnprocessableEntity)
	})

	b.T().Run("case 3 - Should return an error when trying to save a buyer with invalid data", func(t *testing.T) {
		b.SetupTest()

		buyer := internal.Buyer{
			ID:           1,
			CardNumberID: "6544666",
			FirstName:    "Jack",
		}
		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		err := b.sv.Save(&buyer)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Add", 0)
		b.Equal(err, errors.New("couldn't parse buyer"))
	})

}

func (b *BuyerServiceTestSuite) TestBuyerService_Read() {

	b.T().Run("case 1 - Successfully to get all buyers", func(t *testing.T) {
		b.SetupTest()

		buyer := map[int]internal.Buyer{

			1: {ID: 1, CardNumberID: "3445342", FirstName: "Paloma", LastName: "Souza"},
			2: {ID: 2, CardNumberID: "3445343", FirstName: "Brian", LastName: "May"},
			3: {ID: 3, CardNumberID: "3445344", FirstName: "Pah", LastName: "Gabi"},
		}

		b.rp.On("GetAll").Return(buyer, nil)

		buyers := b.sv.GetAll()
		require.Equal(b.T(), buyer, buyers)
	})

	b.T().Run("case 2 - Return an error when trying to get a buyer by a non existent id", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		_, err := b.sv.FindByID(500)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		require.Error(b.T(), err)
		b.Equal(err, service.ErrBuyerNotFound)
	})

	b.T().Run("case 3 - Successfully to get a buyer by id", func(t *testing.T) {
		b.SetupTest()

		buyer := map[int]internal.Buyer{

			1: {ID: 1, CardNumberID: "3445342", FirstName: "Paloma", LastName: "Souza"},
		}

		b.rp.On("GetAll").Return(buyer, nil)

		result, err := b.sv.FindByID(1)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		require.NoError(b.T(), err)
		require.Equal(b.T(), buyer[1], result)
	})

	b.T().Run("case 4 - Return an error when trying to get a buyer by an invalid id", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		_, err := b.sv.FindByID(-20)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		require.Error(b.T(), err)
		b.Equal(err, errors.New("buyer not found"))
	})

}

func (b *BuyerServiceTestSuite) TestBuyerService_Update() {

	b.T().Run("case 1 - Successfully to get an existent buyer, returns updated data", func(t *testing.T) {
		b.SetupTest()

		card := "1111111"
		firstName := "Paloma"
		lastName := "S."

		buyerPatch := internal.BuyerPatch{
			CardNumberID: &card,
			FirstName:    &firstName,
			LastName:     &lastName,
		}

		buyer := internal.Buyer{
			ID:           1,
			CardNumberID: "3445342",
			FirstName:    "Paloma",
			LastName:     "Souza",
		}

		buyerUpdated := internal.Buyer{
			ID:           1,
			CardNumberID: "1111111",
			FirstName:    "Paloma",
			LastName:     "S.",
		}

		b.rp.On("GetAll").Return(map[int]internal.Buyer{
			1: buyer,
		}, nil)

		b.rp.On("Update", mock.AnythingOfType("int"), mock.AnythingOfType("internal.BuyerPatch")).Run(func(args mock.Arguments) {
			id := args.Get(0).(int)
			patch := args.Get(1).(internal.BuyerPatch)
			require.Equal(t, 1, id)
			patch.Patch(&buyer)
		}).Return(nil)

		err := b.sv.Update(1, buyerPatch)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Update", 1)
		require.NoError(b.T(), err)
		require.Equal(b.T(), buyerUpdated, buyer)
	})

	b.T().Run("case 2 - Returns not found error when trying to get non existent buyer", func(t *testing.T) {
		b.SetupTest()

		buyerPatch := internal.BuyerPatch{}
		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)
		err := b.sv.Update(55, buyerPatch)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Update", 0)
		b.Equal(err, service.ErrBuyerNotFound)
	})

	b.T().Run("case 3 - Return conflit error when trying to update with a card number id already in use", func(t *testing.T) {
		b.SetupTest()

		card := "1111111"
		firstName := "Brian"
		lastName := "May"

		buyerPatch := internal.BuyerPatch{
			CardNumberID: &card,
			FirstName:    &firstName,
			LastName:     &lastName,
		}

		buyerWithCardNumberInUse := internal.Buyer{
			ID:           1,
			CardNumberID: "1111111",
			FirstName:    "Paloma",
			LastName:     "Souza",
		}

		b.rp.On("GetAll").Return(map[int]internal.Buyer{
			2: buyerWithCardNumberInUse,
		}, nil)

		err := b.sv.Update(2, buyerPatch)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Update", 0)
		b.Equal(err, service.ErrCardNumberAlreadyInUse)
	})

}

func (b *BuyerServiceTestSuite) TestBuyerService_Delete() {

	b.T().Run("case 1 - Successfully to delete an existent buyer", func(t *testing.T) {
		b.SetupTest()

		buyer := map[int]internal.Buyer{

			1: {ID: 1, CardNumberID: "111111", FirstName: "Paloma", LastName: "Souza"},
		}

		b.rp.On("GetAll").Return(buyer, nil)
		b.rp.On("Delete", 1).Return(int64(0), nil)

		err := b.sv.Delete(1)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Delete", 1)
		require.NoError(b.T(), err)
	})

	b.T().Run("case 2 - Returns not found error when trying to delete non existent buyer", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)

		err := b.sv.Delete(55)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "Delete", 0)
		b.Equal(err, service.ErrBuyerNotFound)
	})
}

func (b *BuyerServiceTestSuite) TestBuyerService_ReportPurchaseOrders() {

	b.T().Run("case 1 - Successfully to get all report purchase orders by buyer", func(t *testing.T) {
		b.SetupTest()

		purchaseOrders := []internal.PurchaseOrdersByBuyer{
			{BuyerID: 1, CardNumberID: "111111", FirstName: "Paloma", LastName: "Souza", PurchaseOrdersCount: 3},
			{BuyerID: 2, CardNumberID: "222222", FirstName: "Pah", LastName: "Gabi", PurchaseOrdersCount: 1},
			{BuyerID: 3, CardNumberID: "333333", FirstName: "Brian", LastName: "May", PurchaseOrdersCount: 5},
		}

		b.rp.On("ReportPurchaseOrders").Return(purchaseOrders, nil)

		allPurchaseOrders, err := b.sv.ReportPurchaseOrders()
		require.Equal(b.T(), purchaseOrders, allPurchaseOrders)
		require.NoError(b.T(), err)

	})

	b.T().Run("case 2 - Return purchase orders not found error when trying to get no buyers records", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("ReportPurchaseOrders").Return([]internal.PurchaseOrdersByBuyer{}, nil)

		_, err := b.sv.ReportPurchaseOrders()

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "ReportPurchaseOrders", 1)
		b.Equal(err, service.ErrPurchaseOrdersNotFound)
	})

	b.T().Run("case 3 - Successfully to get report purchase orders by buyer by Id parameter", func(t *testing.T) {
		b.SetupTest()

		purchaseOrdersByBuyer := []internal.PurchaseOrdersByBuyer{
			{BuyerID: 2, CardNumberID: "222222", FirstName: "Pah", LastName: "Gabi", PurchaseOrdersCount: 1},
		}

		b.rp.On("GetAll").Return(map[int]internal.Buyer{
			1: {ID: 1, CardNumberID: "111111", FirstName: "Paloma", LastName: "Souza"},
			2: {ID: 2, CardNumberID: "222222", FirstName: "Pah", LastName: "Gabi"},
		}, nil)

		b.rp.On("ReportPurchaseOrdersByID", 2).Return(purchaseOrdersByBuyer, nil)

		resultById, err := b.sv.ReportPurchaseOrdersByID(2)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "GetAll", 1)
		b.rp.AssertNumberOfCalls(b.T(), "ReportPurchaseOrdersByID", 1)
		require.NoError(b.T(), err)
		require.Equal(b.T(), purchaseOrdersByBuyer, resultById)

	})

	b.T().Run("case 4 - Return buyer not found error when trying to get an buyer not existent by Id parameter", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("GetAll").Return(map[int]internal.Buyer{}, nil)

		_, err := b.sv.ReportPurchaseOrdersByID(255)

		b.rp.AssertExpectations(b.T())
		b.rp.AssertNumberOfCalls(b.T(), "ReportPurchaseOrdersByID", 0)
		b.Equal(err, service.ErrBuyerNotFound)
	})

	b.T().Run("case 5 - Return an error when trying to get a buyer with no purchase orders", func(t *testing.T) {
		b.SetupTest()

		b.rp.On("GetAll").Return(map[int]internal.Buyer{
			1: {ID: 1, CardNumberID: "111111", FirstName: "Paloma", LastName: "Souza"},
			2: {ID: 2, CardNumberID: "222222", FirstName: "Pah", LastName: "Gabi"},
			3: {ID: 3, CardNumberID: "333333", FirstName: "Brian", LastName: "May"},
			4: {ID: 4, CardNumberID: "444444", FirstName: "Jack", LastName: "Sparrow"},
		}, nil)
		b.rp.On("ReportPurchaseOrdersByID", 4).Return([]internal.PurchaseOrdersByBuyer{}, nil)

		noPurchaseOrders, err := b.sv.ReportPurchaseOrdersByID(4)

		require.Nil(b.T(), noPurchaseOrders)
		require.Error(b.T(), err)
		b.Equal(err, service.ErrPurchaseOrdersByBuyerNotFound)
		b.rp.AssertExpectations(b.T())
	})

}
