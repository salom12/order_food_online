package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"order_food_online/internal/mocks"
	"order_food_online/internal/models"
	"order_food_online/internal/services"
	"order_food_online/pkg/helpers"
	"testing"
)

func TestValidatePromo(t *testing.T) {
	helpers.LoadEnv()
	mockCache := new(mocks.MockPromoCodeCache)

	// Mock cache miss for GetPromoCode
	mockCache.On("GetPromoCode", "PROMO123").Return(nil, nil)

	// Mock cache update for SetPromoCode
	mockCache.On("SetPromoCode", "PROMO123", mock.Anything, mock.Anything).Return(nil)

	service := services.NewPromoCodeService(mockCache)

	// Assuming PROMO123 isn't valid in this test case
	isValid, err := service.ValidatePromo("PROMO123")

	assert.NoError(t, err)
	assert.True(t, isValid)

	// Validate cache interactions
	mockCache.AssertCalled(t, "GetPromoCode", "PROMO123")
}

func TestValidatePromo_WithCacheHit(t *testing.T) {
	mockCache := new(mocks.MockPromoCodeCache)

	// Mock cache to return a valid promo code
	mockCache.On("GetPromoCode", "PROMO123").Return(&models.PromoCode{
		Code:    "PROMO123",
		IsValid: true,
	}, nil)

	service := services.NewPromoCodeService(mockCache)

	isValid, err := service.ValidatePromo("PROMO123")

	assert.NoError(t, err)
	assert.True(t, isValid)

	// Validate that cache methods were called as expected
	mockCache.AssertCalled(t, "GetPromoCode", "PROMO123")
	mockCache.AssertNotCalled(t, "SetPromoCode")
}

func TestValidatePromo_WithCacheMiss(t *testing.T) {
	mockCache := new(mocks.MockPromoCodeCache)

	// Mock cache miss for GetPromoCode
	mockCache.On("GetPromoCode", "PROMO123").Return(nil, nil)

	// Mock cache update for SetPromoCode
	mockCache.On("SetPromoCode", "PROMO123", mock.Anything, mock.Anything).Return(nil)

	service := services.NewPromoCodeService(mockCache)

	// Assuming PROMO123 isn't valid in this test case
	isValid, err := service.ValidatePromo("PROMO123")

	assert.NoError(t, err)
	assert.False(t, isValid)

	// Validate cache interactions
	mockCache.AssertCalled(t, "GetPromoCode", "PROMO123")
	mockCache.AssertCalled(t, "SetPromoCode", "PROMO123", mock.Anything, mock.Anything)
}
