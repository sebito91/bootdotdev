package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

		r.Get("/users", api.getUserByAPIKey)
		r.Post("/users", api.createUser)
	})

	api.Router = r
	return api, nil
}

// createUser will generate a new user in the database with all of the corresponding fields
func (api *APIConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type newUserCheck struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	newUserChk := newUserCheck{}

	if err := decoder.Decode(&newUserChk); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("createUser: could not decode JSON payload: %s", err))
		return
	}

	if newUserChk.Name == "" {
		respondWithError(w, http.StatusBadRequest, "createUser: did not receive value for `name` field")
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createUser: %s", err))
		return
	}

	newTime := time.Now()
	newUser := database.CreateUserParams{
		ID:        newUUID,
		CreatedAt: newTime,
		UpdatedAt: newTime,
		Name:      newUserChk.Name,
	}

	user, err := api.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createUser: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// getUserByAPIKey will fetch the user with the provided API Key from the request header
func (api *APIConfig) getUserByAPIKey(w http.ResponseWriter, r *http.Request) {
	apiKey, err := fetchAPIToken(r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("getUserByAPIKey: %s", err))
		return
	}

	user, err := api.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("getUserByApiKey fetch: %s", err))
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("we're A-OK!")); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write: %s", err))
	}
}
