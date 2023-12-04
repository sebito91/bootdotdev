package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const chirpyAccess = "chirpy-access"
const chirpyRefresh = "chirpy-refresh"

// generateJWT is a helper function to generate a JWT based on the ID of the user and an expiration timeout (in seconds)
func (c *Config) generateJWT(issuer string, expiresInSeconds, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:    issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(expiresInSeconds))),
		Subject:   fmt.Sprintf("%d", id),
	})

	return token.SignedString([]byte(c.jwtSecret))
}

// fetchToken is a helper function to extract the JWT from a given request
func fetchToken(r *http.Request) (string, error) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("expected Authorization token, got %s", bearer)
	}

	return strings.TrimPrefix(bearer, "Bearer "), nil
}

// decodeJWT takes an encoded JWT from a request and parses it into a *jwt.Token object for use within functions;
// will return an error if the token is invalid
func (c *Config) decodeJWT(bearer string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(bearer, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method, expected HMAC, got %v", t.Header["alg"])
		}

		return []byte(c.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// revokeToken will take in a given refresh token from a user and record the token as revoked
// within the database. Any subsequent use of the token will be blocked as Unauthorized.
func (c *Config) revokeToken(w http.ResponseWriter, r *http.Request) {
	bearer, err := fetchToken(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if err := c.db.RevokeToken(bearer); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusOK, nil)
}

// refreshToken will take in a refresh token from a given user, ensure it is valid, and output a new
// access token valid for one hour. We must ensure that the refresh token is not revoked and that it
// is still for a valid user.
func (c *Config) refreshToken(w http.ResponseWriter, r *http.Request) {
	claims, respCode, err := c.fetchClaims(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: respCode,
		}

		errBody.writeErrorToPage(w)
		return
	}

	if issuer, iErr := claims.GetIssuer(); err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("could not fetch issuer from token: %s", iErr),
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	} else if issuer != chirpyRefresh {
		errBody := errorBody{
			Error:     fmt.Sprintf("expected refresh token, got %s", issuer),
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	revokedTokens, err := c.db.GetRevokedTokens()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	bearer, err := fetchToken(r)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusUnauthorized,
		}

		errBody.writeErrorToPage(w)
		return
	}

	for _, revokedToken := range revokedTokens {
		if bearer == revokedToken.Token {
			errBody := errorBody{
				Error:     fmt.Sprintf("provided refresh token was revoked at: %s", revokedToken.RevokedAt),
				errorCode: http.StatusUnauthorized,
			}

			errBody.writeErrorToPage(w)
			return
		}
	}

	// we passed the checks for revoked token, so let's generate a new 60m token
	idString, err := claims.GetSubject()
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("%s", err),
			errorCode: http.StatusBadRequest,
		}

		errBody.writeErrorToPage(w)
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("could not convert userID to string: %s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	token, err := c.generateJWT(chirpyAccess, (60 * 60), id)
	if err != nil {
		errBody := errorBody{
			Error:     fmt.Sprintf("token generate: %s", err),
			errorCode: http.StatusInternalServerError,
		}

		errBody.writeErrorToPage(w)
		return
	}

	writeSuccessToPage(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token})
}
