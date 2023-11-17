package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// adminConfig is a placeholder for our /admin API section
type adminConfig struct {
	API *apiConfig
}

// GetAdminAPI returns the router for the /admin endpoint
func (c *adminConfig) GetAdminAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/metrics", c.metricsEndpoint)
	r.Get("/reset", c.resetEndpoint)

	if c.API == nil {
		panic("Please make sure to initialize the API Config before the Admin Config")
	}

	return r
}

// metricsEndpoint will use a write-enabled middleware to display the number
// of site visits since the start of the server
func (c *adminConfig) metricsEndpoint(w http.ResponseWriter, r *http.Request) {
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
	if _, err := fmt.Fprintf(w, content, c.API.GetFileserverHits()); err != nil {
		panic(err)
	}
}

// resetEndpoint will reset the number of site visits to 0 during a running server
func (c *adminConfig) resetEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.API.ResetFileserverHits()
}
