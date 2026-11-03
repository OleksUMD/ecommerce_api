package server

import "github.com/gin-gonic/gin"

func (s *Server) setApiAndGraphRoutes(router *gin.Engine) {
	s.setUpGraphRoutes(router)
	s.setUpApiRoutes(router)
}

func (s *Server) graphPublicRoutes(router *gin.Engine) {
	// Public
	router.GET("/playground", s.graphQLPlaygroundHandler())

	router.GET("/playground/public", s.graphQLPlaygroundPublicHandler())
	graphqlPublic := router.Group("/graphql/public")
	graphqlPublic.Use(s.graphQLMiddleware())
	graphqlPublic.POST("/", s.graphQLHandler())
}

func (s *Server) graphPrivateRoutes(router *gin.Engine) {
	graphqlProtected := router.Group("/graphql")
	graphqlProtected = s.protectRoutes(graphqlProtected)
	graphqlProtected.Use(s.graphQLMiddleware())
	graphqlProtected.POST("/", s.graphQLHandler())
}

func (s *Server) setUpGraphRoutes(router *gin.Engine) {
	s.graphPublicRoutes(router)
	s.graphPrivateRoutes(router)
}

func (s *Server) protectRoutes(router *gin.RouterGroup) *gin.RouterGroup {
	router.Use(s.authMiddleware())
	return router
}

func (s *Server) setUpApiRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	s.setUpAuthRoutes(api)
	s.setUpApiPrivateRoutes(api)
	s.setUpApiPublicRoutes(api)
}

func (s *Server) setUpAuthRoutes(router *gin.RouterGroup) {
	// Public
	auth := router.Group("auth")
	auth.POST("/register", s.register)
	auth.POST("/login", s.login)
	auth.POST("/logout", s.logout)
	auth.POST("/refresh", s.refreshToken)
	auth.GET("/verify/:token", s.verifyEmail)
}

func (s *Server) setUpApiPrivateRoutes(router *gin.RouterGroup) {
	auth := router.Group("/")
	auth = s.protectRoutes(auth)
	s.setUpUserRoutes(auth)
	s.setUpCategoryRoutes(auth)
	s.setUpProductRoutes(auth)
	s.setUpCartRoutes(auth)
	s.setUpOrderRoutes(auth)
}

func (s *Server) setUpUserRoutes(router *gin.RouterGroup) {
	// Protected
	userRoutes := router.Group("users")
	userRoutes.GET("/profile", s.getProfile)
	userRoutes.PUT("/profile", s.updateProfile)
	userRoutes.POST("/verify-email", s.resendVerificationLink)
}

func (s *Server) setUpCategoryRoutes(router *gin.RouterGroup) {
	// Protected
	categoryRoutes := router.Group("categories")
	categoryRoutes.POST("/", s.adminMiddleware(), s.createCategory)
	categoryRoutes.PUT("/:id", s.adminMiddleware(), s.updateCategory)
	categoryRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteCategory)
}

func (s *Server) setUpProductRoutes(router *gin.RouterGroup) {
	// Protected
	productRoutes := router.Group("products")
	productRoutes.POST("/", s.adminMiddleware(), s.createProduct)
	productRoutes.PUT("/:id", s.adminMiddleware(), s.updateProduct)
	productRoutes.DELETE("/:id", s.adminMiddleware(), s.deleteProduct)
	productRoutes.POST("/:id/images", s.adminMiddleware(), s.uploadProductImage)
}

func (s *Server) setUpCartRoutes(router *gin.RouterGroup) {
	// Protected
	cartRoutes := router.Group("cart")
	cartRoutes.GET("/", s.getCart)
	cartRoutes.POST("/items", s.addToCart)
	cartRoutes.PUT("/items/:id", s.updateCartItem)
	cartRoutes.DELETE("/items/:id", s.removeFromCart)
}

func (s *Server) setUpOrderRoutes(router *gin.RouterGroup) {
	// Protected
	orderRoutes := router.Group("orders")
	orderRoutes.POST("/", s.createOrder)
	orderRoutes.GET("/", s.getOrders)
	orderRoutes.GET("/:id", s.getOrder)
}

func (s *Server) setUpApiPublicRoutes(router *gin.RouterGroup) {
	// Public
	router.GET("/categories", s.getCategories)
	router.GET("/search", s.searchProducts)
	router.GET("/products", s.getProducts)
	router.GET("/products/:id", s.getProduct)
}
