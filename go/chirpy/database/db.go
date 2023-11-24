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
	path string
	mux  sync.RWMutex
}

// Chirp is the default struct for each individual chirp within the system
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// User is the default struct to represent an individual user in the database
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// DBStructure is the interface to render the database
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string) (User, error) {
	var user User

	nextUserID, err := db.getNextUserID()
	if err != nil {
		return user, err
	}

	user.ID = nextUserID
	user.Email = email

	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	if _, ok := dbStructure.Users[user.ID]; ok {
		return user, fmt.Errorf("expected unique userID but found duplicate at %d", user.ID)
	}

	dbStructure.Users[user.ID] = user
	return user, db.writeDB(dbStructure)
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	var chirp Chirp

	nextChirpID, err := db.getNextChirpID()
	if err != nil {
		return chirp, err
	}

	chirp.ID = nextChirpID
	chirp.Body = body

	dbStructure, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	if _, ok := dbStructure.Chirps[chirp.ID]; ok {
		return chirp, fmt.Errorf("expected unique chirpID but found duplicate at %d", chirp.ID)
	}

	dbStructure.Chirps[chirp.ID] = chirp
	return chirp, db.writeDB(dbStructure)
}

// getNextUserID is a helper function to determine the next user's ID from the database
func (db *DB) getNextUserID() (int, error) {
	users, err := db.GetUsers()
	if err != nil {
		return -1, err
	}

	if len(users) == 0 {
		return 1, nil
	}

	ids := make([]int, len(users))
	for _, user := range users {
		ids = append(ids, user.ID)
	}

	// sort the IDs in descending order
	sort.Slice(ids, func(a, b int) bool {
		return ids[a] > ids[b]
	})

	return ids[0] + 1, nil
}

// getNextChirpID is a helper function to determine the next chirp's ID from the database
func (db *DB) getNextChirpID() (int, error) {
	chirps, err := db.GetChirps()
	if err != nil {
		return -1, err
	}

	if len(chirps) == 0 {
		return 1, nil
	}

	ids := make([]int, len(chirps))
	for _, chirp := range chirps {
		ids = append(ids, chirp.ID)
	}

	// sort the IDs in descending order
	sort.Slice(ids, func(a, b int) bool {
		return ids[a] > ids[b]
	})

	return ids[0] + 1, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(db.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if fErr := f.Close(); fErr != nil {
		return fErr
	}

	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	var dbStructure DBStructure

	if err := db.ensureDB(); err != nil {
		return dbStructure, err
	}

	db.mux.RLock()
	defer db.mux.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}

	if len(data) == 0 {
		dbStructure.Chirps = make(map[int]Chirp)
		dbStructure.Users = make(map[int]User)

		return dbStructure, nil
	}

	err = json.Unmarshal(data, &dbStructure)
	return dbStructure, err
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}
