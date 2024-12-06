package cache

import (
	"context"
	"encoding/json"
	"order_food_online/internal/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type ProductCache interface {
	GetAllProducts() ([]models.Product, error)
	SetAllProducts([]models.Product, time.Duration) error
	GetProductByID(int) (*models.Product, error)
	SetProductByID(int, *models.Product, time.Duration) error
}

type redisProductCache struct {
	client *redis.Client
}

func NewProductCache(client *redis.Client) ProductCache {
	return &redisProductCache{client: client}
}

func (c *redisProductCache) GetAllProducts() ([]models.Product, error) {
	data, err := c.client.Get(context.Background(), "products").Result()
	if err != nil {
		return nil, err
	}

	var products []models.Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (c *redisProductCache) SetAllProducts(products []models.Product, ttl time.Duration) error {
	data, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), "products", data, ttl).Err()
}

func (c *redisProductCache) GetProductByID(id int) (*models.Product, error) {
	data, err := c.client.Get(context.Background(), buildProductKey(id)).Result()
	if err != nil {
		return nil, err
	}

	var product models.Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (c *redisProductCache) SetProductByID(id int, product *models.Product, ttl time.Duration) error {
	data, err := json.Marshal(product)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), buildProductKey(id), data, ttl).Err()
}

func buildProductKey(id int) string {
	return "product:" + strconv.Itoa(id)
}
