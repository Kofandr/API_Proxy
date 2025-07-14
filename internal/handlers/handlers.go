package handlers

import (
	"github.com/Kofandr/API_Proxy/internal/logger"
	"github.com/Kofandr/API_Proxy/internal/pathbuilder"
	"github.com/Kofandr/API_Proxy/internal/proxy"
	"github.com/Kofandr/API_Proxy/internal/utils"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Handler struct {
	baseURL     string
	proxyClient *proxy.ProxyClient
}

func New() *Handler {
	return &Handler{
		baseURL:     "https://jsonplaceholder.typicode.com",
		proxyClient: proxy.NewProxyClient(),
	}
}

func (handler *Handler) Proxy(w http.ResponseWriter, r *http.Request) {
	logger := logger.MustLoggerFromCtx(r.Context())

	targetURL, err := pathbuilder.BuildTargetURL(handler.baseURL, r.URL.Path)
	if err != nil {
		logger.Error("Invalid path",
			"path", r.URL.Path,
			"error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	supportedMethods := utils.GetSupportedMethods(targetURL)

	if !utils.Contains(supportedMethods, r.Method) {
		logger.Error("Method not allowed",
			"method", r.Method,
			"path", r.URL.Path)
		w.Header().Set("Allow", strings.Join(supportedMethods, ", "))
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp, err := handler.proxyClient.Do(r.Context(), r.Method, targetURL, r.Body, r.Header)
	if err != nil {
		logger.Error("Upstream error", "url", targetURL, "method", r.Method, "error", err)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

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
