package mocks

import (
	"github.com/stretchr/testify/mock"
	"order_food_online/internal/models"
)

type MockProductRepository struct {
	mock.Mock
}

// GetAllProducts mocks the GetAllProducts method of the repository
func (m *MockProductRepository) GetAllProducts() ([]models.Product, error) {
	args := m.Called()
	return args.Get(0).([]models.Product), args.Error(1)
}

// GetProductByID mocks the GetProductById method of the repository
func (m *MockProductRepository) GetProductByID(id int) (*models.Product, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Product), args.Error(1)
}
