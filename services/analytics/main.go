package main

import (
	"log"
	"context"
	"os"
	"os/signal"
	"syscall"
	"analytics/consumer"
	"analytics/config"
)

func main() {
	// Cancel on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.LoadConfig()

	// Start Kafka trace consumer
	traceConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.TracesTopic,
		cfg.Kafka.GroupID,
	)
	go traceConsumer.Consume(ctx, consumer.HandleTrace)

	// Start Kafka metrics consumer
	metricsConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.MetricsTopic,
		cfg.Kafka.GroupID,
	)
	go metricsConsumer.Consume(ctx, consumer.HandleMetric)

	// Start Kafka logs consumer
	logsConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.LogsTopic,
		cfg.Kafka.GroupID,
	)
	go logsConsumer.Consume(ctx, consumer.HandleLog)

	// Block until context is cancelled
	<-ctx.Done()
	
	// Close consumers
	if err := traceConsumer.Close(); err != nil {
		log.Printf("Error closing trace consumer: %v", err)
	}
	if err := metricsConsumer.Close(); err != nil {
		log.Printf("Error closing metrics consumer: %v", err)
	}
	if err := logsConsumer.Close(); err != nil {
		log.Printf("Error closing logs consumer: %v", err)
	}
}