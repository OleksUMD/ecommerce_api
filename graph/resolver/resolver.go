package resolver

import (
	"strconv"

	"github.com/OleksUMD/ecommerce_api/internal/services"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require
// here.

type Resolver struct {
	authService    services.IAuthService
	userService    services.IUserService
	productService services.IProductService
	cartService    services.ICartService
	orderService   services.IOrderService
}

func NewResolver(
	authService services.IAuthService,
	userService services.IUserService,
	productService services.IProductService,
	cartService services.ICartService,
	orderService services.IOrderService,
) *Resolver {
	return &Resolver{
		authService:    authService,
		userService:    userService,
		productService: productService,
		cartService:    cartService,
		orderService:   orderService,
	}
}

func (r *Resolver) parseID(id string) (uint, error) {
	parsed, err := strconv.ParseUint(id, 10, 32)
	return uint(parsed), err
}
