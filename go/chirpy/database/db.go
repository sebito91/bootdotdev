package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
)

// DB is the struct to point at our database.json file
type DB struct {
	path        string
	dbStructure DBStructure

	mux *sync.RWMutex
}

// Chirp is the default struct for each individual chirp within the system
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// DBStructure is the interface to render the database
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`

	nextChirpID int
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	if path == "" {
		path = "./database.json"
	}

	db := &DB{path: path}
	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if err == os.ErrNotExist {
		return os.WriteFile(db.path, nil, 0666)
	}

	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	var dbStructure DBStructure

	if err := db.ensureDB(); err != nil {
		return dbStructure, err
	}

	data, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}

	err = json.Unmarshal(data, &dbStructure)
	if dbStructure.Chirps != nil {
		ids := make([]int, len(dbStructure.Chirps))

		for _, chirp := range dbStructure.Chirps {
			ids = append(ids, chirp.ID)
		}

		if len(ids) == 0 {
			return dbStructure, fmt.Errorf("could not read any items from the db at %s, found %d records", db.path, len(ids))
		}

		sort.Slice(ids, func(a, b int) bool {
			return ids[a] > ids[b]
		})
		dbStructure.nextChirpID = ids[0]
	}

	return dbStructure, err
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {}
