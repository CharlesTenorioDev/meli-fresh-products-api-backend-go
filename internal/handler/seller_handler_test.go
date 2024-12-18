package handler_test

import (
	"github.com/meli-fresh-products-api-backend-t1/internal"
	"github.com/meli-fresh-products-api-backend-t1/internal/dto"
	"github.com/stretchr/testify/mock"
)

// MockSellerService é uma estrutura mock para internal.SellerService
type MockSellerService struct {
	mock.Mock
}

// FindAll mock
func (m *MockSellerService) FindAll() ([]internal.Seller, error) {
	args := m.Called()
	return args.Get(0).([]internal.Seller), args.Error(1)
}

// FindByID mock
func (m *MockSellerService) FindByID(id int) (internal.Seller, error) {
	args := m.Called(id)
	return args.Get(0).(internal.Seller), args.Error(1)
}

// Save mock
func (m *MockSellerService) Save(seller *internal.Seller) error {
	args := m.Called(seller)
	seller.ID = 1
	return args.Error(0)
}

// Update mock
func (m *MockSellerService) Update(id int, updatedSeller dto.SellersUpdateRequest) (internal.Seller, error) {
	args := m.Called(updatedSeller)
	return args.Get(0).(internal.Seller), args.Error(0)
}

// Delete mock
func (m *MockSellerService) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
