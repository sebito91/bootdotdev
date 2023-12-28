package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// createFeedFollow will add an entry to the 'feed_follows' table for the requesting user and the provided feed_id
func (api *APIConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type newFeedFollowCheck struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	newFeedFollowChk := newFeedFollowCheck{}

	if err := decoder.Decode(&newFeedFollowChk); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("createFeedFollow: could not decode JSON payload: %s", err))
		return
	}

	newFeedFollowUUID, err := uuid.Parse(newFeedFollowChk.FeedID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "createFeedFollow: received invalid value for `feed_id` field")
		return
	}

	if newFeedFollowUUID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "createFeedFollow: `feed_id` field cannot be empty")
		return
	}

	newUUID, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeedFollow: %s", err))
		return
	}

	newTime := time.Now()
	newFeedFollow := database.CreateFeedFollowParams{
		ID:        newUUID,
		CreatedAt: newTime,
		UpdatedAt: newTime,
		FeedID:    newFeedFollowUUID,
		UserID:    user.ID,
	}

	feedFollow, err := api.DB.CreateFeedFollow(r.Context(), newFeedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeedFollow: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}

// getFeedFollows will fetch all feeds that the requesting user follows
func (api *APIConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := api.DB.GetFeedFollowsByUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("getFeedFollows: %s", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollows)
}

// deleteFeedFollow will remove the follow for a given feed_id that the requesting user follows
func (api *APIConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
}
