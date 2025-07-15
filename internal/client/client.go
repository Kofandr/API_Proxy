package client

import (
	"github.com/Kofandr/API_Proxy/internal/logger"
	"github.com/go-resty/resty/v2"
	"time"
)

func NewRestyClient() *resty.Client {
	client := resty.New()

	// решил чтото сам добавить чтобы по разбираться  типа в библиотеке

	client.
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log := logger.MustLoggerFromCtx(resp.Request.Context())
		log.Info("resty response completed",
			"status", resp.StatusCode(),
			"url", resp.Request.URL,
			"duration", resp.Time())

		return nil

	})
	return client
}
