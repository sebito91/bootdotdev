package main

import (
	"fmt"
	"net/http"

	"github.com/sebito91/bootdotdev/go/bloggy/internal/database"
)

// authedHandler is a custom function handler type that only deals with
// HTTP requests if the database.User is authorized (i.e. legit user with entry in system)
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// middlewareAuth defines an http.HandlerFunc that processes the incoming HTTP request
// for a valid user Authorization header providing an ApiKey. If found and acceptable
// the func is successfully returned and the next part of the request is performed.
func (api *APIConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := fetchAPIToken(r)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("middlewareAuth: %s", err))
			return
		}

		user, err := api.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, fmt.Sprintf("middlewareAuth user fetch: %s", err))
			return
		}

		handler(w, r, user)
	}
}
