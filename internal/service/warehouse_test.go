package service_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/mock"
)

func NewWarehouseRepositoryMock() *WarehouseRepositoryMock {
	return &WarehouseRepositoryMock{}
}

type WarehouseRepositoryMock struct {
	mock.Mock
}

func (r *WarehouseRepositoryMock) FindAll() ([]internal.Warehouse, error) {
	args := r.Called()
	return args.Get(0).([]internal.Warehouse), args.Error(1)
}

func (r *WarehouseRepositoryMock) FindByID(id int) (internal.Warehouse, error) {
	args := r.Called(id)
	return args.Get(0).(internal.Warehouse), args.Error(1)
}

func (r *WarehouseRepositoryMock) Save(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *WarehouseRepositoryMock) Update(warehouse *internal.Warehouse) error {
	args := r.Called(warehouse)
	return args.Error(0)
}

func (r *WarehouseRepositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}
