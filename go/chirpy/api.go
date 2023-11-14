package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// apiConfig is a local struct to keep track of site visits
// NOTE: this value is in-memory only and will persist for the duration of the server
type apiConfig struct {
	fileserverHits int
}

// GetAPI returns the router for the /api endpoint
func (c *apiConfig) GetAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", readinessEndpoint)

	return r
}

// GetAdminAPI returns the router for the /admin endpoint
func (c *apiConfig) GetAdminAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/metrics", c.metricsEndpoint)
	r.Get("/reset", c.resetEndpoint)

	return r
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

	w.Header().Set("Content-Type", "text/html")
	content := `
<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
</html>
`
	if _, err := fmt.Fprintf(w, content, c.fileserverHits); err != nil {
		panic(err)
	}
}

// resetEndpoint will reset the number of site visits to 0 during a running server
func (c *apiConfig) resetEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.fileserverHits = 0
}

// readinessEndpoint yields the status and information for the /healthz endpoint
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}
