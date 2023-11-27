package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/sebito91/bootdotdev/go/chirpy/database"
	"golang.org/x/crypto/bcrypt"
)

// APIConfig is a local struct to keep track of site visits
// NOTE: this value is in-memory only and will persist for the duration of the server
type APIConfig struct {
	fileserverHits int
	db             *database.DB
	mux            sync.RWMutex
}

// errorBody is a struct used for returning a JSON-based error code/string
type errorBody struct {
	Error     string `json:"error"`
	errorCode int
}

// NewAPIConfig returns a new instance of the APIConfig
func NewAPIConfig() (*APIConfig, error) {
	db, err := database.NewDB("")
	if err != nil {
		return nil, err
	}

	return &APIConfig{db: db}, nil
}

// GetAPI returns the router for the /api endpoint
func (c *APIConfig) GetAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", readinessEndpoint)

	r.Route("/chirps", func(r chi.Router) {
		r.Get("/", c.getChirps)
		r.Post("/", c.writeChirp)

		r.Route("/{chirpID}", func(r chi.Router) {
			r.Get("/", c.getChirpID)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", c.getUsers)
		r.Post("/", c.writeUser)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", c.getUserByID)
		})
	})

	r.Post("/login", c.loginUser)

	return r
}

// middlewareMetricsInc will use a middleware to increment an in-memory counter
// of the number of site visits during server operation
func (c *APIConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.mux.Lock()
		c.fileserverHits++
		c.mux.Unlock()
		next.ServeHTTP(w, r)
	})
}

// GetFileserverHits returns the current number of site visits since the start of the server
func (c *APIConfig) GetFileserverHits() int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.fileserverHits
}

// ResetFileserverHits will reset the fileserverHits counter as if the server restarted
func (c *APIConfig) ResetFileserverHits() {
	c.mux.Lock()
	c.fileserverHits = 0
	c.mux.Unlock()
}

// getUsers will fetch all of the users stored within the database
func (c *APIConfig) getUsers(w http.ResponseWriter, r *http.Request) {
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
func (c *APIConfig) loginUser(w http.ResponseWriter, r *http.Request) {
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
func (c *APIConfig) getUserByID(w http.ResponseWriter, r *http.Request) {
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
func (c *APIConfig) writeUser(w http.ResponseWriter, r *http.Request) {
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

// getChirps will fetch the chirps from the DB and write to the page
func (c *APIConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.db.GetChirps()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusOK, chirps)
}

// getChirpID will fetch a specific chirp from the database
func (c *APIConfig) getChirpID(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.db.GetChirps()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if chirpID <= 0 {
		errBody := errorBody{
			Error:     fmt.Sprintf("expected valid chirpID (>0), got %d", chirpID),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	for _, chirp := range chirps {
		if chirp.ID == chirpID {
			writeSuccessToPage(w, http.StatusOK, chirp)
			return
		}
	}

	errBody := errorBody{
		Error:     fmt.Sprintf("could not find chirpID %d", chirpID),
		errorCode: http.StatusNotFound,
	}

	errBody.writeErrorToPage(w)
}

// writeChirp will validate the chirp first, and if successful commit to the db
func (c *APIConfig) writeChirp(w http.ResponseWriter, r *http.Request) {
	type bodyCheck struct {
		Body string `json:"body"`
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

	// if chirp is too long (>140 chars), send a 400 error
	if len(bodyChk.Body) > 140 {
		errBody := errorBody{
			Error:     "Chirp is too long",
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	chirp, err := c.db.CreateChirp(cleanedBody(bodyChk.Body))
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusCreated, chirp)
}

func cleanedBody(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	for _, word := range badWords {
		re := regexp.MustCompile(fmt.Sprintf("(?i)%s", word))
		body = re.ReplaceAllString(body, "****")
	}

	return body
}

// readinessEndpoint yields the status and information for the /healthz endpoint
func readinessEndpoint(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := w.Write([]byte("OK")); err != nil {
		panic(err)
	}
}

// writeErrorToPage is a helper function that reuses a JSON-posting for error messages
func (e *errorBody) writeErrorToPage(w http.ResponseWriter) {
	w.WriteHeader(e.errorCode)
	dat, datErr := json.Marshal(e)
	if datErr != nil {
		log.Printf("Error marshaling error JSON: %s", datErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, wErr := w.Write(dat); wErr != nil {
		log.Printf("Error writing error JSON to page: %s", wErr)
		return
	}
}

// writeSuccessToPage is a helper function that reuses a JSON-posting for success messages
func writeSuccessToPage(w http.ResponseWriter, statusCode int, payload interface{}) {
	dat, datErr := json.Marshal(payload)
	if datErr != nil {
		log.Printf("Error marshaling success JSON: %s", datErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, wErr := w.Write(dat); wErr != nil {
		log.Printf("Error writing success JSON to page: %s", wErr)
		return
	}
}
