package services

import (
	"errors"
	"fmt"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/dto"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/utils"
	"gorm.io/gorm"
)

const defaultDateFormat = "2006-01-02T15:04:05Z"

type OrderService struct {
	db     *gorm.DB
	config *config.Config
}

func NewOrderService(db *gorm.DB, config *config.Config) *OrderService {
	return &OrderService{
		db:     db,
		config: config,
	}
}

func (s *OrderService) CreateOrder(userID uint) (*dto.OrderResponse, error) {
	var orderResponse *dto.OrderResponse
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var cart models.Cart
		if err := tx.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
			return errors.New("cart not found")
		}
		if len(cart.CartItems) == 0 {
			return errors.New("cart is empty")
		}
		var totalAmount float64
		var orderItems []models.OrderItem
		for i := range cart.CartItems {
			cartItem := &cart.CartItems[i]
			if cartItem.Product.Stock < cartItem.Quantity {
				return fmt.Errorf("insufficient stock for product: %s", cartItem.Product.Name)
			}
			itemTotal := float64(cartItem.Quantity) * cartItem.Product.Price
			totalAmount += itemTotal
			orderItems = append(orderItems, models.OrderItem{
				ProductID: cartItem.ProductID,
				Quantity:  cartItem.Quantity,
				Price:     cartItem.Product.Price,
			})
			cartItem.Product.Stock -= cartItem.Quantity
			if err := tx.Save(&cartItem.Product).Error; err != nil {
				return err
			}
		}

		order := models.Order{
			UserID:      userID,
			Status:      models.OrderStatusPending,
			TotalAmount: totalAmount,
			OrderItems:  orderItems,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			return err
		}

		response, err := s.getOrderResponse(tx, order.ID)
		if err != nil {
			return err
		}
		orderResponse = response

		return nil
	})
	if err != nil {
		return nil, err
	}

	return orderResponse, nil
}

func (s *OrderService) GetOrders(userID uint, page, limit int) ([]dto.OrderResponse, *utils.PaginationMeta, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit
	var orders []models.Order
	var total int64
	s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)

	if err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).Limit(limit).
		Find(&orders).Error; err != nil {
		return nil, nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i := range orders {
		response[i] = s.convertToOrderResponse(&orders[i])
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
	return response, meta, nil
}

func (s *OrderService) GetOrder(userID, orderID uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := s.db.Preload("OrderItems.Product.Category").
		Where("user_id = ? AND id = ?", userID, orderID).
		First(&order).Error; err != nil {
		return nil, err
	}
	response := s.convertToOrderResponse(&order)
	return &response, nil
}

func (s *OrderService) getOrderResponse(tx *gorm.DB, id uint) (*dto.OrderResponse, error) {
	var order models.Order
	if err := tx.Preload("OrderItems.Product.Category").First(&order, id).Error; err != nil {
		return nil, err
	}

	response := s.convertToOrderResponse(&order)
	return &response, nil
}

func (s *OrderService) convertToOrderResponse(order *models.Order) dto.OrderResponse {
	items := make([]dto.OrderItemResponse, len(order.OrderItems))
	for i := range order.OrderItems {
		items[i] = dto.OrderItemResponse{
			ID:       order.OrderItems[i].ID,
			Quantity: order.OrderItems[i].Quantity,
			Price:    order.OrderItems[i].Price,
			Product: dto.ProductResponse{
				ID:          order.OrderItems[i].Product.ID,
				CategoryID:  order.OrderItems[i].Product.CategoryID,
				Name:        order.OrderItems[i].Product.Name,
				Description: order.OrderItems[i].Product.Description,
				Price:       order.OrderItems[i].Product.Price,
				Stock:       order.OrderItems[i].Product.Stock,
				SKU:         order.OrderItems[i].Product.SKU,
				IsActive:    order.OrderItems[i].Product.IsActive,
				Category: dto.CategoryResponse{
					ID:          order.OrderItems[i].Product.Category.ID,
					Name:        order.OrderItems[i].Product.Category.Name,
					Description: order.OrderItems[i].Product.Category.Description,
					IsActive:    order.OrderItems[i].Product.Category.IsActive,
				},
			},
		}
	}
	return dto.OrderResponse{
		ID:          order.ID,
		UserID:      order.UserID,
		OrderItems:  items,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Format(defaultDateFormat),
	}
}
