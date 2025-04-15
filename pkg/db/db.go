package db

import (
	"fmt"
	"github.com/Egorpalan/medods-api/config"
	"github.com/Egorpalan/medods-api/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func NewDB(cfg *config.Config) *sqlx.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Log.Fatal("NewDB error", zap.Error(err))
	}
	return db
}
