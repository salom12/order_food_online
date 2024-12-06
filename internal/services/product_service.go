package services

import (
	"order_food_online/internal/models"
	"order_food_online/internal/repository"
)

type ProductService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: repo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAllProducts()
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.productRepo.GetProductByID(id)
}
