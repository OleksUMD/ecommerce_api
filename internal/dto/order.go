package dto

// AddToCartRequest when the user adds a particular product to its cart
type AddToCartRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"min=1"`
}

// UpdateCartItemRequest when the user want to update something in the cart related to the particular product in it
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"min=1"`
}

// CartResponse when the user want to get info about the current cart
type CartResponse struct {
	ID        uint               `json:"id"`
	UserID    uint               `json:"user_id"`
	CartItems []CartItemResponse `json:"cart_items"`
	Total     float64            `json:"total"`
}

// CartItemResponse when the user want to get info about the current cart items
type CartItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Subtotal float64         `json:"subtotal"`
}

// OrderResponse when the user want to get info about the particular shop order
type OrderResponse struct {
	ID          uint                `json:"id"`
	UserID      uint                `json:"user_id"`
	Status      string              `json:"status"`
	TotalAmount float64             `json:"total_amount"`
	OrderItems  []OrderItemResponse `json:"order_items"`
	CreatedAt   string              `json:"created_at"`
}

// OrderItemResponse when the user want to get info about the particular shop order item
type OrderItemResponse struct {
	ID       uint            `json:"id"`
	Product  ProductResponse `json:"product"`
	Quantity int             `json:"quantity"`
	Price    float64         `json:"price"`
}
