package handlers

import (
	"github.com/Kofandr/API_Proxy.git/internal/middleware"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Handler struct {
	baseURL string
	client  HTTPClient
}

func New() *Handler {
	return &Handler{baseURL: "https://jsonplaceholder.typicode.com", client: &http.Client{Timeout: 10 * time.Second}}
}

func (handler *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	// Получаем логгер из контекста
	logger := middleware.GetLogger(r.Context())

	path := strings.Trim(strings.TrimPrefix(r.URL.Path, "/api/"), "/")
	if path == "" {
		logger.Error("Empty path",
			"path", r.URL.Path)
		http.Error(w, "Empty path after /api/", http.StatusBadRequest)
		return
	}

	parts := strings.Split(path, "/")
	if parts[0] != "posts" {
		logger.Error("Invalid endpoint",
			"path", r.URL.Path)
		http.Error(w, "Invalid endpoint", http.StatusNotFound)
		return
	}

	targetURL := handler.baseURL + "/posts"
	supportedMethods := []string{http.MethodGet, http.MethodPost}
	if len(parts) == 2 {
		if matched, _ := regexp.MatchString(`^\d+$`, parts[1]); !matched {
			logger.Error("Invalid post ID",
				"id", parts[1])
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		targetURL += "/" + parts[1]
		supportedMethods = []string{http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete}
	} else if len(parts) > 2 {
		logger.Error("Invalid path",
			"path", r.URL.Path)
		http.Error(w, "Invalid path", http.StatusNotFound)
		return
	}

	if !contains(supportedMethods, r.Method) {
		logger.Error("Method not allowed",
			"method", r.Method,
			"path", r.URL.Path)
		w.Header().Set("Allow", strings.Join(supportedMethods, ", "))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		logger.Error("Failed to create request",
			"url", targetURL,
			"error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := handler.client.Do(req)
	if err != nil {
		logger.Error("Upstream error",
			"url", targetURL,
			"method", r.Method,
			"error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		logger.Error("Upstream error",
			"status", resp.StatusCode,
			"url", targetURL,
			"method", r.Method)
		http.Error(w, "Upstream error: "+resp.Status, http.StatusBadGateway)
		return
	}

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	if _, err := io.Copy(w, resp.Body); err != nil {
		logger.Error("Failed to proxy response",
			"error", err)
		http.Error(w, "Failed to proxy response", http.StatusInternalServerError)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
