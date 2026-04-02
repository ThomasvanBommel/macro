package main

import (
	"database/sql"
	"embed"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var testMigrationFiles embed.FS

// SetupTestDB creates an in-memory SQLite database and applies migrations for testing.
func SetupTestDB() *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?cache=shared")
	FatalOnError(err, "Failed to open DB in memory")

	goose.SetBaseFS(testMigrationFiles)
	err = goose.SetDialect("sqlite3")
	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	return db
}

func TestHandleRegisterUser(t *testing.T) {
	db := Database{SetupTestDB()}
	defer db.Close()

	r := gin.New()
	InitLogger(r)
	InitAPI(r, &db)

	t.Run("Missing body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Missing field - password", func(t *testing.T) {
		body := `{"name": "testuser"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("Successful registration", func(t *testing.T) {
		body := `{"name": "testuser", "password": "testpass"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("Duplicate username", func(t *testing.T) {
		body := `{"name": "testuser", "password": "testpass"}`
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}
