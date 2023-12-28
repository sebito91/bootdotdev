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

		r.Get("/users", api.middlewareAuth(api.getUserByAPIKey))
		r.Post("/users", api.createUser)

		r.Get("/feeds", api.getFeeds)
		r.Post("/feeds", api.middlewareAuth(api.createFeed))
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
func (api *APIConfig) getUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, user)
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("we're A-OK!")); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to write: %s", err))
	}
}

// createFeed will generate a new entry in the feeds table using the providing information
func (api *APIConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type newFeedCheck struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	newFeedChk := newFeedCheck{}

	if err := decoder.Decode(&newFeedChk); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("createFeed: could not decode JSON payload: %s", err))
		return
	}

	if newFeedChk.Name == "" {
		respondWithError(w, http.StatusBadRequest, "createFeed: did not receive value for `name` field")
		return
	}

	if newFeedChk.URL == "" {
		respondWithError(w, http.StatusBadRequest, "createFeed: did not receive value for `url` field")
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeed: %s", err))
		return
	}

	newTime := time.Now()
	newFeed := database.CreateFeedParams{
		ID:        newUUID,
		CreatedAt: newTime,
		UpdatedAt: newTime,
		Name:      newFeedChk.Name,
		Url:       newFeedChk.URL,
		UserID:    user.ID,
	}

	feed, err := api.DB.CreateFeed(r.Context(), newFeed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeed: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, feed)
}

// getFeeds will fetch all feeds from the 'feeds' table in the database
func (api *APIConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := api.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("getFeeds: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, feeds)
}
