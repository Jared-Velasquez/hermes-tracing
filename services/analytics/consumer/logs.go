package consumer

import (
	"log"

	otelcollogs "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/protobuf/proto"
)

func HandleLog(data []byte) error {
	var req otelcollogs.ExportLogsServiceRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		log.Printf("Failed to unmarshal logs data: %v", err)
		return err
	}

	log.Printf("Received logs: %v", &req)
	return nil
}
