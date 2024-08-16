package common

import (
	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(broker, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
