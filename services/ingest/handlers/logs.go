package handlers

import (
	"log"
	"context"
	"ingest/config"

	otelcollogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/proto"
)

type LogsServer struct {
	otelcollogs.UnimplementedLogsServiceServer
	config *config.Config
}

func (s *LogsServer) Export(ctx context.Context, req *otelcollogs.ExportLogsServiceRequest) (*otelcollogs.ExportLogsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka logs topic

	// Use protobuf serialization to marshal the request
	data, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal logs request: %v", err)
		return nil, err
	}

	// TODO: Should I use dependency injection for producer?
	producer := NewProducer(s.config.Kafka.Broker, s.config.Kafka.LogsTopic)
	if err := producer.Send(data); err != nil {
		log.Printf("Failed to send logs data to Kafka: %v", err)
		return nil, err
	}

	// 4. Return successful response to client
	log.Println("Logs Export Request: ", req)
	return &otelcollogs.ExportLogsServiceResponse{}, nil
}

func NewLogsServer(config *config.Config) *LogsServer {
	return &LogsServer{
		config: config,
	}
}