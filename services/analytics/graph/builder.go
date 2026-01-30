package graph

import (
	"context"
	"strings"
	"time"

	"analytics/store"
	"analytics/types"

	"github.com/google/uuid"
)

const (
	KeyDelimiter = "->"
	SpanKindInternal = "SPAN_KIND_INTERNAL"
	SpanKindClient   = "SPAN_KIND_CLIENT"
	SpanKindServer   = "SPAN_KIND_SERVER"
	SpanKindProducer = "SPAN_KIND_PRODUCER"
	SpanKindConsumer = "SPAN_KIND_CONSUMER"
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

	nodesMap, edgesMap, err := gb.constructGraph(spans)
	if err != nil {
		return nil, err
	}

	return &types.ServiceGraph{
		GraphID:   uuid.New().String(),
		StartTime: start,
		EndTime:   end,
		Nodes: nodesMap,
		Edges: edgesMap,
		CreatedAt: time.Now(),
	}, nil
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

func (gb *GraphBuilder) constructGraph(spans []types.Span) (nodesMap map[string]types.ServiceNode, edgesMap map[string]types.ServiceEdge, err error) {
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

	spanMap := make(map[string]types.Span)
	for _, span := range spans {
		spanMap[span.SpanID] = span
	}

	nodesMap := make(map[string]types.ServiceNode)

	// edge: parent (caller/producer) -> child (callee/consumer); like Jaeger
	// Jaeger supports self-loops but that's redundant in service graphs
	// Support external services? (e.g. databases, third-party APIs)
	edgesMap := make(map[string]types.ServiceEdge)

	// TODO: implement graph construction
	_ = nodesMap
	_ = edgesMap

	// Switch on Span Kinds
	for _, span := range spans {
		// The service of the span is an internal node to the system; record it
		serviceName := extractServiceName(span)
		addNode(nodesMap, serviceName)

		switch span.Kind {
		case SpanKindClient:
			// Outbound synchronous call via HTTP/RPC to another service/external system

			// May synchronously call DB; DB calls are considered client span kinds
			// https://opentelemetry.io/docs/specs/semconv/db/database-spans/
			if dbTarget := extractDBTarget(span); !dbTarget.IsEmpty() {
				// Client span represents a call to a database system
				addNode(nodesMap, dbTarget.String()) // Want to add DB (e.g. "postgres:orders") as an explicit node in the call graph
				addEdge(edgesMap, serviceName, dbTarget.String())
			} else if target := extractHTTPTarget(span); !target.IsEmpty() {
				// Client span represents an outbound HTTP call to another service/external system
				addNode(nodesMap, target.ServiceName)
				addEdge(edgesMap, serviceName, target.ServiceName)
			} else if target := extractRPCTarget(span); !target.IsEmpty() {
				addNode(nodesMap, target.ServiceName)
				addEdge(edgesMap, serviceName, target.ServiceName)
			}

		case SpanKindServer:
			// Inbound synchronous call via HTTP/RPC from another service/external system
			// Look up parent to find caller
			if parentSpan, ok := spanMap[span.ParentSpanID]; ok {
				parentServiceName := extractServiceName(parentSpan)

				// Ignore self-loops
				if parentServiceName != serviceName {
					addNode(nodesMap, parentServiceName) // Want to record caller even if outside time range / system (e.g. external API)
					addEdge(edgesMap, parentServiceName, serviceName)
				}
			}
		case SpanKindProducer:
			// Outbound asynchronous call via messaging system to another service; publishing a message
			target := extractMessagingTarget(span)
			if !target.IsEmpty() {
				addNode(nodesMap, target.String()) // I want to add messaging system (e.g. "kafka:orders") as an explicit node in the call graph
				addEdge(edgesMap, serviceName, target.String())
			}
		case SpanKindConsumer:
			// Inbound asynchronous call via messaging system from another service; consuming a message
			target := extractMessagingTarget(span)
			if !target.IsEmpty() {
				addNode(nodesMap, target.String()) // See above
				addEdge(edgesMap, target.String(), serviceName)
			}

			// Process span links for batch consumers
			for _, link := range span.Links {
				if parentSpan, ok := spanMap[link.SpanID]; ok {
					parentServiceName := extractServiceName(parentSpan)

					// Ignore self-loops
					if parentServiceName != serviceName {
						addNode(nodesMap, parentServiceName)
						addEdge(edgesMap, parentServiceName, serviceName)
					}
				}
			}
		case SpanKindInternal:
			// Internal operation within a service; no edge to add
		}
	}

	return nodesMap, edgesMap, nil
}

func addNode(nodesMap map[string]types.ServiceNode, serviceName string) {
	if serviceName == "" {
		return
	}
	if _, exists := nodesMap[serviceName]; !exists {
		nodesMap[serviceName] = types.ServiceNode{
			ServiceID:   uuid.New().String(),
			ServiceName: serviceName,
		}
	}
}

func addEdge(edgesMap map[string]types.ServiceEdge, parentService, childService string) {
	if parentService == "" || childService == "" || parentService == childService {
		return
	}
	key := parentService + KeyDelimiter + childService
	if _, exists := edgesMap[key]; !exists {
		edgesMap[key] = types.ServiceEdge{
			Source: parentService,
			Target: childService,
		}
	}
}