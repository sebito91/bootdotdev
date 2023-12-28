package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// APIConfig is a struct to hold references to our database, router, and other components
type APIConfig struct {
	DB     *database.Queries
	Router chi.Router
}

// GetAPI generates the new route for the aggregator and returns a handle to the router
func GetAPI() (*APIConfig, error) {
	api := &APIConfig{}

	err := godotenv.Load()
	if err != nil {
		return api, err
	}

	dbURL := os.Getenv("CONN")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return api, err
	}

	api.DB = database.New(db)

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", mainPage)

		r.Get("/readiness", readinessEndpoint)
		r.Get("/err", errorTester)

		r.Get("/users", api.middlewareAuth(api.getUserByAPIKey))
		r.Post("/users", api.createUser)

		r.Get("/feeds", api.getFeeds)
		r.Post("/feeds", api.middlewareAuth(api.createFeed))

		r.Route("/feed_follows", func(r chi.Router) {
			r.Get("/", api.middlewareAuth(api.getFeedFollows))
			r.Post("/", api.middlewareAuth(api.createFeedFollow))

			r.Route("/{feedFollowID}", func(r chi.Router) {
				r.Delete("/", api.middlewareAuth(api.deleteFeedFollow))
			})
		})
	})

	api.Router = r
	return api, nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("we're A-OK!")); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write: %s", err))
	}
}
