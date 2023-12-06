package api

import (
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/sebito91/bootdotdev/go/chirpy/database"
)

// Config is a local struct to keep track of site visits
// NOTE: this value is in-memory only and will persist for the duration of the server
type Config struct {
	fileserverHits int
	jwtSecret      string
	db             *database.DB
	mux            sync.RWMutex
}

// errorBody is a struct used for returning a JSON-based error code/string
type errorBody struct {
	Error     string `json:"error"`
	errorCode int
}

// NewConfig returns a new instance of the Config
func NewConfig() (*Config, error) {
	db, err := database.NewDB("")
	if err != nil {
		return nil, err
	}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{db: db, jwtSecret: os.Getenv("JWT_SECRET")}, nil
}

// GetAPI returns the router for the /api endpoint
func (c *Config) GetAPI() chi.Router {
	r := chi.NewRouter()

	r.Get("/healthz", readinessEndpoint)

	r.Route("/chirps", func(r chi.Router) {
		r.Get("/", c.getChirps)
		r.Post("/", c.writeChirp)

		r.Route("/{chirpID}", func(r chi.Router) {
			r.Get("/", c.getChirpByID)
			r.Delete("/", c.deleteChirpByID)
		})
	})

	r.Route("/users", func(r chi.Router) {
		r.Get("/", c.getUsers)
		r.Post("/", c.writeUser)
		r.Put("/", c.updateUser)

		r.Route("/{userID}", func(r chi.Router) {
			r.Get("/", c.getUserByID)
		})
	})

	// token-related exercises
	r.Post("/login", c.loginUser)
	r.Post("/refresh", c.refreshToken)
	r.Post("/revoke", c.revokeToken)

	// webhooks
	r.Route("/polka", func(r chi.Router) {
		r.Post("/webhooks", c.processPolkaUpdate)
	})

	return r
}
