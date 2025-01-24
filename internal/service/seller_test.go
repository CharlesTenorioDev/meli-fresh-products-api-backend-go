package service_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type sellerRepositoryMock struct {
	mock.Mock
}

func (r *sellerRepositoryMock) FindAll() ([]internal.Seller, error) {
	args := r.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

func (r *sellerRepositoryMock) FindByID(id int) (internal.Seller, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (r *sellerRepositoryMock) FindByCID(cid int) (internal.Seller, error) {
	args := r.Called(cid)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (r *sellerRepositoryMock) Save(seller *internal.Seller) error {
	args := r.Called(seller)
	return args.Error(0)
}

func (r *sellerRepositoryMock) Update(seller *internal.Seller) error {
	args := r.Called(seller)
	return args.Error(0)
}

func (r *sellerRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func TestSellerServiceDefault_FindAll(t *testing.T) {
	t.Run("should return all sellers", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		expectedSellers := []internal.Seller{{ID: 1}, {ID: 2}}
		repo.On("FindAll").Return(expectedSellers, nil)

		sellers, err := svc.FindAll()

		assert.Nil(t, err)
		assert.Equal(t, expectedSellers, sellers)
	})

	t.Run("should return error when repo fails", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindAll").Return([]internal.Seller(nil), errors.New("repo fails"))

		sellers, err := svc.FindAll()

		assert.NotNil(t, err)
		assert.Equal(t, []internal.Seller(nil), sellers)
	})
}

func TestSellerServiceDefault_FindByID(t *testing.T) {
	t.Run("should return seller by id", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		expectedSeller := internal.Seller{ID: 1}
		repo.On("FindByID", 1).Return(expectedSeller, nil)

		seller, err := svc.FindByID(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedSeller, seller)
	})

	t.Run("should return error if repository fails", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByID", 1).Return(internal.Seller{}, errors.New("repository error"))

		seller, err := svc.FindByID(1)

		assert.NotNil(t, err)
		assert.Equal(t, internal.Seller{}, seller)
	})
}

