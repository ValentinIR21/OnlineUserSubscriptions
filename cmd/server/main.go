package main

import (
	"context"
	"log/slog"
	"net/http"
	"onlineusersub/internal/handler"
	"onlineusersub/internal/repository"
	"onlineusersub/internal/service"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Подключение к PostgreSQL
	postgresURL := getEnv("DB_URL", "postgres://postgres:pass@postgres:5432/usersubscriptions")

	repo, err := repository.NewPostgresRep(ctx, postgresURL)
	if err != nil {
		slog.Error("(main) Не удалось подключиться к PostgreSQL", "err", err)
		os.Exit(1)
	}
	defer repo.Close()

	slog.Info("(main) PostgreSQL подключен")

	// Service
	subService := service.NewSubService(repo)

	// Handler
	h := handler.NewSubHandler(subService)
	router := h.Routes()

	// http server
	addr := getEnv("HTTP_ADDR", ":8081")

	server := http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	slog.Info("Сервер запускается", "addr", addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("(main) ошибка запуска сервера", "err", err)
			os.Exit(1)
		}
	}()

	// Ожидания сигнала завершения
	<-ctx.Done()
	slog.Info("Получен сигнал завершения")

	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Ошибка при остановке сервера")
		os.Exit(1)
	}

	slog.Info("Сервер остановлен")
}

// Чтение перем окружения
func getEnv(value, valueDefualt string) string {

	if val := os.Getenv(value); val != "" {
		return val
	}

	return valueDefualt
}
