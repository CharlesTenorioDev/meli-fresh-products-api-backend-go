package service

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/stretchr/testify/mock"
)

type ProductTypeRepositoryMock struct {
	mock.Mock
}

func (m *ProductTypeRepositoryMock) FindByID(id int) (internal.ProductType, error) {
	args := m.Called(id)
	return args.Get(0).(internal.ProductType), args.Error(1)
}
