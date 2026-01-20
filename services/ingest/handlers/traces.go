package handlers

import (
	"log"
	"context"

	otelcoltrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type TraceServer struct {
	otelcoltrace.UnimplementedTraceServiceServer
}

func (s *TraceServer) Export(ctx context.Context, req *otelcoltrace.ExportTraceServiceRequest) (*otelcoltrace.ExportTraceServiceResponse, error) {
	// TODO: From ingestion pipeline:
	// 1. Validate incoming data
	// 2. Check Kafka backpressure and return codes.ResourceExhausted if overloaded
	// 3. Serialize and enqueue to Kafka traces topic
	// 4. Return successful response to client
	log.Println("Trace Export Request: ", req)

	return &otelcoltrace.ExportTraceServiceResponse{}, nil
}

func NewTraceServer() *TraceServer {
	return &TraceServer{}
}