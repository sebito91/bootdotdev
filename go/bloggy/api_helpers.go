package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// fetchAPIToken is a helper function to extract the APIKey from a given request
func fetchAPIToken(r *http.Request) (string, error) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("expected Authorization token, got %s", bearer)
	}

	return strings.TrimPrefix(bearer, "ApiKey "), nil
}

// respondWithError will write out an error message to the console
func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{
		Error: msg,
	})
}

// respondWithJSON is a helper function that reuses a JSON-posting for success messages
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	dat, datErr := json.Marshal(payload)
	if datErr != nil {
		log.Printf("Error marshaling JSON response: %s", datErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, wErr := w.Write(dat); wErr != nil {
		log.Printf("Error writing JSON to page: %s", wErr)
		return
	}
}

// readinessEndpoint will render a verbal status based on the health + readiness of the webapp
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}

func errorTester(w http.ResponseWriter, r *http.Request) {
	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}
