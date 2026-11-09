// Package server contains server http server information
package server

import (
	"net/http"

	_ "github.com/OleksUMD/ecommerce_api/docs"
	"github.com/OleksUMD/ecommerce_api/internal/cache"
	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Server is an app http server
type Server struct {
	config         *config.Config
	logger         zerolog.Logger
	cache          cache.ICacheProvider
	authService    services.IAuthService
	productService services.IProductService
	userService    services.IUserService
	uploadService  services.IUploadService
	cartService    services.ICartService
	orderService   services.IOrderService
}

// New creates and returns a new instance of app server
func New(
	cfg *config.Config,
	logger zerolog.Logger,
	authService services.IAuthService,
	productService services.IProductService,
	userService services.IUserService,
	uploadService services.IUploadService,
	cartService services.ICartService,
	orderService services.IOrderService,
	cache cache.ICacheProvider,
) *Server {
	return &Server{
		config:         cfg,
		logger:         logger,
		cache:          cache,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

// SetupRoutes initialize gin router
func (s *Server) SetupRoutes() *gin.Engine {
	router := gin.New()
	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleare())
	router.Use(s.rateLimitMiddleware(s.config.Server.RateLimit))
	router.GET("/health", s.healthCheck)

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFile("/api-docs", "./docs/rapidoc.html")

	router.Static("/uploads", "./uploads")

	s.setApiAndGraphRoutes(router)

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
