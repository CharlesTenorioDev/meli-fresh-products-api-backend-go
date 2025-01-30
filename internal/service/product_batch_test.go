package service_test

import (
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func NewProductBatchRepositoryMock() *ProductBatchRepositoryMock {
	return &ProductBatchRepositoryMock{}
}

type ProductBatchRepositoryMock struct {
	mock.Mock
}

func (r *ProductBatchRepositoryMock) FindByID(id int) (internal.ProductBatch, error) {
	args := r.Called(id)
	return args.Get(0).(internal.ProductBatch), args.Error(1)
}

func (r *ProductBatchRepositoryMock) Save(prodBatch *internal.ProductBatch) error {
	args := r.Called(prodBatch)
	return args.Error(0)
}

func (r *ProductBatchRepositoryMock) ProductBatchNumberExists(batchNumber int) (bool, error) {
	args := r.Called(batchNumber)
	return args.Get(0).(bool), args.Error(1)
}

func (r *ProductBatchRepositoryMock) ReportProducts() (prodBatches []internal.ProductBatch, err error) {
	args := r.Called()
	return args.Get(0).([]internal.ProductBatch), args.Error(1)
}

func (r *ProductBatchRepositoryMock) ReportProductsByID(id int) (prodBatches []internal.ProductBatch, err error) {
	args := r.Called(id)
	return args.Get(0).([]internal.ProductBatch), args.Error(1)
}

func newProductBatchService() (*service.ProductBatchService, *ProductBatchRepositoryMock, *SectionRepositoryMock, *RepositoryProductMock) {
	rpProductBatch := NewProductBatchRepositoryMock()
	rpSection := NewSectionRepositoryMock()
	rpProduct := NewRepositoryProductMock()

	return service.NewServiceProductBatch(rpProductBatch, rpSection, rpProduct), rpProductBatch, rpSection, rpProduct
}

func newTestProductBatch(id int, batchNumber int, productID int, prodBatchID int) internal.ProductBatch {
	return internal.ProductBatch{
		ID:                 id,
		BatchNumber:        batchNumber,
		CurrentQuantity:    150,
		CurrentTemperature: 17.4,
		DueDate:            "2023-01-01",
		InitialQuantity:    100,
		ManufacturingDate:  "2023-01-01",
		ManufacturingHour:  17,
		MinumumTemperature: 20.5,
		ProductID:          productID,
		SectionID:          prodBatchID,
	}
}

func TestService_CreateProductBatchUnitTest(t *testing.T) {
	t.Run("successfully create a new product-batch", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 4, 3)

		rpProductBatch.On("ProductBatchNumberExists", prodBatchCreate.BatchNumber).Return(false, nil)
		rpProduct.On("FindByID", prodBatchCreate.ProductID).Return(internal.Product{ID: prodBatchCreate.ProductID}, nil)
		rpSection.On("FindByID", prodBatchCreate.SectionID).Return(internal.Section{ID: prodBatchCreate.SectionID}, nil)
		rpProductBatch.On("Save", &prodBatchCreate).Return(nil)

		err := sv.Save(&prodBatchCreate)

		require.NoError(t, err)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 1)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 1)
	})

	t.Run("return fail error when required field is missing", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 0, 2)

		err := sv.Save(&prodBatchCreate)

		require.ErrorIs(t, err, internal.ErrProductBatchUnprocessableEntity)
		require.Contains(t, err.Error(), "couldn't parse product-batch")

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 0)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("return conflict error when number is already in use", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 4, 3)

		rpProductBatch.On("ProductBatchNumberExists", prodBatchCreate.BatchNumber).Return(true, nil)

		err := sv.Save(&prodBatchCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrProductBatchNumberAlreadyInUse)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 1)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 0)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when product does not exist", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 99, 3)

		rpProductBatch.On("ProductBatchNumberExists", prodBatchCreate.BatchNumber).Return(false, nil)
		rpProduct.On("FindByID", prodBatchCreate.ProductID).Return(internal.Product{}, internal.ErrProductNotFound)

		err := sv.Save(&prodBatchCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrProductNotFound)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 1)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 0)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when product type does not exist", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 3, 99)

		rpProductBatch.On("ProductBatchNumberExists", prodBatchCreate.BatchNumber).Return(false, nil)
		rpProduct.On("FindByID", prodBatchCreate.ProductID).Return(internal.Product{}, nil)
		rpSection.On("FindByID", prodBatchCreate.SectionID).Return(internal.Section{}, internal.ErrSectionNotFound)

		err := sv.Save(&prodBatchCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrSectionNotFound)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 1)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 0)
	})

	t.Run("returns error when product-batch fails to save", func(t *testing.T) {
		sv, rpProductBatch, rpSection, rpProduct := newProductBatchService()

		prodBatchCreate := newTestProductBatch(0, 101, 3, 99)

		rpProductBatch.On("ProductBatchNumberExists", prodBatchCreate.BatchNumber).Return(false, nil)
		rpProduct.On("FindByID", prodBatchCreate.ProductID).Return(internal.Product{}, nil)
		rpSection.On("FindByID", prodBatchCreate.SectionID).Return(internal.Section{}, nil)
		rpProductBatch.On("Save", &prodBatchCreate).Return(internal.ErrProductBatchUnprocessableEntity)

		err := sv.Save(&prodBatchCreate)

		require.Error(t, err)
		require.ErrorIs(t, err, internal.ErrProductBatchUnprocessableEntity)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "ProductBatchNumberExists", 1)
		rpProduct.AssertExpectations(t)
		rpProduct.AssertNumberOfCalls(t, "FindByID", 1)
		rpSection.AssertExpectations(t)
		rpSection.AssertNumberOfCalls(t, "FindByID", 1)
		rpProductBatch.AssertNumberOfCalls(t, "Save", 1)
	})
}

func TestService_ReadProductBatchUnitTest(t *testing.T) {
	t.Run("return error when reading a nonexistent product-batch by ID", func(t *testing.T) {
		sv, rpProductBatch, _, _ := newProductBatchService()
		expectedError := internal.ErrProductBatchNotFound

		rpProductBatch.On("FindByID", 1).Return(internal.ProductBatch{}, expectedError)

		_, err := sv.FindByID(1)

		require.Error(t, err)
		require.ErrorIs(t, err, expectedError)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "FindByID", 1)
	})

	t.Run("successfully read an existing product-batch by ID", func(t *testing.T) {
		sv, rpProductBatch, _, _ := newProductBatchService()
		expectedSection := newTestProductBatch(2, 101, 4, 3)

		rpProductBatch.On("FindByID", 2).Return(expectedSection, nil)

		produBatch, err := sv.FindByID(2)

		require.NoError(t, err)
		require.Equal(t, expectedSection, produBatch)

		rpProductBatch.AssertExpectations(t)
		rpProductBatch.AssertNumberOfCalls(t, "FindByID", 1)
	})
}
