package main

import (
	"github.com/Egorpalan/medods-api/config"
	"github.com/Egorpalan/medods-api/pkg/db"
	"github.com/Egorpalan/medods-api/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitLogger()
	defer logger.Log.Sync()

	database := db.NewDB(cfg)
	defer database.Close()

	logger.Log.Info("Service started successfully!")
	// Дальше — инициализация роутеров, DI, запуск chi и т.д.
}
