package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewPurchaseOrderRepositoryMock() *PurchaseOrderRepositoryMock {
	return &PurchaseOrderRepositoryMock{}
}

type PurchaseOrderRepositoryMock struct {
	mock.Mock
}

func (r *PurchaseOrderRepositoryMock) Save(p *internal.PurchaseOrder) error {
	args := r.Called(p)
	return args.Error(0)
}

func (r *PurchaseOrderRepositoryMock) FindByID(id int) (internal.PurchaseOrder, error) {
	args := r.Called(id)
	return args.Get(0).(internal.PurchaseOrder), args.Error(1)
}

func NewBuyerServiceMock() *BuyerServiceMock {
	return &BuyerServiceMock{}
}

type BuyerServiceMock struct {
	mock.Mock
}

func (bm *BuyerServiceMock) GetAll() map[int]internal.Buyer {
	args := bm.Called()
	return args.Get(0).(map[int]internal.Buyer)
}

func (bm *BuyerServiceMock) FindByID(id int) (internal.Buyer, error) {
	args := bm.Called(id)
	return args.Get(0).(internal.Buyer), args.Error(1)
}

func (bm *BuyerServiceMock) Save(newBuyer *internal.Buyer) (err error) {
	return nil
}

func (bm *BuyerServiceMock) Update(id int, buyerPatch internal.BuyerPatch) (err error) {
	return nil
}

func (bm *BuyerServiceMock) Delete(id int) (err error) {
	return nil
}

func (bm *BuyerServiceMock) ReportPurchaseOrders() ([]internal.PurchaseOrdersByBuyer, error) {
	return nil, nil
}

func (bm *BuyerServiceMock) ReportPurchaseOrdersByID(id int) ([]internal.PurchaseOrdersByBuyer, error) {
	return nil, nil
}

var (
	now  = time.Now()
	date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	po   = internal.PurchaseOrder{
		ID:              1,
		OrderNumber:     "123ABC",
		OrderDate:       date,
		TrackingCode:    "ABC12334",
		BuyerID:         1,
		ProductRecordID: 1,
	}
)

func TestPurchaseOrderService_Save(t *testing.T) {
	t.Run("case 1: success - Should Create a Purchase Order", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, nil)
		svBu.On("FindByID", po.BuyerID).Return(internal.Buyer{}, nil)
		rpPo.On("Save", &po).Return(nil)

		err := sv.Save(&po)

		require.NoError(t, err)
	})

	t.Run("case 2 - error - Should return an error when trying to save a Purchase Order with an duplicated Order Number", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, nil)
		svBu.On("FindByID", po.BuyerID).Return(internal.Buyer{}, nil)
		rpPo.On("Save", &po).Return(internal.ErrPurchaseOrderConflict)

		err := sv.Save(&po)

		require.Error(t, err)
		require.Equal(t, internal.ErrPurchaseOrderConflict, err)
	})

	t.Run("case 3 - error - Should return an error when trying to save a Purchase Order with an non-existent Product Record", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, internal.ErrProductRecordsNotFound)
		svBu.On("FindByID", po.BuyerID).Return(internal.Buyer{}, nil)

		err := sv.Save(&po)

		require.Error(t, err)
		require.Equal(t, internal.ErrProductRecordsNotFound, err)
	})

	t.Run("case 4 - error - Should return an error when trying to save a Purchase Order with an non-existent Buyer", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, nil)
		svBu.On("FindByID", po.BuyerID).Return(internal.Buyer{}, service.ErrBuyerNotFound)

		err := sv.Save(&po)

		require.Error(t, err)
		require.Equal(t, service.ErrBuyerNotFound, err)
	})

	t.Run("case 5 - error - Should return an error when trying to save a Purchase Order with an invalid data", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, nil)
		svBu.On("FindByID", po.BuyerID).Return(internal.Buyer{}, nil)

		err := sv.Save(&internal.PurchaseOrder{})

		require.ErrorContains(t, internal.ErrPurchaseOrderBadRequest, err.Error())
	})

	t.Run("case 4 - error: Should return an error when an internal error occurs", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		rpPr := new(RepositoryProductRecordsMock)
		svBu := NewBuyerServiceMock()
		sv := service.NewPurchaseOrderService(rpPo, rpPr, svBu)

		rpPr.On("FindByID", po.ProductRecordID).Return(internal.ProductRecords{}, errors.New("internal server error"))

		err := sv.Save(&po)

		require.Error(t, err)
		require.EqualError(t, err, "internal server error")
	})
}

func TestPurchaseOrderService_FindByID(t *testing.T) {
	t.Run("case 1: success - Should return a Purchase Order", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		sv := service.NewPurchaseOrderService(rpPo, nil, nil)

		rpPo.On("FindByID", po.ID).Return(po, nil)

		p, err := sv.FindByID(po.ID)

		require.NoError(t, err)
		require.Equal(t, po, p)
	})

	t.Run("case 2 - error - Should return an error when trying to find a non-existent Purchase Order", func(t *testing.T) {
		rpPo := NewPurchaseOrderRepositoryMock()
		sv := service.NewPurchaseOrderService(rpPo, nil, nil)

		rpPo.On("FindByID", po.ID).Return(internal.PurchaseOrder{}, internal.ErrPurchaseOrderNotFound)

		_, err := sv.FindByID(po.ID)

		require.Error(t, err)
		require.Equal(t, internal.ErrPurchaseOrderNotFound, err)
	})
}
