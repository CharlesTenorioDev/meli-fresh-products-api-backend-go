package service_test

import (
	"errors"
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type localityRepositoryMock struct {
	mock.Mock
}

func (l *localityRepositoryMock) Save(locality *internal.Locality) (err error) {
	args := l.Called(locality)
	return args.Error(0)
}

func (l *localityRepositoryMock) ReportSellers() (localities []internal.Locality, err error) {
	args := l.Called()
	return args.Get(0).([]internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) ReportSellersByID(id int) (localities []internal.Locality, err error) {
	args := l.Called(id)
	return args.Get(0).([]internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) FindByID(id int) (locality internal.Locality, err error) {
	args := l.Called(id)
	return args.Get(0).(internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) ReportCarries(localityId int) (amountOfCarries int, e error) {
	args := l.Called(localityId)
	return args.Get(0).(int), args.Error(1)
}

func (l *localityRepositoryMock) GetAmountOfCarriesForEveryLocality() (c []internal.CarriesCountPerLocality, e error) {
	args := l.Called()
	return args.Get(0).([]internal.CarriesCountPerLocality), args.Error(1)
}

func TestLocalityDefault_Save(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		locality := &internal.Locality{
			ID:           2,
			LocalityName: "Test",
			ProvinceName: "Test",
			CountryName:  "Test",
			Sellers:      0,
		}

		mockRepo.On("Save", locality).Return(nil)

		svc := service.NewLocalityDefault(mockRepo)
		err := svc.Save(locality)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro de validação", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		locality := &internal.Locality{} // Crie uma instância inválida de Locality que retorne erros de validação

		svc := service.NewLocalityDefault(mockRepo)
		err := svc.Save(locality)

		assert.Error(t, err)
		assert.IsType(t, internal.DomainError{}, err)
		mockRepo.AssertNotCalled(t, "Save", mock.Anything)
	})

	t.Run("erro ao salvar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		locality := &internal.Locality{
			ID:           2,
			LocalityName: "Test",
			ProvinceName: "Test",
			CountryName:  "Test",
			Sellers:      0,
		}

		mockRepo.On("Save", locality).Return(errors.New("erro ao salvar"))

		svc := service.NewLocalityDefault(mockRepo)
		err := svc.Save(locality)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLocalityDefault_ReportSellers(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		expectedLocalities := []internal.Locality{{}} // Crie uma lista de localidades

		mockRepo.On("ReportSellers").Return(expectedLocalities, nil)

		svc := service.NewLocalityDefault(mockRepo)
		localities, err := svc.ReportSellers()

		assert.NoError(t, err)
		assert.Equal(t, expectedLocalities, localities)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)

		mockRepo.On("ReportSellers").Return([]internal.Locality{}, errors.New("erro ao buscar"))

		svc := service.NewLocalityDefault(mockRepo)
		localities, err := svc.ReportSellers()

		assert.Error(t, err)
		assert.Equal(t, 0, len(localities))
		mockRepo.AssertExpectations(t)
	})
}

func TestLocalityDefault_ReportSellersByID(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		id := 1
		expectedLocalities := []internal.Locality{{}} // Crie uma lista de localidades

		mockRepo.On("ReportSellersByID", id).Return(expectedLocalities, nil)

		svc := service.NewLocalityDefault(mockRepo)
		localities, err := svc.ReportSellersByID(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocalities, localities)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		id := 1

		mockRepo.On("ReportSellersByID", id).Return([]internal.Locality{}, errors.New("erro ao buscar"))

		svc := service.NewLocalityDefault(mockRepo)
		localities, err := svc.ReportSellersByID(id)

		assert.Error(t, err)
		assert.Equal(t, 0, len(localities))
		mockRepo.AssertExpectations(t)
	})
}

func TestLocalityDefault_FindByID(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		id := 1
		expectedLocality := internal.Locality{} // Crie uma instância de Locality

		mockRepo.On("FindByID", id).Return(expectedLocality, nil)

		svc := service.NewLocalityDefault(mockRepo)
		locality, err := svc.FindByID(id)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocality, locality)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		id := 1

		mockRepo.On("FindByID", id).Return(internal.Locality{}, errors.New("erro ao buscar"))

		svc := service.NewLocalityDefault(mockRepo)
		_, err := svc.FindByID(id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLocalityDefault_ReportCarries(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		localityId := 1
		expectedCount := 10

		mockRepo.On("ReportCarries", localityId).Return(expectedCount, nil)

		svc := service.NewLocalityDefault(mockRepo)
		count, err := svc.ReportCarries(localityId)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		localityId := 1

		mockRepo.On("ReportCarries", localityId).Return(0, errors.New("erro ao buscar"))

		svc := service.NewLocalityDefault(mockRepo)
		_, err := svc.ReportCarries(localityId)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestLocalityDefault_GetAmountOfCarriesForEveryLocality(t *testing.T) {
	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)
		expectedCarries := []internal.CarriesCountPerLocality{{}} // Crie uma lista de carries

		mockRepo.On("GetAmountOfCarriesForEveryLocality").Return(expectedCarries, nil)

		svc := service.NewLocalityDefault(mockRepo)
		carries, err := svc.GetAmountOfCarriesForEveryLocality()

		assert.NoError(t, err)
		assert.Equal(t, expectedCarries, carries)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(localityRepositoryMock)

		mockRepo.On("GetAmountOfCarriesForEveryLocality").Return([]internal.CarriesCountPerLocality{}, errors.New("erro ao buscar"))

		svc := service.NewLocalityDefault(mockRepo)
		carries, err := svc.GetAmountOfCarriesForEveryLocality()

		assert.Error(t, err)
		assert.Equal(t, 0, len(carries))
		mockRepo.AssertExpectations(t)
	})
}
