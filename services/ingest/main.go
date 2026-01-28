package main

import (
	"log"
	"net"
	"ingest/handlers"
	"ingest/config"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip" // Register gzip compressor
	otelcollogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	otelcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	otelcoltrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

// Note: likely that we will only need traces to construct call graph

func main() {
	cfg := config.LoadConfig()

	listener, err := net.Listen("tcp", ":4317")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a new gRPC server and connect OTLP services
	grpcServer := grpc.NewServer()

	if cfg.EnableTraces {
		traceServer := handlers.NewTraceServer(cfg)
		otelcoltrace.RegisterTraceServiceServer(grpcServer, traceServer)
		log.Println("Traces endpoint enabled")
	}

	if cfg.EnableMetrics {
		metricsServer := handlers.NewMetricsServer(cfg)
		otelcolmetrics.RegisterMetricsServiceServer(grpcServer, metricsServer)
		log.Println("Metrics endpoint enabled")
	}

	if cfg.EnableLogs {
		logsServer := handlers.NewLogsServer(cfg)
		otelcollogs.RegisterLogsServiceServer(grpcServer, logsServer)
		log.Println("Logs endpoint enabled")
	}

	log.Println("Starting OTLP gRPC server on :4317")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}