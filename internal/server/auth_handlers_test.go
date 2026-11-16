package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthcheck(t *testing.T) {
	req1 := httptest.NewRequest(http.MethodGet, "/health", nil)
	w1 := httptest.NewRecorder()
	testServer.Handler.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)

}
func TestRegisterNewUser_Success(t *testing.T) {
	defer cleanupTestData(t)
}

func TestRegisterNewUser_BadRequest(t *testing.T) {
	defer cleanupTestData(t)
}
