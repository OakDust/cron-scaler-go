package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"scale-handler/internal/app"
	"scale-handler/internal/config"
	"scale-handler/internal/controller"
	"scale-handler/internal/repository/postgres"
	"scale-handler/internal/usecase"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Настраиваем логгер
	logger := setupLogger()

	// Подключаемся к базе данных
	db, err := connectToDatabase(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("Connected to database",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"dbname", cfg.Database.DBName)

	// Инициализируем репозиторий
	scheduleRepo := postgres.NewScheduleRepository(db, logger)

	// Проверим подключение и таблицу
	if err := scheduleRepo.CheckConnection(context.Background()); err != nil {
		logger.Error("Database check failed", "error", err)
		os.Exit(1)
	}
	logger.Info("Database check passed, table exists")

	// Инициализируем usecase
	scheduleUC := usecase.NewScheduleUseCase(scheduleRepo, logger)

	// Инициализируем контроллер
	ctrl := controller.NewController(scheduleUC, logger)

	// Создаем gRPC сервер
	grpcServer, err := app.NewGRPCServer(cfg.GRPCPort, ctrl, logger)
	if err != nil {
		logger.Error("Failed to create gRPC server", "error", err)
		os.Exit(1)
	}

	// Запускаем gRPC сервер в горутине
	go func() {
		if err := grpcServer.Start(); err != nil {
			logger.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Scale-handler service started", "grpc_port", cfg.GRPCPort)

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down service...")

	// Graceful shutdown
	grpcServer.Stop()
	logger.Info("Service stopped gracefully")
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func connectToDatabase(dbConfig config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s search_path=public",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
