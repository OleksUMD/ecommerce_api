package server

import (
	"strconv"

	_ "github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

// @Summary Get user's orders
// @Description Retrieve paginated list of user's orders
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} utils.PaginatedResponse{data=[]dto.OrderResponse} "Orders retrieved successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal server error"
// @Router /orders [get]
func (s *Server) getOrders(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		utils.BadRequestResponse(c, "Get orders failed", err)
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		utils.BadRequestResponse(c, "Get orders failed", err)
		return
	}
	userID := c.GetUint("user_id")
	response, meta, err := s.orderService.GetOrders(userID, page, limit)
	if err != nil {
		utils.BadRequestResponse(c, "Get orders failed", err)
		return
	}
	utils.PaginatedSuccessResponse(c, "Orders retrieved successfully", response, *meta)
}

// @Summary Get order by ID
// @Description Retrieve detailed information about a specific order
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Param id path int true "Order ID"
// @Success 200 {object} utils.Response{data=dto.OrderResponse} "Order retrieved successfully"
// @Failure 400 {object} utils.Response "Invalid order ID"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "Order not found"
// @Router /orders/{id} [get]
func (s *Server) getOrder(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		utils.BadRequestResponse(c, "Get order failed", err)
		return
	}
	userID := c.GetUint("user_id")
	response, err := s.orderService.GetOrder(userID, uint(orderID))
	if err != nil {
		utils.BadRequestResponse(c, "Get order failed", err)
		return
	}
	utils.SuccessResponse(c, "Order retrieved successfully", response)
}

// @Summary Create an order
// @Description Create an order from the current user's cart
// @Tags Orders
// @Produce json
// @Security BearerAuth
// @Success 201 {object} utils.Response{data=dto.OrderResponse} "Order created successfully"
// @Failure 400 {object} utils.Response "Cart is empty or insufficient stock"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /orders [post]
func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create an order", err)
		return
	}
	utils.CreatedResponse(c, "Order created successfully", order)
}
