package handlers

import (
	"context"
	"net/http"
	"time"
	"encoding/json"

	"analytics/store"
	"analytics/types"
	"analytics/graph"
)

type AnalyticsHandler struct {
	graphBuilder *graph.GraphBuilder
}

func NewAnalyticsHandler(graphBuilder *graph.GraphBuilder) *AnalyticsHandler {
	return &AnalyticsHandler{graphBuilder: graphBuilder}
}

func (ah *AnalyticsHandler) ConstructServiceGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	start, end, err := extractTimeRange(r)
	if err != nil {
		http.Error(w, "Invalid time range", http.StatusBadRequest)
		return
	}

	serviceGraph, err := ah.graphBuilder.Build(ctx, start, end)
	if err != nil {
		http.Error(w, "Failed to build service graph", http.StatusInternalServerError)
		return
	}
	writeGraphJSONResponse(w, serviceGraph)
}

func extractTimeRange(r *http.Request) (time.Time, time.Time, error) {
	// start and end are in the query parameters
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func writeGraphJSONResponse(w http.ResponseWriter, graph *types.ServiceGraph) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(graph)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}