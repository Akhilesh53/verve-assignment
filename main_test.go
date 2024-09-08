package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup a test router
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/api/verve/accept", handleRequest)
	return router
}

// Test: Valid request with unique ID
func TestValidRequest(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/verve/accept?id=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}

// Test: Invalid request (missing ID)
func TestInvalidRequest(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/verve/accept", nil) // Missing ID
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "failed", w.Body.String())
}

// Test: Deduplication of IDs (same ID sent twice)
func TestDuplicateRequest(t *testing.T) {
	router := setupRouter()

	req1, _ := http.NewRequest("GET", "/api/verve/accept?id=2", nil)
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "ok", w1.Body.String())

	req2, _ := http.NewRequest("GET", "/api/verve/accept?id=2", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code) // Deduplication should still return "ok"
	assert.Equal(t, "ok", w2.Body.String())
}

// Test: External endpoint with GET request
func TestExternalGETRequest(t *testing.T) {
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/verve/accept?id=3&endpoint=http://httpbin.org/get", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "ok", w.Body.String())
}
