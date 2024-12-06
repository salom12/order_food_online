package tests

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"order_food_online/internal/handlers"
	"order_food_online/internal/mocks"
	"order_food_online/internal/models"
	"order_food_online/internal/services"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetProductsHandler(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)

	// Mock the service behavior
	mockRepo.On("GetAllProducts").Return([]models.Product{
		{ID: 1, Name: "Mock Product 1", Price: 19.99},
		{ID: 2, Name: "Mock Product 2", Price: 29.99},
	}, nil)

	service := services.NewProductService(mockRepo)
	handler := handlers.NewProductHandler(service, slog.Default())

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	if assert.NoError(t, handler.GetProducts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Mock Product 1")
		assert.Contains(t, rec.Body.String(), "Mock Product 2")
	}

	// Assert that the mock service was called
	mockRepo.AssertExpectations(t)
}
