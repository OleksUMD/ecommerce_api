// Package server contains server http server information
package server

import (
	"net/http"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Server is an app http server
type Server struct {
	config *config.Config
	db     *gorm.DB
	logger zerolog.Logger
}

// New creates and returns a new instance of app server
func New(cfg *config.Config, db *gorm.DB, logger zerolog.Logger) *Server {
	return &Server{
		config: cfg,
		db:     db,
		logger: logger,
	}
}

// SetupRoutes initialize gin router
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleare())
	router.GET("/health", s.healthCheck)

	api := router.Group("/api/v1")
	auth := api.Group("auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)
	auth.POST("/logout", s.logout)
	auth.POST("/refresh", s.refreshToken)
	userRoutes := api.Group("users")
	userRoutes.Use(s.authMiddleware())
	userRoutes.GET("/profile", s.getProfile)
	userRoutes.PUT("/profile", s.updateProfile)

	return router
}

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) corsMiddleare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-PIN")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
