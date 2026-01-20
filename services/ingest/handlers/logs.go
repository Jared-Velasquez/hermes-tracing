package handlers

import (
	"log"
	"context"

	otelcollogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
)

type LogsServer struct {
	otelcollogs.UnimplementedLogsServiceServer
}

func (s *LogsServer) Export(ctx context.Context, req *otelcollogs.ExportLogsServiceRequest) (*otelcollogs.ExportLogsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka logs topic
	// 4. Return successful response to client
	log.Println("Logs Export Request: ", req)
	return &otelcollogs.ExportLogsServiceResponse{}, nil
}

func NewLogsServer() *LogsServer {
	return &LogsServer{}
}