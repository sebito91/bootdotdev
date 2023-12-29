package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// Feed is an API-version of the struct the comes from the database. This helps to delineate HTTP from DB API
type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

// createFeed will generate a new entry in the feeds table using the providing information
func (ac *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
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

	newFeedFollowUUID, err := uuid.NewRandom()
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

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        newFeedFollowUUID,
		CreatedAt: newTime,
		UpdatedAt: newTime,
		FeedID:    newUUID,
		UserID:    user.ID,
	}

	dbFeed, err := ac.DB.CreateFeed(r.Context(), newFeed)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeed: %s", err))
		return
	}

	dbFeedFollow, err := ac.DB.CreateFeedFollow(r.Context(), newFeedFollow)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createFeed: createFeedFollow: %s", err))
		return
	}

	feed := Feed{
		ID:        dbFeed.ID,
		Name:      dbFeed.Name,
		URL:       dbFeed.Url,
		UserID:    dbFeed.UserID,
		CreatedAt: dbFeed.CreatedAt,
		UpdatedAt: dbFeed.UpdatedAt,
	}

	feedFollow := FeedFollow{
		ID:        dbFeedFollow.ID,
		CreatedAt: dbFeedFollow.CreatedAt,
		UpdatedAt: dbFeedFollow.UpdatedAt,
		FeedID:    dbFeedFollow.FeedID,
		UserID:    dbFeedFollow.UserID,
	}

	respondWithJSON(w, http.StatusCreated, struct {
		Feed       Feed       `json:"feed"`
		FeedFollow FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feedFollow,
	})
}

// getFeeds will fetch all feeds from the 'feeds' table in the database
func (ac *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	dbFeeds, err := ac.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("getFeeds: %s", err))
		return
	}

	feeds := make([]Feed, len(dbFeeds))
	for idx, dbFeed := range dbFeeds {
		feeds[idx] = Feed{
			ID:        dbFeed.ID,
			Name:      dbFeed.Name,
			URL:       dbFeed.Url,
			UserID:    dbFeed.UserID,
			CreatedAt: dbFeed.CreatedAt,
			UpdatedAt: dbFeed.UpdatedAt,
		}
	}

	respondWithJSON(w, http.StatusCreated, feeds)
}
