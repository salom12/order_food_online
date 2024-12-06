package services

import (
	"order_food_online/internal/models"
	"order_food_online/internal/repository"
)

type OrderService struct {
	orderRepo repository.OrderRepository
}

func NewOrderService(repo repository.OrderRepository) *OrderService {
	return &OrderService{orderRepo: repo}
}

func (s *OrderService) GetAllOrders() ([]models.Order, error) {
	return s.orderRepo.GetAllOrders()
}

func (s *OrderService) GetOrderByID(id int) (*models.Order, error) {
	return s.orderRepo.GetOrderByID(id)
}

func (s *OrderService) PlaceOrder(orderReq models.OrderRequest) (*models.Order, error) {
	return s.orderRepo.PlaceOrder(orderReq)
}

func (s *OrderService) CheckProductExists(productID int) (bool, error) {
	return s.orderRepo.CheckProductExists(productID)
}
