package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"
	"vk-test-task/internal/configs"
	"vk-test-task/internal/http-server/handlers"
	"vk-test-task/internal/service"
	"vk-test-task/internal/storage"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {

	pg_cfg, err := configs.PgConfigFromEnv()
	logger := setupLogger(envLocal)

	db, err := storage.MustStorage(pg_cfg)
	if err != nil {
		logger.Info(err.Error())
		panic(err.Error())
	}

	srvc, err := service.MustService(db)
	if err != nil {
		panic(err.Error())
	}

	handl, err := handlers.MustHandlers(srvc)
	if err != nil {
		panic(err.Error())
	}

	router := handl.InitRouts(logger)

	server := http.Server{
		Addr:         ":8082",
		Handler:      router,
		ReadTimeout:  time.Millisecond * 50,
		WriteTimeout: time.Millisecond * 50,
	}

	server.ListenAndServe()
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
