package main

import (
	"fmt"
	"net/http"
)

// apiConfig is a local struct to keep track of site visits
// NOTE: this value is in-memory only and will persist for the duration of the server
type apiConfig struct {
	fileserverHits int
}

// middlewareMetricsInc will use a middleware to increment an in-memory counter
// of the number of site visits during server operation
func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

// metricsEndpoint will use a write-enabled middleware to display the number
// of site visits since the start of the server
func (c *apiConfig) metricsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := fmt.Fprintf(w, "Hits: %d", c.fileserverHits); err != nil {
		panic(err)
	}
}

// resetEndpoint will reset the number of site visits to 0 during a running server
func (c *apiConfig) resetEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.fileserverHits = 0
}

// middlewareCors enables the cross-origin features required to run via boot.dev test servers
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// readinessEndpoint yields the status and information for the /healthz endpoint
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}
