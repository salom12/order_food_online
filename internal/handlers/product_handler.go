package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order_food_online/internal/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Custom error definitions
var (
	errFailedToFetchProducts = errors.New("failed to fetch products")
	errInvalidProductID      = errors.New("invalid product ID")
	errProductNotFound       = errors.New("product not found")
)

// ProductHandler handles HTTP requests related to products
type ProductHandler struct {
	service *services.ProductService
	logger  *slog.Logger
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(service *services.ProductService, logger *slog.Logger) *ProductHandler {
	return &ProductHandler{service: service, logger: logger}
}

// RegisterProductRoutes sets up the routes for product-related endpoints
func (h *ProductHandler) RegisterProductRoutes(e *echo.Echo) {
	e.GET("/products", h.GetProducts)
	e.GET("/products/:id", h.GetProductByID)
}

// GetProducts handles the GET /products request
func (h *ProductHandler) GetProducts(c echo.Context) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		err := fmt.Errorf("%w: %v", errFailedToFetchProducts, err)
		h.logger.Error(err.Error(), "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errFailedToFetchProducts.Error()})
	}
	return c.JSON(http.StatusOK, products)
}

// GetProductByID handles the GET /products/:id request
func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("%w: %v", errInvalidProductID, err)
		h.logger.Error(err.Error(), slog.String("param", c.Param("id")), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errInvalidProductID.Error()})
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		err := fmt.Errorf("%w: %v", errProductNotFound, err)
		h.logger.Error(err.Error(), slog.Int("productID", id), "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": errProductNotFound.Error()})
	}
	return c.JSON(http.StatusOK, product)
}
