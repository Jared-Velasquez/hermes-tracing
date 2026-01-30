package graph

import (
	"time"
	"context"
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

	// edge: parent (caller/producer) -> child (callee/consumer)
	edgesMap := make(map[string]types.ServiceEdge)
}