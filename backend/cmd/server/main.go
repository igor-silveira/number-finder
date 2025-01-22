package main

import (
	"fmt"
	"log/slog"
	"number-finder-api/internal/api"
	"number-finder-api/internal/config"
	"number-finder-api/internal/service"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger := setupLogger(cfg.LogLevel)

	finder, err := service.NewFinder(cfg.DataPath)

	if err != nil {
		logger.Error("Failed to initialize finder service", "error", err)
		os.Exit(1)
	}

	handler := api.NewHandler(finder, logger)
	server := api.NewServer(handler, logger)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err = server.Start(cfg.Port); err != nil {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	<-quit

	if err = server.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}

func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
