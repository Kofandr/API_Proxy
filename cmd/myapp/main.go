package main

import (
	"context"
	"github.com/Kofandr/API_Proxy.git/internal/logger"
	"github.com/Kofandr/API_Proxy.git/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Ну как заебок???
func main() {

	log := logger.New("INFO")

	mainServer := server.New(log)

	go func() {
		if err := mainServer.Start(); err != nil && err != http.ErrServerClosed {
			log.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Shutting down...")
	if err := mainServer.Shutdown(ctx); err != nil {
		log.Error("Shutdown failed", "error", err)
	} else {
		log.Info("Server stopped")
	}

}
