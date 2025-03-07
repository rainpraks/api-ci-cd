package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

// Metrics struct to track both latency and request duration
type Metrics struct {
	mu                sync.Mutex
	totalTime         time.Duration
	totalLatency      time.Duration
	requests          int
	endpointMetrics   map[string]time.Duration
	endpointLatencies map[string]time.Duration
}

// Global metrics instance
var metrics = Metrics{
	endpointMetrics:   make(map[string]time.Duration),
	endpointLatencies: make(map[string]time.Duration),
}

// Middleware to log requests and measure latency/duration
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Simulate Latency Measurement (time until processing starts)
		latencyStart := time.Now()
		time.Sleep(5 * time.Millisecond) // Simulated network delay (for real-world, remove)
		latency := time.Since(latencyStart)

		// Log request
		log.Printf(" %s %s", r.Method, r.URL.Path)

		// Process request
		next.ServeHTTP(w, r)

		// Compute request duration
		duration := time.Since(start)

		// Store metrics safely
		metrics.mu.Lock()
		metrics.requests++
		metrics.totalTime += duration
		metrics.totalLatency += latency
		metrics.endpointMetrics[r.URL.Path] += duration
		metrics.endpointLatencies[r.URL.Path] += latency
		metrics.mu.Unlock()

		// Log response time
		log.Printf(" %s %s took %v (Latency: %v, Processing: %v)",
			r.Method, r.URL.Path, duration, latency, duration-latency)
	})
}

// Handler for /metrics endpoint
func getMetrics(w http.ResponseWriter, r *http.Request) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	// Compute mean request time and latency
	meanDuration := time.Duration(0)
	meanLatency := time.Duration(0)

	if metrics.requests > 0 {
		meanDuration = metrics.totalTime / time.Duration(metrics.requests)
		meanLatency = metrics.totalLatency / time.Duration(metrics.requests)
	}

	// Build response JSON
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
