// This is a test package
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/database"
	"github.com/OleksUMD/ecommerce_api/internal/interfaces"
	"github.com/OleksUMD/ecommerce_api/internal/logger"
	"github.com/OleksUMD/ecommerce_api/internal/providers"
	"github.com/OleksUMD/ecommerce_api/internal/server"
	"github.com/OleksUMD/ecommerce_api/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	db, err := database.New(&cfg.Database)
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
	authService := services.NewAuthService(db, cfg)
	userService := services.NewUserService(db, cfg)
	productService := services.NewProductService(db, cfg)

	var uploadProvider interfaces.UploadProvider
	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}

	uploadService := services.NewUploadService(uploadProvider)
	cartService := services.NewCartService(db, cfg)
	orderService := services.NewOrderService(db, cfg)
	srv := server.New(
		cfg,
		db,
		log,
		authService,
		productService,
		userService,
		uploadService,
		cartService,
		orderService,
	)
	router := srv.SetupRoutes()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting the http server...")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info().Msg("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}
	log.Info().Msg("shutting down database")
}
