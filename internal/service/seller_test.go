package service_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"testing"

	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) FindAll() ([]internal.Seller, error) {
	args := r.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

func (r *repositoryMock) FindByID(id int) (internal.Seller, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (r *repositoryMock) FindByCID(cid int) (internal.Seller, error) {
	args := r.Called(cid)
	return args.Get(0).(internal.Seller), args.Error(1)
}

func (r *repositoryMock) Save(seller *internal.Seller) error {
	args := r.Called(seller)
	return args.Error(0)
}

func (r *repositoryMock) Update(seller *internal.Seller) error {
	args := r.Called(seller)
	return args.Error(0)
}

func (r *repositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func TestSellerServiceDefault_FindAll(t *testing.T) {
	t.Run("should return all sellers", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		expectedSellers := []internal.Seller{{ID: 1}, {ID: 2}}
		repo.On("FindAll").Return(expectedSellers, nil)

		sellers, err := svc.FindAll()

		assert.Nil(t, err)
		assert.Equal(t, expectedSellers, sellers)
	})
}

func TestSellerServiceDefault_FindByID(t *testing.T) {
	t.Run("should return seller by id", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		expectedSeller := internal.Seller{ID: 1}
		repo.On("FindByID", 1).Return(expectedSeller, nil)

		seller, err := svc.FindByID(1)

		assert.Nil(t, err)
		assert.Equal(t, expectedSeller, seller)
	})

	t.Run("should return error if repository fails", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByID", 1).Return(internal.Seller{}, errors.New("repository error"))

		seller, err := svc.FindByID(1)

		assert.NotNil(t, err)
		assert.Equal(t, internal.Seller{}, seller)
	})
}

func TestSellerServiceDefault_Save(t *testing.T) {
	seller := &internal.Seller{
		CID:         12345,
		CompanyName: "Test Seller",
		Address:     "Rua 1",
		Telephone:   "1234567890",
		Locality:    1,
	}

	t.Run("should save seller successfully", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{}, nil)
		repo.On("Save", seller).Return(nil)

		err := svc.Save(seller)

		assert.Nil(t, err)
	})

	t.Run("should return error if seller CID already exists", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{CID: seller.CID}, nil)

		err := svc.Save(seller)

		assert.Equal(t, internal.ErrSellerCIDAlreadyExists, err)
	})

	t.Run("should return error if locality does not exist", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{}, errors.New("locality not found"))

		err := svc.Save(seller)

		assert.NotNil(t, err)
	})

	t.Run("should return error if repository fails to save", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("FindByCID", seller.CID).Return(internal.Seller{}, internal.ErrSellerNotFound)
		localityRepo.On("FindByID", seller.Locality).Return(internal.Locality{}, nil)
		repo.On("Save", seller).Return(errors.New("repository error"))

		err := svc.Save(seller)

		assert.NotNil(t, err)
	})
}

func TestSellerServiceDefault_Delete(t *testing.T) {
	t.Run("should delete seller successfully", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("Delete", 1).Return(nil)

		err := svc.Delete(1)

		assert.Nil(t, err)
	})

	t.Run("should return error if repository fails to delete", func(t *testing.T) {
		repo := new(repositoryMock)
		localityRepo := new(localityRepositoryMock)
		svc := service.NewSellerServiceDefault(repo, localityRepo)

		repo.On("Delete", 1).Return(errors.New("repository error"))

		err := svc.Delete(1)

		assert.NotNil(t, err)
	})
}
