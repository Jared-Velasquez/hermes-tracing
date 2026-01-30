package graph

import (
	"analytics/types"
)

// OpenTelemetry semantic convention attribute names
// https://opentelemetry.io/docs/specs/semconv/
const (
	// Resource attributes
	AttrServiceName = "service.name"

	// Database attributes
	AttrDBSystem     = "db.system"
	AttrDBSystemName = "db.system.name" // Current
	AttrDBName       = "db.name"
	AttrDBNamespace  = "db.namespace" // Current

	// Messaging attributes
	AttrMessagingSystem      = "messaging.system"
	AttrMessagingDestination = "messaging.destination.name"

	// HTTP attributes
	AttrHTTPHost = "http.host" // Legacy
	AttrHTTPURL  = "http.url"  // Legacy
	AttrURLFull  = "url.full"  // Current

	// RPC attributes
	AttrRPCSystem = "rpc.system"

	// Network/address attributes
	AttrPeerService   = "peer.service"
	AttrServerAddress = "server.address"
	AttrClientAddress = "client.address"
)

// extractAttribute returns the first non-empty string value found for the given keys
func extractAttribute(m map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if val, ok := m[key].(string); ok && val != "" {
			return val
		}
	}
	return ""
}

func extractServiceName(span types.Span) string {
	// Extract service.name from resource attributes
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/service/

	// TODO: if service.name is not specified, SDK falls back to "unknown_service"
	// concatenated with process.executable.name (if available). Should I handle
	// this case differently?

	serviceName := extractAttribute(span.Attributes, AttrPeerService)
	return serviceName
}

func extractDBTarget(span types.Span) types.DBTarget {
	// Extract DB system name and database from attributes:
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/db/

	systemName := extractAttribute(span.Attributes, AttrDBSystemName, AttrDBSystem)
	if systemName == "" {
		return types.DBTarget{}
	}

	databaseName := extractAttribute(span.Attributes, AttrDBNamespace, AttrDBName)

	return types.DBTarget{System: systemName, DatabaseName: databaseName}
}

func extractMessagingTarget(span types.Span) types.MessagingTarget {
	// Extract messaging system name and destination from attributes:
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/messaging

	systemName := extractAttribute(span.Attributes, AttrMessagingSystem)
	if systemName == "" {
		return types.MessagingTarget{}
	}

	destinationName := extractAttribute(span.Attributes, AttrMessagingDestination)

	return types.MessagingTarget{System: systemName, DestinationName: destinationName}
}

// Note: how to distinguish internal vs. external HTTP/RPC synchronous calls?
// Internal calls: connects two ServiceNodes within the system
// External calls: connects a ServiceNode to an external system (e.g. third-party API)

func extractHTTPTarget(span types.Span) types.SyncCallTarget {
	// Extract HTTP target from attributes:
	// HTTP: https://opentelemetry.io/docs/specs/semconv/registry/attributes/http/
	// Client (source): https://opentelemetry.io/docs/specs/semconv/registry/attributes/client/
	// Server (destination): https://opentelemetry.io/docs/specs/semconv/registry/attributes/server/

	// Extract destination from span attributes
	destination := extractAttribute(span.Attributes, AttrPeerService, AttrServerAddress, AttrHTTPHost)

	// Fallback: parse from URL
	if destination == "" {
		if urlStr := extractAttribute(span.Attributes, AttrURLFull, AttrHTTPURL); urlStr != "" {
			destination = extractHostFromURL(urlStr)
		}
	}

	// Extract source: service.name from resource, fallback to client.address from span
	source := extractAttribute(span.Resource, AttrServiceName)
	if source == "" {
		source = extractAttribute(span.Attributes, AttrClientAddress)
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

func extractRPCTarget(span types.Span) types.SyncCallTarget {
	// Extract RPC target from attributes:
	// https://opentelemetry.io/docs/specs/semconv/registry/attributes/rpc/

	// Extract RPC system (grpc, thrift, etc.) - required
	protocol := extractAttribute(span.Attributes, AttrRPCSystem)
	if protocol == "" {
		return types.SyncCallTarget{}
	}

	// Extract destination from span attributes
	destination := extractAttribute(span.Attributes, AttrPeerService, AttrServerAddress)

	// Extract source: service.name from resource, fallback to client.address from span
	source := extractAttribute(span.Resource, AttrServiceName)
	if source == "" {
		source = extractAttribute(span.Attributes, AttrClientAddress)
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
