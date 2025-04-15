package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Egorpalan/medods-api/config"
	"github.com/Egorpalan/medods-api/internal/handler"
	"github.com/Egorpalan/medods-api/internal/repository"
	"github.com/Egorpalan/medods-api/internal/usecase"
	"github.com/Egorpalan/medods-api/pkg/db"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func main() {
	// Инициализация конфига и логгера
	cfg := config.LoadConfig()
	logger.InitLogger()
	defer logger.Log.Sync()

	// Подключение к БД
	database := db.NewDB(cfg)
	defer database.Close()

	// DI: репозиторий, usecase, handler
	refreshTokenRepo := repository.NewRefreshTokenRepo(database)
	authUsecase := usecase.NewAuthUsecase(refreshTokenRepo, []byte(cfg.JWTSecret))
	authHandler := handler.NewAuthHandler(authUsecase)

	// Роутер
	r := chi.NewRouter()
	r.Post("/token", authHandler.GenerateTokenPair)
	r.Post("/refresh", authHandler.RefreshTokenPair)

	// HTTP сервер
	server := &http.Server{
		Addr:    ":" + cfg.AppPort,
		Handler: r,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		logger.Log.Info("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Log.Error("Server forced to shutdown", zap.Error(err))
		}
	}()

	logger.Log.Info("Starting server...", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Fatal("ListenAndServe error", zap.Error(err))
	}
}
