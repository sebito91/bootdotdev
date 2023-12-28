package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

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
