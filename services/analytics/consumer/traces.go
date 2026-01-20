package consumer

import (
	"log"

	otelcoltrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
)

func HandleTrace(data []byte) error {
	var req otelcoltrace.ExportTraceServiceRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		log.Printf("Failed to unmarshal trace data: %v", err)
		return err
	}

	log.Printf("Received trace: %v", &req)
	return nil
}
