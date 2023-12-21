package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GetAPI generates the new route for the aggregator and returns a handle to the router
func GetAPI() (chi.Router, error) {
	r := chi.NewRouter()

	//	r.Get("/healthz", readinessEndpoint)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", mainPage)
	})

	return r, nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("we're A-OK!")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
