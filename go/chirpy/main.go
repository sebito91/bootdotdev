package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sebito91/bootdotdev/go/chirpy/admin"
	"github.com/sebito91/bootdotdev/go/chirpy/api"
)

func main() {
	fmt.Println("Welcome to Chirpy!")

	appPrefix := "/"
	apiCfg, apiErr := api.NewConfig()
	if apiErr != nil {
		panic(apiErr)
	}

	adminCfg, adminErr := admin.NewConfig(apiCfg)
	if adminErr != nil {
		panic(adminErr)
	}

	mainHandler := http.StripPrefix(appPrefix, http.FileServer(http.Dir(".")))
	fsHandler := apiCfg.MiddlewareMetricsInc(mainHandler)

	// kick off the new multiplexer
	r := chi.NewRouter()
	r.Handle(appPrefix, fsHandler)
	r.Handle(appPrefix+"/*", fsHandler)

	r.Mount("/api", apiCfg.MiddlewareMetricsInc(apiCfg.GetAPI()))
	r.Mount("/admin", adminCfg.GetAdminAPI())

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
