package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"order_food_online/internal/models"
	"time"
)

type PromoCodeCache interface {
	GetPromoCode(string) (*models.PromoCode, error)
	SetPromoCode(string, *models.PromoCode, time.Duration) error
}

type redisPromoCodeCache struct {
	client *redis.Client
}

func NewPromoCodeCache(client *redis.Client) PromoCodeCache {
	return &redisPromoCodeCache{client: client}
}

func (c *redisPromoCodeCache) GetPromoCode(code string) (*models.PromoCode, error) {
	data, err := c.client.Get(context.Background(), buildPromoCodeKey(code)).Result()
	if err != nil {
		return nil, err
	}

	var PromoCode models.PromoCode
	if err := json.Unmarshal([]byte(data), &PromoCode); err != nil {
		return nil, err
	}
	return &PromoCode, nil
}

func (c *redisPromoCodeCache) SetPromoCode(code string, PromoCode *models.PromoCode, ttl time.Duration) error {
	data, err := json.Marshal(PromoCode)
	if err != nil {
		return err
	}
	return c.client.Set(context.Background(), buildPromoCodeKey(code), data, ttl).Err()
}

func buildPromoCodeKey(code string) string {
	return fmt.Sprintf("PromoCode:%s", code)
}
