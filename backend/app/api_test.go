package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

func TestBindInput(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	type TestInput struct {
		Name   string      `json:"name" binding:"required"`
		Number json.Number `json:"number" binding:"required"`
	}

	t.Run("valid-input", func(t *testing.T) {
		body := `{"name": "test", "number": "123"}`
		c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		var input TestInput
		res := bindInput(c, &input)

		if !res {
			t.Errorf("Expected successful bind, got error: %v", c.Errors)
		}

		if input.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", input.Name)
		}

		if input.Number != "123" {
			t.Errorf("Expected number '123', got '%s'", input.Number)
		}
	})

	t.Run("missing-name", func(t *testing.T) {
		c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(`{"number": "123"}`))
		c.Request.Header.Set("Content-Type", "application/json")

		var input TestInput
		res := bindInput(c, &input)

		if res {
			t.Errorf("Expected bind error for missing name, got success")
		}

		if len(c.Errors) == 0 {
			t.Errorf("Expected error for missing name, got none")
		}
	})

	t.Run("parse-decimal", func(t *testing.T) {
		body := `{"name": "test", "number": "123.45"}`
		c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
		c.Request.Header.Set("Content-Type", "application/json")

		var input TestInput
		res := bindInput(c, &input)

		if !res {
			t.Errorf("Expected successful bind, got error: %v", c.Errors)
		}

		if input.Number != "123.45" {
			t.Errorf("Expected number '123.45', got '%s'", input.Number)
		}
	})
}

func sessionTestSetup() (*httptest.ResponseRecorder, *gin.Context) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", nil)
	sessions.Sessions("macro_session", cookie.NewStore([]byte("super-duper-secret")))(c)

	return w, c
}

func TestInitSession(t *testing.T) {
	_, c := sessionTestSetup()

	s := Session{
		UserName: "testuser",
		Token:    "testtoken",
	}

	initSession(&s, c)

	session := sessions.Default(c)
	if session.Get("username") != s.UserName {
		t.Errorf("Expected session username '%s', got '%s'", s.UserName, session.Get("username"))
	}

	if session.Get("token") != s.Token {
		t.Errorf("Expected session token '%s', got '%s'", s.Token, session.Get("token"))
	}
}

func TestClearSession(t *testing.T) {
	_, c := sessionTestSetup()

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	clearSession(c)

	session := sessions.Default(c)
	if session.Get("username") != nil {
		t.Errorf("Expected session username to be cleared, got '%s'", session.Get("username"))
	}

	if session.Get("token") != nil {
		t.Errorf("Expected session token to be cleared, got '%s'", session.Get("token"))
	}
}

func TestGetSessionToken(t *testing.T) {
	_, c := sessionTestSetup()

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	token, success := getSessionToken(c)
	if !success {
		t.Errorf("Expected successful session token retrieval, got failure")
	}

	if token != "testtoken" {
		t.Errorf("Expected session token 'testtoken', got '%s'", token)
	}

	s.Clear()
	s.Save()

	_, success = getSessionToken(c)
	if success {
		t.Errorf("Expected failure retrieving token from cleared session, got success")
	}
}

func apiTestSetup() (*API, *httptest.ResponseRecorder, *gin.Context) {
	db := Database{SetupTestDB()}
	api := API{db: &db}
	w, c := sessionTestSetup()

	return &api, w, c
}

func TestCreateSession(t *testing.T) {
	api, _, c := apiTestSetup()
	defer api.db.Close()

	api.createSession("testuser", c)

	session := sessions.Default(c)
	if session.Get("username") != "testuser" {
		t.Errorf("Expected session username 'testuser', got '%s'", session.Get("username"))
	}

	if session.Get("token") == nil {
		t.Errorf("Expected session token to be set, got nil")
	}
}

func TestHandleSessionValidation(t *testing.T) {

	// Test with no session
	t.Run("no-session", func(t *testing.T) {
		api, _, c := apiTestSetup()
		defer api.db.Close()

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status %d for no session, got %d", http.StatusUnauthorized,
				c.Writer.Status())
		}
	})

	// Test with invalid session token
	t.Run("invalid-session", func(t *testing.T) {
		api, _, c := apiTestSetup()
		defer api.db.Close()

		session := sessions.Default(c)
		session.Set("username", "testuser")
		session.Set("token", "invalidtoken")
		session.Save()

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status %d for invalid session, got %d", http.StatusUnauthorized,
				c.Writer.Status())
		}
	})

	// Test with valid session
	t.Run("valid-session", func(t *testing.T) {
		api, _, c := apiTestSetup()
		defer api.db.Close()

		api.createSession("testuser123", c)

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status %d for valid session, got %d", http.StatusOK, c.Writer.Status())
		}
	})
}
