package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/redis/go-redis/v9"
)

var ErrRedisNotFound = errors.New("key not found")

type CacheProvider struct {
	CacheClient   *redis.Client
	SessionClient *redis.Client
	Config        *config.CacheConfig
}

func NewCache(
	config *config.CacheConfig,
) *CacheProvider {
	addr := net.JoinHostPort(config.Host, config.Port)
	cacheClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.CacheDB,
	})
	sessionClient := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.AuthDB,
	})
	return &CacheProvider{
		Config:        config,
		CacheClient:   cacheClient,
		SessionClient: sessionClient,
	}
}

func (c *CacheProvider) Close() error {
	if err := c.CacheClient.Close(); err != nil {
		return fmt.Errorf("failed to close cache client: %w", err)
	}
	if err := c.SessionClient.Close(); err != nil {
		return fmt.Errorf("failed to close session client: %w", err)
	}
	return nil
}

func (c *CacheProvider) GetAuthToken(req *http.Request) (string, error) {
	return "", nil
}

func (c *CacheProvider) GetCategories(ctx context.Context) ([]dto.CategoryResponse, error) {
	var response []dto.CategoryResponse
	rawString, err := c.CacheClient.Get(ctx, "categories").Result()
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(rawString), &response); err != nil {
		return nil, err
	}
	return response, nil
}

func (c *CacheProvider) SetCategories(ctx context.Context, data []dto.CategoryResponse) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := c.CacheClient.Set(ctx, "categories", string(bytes), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (c *CacheProvider) CleanCategories(ctx context.Context) error {
	_, err := c.CacheClient.Del(ctx, "categories").Result()
	if err != nil {
		return err
	}
	return nil
}

func (c *CacheProvider) SetTempToken(token string, data *models.UnverifiedUser) error {
	ctx := context.Background()
	_, err := c.SessionClient.Set(ctx, token, data, c.Config.AuthTTL*time.Minute).Result()
	return err
}

func (c *CacheProvider) GetUnverifiedUser(token string) (*models.UnverifiedUser, error) {
	ctx := context.Background()
	res, err := c.SessionClient.Get(ctx, token).Bytes()
	if err != nil {
		return nil, ErrRedisNotFound
	}
	var unverifiedUser models.UnverifiedUser
	err = json.Unmarshal(res, &unverifiedUser)
	if err != nil {
		return nil, err
	}
	return &unverifiedUser, nil
}
