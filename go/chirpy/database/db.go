package database

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

// DB is the struct to point at our database.json file
type DB struct {
	path string
	mux  sync.RWMutex
}

// Chirp is the default struct for each individual chirp within the system
type Chirp struct {
	ID       int    `json:"id"`
	AuthorID int    `json:"author_id"`
	Body     string `json:"body"`
}

// User is the default struct to represent an individual user in the database
type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// UserWithPassword is a superset struct for a given user that appends their password (bcrypt-hashed)
type UserWithPassword struct {
	User
	PasswordHash []byte `json:"password"`
}

// RevokedToken is the struct to consume the revoked tokens within the database
type RevokedToken struct {
	RevokedAt time.Time `json:"revoked_at"`
	Token     string    `json:"token"`
}

// DBStructure is the interface to render the database
type DBStructure struct {
	Chirps        map[int]Chirp            `json:"chirps"`
	Users         map[int]UserWithPassword `json:"users"`
	RevokedTokens map[int]RevokedToken     `json:"revoked_tokens"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	if path == "" {
		path = "./database.json"
	}

	db := &DB{path: path}
	if err := db.reassureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email string, password []byte) (User, error) {
	var user UserWithPassword

	nextUserID, err := db.getNextUserID()
	if err != nil {
		return user.User, err
	}

	user.ID = nextUserID
	user.Email = email
	user.PasswordHash = password

	dbStructure, err := db.loadDB()
	if err != nil {
		return user.User, err
	}

	if _, ok := dbStructure.Users[user.ID]; ok {
		return user.User, fmt.Errorf("expected unique userID but found duplicate at %d", user.ID)
	}

	// check if user already exists and throw an error if they do
	for _, existingUser := range dbStructure.Users {
		if user.Email == existingUser.Email {
			return user.User, fmt.Errorf("found duplicate user with email %s", user.Email)
		}
	}

	dbStructure.Users[user.ID] = user
	return user.User, db.writeDB(dbStructure)
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
	var chirp Chirp

	nextChirpID, err := db.getNextChirpID()
	if err != nil {
		return chirp, err
	}

	chirp.ID = nextChirpID
	chirp.AuthorID = authorID
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

// getNextRevokedTokenID is a helper function to determine the next revoked token ID from the database
func (db *DB) getNextRevokedTokenID() (int, error) {
	revokedTokens, err := db.GetRevokedTokens()
	if err != nil {
		return -1, err
	}

	if len(revokedTokens) == 0 {
		return 1, nil
	}

	ids := make([]int, len(revokedTokens))
	for idx, _ := range revokedTokens {
		ids = append(ids, idx)
	}

	// sort the IDs in descending order
	sort.Slice(ids, func(a, b int) bool {
		return ids[a] > ids[b]
	})

	return ids[0] + 1, nil
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

// GetUsersFull return all users in the database with hashed passwords
func (db *DB) GetUsersFull() ([]UserWithPassword, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]UserWithPassword, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	return users, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user.User)
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

// DeleteChirp will remove the provided chirp from the database
func (db *DB) DeleteChirp(chirpToDelete Chirp) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	for idx, chirp := range dbStructure.Chirps {
		if chirp.ID == chirpToDelete.ID {
			delete(dbStructure.Chirps, idx)
		}
	}

	return db.writeDB(dbStructure)
}

// GetRevokedTokens retrieves the set of revoked tokens from the database
func (db *DB) GetRevokedTokens() ([]RevokedToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	revokedTokens := make([]RevokedToken, 0, len(dbStructure.RevokedTokens))
	for _, revokedToken := range dbStructure.RevokedTokens {
		revokedTokens = append(revokedTokens, revokedToken)
	}

	return revokedTokens, nil
}

// RevokeToken will revoke the provided token from the database
func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	nextID, err := db.getNextRevokedTokenID()
	if err != nil {
		return err
	}

	dbStructure.RevokedTokens[nextID] = RevokedToken{
		RevokedAt: time.Now(),
		Token:     token,
	}

	return db.writeDB(dbStructure)
}

// UpdateUser will update the existing user at userID with a new email/password combination
func (db *DB) UpdateUser(userID int, email string, passwordHash []byte) (User, error) {
	var user UserWithPassword

	user.ID = userID
	user.Email = email
	user.PasswordHash = passwordHash

	dbStructure, err := db.loadDB()
	if err != nil {
		return user.User, err
	}

	dbStructure.Users[user.ID] = user
	return user.User, db.writeDB(dbStructure)
}

// reassureDB creates a new database file if it doesn't exist
func (db *DB) reassureDB() error {
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

	if err := db.reassureDB(); err != nil {
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
		dbStructure.Users = make(map[int]UserWithPassword)
		dbStructure.RevokedTokens = make(map[int]RevokedToken)

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
