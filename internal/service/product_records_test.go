package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryProductRecordsMock struct {
	mock.Mock
}

func (m *RepositoryProductRecordsMock) Save(productRec internal.ProductRecords) (internal.ProductRecords, error) {
	args := m.Called(productRec)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

func (m *RepositoryProductRecordsMock) FindAll() ([]internal.ProductRecords, error) {
	args := m.Called()
	return args.Get(0).([]internal.ProductRecords), args.Error(1)
}

func (m *RepositoryProductRecordsMock) FindByID(id int) (internal.ProductRecords, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductRecords), args.Error(1)
}

func TestProductRecords_Create(t *testing.T) {
	product := internal.Product{ID: 1, ProductCode: "code-1"}
	productRec := internal.ProductRecords{
		ID:             1,
		LastUpdateDate: time.Now(),
		ProductID:      product.ID,
		SalePrice:      100.0,
		PurchasePrice:  50.0,
	}

	t.Run("successfully create product record", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		productRepo := new(RepositoryProductMock)

		serv := service.NewProductRecordsDefault(productRecRepo, productRepo)

		productRepo.On("FindByID", product.ID).Return(product, nil)
		productRecRepo.On("Save", productRec).Return(productRec, nil)

		result, err := serv.Create(productRec)

		assert.Nil(t, err)
		assert.Equal(t, productRec, result)
	})

	t.Run("error: product not found", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		productRepo := new(RepositoryProductMock)

		serv := service.NewProductRecordsDefault(productRecRepo, productRepo)
		productRepo.On("FindByID", productRec.ProductID).Return(internal.Product{}, internal.ErrProductNotFound)
		productRecRepo.On("Save", productRec).Return(productRec, nil)

		_, err := serv.Create(productRec)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductIdNotFound, err)
	})

	t.Run("error: validation error", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		productRepo := new(RepositoryProductMock)

		serv := service.NewProductRecordsDefault(productRecRepo, productRepo)
		invalidProductRec := internal.ProductRecords{
			ID:        1,
			ProductID: 0,
			SalePrice: 0,
		}

		_, err := serv.Create(invalidProductRec)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductUnprocessableEntity, err)
	})

	t.Run("error: failed to save product record", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		productRepo := new(RepositoryProductMock)

		serv := service.NewProductRecordsDefault(productRecRepo, productRepo)
		productRepo.On("FindByID", product.ID).Return(product, nil)
		productRecRepo.On("Save", productRec).Return(internal.ProductRecords{}, internal.ErrProductUnprocessableEntity)

		_, err := serv.Create(productRec)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrProductUnprocessableEntity, err)
	})
}

func TestProductRecords_GetAll(t *testing.T) {
	productRecords := []internal.ProductRecords{
		{ID: 1, ProductID: 1, SalePrice: 100.0, PurchasePrice: 50.0},
		{ID: 2, ProductID: 2, SalePrice: 200.0, PurchasePrice: 100.0},
	}

	t.Run("successfully retrieve all product records", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		serv := service.NewProductRecordsDefault(productRecRepo, nil)
		productRecRepo.On("FindAll").Return(productRecords, nil)
		result, err := serv.GetAll()

		assert.Nil(t, err)
		assert.Equal(t, productRecords, result)
	})

	t.Run("error retrieving product records", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		serv := service.NewProductRecordsDefault(productRecRepo, nil)
		productRecRepo.On("FindAll").Return([]internal.ProductRecords{}, errors.New("error retrieving product records"))
		result, err := serv.GetAll()

		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.New("error retrieving product records"), err)
	})
}

func TestProductRecords_GetByID(t *testing.T) {
	productRecords := internal.ProductRecords{ID: 1, ProductID: 1, SalePrice: 100.0, PurchasePrice: 50.0}

	t.Run("successfully retrieve by ID product records", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		serv := service.NewProductRecordsDefault(productRecRepo, nil)
		productRecRepo.On("FindByID", productRecords.ProductID).Return(productRecords, nil)
		result, err := serv.GetByID(productRecords.ProductID)

		assert.Nil(t, err)
		assert.Equal(t, productRecords, result)
	})

	t.Run("error by ID product records", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		serv := service.NewProductRecordsDefault(productRecRepo, nil)
		productRecRepo.On("FindByID", productRecords.ProductID).Return(productRecords, errors.New("ID invalid"))
		_, err := serv.GetByID(productRecords.ProductID)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("ID invalid"), err)
	})

	t.Run("error ID 0 product records", func(t *testing.T) {
		productRecRepo := new(RepositoryProductRecordsMock)
		serv := service.NewProductRecordsDefault(productRecRepo, nil)
		productRecRepo.On("FindByID", 0).Return(productRecords, errors.New("ID invalid"))
		_, err := serv.GetByID(0)

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("ID invalid"), err)
	})

}
