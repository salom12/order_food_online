package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order_food_online/internal/models"
	"order_food_online/internal/services"
	"strconv"

	"github.com/labstack/echo/v4"
)

// Custom error definitions
var (
	errFailedToFetchOrders = errors.New("failed to fetch orders")
	errInvalidOrderID      = errors.New("invalid order ID")
	errOrderNotFound       = errors.New("order not found")
)

// OrderHandler handles HTTP requests related to Orders
type OrderHandler struct {
	service          *services.OrderService
	promoCodeService *services.PromoCodeService
	logger           *slog.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(service *services.OrderService, logger *slog.Logger) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

// RegisterOrderRoutes sets up the routes for Order-related endpoints
func (h *OrderHandler) RegisterOrderRoutes(e *echo.Echo) {
	e.GET("/orders", h.GetOrders)
	e.POST("/orders", h.PlaceOrder)
	e.GET("/orders/:id", h.GetOrderByID)
}

// GetOrders handles the GET /Orders request
func (h *OrderHandler) GetOrders(c echo.Context) error {
	Orders, err := h.service.GetAllOrders()
	if err != nil {
		err := fmt.Errorf("%w: %v", errFailedToFetchOrders, err)
		h.logger.Error(err.Error(), "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errFailedToFetchOrders.Error()})
	}
	return c.JSON(http.StatusOK, Orders)
}

// GetOrderByID handles the GET /Orders/:id request
func (h *OrderHandler) GetOrderByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		err := fmt.Errorf("%w: %v", errInvalidOrderID, err)
		h.logger.Error(err.Error(), slog.String("param", c.Param("id")), "error", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": errInvalidOrderID.Error()})
	}

	Order, err := h.service.GetOrderByID(id)
	if err != nil {
		err := fmt.Errorf("%w: %v", errOrderNotFound, err)
		h.logger.Error(err.Error(), slog.Int("OrderID", id), "error", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": errOrderNotFound.Error()})
	}
	return c.JSON(http.StatusOK, Order)
}

// PlaceOrder handles the POST /Orders request
func (h *OrderHandler) PlaceOrder(c echo.Context) error {
	var orderReq models.OrderRequest

	// Bind the request body to OrderRequest
	if err := c.Bind(&orderReq); err != nil {
		h.logger.Error("Invalid request payload", slog.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// check promo code
	if orderReq.CouponCode.Valid && orderReq.CouponCode.String != "" {
		ok, err := h.promoCodeService.ValidatePromo(orderReq.CouponCode.String)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		}
		if !ok {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid coupon"})
		}
	}

	// Validate that products exist
	for _, item := range orderReq.Items {
		productExists, err := h.service.CheckProductExists(item.ProductID)
		if err != nil {
			h.logger.Error("Failed to check product existence", slog.Int("productID", item.ProductID), "error", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
		if !productExists {
			h.logger.Warn("Product does not exist", slog.Int("productID", item.ProductID))
			return c.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Product with ID %d does not exist", item.ProductID)})
		}
	}

	// Place the order
	order, err := h.service.PlaceOrder(orderReq)
	if err != nil {
		h.logger.Error("Failed to place order", slog.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to place order"})
	}

	return c.JSON(http.StatusCreated, order)
}
