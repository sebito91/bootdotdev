package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// User is an API-version of the struct the comes from the database. This helps to delineate HTTP from DB API
type User struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// createUser will generate a new user in the database with all of the corresponding fields
func (ac *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
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

	dbUser, err := ac.DB.CreateUser(r.Context(), newUser)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("createUser: %s", err))
		return
	}

	user := User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// getUserByAPIKey will fetch the user with the provided API Key from the request header
func (ac *apiConfig) getUserByAPIKey(w http.ResponseWriter, r *http.Request, dbUser database.User) {
	user := User{
		ID:        dbUser.ID,
		Name:      dbUser.Name,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	respondWithJSON(w, http.StatusOK, user)
}