func TestSellerServiceDefault_Save(t *testing.T) {

	t.Run("should save seller successfully", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         12345,
			CompanyName: "Test Seller",
			Address:     "Rua 1",
			Telephone:   "12 94562-7890",
			Locality:    1,
		}

		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{ID: 1}, nil)
		repo.On("Save", seller).Return(nil)

		err := svc.Save(seller)

		assert.Nil(t, err)
	})

	t.Run("should return error if seller CID already exists", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         12345,
			CompanyName: "Test Seller",
			Address:     "Rua 1",
			Telephone:   "12 94562-7890",
			Locality:    1,
		}

		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{CID: seller.CID}, nil)

		err := svc.Save(seller)

		assert.Equal(t, internal.ErrSellerCIDAlreadyExists, err)
	})

	t.Run("should return error if locality does not exist", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         12345,
			CompanyName: "Test Seller",
			Address:     "Rua 1",
			Telephone:   "12 94562-7890",
			Locality:    1,
		}

		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{}, errors.New("locality not found"))

		err := svc.Save(seller)

		assert.NotNil(t, err)
	})

	t.Run("should return error if repository fails to save", func(t *testing.T) {
		seller := &internal.Seller{
			ID:          1,
			CID:         12345,
			CompanyName: "Test Seller",
			Address:     "Rua 1",
			Telephone:   "12 94562-7890",
			Locality:    1,
		}

		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{}, nil)
		repo.On("Save", seller).Return(errors.New("repository error"))

		err := svc.Save(seller)

		assert.NotNil(t, err)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", mock.Anything).Return(internal.Seller{}, errors.New("repository error"))

		err := svc.Save(&internal.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Red Store",
			Address:     "Avenida Augusta",
			Telephone:   "11 95173-3133",
			Locality:    1,
		})

		assert.NotNil(t, err)
		assert.Equal(t, errors.New("repository error"), err)
	})

	t.Run("should return domain error with invalid fields", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		var expectedError = internal.DomainError{
			Message: "Seller fields invalid",
			Causes: []internal.Causes{
				{
					Field:   "cid",
					Message: "Company ID is required",
				},
				{
					Field:   "company_name",
					Message: "Company name is required",
				},
			},
		}

		seller := &internal.Seller{
			ID:          1,
			CID:         0,
			CompanyName: "",
			Address:     "Avenida Augusta",
			Telephone:   "11 95173-3133",
			Locality:    1,
		}

		err := svc.Save(seller)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSellerServiceDefault_Update(t *testing.T) {
	t.Run("should update seller successfully", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		actualSeller := internal.Seller{
			ID:          1,
			CID:         2,
			CompanyName: "Blue Store",
			Address:     "Avenida Paulista",
			Telephone:   "11 91332-3232",
			Locality:    1,
		}

		sellerPatch := internal.SellerPatch{
			CID:         intPtr(2),
			CompanyName: stringPtr("Red Store"),
			Address:     stringPtr("Avenida Augusta"),
			Telephone:   stringPtr("11 91332-3134"),
			Locality:    intPtr(2),
		}

		expectedSeller := internal.Seller{
			ID:          1,
			CID:         2,
			CompanyName: "Red Store",
			Address:     "Avenida Augusta",
			Telephone:   "11 91332-3134",
			Locality:    2,
		}

		repo.On("FindByID", mock.Anything).Return(actualSeller, nil)
		repo.On("FindByCID", mock.Anything).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", mock.Anything).Return(internal.Locality{ID: 2}, nil)
		repo.On("Update", mock.Anything).Return(nil)

		seller, err := svc.Update(expectedSeller.ID, sellerPatch)

		assert.Nil(t, err)
		assert.Equal(t, expectedSeller, seller)
	})

	t.Run("should return error if seller CID exists", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		actualSeller := internal.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Blue Store",
			Address:     "Avenida Paulista",
			Telephone:   "11 91332-3232",
			Locality:    1,
		}

		updatedSeller := internal.SellerPatch{
			CID:         intPtr(2),
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
			Locality:    nil,
		}

		repo.On("FindByID", actualSeller.ID).Return(actualSeller, nil)
		repo.On("FindByCID", *updatedSeller.CID).Return(internal.Seller{ID: 3, CID: *updatedSeller.CID}, nil)

		_, err := svc.Update(actualSeller.ID, updatedSeller)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrSellerCIDAlreadyExists, err)
	})

	t.Run("should return error if locality does not exist", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		actualSeller := internal.Seller{
			ID:          1,
			CID:         1,
			CompanyName: "Blue Store",
			Address:     "Avenida Paulista",
			Telephone:   "11 91332-3232",
			Locality:    1,
		}

		updatedSeller := internal.SellerPatch{
			CID:         nil,
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
			Locality:    intPtr(2),
		}

		repo.On("FindByID", actualSeller.ID).Return(actualSeller, nil)
		localityRepo.On("FindByID", *updatedSeller.Locality).Return(internal.Locality{}, internal.ErrLocalityNotFound)

		_, err := svc.Update(actualSeller.ID, updatedSeller)

		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrLocalityNotFound, err)
	})

	t.Run("should return error if repository does not find seller", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		sellerPatch := internal.SellerPatch{
			CID:         intPtr(2),
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
			Locality:    nil,
		}

		repo.On("FindByID", mock.Anything).Return(internal.Seller{}, internal.ErrSellerNotFound)

		_, err := svc.Update(1, sellerPatch)
		assert.NotNil(t, err)
		assert.Equal(t, internal.ErrSellerNotFound, err)
	})

	t.Run("should return error if repo fails", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		sellerPatch := internal.SellerPatch{
			CID:         intPtr(2),
			CompanyName: nil,
			Address:     nil,
			Telephone:   nil,
			Locality:    nil,
		}

		repo.On("FindByID", mock.Anything).Return(internal.Seller{}, nil)
		repo.On("FindByCID", mock.Anything).Return(internal.Seller{}, errors.New("repository error"))

		_, err := svc.Update(1, sellerPatch)
		assert.NotNil(t, err)
		assert.Equal(t, errors.New("repository error"), err)
	})
}

func TestSellerServiceDefault_Delete(t *testing.T) {
	t.Run("should delete seller successfully", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("Delete", 1).Return(nil)

		err := svc.Delete(1)

		assert.Nil(t, err)
	})

	t.Run("should return error if repository fails to delete", func(t *testing.T) {
		repo := new(sellerRepositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("Delete", 1).Return(errors.New("repository error"))

		err := svc.Delete(1)

		assert.NotNil(t, err)
	})
}
