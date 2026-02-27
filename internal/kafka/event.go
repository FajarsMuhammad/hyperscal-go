package kafka

// Topic names used across the application.
// Define all Kafka topic constants here so producers and consumers
// reference the same string without typos.
const (
	TopicUserCreated = "user.created"
	TopicOrderPlaced = "order.placed"
)

// UserCreatedEvent is published when a new user registers.
type UserCreatedEvent struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

// OrderPlacedEvent is published when an order is placed.
type OrderPlacedEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    uint    `json:"user_id"`
	Total     float64 `json:"total"`
	CreatedAt string  `json:"created_at"`
}
