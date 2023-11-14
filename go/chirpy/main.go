package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	fmt.Println("vim-go")

	appPrefix := "/app"
	apiCfg := &apiConfig{}

	mainHandler := http.StripPrefix(appPrefix, http.FileServer(http.Dir(".")))

	// kick off the new multiplexer
	mux := http.NewServeMux()
	mux.Handle(appPrefix+"/", apiCfg.middlewareMetricsInc(mainHandler))
	mux.HandleFunc("/healthz", readinessEndpoint)
	mux.HandleFunc("/metrics", apiCfg.metricsEndpoint)
	mux.HandleFunc("/reset", apiCfg.resetEndpoint)

	// wrap the mux in a custom middleware for CORS headers
	corsMux := middlewareCors(mux)

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
