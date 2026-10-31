package cache

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
	config *config.Config
}

func NewCache(
	config *config.Config,
) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.Cache.Password,
		DB:       config.Cache.CacheDB,
	})
	return &Cache{
		config: config,
		client: client,
	}
}
func (c *Cache) Close() error {
	return c.client.Close()
}
func (c *Cache) GetAuthToken(req *http.Request) (string, error) {
	return "", nil
}

func (c *Cache) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	var response []dto.CategoryResponse
	rawString, err := c.client.Get(ctx, "categories").Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(rawString), &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *Cache) SetCategories(ctx context.Context, data []dto.CategoryResponse) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := c.client.Set(ctx, "categories", string(bytes), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cache) CleanCategories(ctx context.Context) error {
	_, err := c.client.Del(ctx, "categories").Result()
	if err != nil {
		return err
	}
	return nil
}
