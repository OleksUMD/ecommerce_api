package cache

import (
	"context"
	"net/http"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
)

type CacheInterface interface {
	GetAuthToken(req *http.Request) (string, error)
	GetCategories(ctx context.Context) ([]dto.CategoryResponse, error)
	SetCategories(ctx context.Context, data []dto.CategoryResponse) error
	CleanCategories(ctx context.Context) error
}
