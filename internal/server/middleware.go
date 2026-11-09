package server

import (
	"strings"

	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header required", nil)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.config.JWT.Secret)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid authorization header format", err)
			c.Abort()
			return
		}
		if !s.userService.UserIsVerified(claims.UserID) {
			utils.UnauthorizedResponse(c, "User verification failed", err)
			c.Abort()
			return
		}

		c.Set(utils.UserIDKey, claims.UserID)
		c.Set(utils.UserEmailKey, claims.Email)
		c.Set(utils.UserRoleKey, claims.Role)
		c.Next()
	}
}

func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(utils.UserRoleKey)
		if !exists {
			utils.ForbiddenResponse(c, "Forbidden", nil)
			c.Abort()
			return
		}

		if role.(string) != string(models.UserRoleAdmin) {
			utils.ForbiddenResponse(c, "Forbidden", nil)
			c.Abort()
			return
		}
		userID, exists := c.Get(utils.UserIDKey)
		if !exists {
			utils.ForbiddenResponse(c, "Forbidden", nil)
			c.Abort()
			return
		}
		if !s.userService.UserIsVerified(userID.(uint)) {
			utils.UnauthorizedResponse(c, "User verification failed", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Server) rateLimitMiddleware(limit int) gin.HandlerFunc {
	if limit <= 0 {
		limit = 10
	}
	var limiter = rate.NewLimiter(rate.Limit(limit), 100)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			utils.TooManyRequestResponse(c, "Too many requests")
			c.Abort()
			return
		}
		c.Next()
	}
}
