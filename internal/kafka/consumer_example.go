package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	kafkapkg "hyperscal-go/pkg/kafka"
)

// UserEventConsumer handles consuming user-related events from Kafka.
type UserEventConsumer struct {
	consumer *kafkapkg.Consumer
}

// NewUserEventConsumer creates a UserEventConsumer subscribed to the user.created topic.
//
// Usage:
//
//	kafkaClient := kafka.NewKafkaClient(&cfg.Kafka)
//	c, err := kafkainternal.NewUserEventConsumer(kafkaClient)
//	if err != nil { log.Fatal(err) }
//	defer c.Close()
//	go c.Start(ctx)
func NewUserEventConsumer(client *kafkapkg.KafkaClient) (*UserEventConsumer, error) {
	cons, err := client.NewConsumer(TopicUserCreated)
	if err != nil {
		return nil, fmt.Errorf("failed to create user event consumer: %w", err)
	}
	return &UserEventConsumer{
		consumer: cons,
	}, nil
}

// Start begins consuming messages in a blocking loop.
// Cancel ctx to stop gracefully.
func (c *UserEventConsumer) Start(ctx context.Context) error {
	return c.consumer.Consume(ctx, c.handleMessage)
}

// handleMessage is the internal handler that decodes the raw Kafka message
// and delegates to business logic.
func (c *UserEventConsumer) handleMessage(ctx context.Context, msg kafkapkg.Message) error {
	var event UserCreatedEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		log.Printf("[Kafka Consumer] Failed to unmarshal UserCreatedEvent: %v | raw=%s", err, string(msg.Value))
		return err
	}

	// --- Business logic goes here ---
	log.Printf("[Kafka Consumer] Received %s | user_id=%d email=%s",
		TopicUserCreated, event.UserID, event.Email)

	// Example: send welcome email, update read-model, notify downstream service, etc.
	return processUserCreated(ctx, event)
}

// processUserCreated contains the actual business logic after receiving the event.
// Replace this with your real implementation.
func processUserCreated(_ context.Context, event UserCreatedEvent) error {
	// TODO: implement real downstream logic (e.g., send welcome email)
	log.Printf("[Kafka Consumer] Processing user.created for user_id=%d email=%s",
		event.UserID, event.Email)
	return nil
}

// Close releases the underlying reader.
func (c *UserEventConsumer) Close() error {
	return c.consumer.Close()
}
