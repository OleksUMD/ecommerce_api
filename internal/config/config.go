// Package config contains all structs and functions for the app configuration
package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config is a main container for other type of settings
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

// ServerConfig is a container for server configuration
type ServerConfig struct {
	Port    string
	GinMode string
}

// DatabaseConfig is a container for database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// JWTConfig is a container for authentication settings
type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshtokenExpires time.Duration
}

// AWSConfig contains settings for AWS services
type AWSConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
}

// UploadConfig contains settings for file storage service
type UploadConfig struct {
	Path           string
	MaxFileSize    int64
	UploadProvider string
}

// Load get configuration from environment variables and create corresponding container(struct) for each of them
func Load() (*Config, error) {
	var err error

	_ = godotenv.Load()

	jwtExpiresIn, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		return nil, err
	}
	refreshTokenExpires, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "72h"))
	if err != nil {
		return nil, err
	}
	maxUploadSize, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	if err != nil {
		return nil, err
	}
	server := ServerConfig{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "8080"),
	}
	database := DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		Name:     getEnv("DB_NAME", "ecommerce_api"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	aws := AWSConfig{
		Region:          getEnv("AWS_REGION", "us-east-1"),
		AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", "test"),
		SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
		S3Bucket:        getEnv("AWS_S3_BUCKET", "ecommerce-uploads"),
		S3Endpoint:      getEnv("AWS_S3_ENDPOINT", "http://localhost:9000"),
	}

	jwt := JWTConfig{
		Secret:              getEnv("JWT_SECRET", "super-long-secret"),
		ExpiresIn:           jwtExpiresIn,
		RefreshtokenExpires: refreshTokenExpires,
	}
	upload := UploadConfig{
		Path:           getEnv("UPLOAD_PATH", "8080"),
		MaxFileSize:    maxUploadSize,
		UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
	}
	config := &Config{}
	config.Server = server
	config.Database = database
	config.AWS = aws
	config.JWT = jwt
	config.Upload = upload
	return config, err
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
