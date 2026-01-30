package consumer

import (
	"context"
	"encoding/hex"
	"log"
	"analytics/store"

	commonv1 "go.opentelemetry.io/proto/otlp/common/v1"
	otelcoltrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	"google.golang.org/protobuf/proto"
)

type TraceHandler struct {
	store *store.Store
}

func NewTraceHandler(store *store.Store) *TraceHandler {
	return &TraceHandler{store: store}
}

func (h *TraceHandler) Handle(data []byte) error {
	var req otelcoltrace.ExportTraceServiceRequest
	if err := proto.Unmarshal(data, &req); err != nil {
		log.Printf("Failed to unmarshal trace data: %v", err)
		return err
	}

	ctx := context.Background()

	// ResourceSpans -> ScopeSpans -> Spans
	// Docs: https://pkg.go.dev/go.opentelemetry.io/proto/otlp@v1.9.0/collector/trace/v1
	for _, resourceSpan := range req.GetResourceSpans() {
		resource := resourceSpan.GetResource()

		for _, scopeSpan := range resourceSpan.GetScopeSpans() {
			scope := scopeSpan.GetScope()

			for _, span := range scopeSpan.GetSpans() {
				doc := map[string]interface{}{
					"trace_id":       hex.EncodeToString(span.GetTraceId()),
					"span_id":        hex.EncodeToString(span.GetSpanId()),
					"parent_span_id": hex.EncodeToString(span.GetParentSpanId()),
					"name":           span.GetName(),
					"kind":           span.GetKind().String(),
					"links": 		  func() []map[string]string {
						links := make([]map[string]string, 0, len(span.GetLinks()))
						for _, link := range span.GetLinks() {
							links = append(links, map[string]string{
								"trace_id": hex.EncodeToString(link.GetTraceId()),
								"span_id":  hex.EncodeToString(link.GetSpanId()),
							})
						}
						return links
					}(),
					"start_time":     span.GetStartTimeUnixNano(),
					"end_time":       span.GetEndTimeUnixNano(),
					"status":         span.GetStatus().GetCode().String(),

					// We need to use flattenMap for the attributes and resource fields since Elasticsearch
					// sometimes indexes attributes as objects but attributes returns as a string
					"attributes":     flattenMap(attributesToMap(span.GetAttributes()), ""),
					"resource":       flattenMap(attributesToMap(resource.GetAttributes()), ""),
					"scope_name":     scope.GetName(),
				}

				log.Printf("Indexing span: %v", doc)

				if err := h.store.IndexSpan(ctx, doc); err != nil {
					log.Printf("Failed to index span %s: %v", doc["span_id"], err)
				}
			}
		}
	}
	return nil
}

func flattenMap(m map[string]interface{}, prefix string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		if nested, ok := v.(map[string]interface{}); ok {
			for nk, nv := range flattenMap(nested, key) {
				result[nk] = nv
			}
		} else {
			result[key] = v
		}
	}
	return result
}

func attributesToMap(attrs []*commonv1.KeyValue) map[string]interface{} {
	m := make(map[string]interface{})
	for _, kv := range attrs {
		m[kv.GetKey()] = extractValue(kv.GetValue())
	}
	return m
}

func extractValue(v *commonv1.AnyValue) interface{} {
	if v == nil {
		return nil
	}
	switch val := v.GetValue().(type) {
	case *commonv1.AnyValue_StringValue:
		return val.StringValue
	case *commonv1.AnyValue_IntValue:
		return val.IntValue
	case *commonv1.AnyValue_DoubleValue:
		return val.DoubleValue
	case *commonv1.AnyValue_BoolValue:
		return val.BoolValue
	case *commonv1.AnyValue_ArrayValue:
		arr := make([]interface{}, 0, len(val.ArrayValue.GetValues()))
		for _, elem := range val.ArrayValue.GetValues() {
			arr = append(arr, extractValue(elem))
		}
		return arr
	case *commonv1.AnyValue_KvlistValue:
		return attributesToMap(val.KvlistValue.GetValues())
	default:
		return nil
	}
}
