package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("vim-go")

	appPrefix := "/app"
	apiCfg := &apiConfig{}

	mainHandler := http.StripPrefix(appPrefix, http.FileServer(http.Dir(".")))
	fsHandler := apiCfg.middlewareMetricsInc(mainHandler)

	// kick off the new multiplexer
	r := chi.NewRouter()
	r.Handle(appPrefix, fsHandler)
	r.Handle(appPrefix+"/*", fsHandler)

	r.Get("/healthz", readinessEndpoint)
	r.Get("/metrics", apiCfg.metricsEndpoint)
	r.Get("/reset", apiCfg.resetEndpoint)

	// wrap the mux in a custom middleware for CORS headers
	corsMux := middlewareCors(r)

	// create the server struct
	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           corsMux,
		ReadHeaderTimeout: time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
