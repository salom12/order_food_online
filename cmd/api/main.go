package main

import (
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"go.uber.org/dig"
	"log"
	"log/slog"
	"order_food_online/config"
	"order_food_online/internal/cache"
	"order_food_online/internal/handlers"
	"order_food_online/internal/repository"
	"order_food_online/internal/server"
	"order_food_online/internal/services"
	"os"
)

func provideDependencies(container *dig.Container) error {
	// Provide the SQL database
	if err := container.Provide(func() (*sql.DB, error) {
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}
		if err = db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}); err != nil {
		return err
	}

	// Provide the Redis client
	if err := container.Provide(func() (*redis.Client, error) {
		client := redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_URL"),
		})
		if _, err := client.Ping(client.Context()).Result(); err != nil {
			return nil, err
		}
		return client, nil
	}); err != nil {
		return err
	}

	// Provide logger
	if err := container.Provide(slog.Default); err != nil {
		return err
	}

	// Provide cache
	if err := container.Provide(cache.NewProductCache); err != nil {
		return err
	}
	if err := container.Provide(cache.NewOrderCache); err != nil {
		return err
	}
	if err := container.Provide(cache.NewPromoCodeCache); err != nil {
		return err
	}

	// Provide repositories
	if err := container.Provide(repository.NewProductRepository); err != nil {
		return err
	}
	if err := container.Provide(repository.NewOrderRepository); err != nil {
		return err
	}

	// Provide services
	if err := container.Provide(services.NewProductService); err != nil {
		return err
	}
	if err := container.Provide(services.NewOrderService); err != nil {
		return err
	}
	if err := container.Provide(services.NewPromoCodeService); err != nil {
		return err
	}

	// Provide handlers
	if err := container.Provide(handlers.NewProductHandler); err != nil {
		return err
	}
	if err := container.Provide(handlers.NewOrderHandler); err != nil {
		return err
	}

	// Provide the Echo instance
	if err := container.Provide(func() *echo.Echo {
		return echo.New()
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	// Create a new Dig container
	container := dig.New()

	// Provide dependencies to the container
	if err := provideDependencies(container); err != nil {
		log.Fatalf("Error providing dependencies: %v", err)
	}

	// Invoke the API server setup
	if err := container.Invoke(server.StartServer); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
