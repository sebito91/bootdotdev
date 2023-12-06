package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// processPolkaUpdate is a webhook that receives updates for users from the Polka payment system.
// if the user has upgraded their service to ChirpyRed, and payment is verified by Polka, we must
// update that user in the database.
// all other events are ignored
func (c *Config) processPolkaUpdate(w http.ResponseWriter, r *http.Request) {
	type eventUser struct {
		UserID int `json:"user_id"`
	}

	type event struct {
		Event     string    `json:"event"`
		EventUser eventUser `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	eventData := event{}

	apiKey, err := fetchAPIToken(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	} else if apiKey != c.polkaAPIKey {
		errBody := errorBody{
			Error:     "received incorrect ApiKey",
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	// handle a decode error
	if err := decoder.Decode(&eventData); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if eventData.Event != "user.upgraded" {
		writeSuccessToPage(w, http.StatusOK, nil)
		return
	}

	if err := c.db.UpdateUserToRed(eventData.EventUser.UserID); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusNotFound,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusOK, nil)
}
