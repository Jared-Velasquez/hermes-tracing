package graph

import (
	"context"
	"strings"
	"time"

	"analytics/store"
	"analytics/types"
)

type GraphBuilder struct {
	store *store.Store
}

func NewGraphBuilder(store *store.Store) *GraphBuilder {
	return &GraphBuilder{store: store}
}

func (gb *GraphBuilder) Build(ctx context.Context, start, end time.Time) (*types.ServiceGraph, error) {
	spans, err := gb.querySpans(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return gb.constructGraph(spans)
}

func (gb *GraphBuilder) querySpans(ctx context.Context, start, end time.Time) ([]types.Span, error) {
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"range": map[string]interface{}{
				"start_time": map[string]interface{}{
					"gte": start.UnixNano(),
					"lte": end.UnixNano(),
				},
			},
		},
		// Note: default parameters in Go?
		// "size": 10000,
	}

	return gb.store.QuerySpans(ctx, query)
}

func (gb *GraphBuilder) constructGraph(spans []types.Span) (*types.ServiceGraph, error) {
	// Note: once constructGraph is implemented I'll move these notes to the analytics README

	// Notes about traces and spans:
	
	// A single trace is made up of one or more spans

	// For synchronous calls (e.g. HTTP, gRPC):
	// Parent span represents the caller service, child span represents the callee service
	// Span Kinds: Client, Server, Internal

	// For asynchronous calls (e.g. message queues):
	// Parent span represents the producer service, child span represents the consumer service
	// Span Kinds: Producer, Consumer

	// Note: consumers most commonly have span links, but any span (internal, producer, server, etc.) can have links
	
	// Span Links may be used to link spans across (usually async) boundaries when the parent-child relationship is not sufficient. For
	// example, Producer A, B, and C all produce messages to a queue, and Consumer D consumes messages in batches from the queue.
	// The spans for Producer A, B, and C can be linked to the span for Consumer D using Span Links. Not certain but it seems like
	// a span can have both a parent span and span links; must consider that when constructing the graph

	// Notes about graph construction:
	// What happens if we process a span with a parent_span_id that we haven't seen yet?
	// What happens if we process a span with a parent_span_id that is not within the time range?

	nodesMap := make(map[string]types.ServiceNode)

	// edge: parent (caller/producer) -> child (callee/consumer); like Jaeger
	// Jaeger supports self-loops but that's redundant in service graphs
	// Support external services? (e.g. databases, third-party APIs)
	edgesMap := make(map[string]types.ServiceEdge)

	// TODO: implement graph construction
	_ = nodesMap
	_ = edgesMap

	return nil, nil
}

func extractDBTarget(attrs map[string]interface{}) types.DBTarget {
	// Extract DB system name and database from attributes:
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/db/

	systemName := ""
	databaseName := ""
	if dbSystem, ok := attrs["db.system.name"].(string); ok && dbSystem != "" { // Current
		systemName = dbSystem
	} else if dbSystem, ok := attrs["db.system"].(string); ok && dbSystem != "" { // Legacy/Deprecated
		systemName = dbSystem
	} else {
		return types.DBTarget{}
	}

	if dbName, ok := attrs["db.name"].(string); ok && dbName != "" { // Current
		databaseName = dbName
	}

	return types.DBTarget{System: systemName, DatabaseName: databaseName}
}

func extractMessagingTarget(attrs map[string]interface{}) types.MessagingTarget {
	// Extract messaging system name and destination from attributes:
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/messaging

	systemName := ""
	destinationName := ""
	if msgSystem, ok := attrs["messaging.system"].(string); ok && msgSystem != "" {
		systemName = msgSystem
	} else {
		return types.MessagingTarget{}
	}

	if destName, ok := attrs["messaging.destination.name"].(string); ok && destName != "" {
		destinationName = destName
	}

	return types.MessagingTarget{System: systemName, DestinationName: destinationName}
}

// Note: how to distinguish internal vs. external HTTP/RPC synchronous calls?
// Internal calls: connects two ServiceNodes within the system
// External calls: connects a ServiceNode to an external system (e.g. third-party API)

