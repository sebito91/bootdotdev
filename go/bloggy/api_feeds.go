package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

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
