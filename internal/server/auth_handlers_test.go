package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

}
func TestRegisterNewUser_Success(t *testing.T) {
	defer cleanupTestData(t)
	formData := dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "Testpassword123!",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+1111111111",
	}
	reqBody, err := json.Marshal(&formData)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegisterNewUser_BadRequest(t *testing.T) {
	defer cleanupTestData(t)
	formDataMisEmail := dto.RegisterRequest{
		Password:  "Testpassword123!",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+1111111111",
	}
	reqBody, err := json.Marshal(&formDataMisEmail)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
	w := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	formDataSimplePassword := dto.RegisterRequest{
		Email:     "test@test.com",
		Password:  "toosimple",
		FirstName: "John",
		LastName:  "Doe",
		Phone:     "+1111111111",
	}
	reqBody, err = json.Marshal(&formDataSimplePassword)
	assert.NoError(t, err)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(reqBody))
	w = httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
