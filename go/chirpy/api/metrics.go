package api

import "net/http"

// MiddlewareMetricsInc will use a middleware to increment an in-memory counter
// of the number of site visits during server operation
func (c *Config) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.mux.Lock()
		c.fileserverHits++
		c.mux.Unlock()
		next.ServeHTTP(w, r)
	})
}

// GetFileserverHits returns the current number of site visits since the start of the server
func (c *Config) GetFileserverHits() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.fileserverHits
}

// ResetFileserverHits will reset the fileserverHits counter as if the server restarted
func (c *Config) ResetFileserverHits() {
	c.mux.Lock()
	c.fileserverHits = 0
	c.mux.Unlock()
}

// readinessEndpoint yields the status and information for the /healthz endpoint
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}
