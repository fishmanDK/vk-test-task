package main

import (
	"context"
	"fmt"
	"github.com/fishmanDK/vk-test-task/config"
	"github.com/fishmanDK/vk-test-task/internal/http-server/handlers"
	"github.com/fishmanDK/vk-test-task/internal/service"
	"github.com/fishmanDK/vk-test-task/internal/storage"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

// @title VK-test-task-API
// @version 1.0
// @description Api Server for VK-test Application

// @host localhost:8081
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @description Enter your bearer token in the format Bearer <token>
// @in header
// @name Authorization
func main() {
	cfg := config.MustLoadConfigHTTP()
	pg_cfg := config.MustLoadConfigPostgresDB()
	logger := setupLogger(cfg.Env)

	db, err := storage.MustStorage(pg_cfg)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	srvc, err := service.MustService(db)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	handl, err := handlers.MustHandlers(srvc)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	router := handl.InitRouts(logger)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server error", slog.String("err", err.Error()))
			os.Exit(1)
		}
	}()

	logger.Info("VK-test-task Started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	sig := <-stop
	fmt.Printf("Received signal: %v\n", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", slog.String("err", err.Error()))
		os.Exit(1)
	}

	logger.Info("Server gracefully stopped")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewTextHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	case envDev:
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		slogHandler := slog.NewJSONHandler(os.Stdout, opts)
		logger = slog.New(slogHandler)
	}

	return logger
}
