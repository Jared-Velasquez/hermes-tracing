package consumer

import (
	"context"
	"log"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
	}
}

func (c *Consumer) Consume(ctx context.Context, handler func([]byte) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			log.Printf("Failed to read message: %v", err)
			time.Sleep(time.Second)
			continue
		}

		if err := handler(msg.Value); err != nil {
			log.Printf("Failed to handle message: %v", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
