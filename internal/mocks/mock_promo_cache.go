package mocks

import (
	"github.com/stretchr/testify/mock"
	"order_food_online/internal/models"
	"time"
)

type MockPromoCodeCache struct {
	mock.Mock
}

func (m *MockPromoCodeCache) GetPromoCode(code string) (*models.PromoCode, error) {
	args := m.Called(code)

	// Handle nil return safely
	if promoCode, ok := args.Get(0).(*models.PromoCode); ok {
		return promoCode, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPromoCodeCache) SetPromoCode(code string, PromoCode *models.PromoCode, ttl time.Duration) error {
	args := m.Called(code, PromoCode, ttl)
	return args.Error(0)
}
