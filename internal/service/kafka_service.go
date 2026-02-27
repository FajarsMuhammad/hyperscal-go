package service

import (
	"context"
	"fmt"
	"time"

	kafkainternal "hyperscal-go/internal/kafka"
	kafkapkg "hyperscal-go/pkg/kafka"
)

// KafkaService exposes Kafka publish operations for use by controllers.
type KafkaService struct {
	userProducer  *kafkainternal.UserEventProducer
	orderProducer *kafkapkg.Producer
}

// NewKafkaService wires up producers for all supported topics.
func NewKafkaService(client *kafkapkg.KafkaClient) (*KafkaService, error) {
	userProducer, err := kafkainternal.NewUserEventProducer(client)
	if err != nil {
		return nil, fmt.Errorf("kafka service: %w", err)
	}

	orderProducer, err := client.NewProducer(kafkainternal.TopicOrderPlaced)
	if err != nil {
		return nil, fmt.Errorf("kafka service: %w", err)
	}

	return &KafkaService{
		userProducer:  userProducer,
		orderProducer: orderProducer,
	}, nil
}

// PublishUserCreated publishes a user.created event.
func (s *KafkaService) PublishUserCreated(ctx context.Context, userID uint, email string) error {
	event := kafkainternal.UserCreatedEvent{
		UserID:    userID,
		Email:     email,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	return s.userProducer.PublishUserCreated(ctx, event)
}

// PublishOrderPlaced publishes an order.placed event.
func (s *KafkaService) PublishOrderPlaced(ctx context.Context, orderID string, userID uint, total float64) error {
	event := kafkainternal.OrderPlacedEvent{
		OrderID:   orderID,
		UserID:    userID,
		Total:     total,
		CreatedAt: time.Now().Format(time.RFC3339),
	}
	return s.orderProducer.PublishJSON(ctx, orderID, event)
}

// Close releases all underlying producer connections.
func (s *KafkaService) Close() {
	s.userProducer.Close()
	s.orderProducer.Close()
}
