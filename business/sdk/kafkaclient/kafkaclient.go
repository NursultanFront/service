// Package kafkaclient provides support for accessing Kafka.
package kafkaclient

import (
	"github.com/segmentio/kafka-go"
)

// Config is the required properties to use Kafka.
type Config struct {
	Brokers []string
}

// NewWriter constructs a Kafka writer (producer) for the specified topic.
// The caller is responsible for calling Close() on the returned writer.
func NewWriter(cfg Config, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

// NewReader constructs a Kafka reader (consumer) for the specified topic and
// consumer group. The caller is responsible for calling Close() on the
// returned reader.
func NewReader(cfg Config, topic string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   topic,
		GroupID: groupID,
	})
}
