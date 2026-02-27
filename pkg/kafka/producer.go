package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Producer wraps kgo.Client to publish messages to a single topic.
type Producer struct {
	client *kgo.Client
	topic  string
}

// newProducer creates an internal Producer. Use KafkaClient.NewProducer instead.
func newProducer(brokers []string, topic string) (*Producer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ProduceRequestTimeout(10 * time.Second),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client for producer: %w", err)
	}

	return &Producer{
		client: client,
		topic:  topic,
	}, nil
}

// Publish sends a raw byte message to Kafka.
// Set key to nil if you don't need key-based partitioning.
func (p *Producer) Publish(ctx context.Context, key, value []byte) error {
	msg := &kgo.Record{
		Topic: p.topic,
		Key:   key,
		Value: value,
	}

	// ProduceSync blocks until the message is successfully produced or fails
	res := p.client.ProduceSync(ctx, msg)
	if err := res.FirstErr(); err != nil {
		return fmt.Errorf("kafka producer: failed to write message: %w", err)
	}
	return nil
}

// PublishJSON marshals payload to JSON and publishes it.
// key is used for partition routing (e.g., user ID, order ID).
func (p *Producer) PublishJSON(ctx context.Context, key string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("kafka producer: failed to marshal payload: %w", err)
	}
	var kBytes []byte
	if key != "" {
		kBytes = []byte(key)
	}
	return p.Publish(ctx, kBytes, data)
}

// Close flushes pending writes and releases the client connection.
func (p *Producer) Close() error {
	p.client.Close()
	return nil
}
