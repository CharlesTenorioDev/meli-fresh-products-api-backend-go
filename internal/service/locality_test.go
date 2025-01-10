package service_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/mock"
)

type localityRepositoryMock struct {
	mock.Mock
}

func (l *localityRepositoryMock) Save(locality *internal.Locality) (err error) {
	args := l.Called()
	return args.Error(0)
}

func (l *localityRepositoryMock) ReportSellers() (localities []internal.Locality, err error) {
	args := l.Called()
	return args.Get(0).([]internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) ReportSellersByID(id int) (localities []internal.Locality, err error) {
	args := l.Called()
	return args.Get(0).([]internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) FindByID(id int) (locality internal.Locality, err error) {
	args := l.Called()
	return args.Get(0).(internal.Locality), args.Error(1)
}

func (l *localityRepositoryMock) ReportCarries(localityId int) (amountOfCarries int, e error) {
	args := l.Called()
	return args.Get(0).(int), args.Error(1)
}

func (l *localityRepositoryMock) GetAmountOfCarriesForEveryLocality() (c []internal.CarriesCountPerLocality, e error) {
	args := l.Called()
	return args.Get(0).([]internal.CarriesCountPerLocality), args.Error(1)
}
