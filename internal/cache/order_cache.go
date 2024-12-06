package cache

import (
	"context"
	"encoding/json"
	"order_food_online/internal/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type OrderCache interface {
	GetAllOrders() ([]models.Order, error)
	SetAllOrders([]models.Order, time.Duration) error
	GetOrderByID(int) (*models.Order, error)
	SetOrderByID(int, *models.Order, time.Duration) error
}

type redisOrderCache struct {
	client *redis.Client
}

func NewOrderCache(client *redis.Client) OrderCache {
	return &redisOrderCache{client: client}
}

func (c *redisOrderCache) GetAllOrders() ([]models.Order, error) {
	data, err := c.client.Get(context.Background(), "Orders").Result()
	if err != nil {
		return nil, err
	}

	var Orders []models.Order
	if err := json.Unmarshal([]byte(data), &Orders); err != nil {
		return nil, err
	}
	return Orders, nil
}

func (c *redisOrderCache) SetAllOrders(Orders []models.Order, ttl time.Duration) error {
	data, err := json.Marshal(Orders)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), "Orders", data, ttl).Err()
}

func (c *redisOrderCache) GetOrderByID(id int) (*models.Order, error) {
	data, err := c.client.Get(context.Background(), buildOrderKey(id)).Result()
	if err != nil {
		return nil, err
	}

	var Order models.Order
	if err := json.Unmarshal([]byte(data), &Order); err != nil {
		return nil, err
	}
	return &Order, nil
}

func (c *redisOrderCache) SetOrderByID(id int, Order *models.Order, ttl time.Duration) error {
	data, err := json.Marshal(Order)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), buildOrderKey(id), data, ttl).Err()
}

func buildOrderKey(id int) string {
	return "Order:" + strconv.Itoa(id)
}
