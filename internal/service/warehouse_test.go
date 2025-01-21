package service_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/mock"
)

type warehouseRepositoryMock struct {
	mock.Mock
}

func (r *warehouseRepositoryMock) FindAll() ([]internal.Warehouse, error) {
	args := r.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

func (r *warehouseRepositoryMock) FindByID(id int) (internal.Warehouse, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (r *warehouseRepositoryMock) Save(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *warehouseRepositoryMock) Update(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *warehouseRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}
