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

type DBTarget struct {
	System 		 string
	DatabaseName string
}

func (db DBTarget) String() string {
	return db.System + ":" + db.DatabaseName
}

func (db DBTarget) IsEmpty() bool {
	return db.System == ""
}

type MessagingTarget struct {
	System   		string
	DestinationName string
}

func (mt MessagingTarget) String() string {
	return mt.System + ":" + mt.DestinationName
}

func (mt MessagingTarget) IsEmpty() bool {
	return mt.System == ""
}

type SyncCallTarget struct {
	// Is protocol the right word?
	Protocol 	string // Example: http, grpc, thrift
	Source 	 	string // Example: service.name (internal)
	Destination string // Example: service.name (internal), api.stripe.com (external)
}

func (sct SyncCallTarget) String() string {
	// TODO: how is this related with ServiceNodes?
	return sct.Source + "->" + sct.Destination
}

func (sct SyncCallTarget) IsEmpty() bool {
	return sct.Source == "" && sct.Destination == ""
}