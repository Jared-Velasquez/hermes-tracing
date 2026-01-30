package types

import (
	"time"
)

// Span document stored in Elasticsearch
// How to handle Span Links?
type Span struct {
	TraceID      string                 `json:"trace_id"`
	SpanID       string                 `json:"span_id"`
	ParentSpanID string                 `json:"parent_span_id"`
	Name         string                 `json:"name"`
	Kind         string                 `json:"kind"`
	Links        []SpanLink             `json:"links"`
	StartTime    uint64                 `json:"start_time"`
	EndTime      uint64                 `json:"end_time"`
	Status       string                 `json:"status"`
	Attributes   map[string]interface{} `json:"attributes"`
	Resource     map[string]interface{} `json:"resource"`
	ScopeName    string                 `json:"scope_name"`
}

type SpanLink struct {
	TraceID string `json:"trace_id"`
	SpanID  string `json:"span_id"`
}

type ServiceNode struct {
	ServiceID  	string `json:"service_id"`
	ServiceName string `json:"service_name"`
}

type ServiceEdge struct {
	Source     	string `json:"source"`
	Target     	string `json:"target"`
}

type ServiceGraph struct {
	GraphID   string         `json:"graph_id"`
	StartTime time.Time      `json:"start_time"`
	EndTime   time.Time      `json:"end_time"`
	Nodes     []ServiceNode  `json:"nodes"`
	Edges     []ServiceEdge  `json:"edges"`
	CreatedAt time.Time      `json:"created_at"`
}