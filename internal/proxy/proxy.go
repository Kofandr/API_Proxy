package proxy

import (
	"context"
	"io"
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ProxyClient struct {
	Client HTTPClient // Используем интерфейс вместо *http.Client
}

func NewProxyClient() *ProxyClient {
	return &ProxyClient{
		Client: &http.Client{},
	}
}

func (p *ProxyClient) Do(ctx context.Context, method, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return p.Client.Do(req)
}
