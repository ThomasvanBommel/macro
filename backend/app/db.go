package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"time"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type Handler struct {
	db *sql.DB
}

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

	ticker := time.NewTicker(10 * time.Minute)
	go func() {
		for range ticker.C {
			log.Println("Cleaning up expired sessions...")
			if err := cleanupExpiredSessions(db); err != nil {
				log.Printf("Failed to clean up expired sessions: %v", err)
			}
		}
	}()

	return db, nil
}

func createUser(db *sql.DB, username string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users (username, password_hash) VALUES (?, ?);", username, string(hash))
	if err != nil {
		return err
	}

	return nil
}

func createSession(db *sql.DB, user_id int) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	q := `
		INSERT INTO sessions (user_id, token, expires_at)
		VALUES (?, ?, datetime('now', '+1 hour'));
	`

	_, err := db.Exec(q, user_id, token)
	if err != nil {
		return "", err
	}

	return token, nil
}

func cleanupExpiredSessions(db *sql.DB) error {
	q := `
		DELETE
		FROM sessions
		WHERE expires_at <= datetime('now');
	`
	_, err := db.Exec(q)
	return err
}

func loginUser(db *sql.DB, username string, password string) (string, error) {
	var user_id int
	var hash string
	err := db.QueryRow(
		"SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&user_id, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", nil
	}

	return createSession(db, user_id)
}

func selectSession(db *sql.DB, token string) (*Session, error) {
	var s Session
	q := `
		SELECT user_id, username, expires_at
		FROM users
		JOIN sessions
		  ON users.id = user_id
		WHERE token = ?
		  AND expires_at > datetime('now');
	`
	err := db.QueryRow(q, token).Scan(&s.User_ID, &s.Username, &s.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func deleteSession(db *sql.DB, token string) error {
	q := `
		DELETE
		FROM sessions
		WHERE token = ?;
	`
	_, err := db.Exec(q, token)
	return err
}

func insertFood(db *sql.DB, food Food) error {
	q := `
		INSERT INTO food (name, brand, created_by, calories, carbs, protein, fat)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`
	_, err := db.Exec(q, food.Name, food.Brand, food.CreatedBy, food.Calories, food.Carbs, 
		food.Protein, food.Fat)
	return err
}