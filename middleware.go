package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Metrics struct {
	mu                sync.Mutex
	totalTime         time.Duration
	totalLatency      time.Duration
	requests          int
	endpointMetrics   map[string]time.Duration
	endpointLatencies map[string]time.Duration
}

var metrics = Metrics{
	endpointMetrics:   make(map[string]time.Duration),
	endpointLatencies: make(map[string]time.Duration),
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		metrics.mu.Lock()
		metrics.requests++
		metrics.totalTime += duration
		metrics.endpointMetrics[r.URL.Path] += duration
		metrics.mu.Unlock()

		log.Printf("%s | %s | %v", r.Method, r.URL.Path, duration)
	})
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	meanDuration := time.Duration(0)
	meanLatency := time.Duration(0)

	if metrics.requests > 0 {
		meanDuration = metrics.totalTime / time.Duration(metrics.requests)
		meanLatency = metrics.totalLatency / time.Duration(metrics.requests)
	}

	response := map[string]interface{}{
		"total_requests":        metrics.requests,
		"mean_request_duration": meanDuration.String(),
		"mean_request_latency":  meanLatency.String(),
		"endpoint_metrics":      metrics.endpointMetrics,
		"endpoint_latencies":    metrics.endpointLatencies,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
