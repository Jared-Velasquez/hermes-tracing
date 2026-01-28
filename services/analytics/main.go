package main

import (
	"log"
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"analytics/consumer"
	"analytics/config"
)

func main() {
	// Cancel on interrupt signal
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg := config.LoadConfig()

	var wg sync.WaitGroup

	// Start Kafka trace consumer
	traceConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.TracesTopic,
		cfg.Kafka.GroupID,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		traceConsumer.Consume(ctx, consumer.HandleTrace)
	}()

	// Start Kafka metrics consumer
	metricsConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.MetricsTopic,
		cfg.Kafka.GroupID,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		metricsConsumer.Consume(ctx, consumer.HandleMetric)
	}()

	// Start Kafka logs consumer
	logsConsumer := consumer.NewConsumer(
		cfg.Kafka.Broker,
		cfg.Kafka.LogsTopic,
		cfg.Kafka.GroupID,
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		logsConsumer.Consume(ctx, consumer.HandleLog)
	}()

	// Block until context is cancelled
	<-ctx.Done()
	log.Println("Shutting down...")

	// Wait for all consumers to finish processing
	wg.Wait()

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
	log.Println("Shutdown complete")
}