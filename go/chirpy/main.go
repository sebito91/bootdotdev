package main

import (
	"fmt"
	"net/http"
	"time"
)

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

func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("vim-go")

	appPrefix := "/app"

	// kick off the new multiplexer
	mux := http.NewServeMux()
	mux.Handle(appPrefix+"/", http.StripPrefix(appPrefix, http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", readinessEndpoint)

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
