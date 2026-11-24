package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UrielMT343/podtank/internal/server/signaling"
)

func main() {
	opts := slog.HandlerOptions{Level: slog.LevelInfo}
	handler := slog.NewJSONHandler(os.Stdout, &opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("Starting game server")

	rootCtx, stopCtx := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopCtx()

	signalingHandler := signaling.NewHandler()

	h := &Handlers{
		signalingHandler: signalingHandler,
	}

	basePrefix := "/api"

	routeHandler := routes(h, basePrefix)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           routeHandler,
		BaseContext:       func(net.Listener) context.Context { return rootCtx },
		ReadHeaderTimeout: 3 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       120 * time.Second,
	}

	slog.Info("Server started", "Port", "8080")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Critical server error", "error", err)
		}
	}()
	<-rootCtx.Done()

	slog.Info("Shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("Server exited properly")
}
