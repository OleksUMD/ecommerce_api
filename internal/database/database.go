// Package database contains tools to setup the database connection
package database

import (
	"fmt"
	"time"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New initializes and configures the gorm database struct
func New(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode, time.UTC,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s", err)
	}
	return db, nil
}
