package consumer

import (
	"log"

	otelcolmetrics "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/protobuf/proto"
)

func HandleMetric(data []byte) error {
	var req otelcolmetrics.ExportMetricsServiceRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		log.Printf("Failed to unmarshal metrics data: %v", err)
		return err
	}

	log.Printf("Received metrics: %v", &req)
	return nil
}
