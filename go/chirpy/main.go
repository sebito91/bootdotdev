package main

import (
	"fmt"
	"net/http"
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

func main() {
	fmt.Println("vim-go")

	// kick off the new multiplexer
	mux := http.NewServeMux()

	// wrap the mux in a custom middleware for CORS headers
	corsMux := middlewareCors(mux)

	// create the server struct
	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: corsMux,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
