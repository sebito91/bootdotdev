package api

import (
	"encoding/json"
	"log"
	"net/http"
)

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
