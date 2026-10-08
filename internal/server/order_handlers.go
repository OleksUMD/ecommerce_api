package server

import (
	"strconv"

	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"github.com/gin-gonic/gin"
)

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

func (s *Server) createOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	order, err := s.orderService.CreateOrder(userID)
	if err != nil {
		utils.BadRequestResponse(c, "Failed to create an order", err)
		return
	}
	utils.CreatedResponse(c, "Order created successfully", order)
}
