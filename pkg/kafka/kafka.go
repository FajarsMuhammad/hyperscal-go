package kafka

import (
	"hyperscal-go/config"
	"time"
)

// KafkaClient is the main entry point to create producers and consumers.
// It wraps the config so all derived writers/readers share the same broker list.
type KafkaClient struct {
	cfg *config.KafkaConfig
}

// NewKafkaClient creates a new KafkaClient from the application config.
func NewKafkaClient(cfg *config.KafkaConfig) *KafkaClient {
	return &KafkaClient{cfg: cfg}
}

// Config returns the underlying Kafka configuration.
func (k *KafkaClient) Config() *config.KafkaConfig {
	return k.cfg
}

// NewProducer creates a new Producer for the given topic.
func (k *KafkaClient) NewProducer(topic string) (*Producer, error) {
	return newProducer(k.cfg.Brokers, topic)
}

// NewConsumer creates a new Consumer that subscribes to the given topic.
// It uses the GroupID from config to enable consumer group load-balancing.
func (k *KafkaClient) NewConsumer(topic string) (*Consumer, error) {
	return newConsumer(k.cfg.Brokers, topic, k.cfg.GroupID, k.cfg.AutoOffsetReset)
}

// dialTimeout is the default TCP dial timeout used by internal helpers.
const dialTimeout = 10 * time.Second
