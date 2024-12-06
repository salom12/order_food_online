package mocks

import (
	"github.com/stretchr/testify/mock"
	"order_food_online/internal/models"
)

// MockOrderService is a mock implementation of the OrderService interface
type MockOrderService struct {
	mock.Mock
}

// PlaceOrder mocks the CreateOrder method
func (m *MockOrderService) PlaceOrder(orderReq models.OrderRequest) (*models.Order, error) {
	args := m.Called(orderReq)
	return args.Get(0).(*models.Order), args.Error(1)
}

// GetOrderByID mocks the GetOrderByID method
func (m *MockOrderService) GetOrderByID(id int) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

// GetAllOrders mocks the GetAllOrders method
func (m *MockOrderService) GetAllOrders() ([]models.Order, error) {
	args := m.Called()
	return args.Get(0).([]models.Order), args.Error(1)
}

// CheckProductExists mocks the CheckProductExists method
func (m *MockOrderService) CheckProductExists(productID int) (bool, error) {
	args := m.Called(productID)
	return args.Get(0).(bool), args.Error(1)
}
