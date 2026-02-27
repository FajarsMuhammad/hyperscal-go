package dto

// PublishUserCreatedRequest is the request body for POST /api/kafka/user-created
type PublishUserCreatedRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
}

// PublishOrderPlacedRequest is the request body for POST /api/kafka/order-placed
type PublishOrderPlacedRequest struct {
	OrderID string  `json:"order_id" binding:"required"`
	UserID  uint    `json:"user_id" binding:"required"`
	Total   float64 `json:"total" binding:"required,gt=0"`
}
