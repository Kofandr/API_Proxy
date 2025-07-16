package handler

import (
	"github.com/Kofandr/API_Proxy/config"
	"github.com/Kofandr/API_Proxy/internal/logger"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type Handler struct {
	baseURL     string
	restyClient *resty.Client
}

func New(client *resty.Client, cfg *config.Configuration) *Handler {
	return &Handler{
		baseURL:     cfg.PathProxy,
		restyClient: client,
	}
}

func (handler *Handler) Get(w http.ResponseWriter, r *http.Request) {
	logger := logger.MustLoggerFromCtx(r.Context())

	if r.Method != http.MethodGet {
		logger.Info("Invalid method", "method", r.Method)
		http.Error(w, "Invalid method", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/posts/")
	resp, err := handler.restyClient.R().
		SetHeader("Accept", "application/json").
		SetContext(r.Context()).
		Get(handler.baseURL + id)
	if err != nil {
		logger.Error("Upstream error", "url", (handler.baseURL + id), "method", r.Method, "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	for key, values := range resp.Header() {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode())

	if _, err := w.Write(resp.Body()); err != nil {
		logger.Error("Failed to write response", "error", err)
		http.Error(w, "Failed to proxy response", http.StatusInternalServerError)
	}

}
