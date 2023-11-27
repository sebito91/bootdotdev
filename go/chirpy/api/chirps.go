package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// getChirps will fetch the chirps from the DB and write to the page
func (c *Config) getChirps(w http.ResponseWriter, r *http.Request) {
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
func (c *Config) getChirpID(w http.ResponseWriter, r *http.Request) {
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
func (c *Config) writeChirp(w http.ResponseWriter, r *http.Request) {
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
