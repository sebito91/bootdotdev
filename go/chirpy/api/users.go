package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sebito91/bootdotdev/go/chirpy/database"
	"golang.org/x/crypto/bcrypt"
)

// getUsers will fetch all of the users stored within the database
func (c *Config) getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := c.db.GetUsers()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusOK, users)
}

// loginUser will check if a given user is stored in the database and the credentials provided are
// correct/matching. If all matches as expected, a success is sent; if anything is a mismatch or the user
// doesn't exist, an error is sent
func (c *Config) loginUser(w http.ResponseWriter, r *http.Request) {
	type bodyCheck struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	bodyChk := bodyCheck{}

	// handle a decode error
	if err := decoder.Decode(&bodyChk); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	users, err := c.db.GetUsersFull()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if bodyChk.Email == "" || bodyChk.Password == "" {
		errBody := errorBody{
			Error:     "login expected valid user email adddress and password",
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	for _, user := range users {
		passErr := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(bodyChk.Password))
		if user.Email == bodyChk.Email {
			if passErr == nil {
				writeSuccessToPage(w, http.StatusOK, database.User{ID: user.ID, Email: user.Email})
				return
			}

			errBody := errorBody{
				Error:     fmt.Sprintf("could not authenticate user with email %s", bodyChk.Email),
				errorCode: http.StatusUnauthorized,
			}

			errBody.writeErrorToPage(w)
			return
		}
	}

	errBody := errorBody{
		Error:     fmt.Sprintf("could not find user with email %s", bodyChk.Email),
		errorCode: http.StatusNotFound,
	}

	errBody.writeErrorToPage(w)
}

// getUserByID will fetch the specific user with the provided userID from the database
func (c *Config) getUserByID(w http.ResponseWriter, r *http.Request) {
	users, err := c.db.GetUsers()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if userID <= 0 {
		errBody := errorBody{
			Error:     fmt.Sprintf("expected valid userID (>0), got %d", userID),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	for _, user := range users {
		if user.ID == userID {
			writeSuccessToPage(w, http.StatusOK, user)
			return
		}
	}

	errBody := errorBody{
		Error:     fmt.Sprintf("could not find userID %d", userID),
		errorCode: http.StatusNotFound,
	}

	errBody.writeErrorToPage(w)
}

// writeUser will persist the user to the database, if the user does not exist
func (c *Config) writeUser(w http.ResponseWriter, r *http.Request) {
	type bodyCheck struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	bodyChk := bodyCheck{}

	// handle a decode error

	if err := decoder.Decode(&bodyChk); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if bodyChk.Email == "" || bodyChk.Password == "" {
		errBody := errorBody{
			Error:     "system requires both a valid email and password",
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(bodyChk.Password), bcrypt.DefaultCost)
	if err != nil {
		errBody := errorBody{
			Error:     "could not encode password, please send valid string",
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	user, err := c.db.CreateUser(bodyChk.Email, passHash)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusCreated, user)
}
