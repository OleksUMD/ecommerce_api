package server

import (
	"strconv"

	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	cart, err := s.cartService.GetCart(userID)
	if err != nil {
		utils.NotFoundResponse(c, "Cart not found", nil)
		return
	}
	utils.SuccessResponse(c, "Cart retrieved successfully", cart)
}

func (s *Server) addToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	response, err := s.cartService.AddToCart(userID, &req)
	if err != nil {
		utils.BadRequestResponse(c, "Add product to cart failed", err)
		return
	}
	utils.SuccessResponse(c, "Product added to cart successfully", response)
}

func (s *Server) updateCartItem(c *gin.Context) {
	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "Invalid request data", err)
		return
	}
	userID := c.GetUint("user_id")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Update cart item failed", err)
		return
	}
	response, err := s.cartService.UpdateToCart(userID, uint(itemID), &req)
	if err != nil {
		utils.BadRequestResponse(c, "Update the cart item failed", err)
		return
	}
	utils.SuccessResponse(c, "The cart item updated successfully", response)
}

func (s *Server) removeFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Remove cart item failed", err)
		return
	}

	err = s.cartService.RemoveFromCart(userID, uint(itemID))
	if err != nil {
		utils.BadRequestResponse(c, "Failed to remove item from the cart", err)
		return
	}
	utils.SuccessNoContentResponse(c, "Cart item removed from the cart successfully", nil)
}
