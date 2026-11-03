package cache

import (
	"context"
	"net/http"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/models"
)

type ICacheProvider interface {
	Close() error
	GetAuthToken(req *http.Request) (string, error)
	GetCategories(ctx context.Context) ([]dto.CategoryResponse, error)
	SetCategories(ctx context.Context, data []dto.CategoryResponse) error
	CleanCategories(ctx context.Context) error
	SetTempToken(token string, data *models.UnverifiedUser) error
	GetUnverifiedUser(token string) (*models.UnverifiedUser, error)
}
