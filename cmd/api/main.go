// This is a test package
package main

import (
	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/database"
	"github.com/OleksUMD/ecommerce_api/internal/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}
	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer func() {
		if err := mainDB.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close database connection")
		}
	}()
	gin.SetMode(cfg.Server.GinMode)
	log.Info().Msg("starting server")
}
