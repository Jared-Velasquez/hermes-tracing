package handlers

import (
	"log"
	"context"
	"ingest/config"

	otelcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/protobuf/proto"
)

type MetricsServer struct {
	otelcolmetrics.UnimplementedMetricsServiceServer
	config *config.Config
}

func (s *MetricsServer) Export(ctx context.Context, req *otelcolmetrics.ExportMetricsServiceRequest) (*otelcolmetrics.ExportMetricsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka metrics topic

	// Use protobuf serialization to marshal the request
	data, err := proto.Marshal(req)
	if err != nil {
		log.Printf("Failed to marshal metrics request: %v", err)
		return nil, err
	}

	// TODO: Should I use dependency injection for producer?
	producer := NewProducer(s.config.Kafka.Broker, s.config.Kafka.MetricsTopic)
	if err := producer.Send(data); err != nil {
		log.Printf("Failed to send metrics data to Kafka: %v", err)
		return nil, err
	}

	// 4. Return successful response to client
	log.Println("Metrics Export Request: ", req)
	return &otelcolmetrics.ExportMetricsServiceResponse{}, nil
}

func NewMetricsServer(config *config.Config) *MetricsServer {
	return &MetricsServer{
		config: config,
	}
}