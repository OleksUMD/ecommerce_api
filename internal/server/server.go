// Package server contains server http server information
package server

import (
	"net/http"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// Server is an app http server
type Server struct {
	config         *config.Config
	db             *gorm.DB
	logger         zerolog.Logger
	authService    *services.AuthService
	productService *services.ProductService
	userService    *services.UserService
	uploadService  *services.UploadService
	cartService    *services.CartService
}

// New creates and returns a new instance of app server
func New(
	cfg *config.Config,
	db *gorm.DB,
	logger zerolog.Logger,
	authService *services.AuthService,
	productService *services.ProductService,
	userService *services.UserService,
	uploadService *services.UploadService,
	cartService *services.CartService,
) *Server {
	return &Server{
		config:         cfg,
		db:             db,
		logger:         logger,
		authService:    authService,
		productService: productService,
		userService:    userService,
		uploadService:  uploadService,
		cartService:    cartService,
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
	router.Static("/uploads", "./uploads")

	api := router.Group("/api/v1")
	auth := api.Group("auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)
	auth.POST("/logout", s.logout)
	auth.POST("/refresh", s.refreshToken)

	// Private routes
	protected := api.Group("/")
	protected.Use(s.authMiddleware())

	userRoutes := protected.Group("users")
	userRoutes.GET("/profile", s.getProfile)
	userRoutes.PUT("/profile", s.updateProfile)

	categoryRoutes := protected.Group("categories")
	categoryRoutes.POST("/", s.adminMiddleware(), s.createCategory)
	categoryRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
	categoryRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)

	productRoutes := protected.Group("products")
	productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
	productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
	productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
	productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)

	cartRoutes := protected.Group("cart")
	cartRoutes.GET("/", s.getCart)
	cartRoutes.POST("/items", s.addToCart)
	cartRoutes.PUT("/items/:id", s.updateCartItem)
	cartRoutes.DELETE("/items/:id", s.removeFromCart)

	// Public routes
	api.GET("/categories", s.getCategories)
	api.GET("/products", s.getProducts)
	api.GET("/products/:id", s.getProduct)

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
