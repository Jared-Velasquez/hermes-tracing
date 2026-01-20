## Ingestion Server
The Ingestion Server is a gRPC server that implements OTLP over gRPC:

TraceService: `rpc Export(ExportTraceServiceRequest) returns (ExportTraceServiceResponse)`

MetricsService: `rpc Export(ExportMetricsServiceRequest) returns (ExportMetricsServiceResponse)`

LogsService: `rpc Export(ExportLogsServiceRequest) returns (ExportLogsServiceResponse)`

### Ingestion Pipeline
1. Receive Export gRPC request from OpenTelemetry Collector (configured to export to Hermes)
2. Validate request schema
3. Apply message queue backpressure; return gRPC errors if queue is overloaded
4. Enqueue to Kafka for downstream processing (persistent storage)
5. Return ExportResponse