package server

import (
	"context"
	"database/sql"
	"github.com/labstack/echo/v4"
	echo_middleware "github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"order_food_online/internal/handlers"
	"order_food_online/pkg/middleware"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func StartServer(
	e *echo.Echo,
	productHandler *handlers.ProductHandler,
	orderHandler *handlers.OrderHandler,
	db *sql.DB,
) {
	// Register routes
	e.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	productHandler.RegisterProductRoutes(e)
	orderHandler.RegisterOrderRoutes(e)

	e.Use(echo_middleware.Logger())
	e.Use(echo_middleware.Recover())
	e.Use(middleware.AuthMiddleware())

	// Start the server
	port := getPort()
	go func() {
		if err := e.Start(port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start the server: %v", err)
		}
	}()

	// Handle OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shut down the server: %v", err)
	}

	// Cleanup resources
	if err := db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}
