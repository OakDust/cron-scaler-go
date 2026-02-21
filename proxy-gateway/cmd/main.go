// Package main - proxy-gateway API
// @title Cron Scaler Proxy Gateway API
// @version 1.0
// @description REST API для управления расписаниями масштабирования
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxy-gateway/controller"
	"proxy-gateway/internal/config"

	_ "proxy-gateway/docs" // swagger docs

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Настраиваем логгер
	logger := setupLogger()

	// Создаем контроллер
	ctrl, err := controller.NewController(logger, cfg.GRPCServerAddr)
	if err != nil {
		logger.Error("Failed to create controller", "error", err)
		os.Exit(1)
	}
	defer ctrl.Close()

	apiRouter := controller.NewRouter(ctrl)
	mux := http.NewServeMux()
	mux.Handle("/", apiRouter)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	logger.Info("Starting proxy-gateway",
		"port", cfg.HTTPPort,
		"grpc_addr", cfg.GRPCServerAddr)

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Запускаем сервер в горутине
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("Server stopped")
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}
