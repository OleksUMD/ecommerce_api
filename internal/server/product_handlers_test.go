package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/stretchr/testify/assert"
)

type TestCategoryResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    []dto.CategoryResponse `json:"data"`
	Error   string                 `json:"error"`
}

type TestProductResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []dto.ProductResponse `json:"data"`
	Error   string                `json:"error"`
	Meta    utils.PaginationMeta  `json:"meta"`
}

func TestGetCategories(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp TestCategoryResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp.Data, 2)
	category := resp.Data[0]
	assert.NotNil(t, category)
	assert.NotZero(t, category.ID)
}

func TestGetProducts(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/products", nil)
	w := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp TestProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), resp.Meta.Total)
	assert.Len(t, resp.Data, 2)
	assert.NotNil(t, resp.Data[0])
	assert.NotZero(t, resp.Data[0].ID)
}
