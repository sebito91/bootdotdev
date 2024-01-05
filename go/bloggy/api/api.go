package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// apiConfig is a struct to hold references to our database, router, and other components
type apiConfig struct {
	DB            *database.Queries
	Router        chi.Router
	concurrency   int
	sleepInterval time.Duration
}

// GetAPI generates the new route for the aggregator and returns a handle to the router
func GetAPI(concurrency int, sleepInterval time.Duration) (*apiConfig, error) {
	apiCfg := &apiConfig{
		concurrency:   concurrency,
		sleepInterval: sleepInterval,
	}

	err := godotenv.Load()
	if err != nil {
		return apiCfg, err
	}

	dbURL := os.Getenv("CONN")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return apiCfg, err
	}

	apiCfg.DB = database.New(db)

	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", mainPage)

		r.Get("/readiness", readinessEndpoint)
		r.Get("/err", errorTester)

		r.Get("/users", apiCfg.middlewareAuth(apiCfg.getUserByAPIKey))
		r.Post("/users", apiCfg.createUser)

		r.Get("/feeds", apiCfg.getFeeds)
		r.Post("/feeds", apiCfg.middlewareAuth(apiCfg.createFeed))

		r.Route("/feed_follows", func(r chi.Router) {
			r.Get("/", apiCfg.middlewareAuth(apiCfg.getFeedFollows))
			r.Post("/", apiCfg.middlewareAuth(apiCfg.createFeedFollow))

			r.Route("/{feedFollowID}", func(r chi.Router) {
				r.Delete("/", apiCfg.middlewareAuth(apiCfg.deleteFeedFollow))
			})
		})
	})

	apiCfg.Router = r
	return apiCfg, nil
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("we're A-OK!")); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write: %s", err))
	}
}
