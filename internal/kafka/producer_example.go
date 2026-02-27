package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	kafkapkg "hyperscal-go/pkg/kafka"
)

// UserEventProducer handles publishing user-related events to Kafka.
type UserEventProducer struct {
	producer *kafkapkg.Producer
}

// NewUserEventProducer creates a UserEventProducer for the user.created topic.
//
// Usage:
//
//	kafkaClient := kafka.NewKafkaClient(&cfg.Kafka)
//	p, err := kafkainternal.NewUserEventProducer(kafkaClient)
//	if err != nil { log.Fatal(err) }
//	defer p.Close()
func NewUserEventProducer(client *kafkapkg.KafkaClient) (*UserEventProducer, error) {
	prod, err := client.NewProducer(TopicUserCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create user event producer: %w", err)
	}
	return &UserEventProducer{
		producer: prod,
	}, nil
}

// PublishUserCreated publishes a UserCreatedEvent.
// key is used for partition routing – typically the user ID string.
func (p *UserEventProducer) PublishUserCreated(ctx context.Context, event UserCreatedEvent) error {
	key := fmt.Sprintf("%d", event.UserID)
	if err := p.producer.PublishJSON(ctx, key, event); err != nil {
		return fmt.Errorf("UserEventProducer: %w", err)
	}
	log.Printf("[Kafka Producer] Published %s for user_id=%d", TopicUserCreated, event.UserID)
	return nil
}

// Close releases the underlying writer.
func (p *UserEventProducer) Close() error {
	return p.producer.Close()
}

// --- Example standalone function (for demonstration / quick testing) ---

// ExamplePublishUserCreated shows how to publish a single event
// without the wrapper struct. Useful as a reference when wiring into a service.
func ExamplePublishUserCreated(client *kafkapkg.KafkaClient, userID uint, email string) error {
	p, err := client.NewProducer(TopicUserCreated)
	if err != nil {
		return fmt.Errorf("failed to create producer: %w", err)
	}
	defer p.Close()

	event := UserCreatedEvent{
		UserID:    userID,
		Email:     email,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	return p.PublishJSON(context.Background(), fmt.Sprintf("%d", userID), event)
}
