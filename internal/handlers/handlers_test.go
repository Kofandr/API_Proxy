package handlers

import (
	"context"
	"errors"
	"github.com/Kofandr/API_Proxy/internal/logger"
	"github.com/Kofandr/API_Proxy/internal/middleware"
	"github.com/Kofandr/API_Proxy/internal/proxy"

	"github.com/stretchr/testify/assert"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func setupRequest(t *testing.T, method, path string) (*http.Request, *slog.Logger) {
	logger := logger.New("INFO")
	ctx := context.WithValue(context.Background(), "logger", logger)
	req := httptest.NewRequest(method, path, nil).WithContext(ctx)
	return req, logger
}

// TestProxy_Success тестирует успешный сценарий проксирования запроса
func TestProxy_Success(t *testing.T) {
	// Arrange
	req, logger := setupRequest(t, "GET", "/posts")
	w := httptest.NewRecorder()

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"message": "success"}`)),
				Header:     http.Header{"Content-Type": []string{"application/json"}},
			}, nil
		},
	}

	handler := &Handler{
		baseURL: "https://jsonplaceholder.typicode.com",
		proxyClient: &proxy.ProxyClient{
			Client: mockClient,
		},
	}

	middlewareHandler := middleware.LoggerMiddleware(logger, http.HandlerFunc(handler.Proxy))

	middlewareHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	body, _ := io.ReadAll(resp.Body)
	assert.JSONEq(t, `{"message": "success"}`, string(body))
}

func TestProxy_InvalidPath(t *testing.T) {

	req, logger := setupRequest(t, "GET", "/invalid&path")
	w := httptest.NewRecorder()

	handler := New()

	middlewareHandler := middleware.LoggerMiddleware(logger, http.HandlerFunc(handler.Proxy))

	middlewareHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestProxy_UpstreamError(t *testing.T) {
	// Arrange
	req, logger := setupRequest(t, "GET", "/posts")
	w := httptest.NewRecorder()

	mockClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return nil, errors.New("upstream error")
		},
	}

	handler := &Handler{
		baseURL: "https://jsonplaceholder.typicode.com",
		proxyClient: &proxy.ProxyClient{
			Client: mockClient,
		},
	}

	middlewareHandler := middleware.LoggerMiddleware(logger, http.HandlerFunc(handler.Proxy))

	middlewareHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
