package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"log/slog"
)

type Server struct {
	echo *echo.Echo
	addr string
	log  *slog.Logger
}

func New(log *slog.Logger) *Server {
	serverEcho := echo.New()

	return &Server{echo: serverEcho, addr: ":8080", log: log}

}

func (server *Server) Start() error {
	server.log.Info("Starting server", "addr", server.addr)
	return server.echo.Start(server.addr)
}

func (server *Server) Shutdown(ctx context.Context) error {
	server.log.Info("Shutting down server")
	return server.echo.Shutdown(ctx)
}
