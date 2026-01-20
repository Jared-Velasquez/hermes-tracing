package handlers

import (
	"log"
	"context"

	otelcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
)

type MetricsServer struct {
	otelcolmetrics.UnimplementedMetricsServiceServer
}

func (s *MetricsServer) Export(ctx context.Context, req *otelcolmetrics.ExportMetricsServiceRequest) (*otelcolmetrics.ExportMetricsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka metrics topic
	// 4. Return successful response to client
	log.Println("Metrics Export Request: ", req)
	return &otelcolmetrics.ExportMetricsServiceResponse{}, nil
}

func NewMetricsServer() *MetricsServer {
	return &MetricsServer{}
}