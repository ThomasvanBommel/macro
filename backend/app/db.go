package main

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "/app/macro.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	log.Println("Running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	return db, nil
}

func registerUser(db *sql.DB, req RegisterUserRequest) (RegisterUserResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterUserResponse{}, err
	}

	result, err := db.Exec("INSERT INTO users (username, password_hash) VALUES (?, ?)", req.Username, string(hash))
	if err != nil {
		return RegisterUserResponse{}, err
	}

	id, _ := result.LastInsertId()
	return RegisterUserResponse{
		ID:       int(id),
		Username: req.Username,
		Message:  "User created successfully",
	}, nil
}
