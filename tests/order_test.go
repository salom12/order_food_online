package tests

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"order_food_online/internal/handlers"
	"order_food_online/internal/mocks"
	"order_food_online/internal/models"
	"order_food_online/internal/services"
	"testing"
)

func TestGetOrdersHandler(t *testing.T) {
	mockRepo := new(mocks.MockOrderService)

	// Mock the service behavior
	mockRepo.On("GetAllOrders").Return([]models.Order{
		{ID: 1, CouponCode: "test", FinalPrice: 100},
	}, nil)

	service := services.NewOrderService(mockRepo)
	handler := handlers.NewOrderHandler(service, slog.Default())

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if assert.NoError(t, handler.GetOrders(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "test")
	}

	// Assert that the mock service was called
	mockRepo.AssertExpectations(t)
}
