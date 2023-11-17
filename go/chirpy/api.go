package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"

	"github.com/go-chi/chi/v5"
)

// apiConfig is a local struct to keep track of site visits
// NOTE: this value is in-memory only and will persist for the duration of the server
type apiConfig struct {
	fileserverHits int
	mux            sync.RWMutex
}

// errorBody is a struct used for returning a JSON-based error code/string
type errorBody struct {
	Error     string `json:"error"`
	errorCode int
}

// GetAPI returns the router for the /api endpoint
func (c *apiConfig) GetAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", readinessEndpoint)

	r.Post("/validate_chirp", c.chirps)

	return r
}

// middlewareMetricsInc will use a middleware to increment an in-memory counter
// of the number of site visits during server operation
func (c *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.mux.Lock()
		c.fileserverHits++
		c.mux.Unlock()
		next.ServeHTTP(w, r)
	})
}

// GetFileserverHits returns the current number of site visits since the start of the server
func (c *apiConfig) GetFileserverHits() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.fileserverHits
}

// ResetFileserverHits will reset the fileserverHits counter as if the server restarted
func (c *apiConfig) ResetFileserverHits() {
	c.mux.Lock()
	c.fileserverHits = 0
	c.mux.Unlock()
}

func (c *apiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	type bodyCheck struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	bodyChk := bodyCheck{}

	// handle a decode error
	if err := decoder.Decode(&bodyChk); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	// if chirp is too long (>140 chars), send a 400 error
	if len(bodyChk.Body) > 140 {
		errBody := errorBody{
			Error:     "Chirp is too long",
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	payload := struct {
		Body string `json:"cleaned_body"`
	}{
		Body: cleanedBody(bodyChk.Body),
	}
	writeSuccessToPage(w, http.StatusOK, payload)
}

func cleanedBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range badWords {
		re := regexp.MustCompile(fmt.Sprintf("(?i)%s", word))
		body = re.ReplaceAllString(body, "****")
	}

	return body
}

// readinessEndpoint yields the status and information for the /healthz endpoint
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}

// writeErrorToPage is a helper function that reuses a JSON-posting for error messages
func (e *errorBody) writeErrorToPage(w http.ResponseWriter) {
	w.WriteHeader(e.errorCode)
	dat, datErr := json.Marshal(e)
	if datErr != nil {
		log.Printf("Error marshaling error JSON: %s", datErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, wErr := w.Write(dat); wErr != nil {
		log.Printf("Error writing error JSON to page: %s", wErr)
		return
	}
}

// writeSuccessToPage is a helper function that reuses a JSON-posting for success messages
func writeSuccessToPage(w http.ResponseWriter, statusCode int, payload interface{}) {
	dat, datErr := json.Marshal(payload)
	if datErr != nil {
		log.Printf("Error marshaling success JSON: %s", datErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, wErr := w.Write(dat); wErr != nil {
		log.Printf("Error writing success JSON to page: %s", wErr)
		return
	}
}
