package main

import (
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func connectDB(dbFile string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", dbFile)

	if err != nil {
		return nil, err
	}

	if db.Ping() != nil {
		return nil, err
	}

	return db, nil
}

func createUsersTable(db *sql.DB) (sql.Result, error) {
	query := `
	  CREATE TABLE users (
		user_id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	  )`

	result, err := db.Exec(query)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func generateHashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func insertUser(db *sql.DB, email string, password string) (sql.Result, error) {
	user_id, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	hash, err := generateHashPassword(password)

	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO users
		(user_id, email, password)
		VALUES(?, ?, ?)
	`

	result, err := db.Exec(query, user_id, email, hash)

	if err != nil {
		return nil, err
	}

	return result, nil
}
