package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	otelcollogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otelcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	otelcoltrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

// Servers implement the OTLP gRPC services for traces, metrics, and logs
type TraceServer struct {
	otelcoltrace.UnimplementedTraceServiceServer
}

func (s *TraceServer) Export(ctx context.Context, req *otelcoltrace.ExportTraceServiceRequest) (*otelcoltrace.ExportTraceServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka traces topic
	// 4. Return successful response to client

	return &otelcoltrace.ExportTraceServiceResponse{}, nil
}

type MetricsServer struct {
	otelcolmetrics.UnimplementedMetricsServiceServer
}

func (s *MetricsServer) Export(ctx context.Context, req *otelcolmetrics.ExportMetricsServiceRequest) (*otelcolmetrics.ExportMetricsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka metrics topic
	// 4. Return successful response to client

	return &otelcolmetrics.ExportMetricsServiceResponse{}, nil
}

type LogsServer struct {
	otelcollogs.UnimplementedLogsServiceServer
}

func (s *LogsServer) Export(ctx context.Context, req *otelcollogs.ExportLogsServiceRequest) (*otelcollogs.ExportLogsServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka logs topic
	// 4. Return successful response to client

	return &otelcollogs.ExportLogsServiceResponse{}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":4317")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server and connect OTLP services
	grpcServer := grpc.NewServer()

	traceServer := TraceServer{}
	metricsServer := MetricsServer{}
	logsServer := LogsServer{}

	otelcoltrace.RegisterTraceServiceServer(grpcServer, &traceServer)
	otelcolmetrics.RegisterMetricsServiceServer(grpcServer, &metricsServer)
	otelcollogs.RegisterLogsServiceServer(grpcServer, &logsServer)

	log.Println("Starting OTLP gRPC server on :4317")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}