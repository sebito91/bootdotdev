package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"sort"
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

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		// sort the IDs in descending order
		sort.Slice(chirps, func(a, b int) bool {
			return chirps[a].ID > chirps[b].ID
		})
	} else {
		// sort the IDs in ascending order
		sort.Slice(chirps, func(a, b int) bool {
			return chirps[a].ID < chirps[b].ID
		})
	}

	authorIDParam := r.URL.Query().Get("author_id")
	if authorIDParam == "" {
		writeSuccessToPage(w, http.StatusOK, chirps)
		return
	}

	authorID, err := strconv.Atoi(authorIDParam)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	chirpsByAuthor, err := c.db.GetChirpsByAuthorID(authorID)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if sortOrder == "desc" {
		// sort the IDs in descending order
		sort.Slice(chirpsByAuthor, func(a, b int) bool {
			return chirpsByAuthor[a].ID > chirpsByAuthor[b].ID
		})
	} else {
		// sort the IDs in ascending order
		sort.Slice(chirpsByAuthor, func(a, b int) bool {
			return chirpsByAuthor[a].ID < chirpsByAuthor[b].ID
		})
	}

	writeSuccessToPage(w, http.StatusOK, chirpsByAuthor)
}

// getChirpByID will fetch a specific chirp from the database
func (c *Config) getChirpByID(w http.ResponseWriter, r *http.Request) {
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

// deleteChirpByID will delete a specific chirp from the database if the user is authorized to do so
// the JWT from the request is validated, and if the authenticated user matches the author_id of the given
// chirp, then the system will remove the chirp from the database.
func (c *Config) deleteChirpByID(w http.ResponseWriter, r *http.Request) {
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

	claims, respCode, err := c.fetchClaims(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: respCode,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if issuer, claimErr := claims.GetIssuer(); claimErr != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", claimErr),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	} else if issuer == chirpyRefresh {
		errBody := errorBody{
			Error:     "cannot use refresh token for delete chirp request, please provide valid access token",
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	idString, err := claims.GetSubject()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	authorID, err := strconv.Atoi(idString)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("could not convert userID to string: %s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	for _, chirp := range chirps {
		if chirp.ID == chirpID && chirp.AuthorID == authorID {
			if err := c.db.DeleteChirp(chirp); err != nil {
				errBody := errorBody{
					Error:     fmt.Sprintf("could not delete chirpID %d: %s", chirpID, err),
					errorCode: http.StatusInternalServerError,
				}

				errBody.writeErrorToPage(w)
				return
			}

			writeSuccessToPage(w, http.StatusOK, nil)
			return
		} else if chirp.ID == chirpID && chirp.AuthorID != authorID {
			errBody := errorBody{
				Error:     fmt.Sprintf("cannot delete chirpID %d by authorID %d, unauthorized", chirpID, authorID),
				errorCode: http.StatusForbidden,
			}

			errBody.writeErrorToPage(w)
			return
		}
	}

	errBody := errorBody{
		Error:     fmt.Sprintf("could not find chirpID %d by author %d", chirpID, authorID),
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

	claims, respCode, err := c.fetchClaims(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: respCode,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if issuer, claimErr := claims.GetIssuer(); claimErr != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", claimErr),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	} else if issuer == chirpyRefresh {
		errBody := errorBody{
			Error:     "cannot use refresh token for chirp request, please provide valid access token",
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	idString, err := claims.GetSubject()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	authorID, err := strconv.Atoi(idString)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("could not convert userID to string: %s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	chirp, err := c.db.CreateChirp(authorID, cleanedBody(bodyChk.Body))
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