func extractHTTPTarget(attrs map[string]interface{}) types.SyncCallTarget {
	// Extract HTTP target from attributes:
	// HTTP: https://opentelemetry.io/docs/specs/semconv/registry/attributes/http/
	// Client (source): https://opentelemetry.io/docs/specs/semconv/registry/attributes/client/
	// Server (destination): https://opentelemetry.io/docs/specs/semconv/registry/attributes/server/

	// TODO: need both span and resource attributes to get service.name for source/destination?

	destination := ""
	source := ""

	if peer, ok := attrs["service.peer.name"].(string); ok && peer != "" { // Current
		destination = peer
	} else if peer, ok := attrs["peer.service"].(string); ok && peer != "" { // Legacy/deprecated
		destination = peer
	} else if serverAddr, ok := attrs["server.address"].(string); ok && serverAddr != "" { // Current (for extracting address)
		destination = serverAddr
	} else if httpHost, ok := attrs["http.host"].(string); ok && httpHost != "" { // Legacy/deprecated (for extracting address); TODO: docs don't make this clear if httpHost is client or server
		destination = httpHost
	}

	// Fallback: parse from URL (current: url.full, legacy: http.url)
	if destination == "" {
		urlStr := ""
		if u, ok := attrs["url.full"].(string); ok && u != "" {
			urlStr = u
		} else if u, ok := attrs["http.url"].(string); ok && u != "" {
			urlStr = u
		}
		if urlStr != "" {
			destination = extractHostFromURL(urlStr)
		}
	}

	if service, ok := attrs["service.name"].(string); ok && service != "" { // Current
		source = service
	} else if clientAddr, ok := attrs["client.address"].(string); ok && clientAddr != "" { // Current (for extracting address)
		source = clientAddr
	}

	if source == "" || destination == "" {
		return types.SyncCallTarget{}
	}

	return types.SyncCallTarget{
		Protocol:    "http",
		Source:      source,
		Destination: destination,
	}
}

func extractRPCTarget(attrs map[string]interface{}) types.SyncCallTarget {
	// Extract RPC target from attributes:
	// RPC: https://opentelemetry.io/docs/specs/semconv/registry/attributes/rpc/
	// Client (source): https://opentelemetry.io/docs/specs/semconv/registry/attributes/client/
	// Server (destination): https://opentelemetry.io/docs/specs/semconv/registry/attributes/server/

	// TODO: need both span and resource attributes to get service.name for source/destination?

	destination := ""
	source := ""
	protocol := ""

	// Extract RPC system (grpc, thrift, etc.)
	if rpcSystem, ok := attrs["rpc.system.name"].(string); ok && rpcSystem != "" { // Current
		protocol = rpcSystem
	} else if rpcSystem, ok := attrs["rpc.system"].(string); ok && rpcSystem != "" { // Legacy/deprecated
		protocol = rpcSystem
	} else {
		// No rpc.system means this isn't an RPC span
		return types.SyncCallTarget{}
	}

	// Extract service destination
	if peer, ok := attrs["peer.service"].(string); ok && peer != "" { // Current
		destination = peer
	} else if serverAddr, ok := attrs["server.address"].(string); ok && serverAddr != "" { // Current (for extracting address)
		destination = serverAddr
	}

	// Extract service source
	if service, ok := attrs["service.name"].(string); ok && service != "" { // Current
		source = service
	} else if clientAddr, ok := attrs["client.address"].(string); ok && clientAddr != "" { // Current (for extracting address)
		source = clientAddr
	}

	if source == "" || destination == "" {
		return types.SyncCallTarget{}
	}

	return types.SyncCallTarget{
		Protocol:    protocol,
		Source:      source,
		Destination: destination,
	}
}

// extractHostFromURL parses a URL string and returns the host (without port)
func extractHostFromURL(urlStr string) string {
	// Simple parsing without importing net/url to keep it lightweight
	// Handles: "https://api.stripe.com:443/v1/charges" -> "api.stripe.com"

	// Remove scheme
	if idx := strings.Index(urlStr, "://"); idx != -1 {
		urlStr = urlStr[idx+3:]
	}

	// Remove path
	if idx := strings.Index(urlStr, "/"); idx != -1 {
		urlStr = urlStr[:idx]
	}

	// Remove port
	if idx := strings.LastIndex(urlStr, ":"); idx != -1 {
		// Make sure it's a port, not part of IPv6
		if !strings.Contains(urlStr[idx:], "]") {
			urlStr = urlStr[:idx]
		}
	}

	// Remove userinfo if present
	if idx := strings.Index(urlStr, "@"); idx != -1 {
		urlStr = urlStr[idx+1:]
	}

	return urlStr
}