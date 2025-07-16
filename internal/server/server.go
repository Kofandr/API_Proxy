package server

import (
	"context"
	"github.com/Kofandr/API_Proxy/internal/client"
	"github.com/Kofandr/API_Proxy/internal/handler"
	"github.com/Kofandr/API_Proxy/internal/middleware"
	"time"

	"log/slog"
	"net/http"
)

type Server struct {
	Http *http.Server
	log  *slog.Logger
}

func New(log *slog.Logger) *Server {

	restyClient := client.NewRestyClient()
	handler := handler.New(restyClient)

	mux := http.NewServeMux()

	mux.Handle("/posts/", middleware.LoggerMiddleware(log, http.HandlerFunc(handler.Get)))

	Http := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return &Server{Http, log}

}

func (server *Server) Start() error {
	server.log.Info("Starting server", "addr", server.Http.Addr)
	return server.Http.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.log.Info("Shutting down server")
	return server.Http.Shutdown(ctx)
}
