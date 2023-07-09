package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type user_raw struct {
	user_id  string
	email    string
	password string
}

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

func getUser(db *sql.DB) (*user_raw, error) {
	var raw user_raw

	query := `SELECT * FROM users`
	err := db.QueryRow(query).Scan(&raw.user_id, &raw.email, &raw.password)

	if err != nil {
		return nil, err
	} else {
		return &raw, nil
	}
}

func main() {
	var db *sql.DB
	var err error
	var raw *user_raw

	db, err = connectDB(":memory:")

	if err != nil {
		panic("데이터베이스가 연결되지 않았습니다.")
	}

	fmt.Println("DB Ready.")

	_, err = createUsersTable(db)

	if err != nil {
		panic("유저 테이블이 생성되지 않았습니다.")
	}

	fmt.Println("Table Created.")

	_, err = insertUser(db, "abc@gmail.com", "12345678765")

	if err != nil {
		panic("유저가 생성되지 않았습니다.")
	}
	fmt.Println("User Created.")

	raw, err = getUser(db)

	if err != nil {
		panic("유저를 불러오는데 실패했습니다.")
	}

	fmt.Println()
	fmt.Println("user_id", raw.user_id)
	fmt.Println("email", raw.email)
	fmt.Println("password", raw.password)

}
