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

// DBStructure is the interface to render the database
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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
	fmt.Printf("inside CreateChirp, going to post: %s...\n", body)
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

// getNextChirpID is a helper function to determine the next chirp's ID from the database
func (db *DB) getNextChirpID() (int, error) {
	fmt.Println("inside getNextChirpID...")
	chirps, err := db.GetChirps()
	if err != nil {
		fmt.Println("we had an error in getNextChirpID")
		return -1, err
	}

	fmt.Printf("got chirps: %+v\n", chirps)

	if len(chirps) == 0 {
		fmt.Println("this is our first chirp!")
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

	fmt.Printf("we're sending back chirpID %d\n", ids[0]+1)

	return ids[0] + 1, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	fmt.Println("inside GetChirps...")
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	fmt.Println("inside ensureDB...")
	db.mux.Lock()
	defer fmt.Println("leaving ensureDB...")
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
	fmt.Println("inside loadDB...")
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
		return dbStructure, nil
	}

	err = json.Unmarshal(data, &dbStructure)
	fmt.Printf("what is the error in loadDB? %s\n", err)
	return dbStructure, err
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	fmt.Println("inside writeDB...")
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	return os.WriteFile(db.path, data, 0644)
}
