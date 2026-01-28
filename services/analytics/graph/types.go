import (
	"time"
)

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