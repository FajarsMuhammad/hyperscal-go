package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Message is a simplified representation of a Kafka message.
type Message struct {
	Topic     string
	Partition int32
	Offset    int64
	Key       []byte
	Value     []byte
}

// Handler is a function that processes a single Kafka message.
// Return an error to log the failure; the consumer will continue to the next message.
type Handler func(ctx context.Context, msg Message) error

// Consumer wraps kgo.Client to consume messages from a single topic.
type Consumer struct {
	client *kgo.Client
	topic  string
}

// newConsumer creates an internal Consumer. Use KafkaClient.NewConsumer instead.
func newConsumer(brokers []string, topic, groupID, autoOffsetReset string) (*Consumer, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(brokers...),
		kgo.ConsumeTopics(topic),
		kgo.ConsumerGroup(groupID),
		kgo.DisableAutoCommit(), // We will commit manually after processing
	}

	if autoOffsetReset == "latest" {
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtEnd()))
	} else {
		// earliest is the default if not specified
		opts = append(opts, kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client for consumer: %w", err)
	}

	return &Consumer{
		client: client,
		topic:  topic,
	}, nil
}

// Consume starts a blocking loop that reads messages and calls handler for each one.
// It commits the offset automatically after the handler returns (even on error).
// Cancel ctx to stop consuming gracefully.
func (c *Consumer) Consume(ctx context.Context, handler Handler) error {
	log.Printf("[Kafka Consumer] Starting consumer on topic: %s", c.topic)
	for {
		fetches := c.client.PollFetches(ctx)
		if errMsgs := fetches.Errors(); len(errMsgs) > 0 {
			if ctx.Err() != nil {
				log.Println("[Kafka Consumer] Context cancelled, shutting down...")
				return nil
			}
			// Just log poll/fetch errors and continue retrying
			for _, err := range errMsgs {
				log.Printf("[Kafka Consumer] fetch error from topic %s, partition %d: %v",
					err.Topic, err.Partition, err.Err)
			}
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			msg := Message{
				Topic:     record.Topic,
				Partition: record.Partition,
				Offset:    record.Offset,
				Key:       record.Key,
				Value:     record.Value,
			}

			if err := handler(ctx, msg); err != nil {
				log.Printf("[Kafka Consumer] Handler error on topic=%s partition=%d offset=%d: %v",
					record.Topic, record.Partition, record.Offset, err)
			}

			// Commit offset regardless of handler success so we don't re-process on crash
			if err := c.client.CommitRecords(ctx, record); err != nil {
				log.Printf("[Kafka Consumer] Failed to commit offset: %v", err)
			}
		}
	}
}

// Close releases the client connection.
func (c *Consumer) Close() error {
	c.client.Close()
	return nil
}
