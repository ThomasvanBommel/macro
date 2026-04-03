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

// setupTestDB creates an in-memory SQLite database and applies migrations for testing.
// func setupTestDB() *sql.DB {
// 	db, err := sql.Open("sqlite", ":memory:?cache=shared")
// 	FatalOnError(err, "Failed to open DB in memory")

// 	goose.SetBaseFS(testMigrationFiles)
// 	err = goose.SetDialect("sqlite3")
// 	FatalOnError(err, "Failed to set goose dialect")
// 	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

// 	return db
// }

// newAPI creates a new API instance with an in-memory database for testing.
func newAPI() *API {
	db, err := sql.Open("sqlite", ":memory:?cache=shared")
	FatalOnError(err, "Failed to open DB in memory")

	goose.SetBaseFS(testMigrationFiles)
	err = goose.SetDialect("sqlite3")
	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	return &API{db: &Database{db}}
}

// newContext creates a test Gin context with the given JSON body and session middleware applied.
func newContext(t *testing.T, body string) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	sessions.Sessions("macro_session", cookie.NewStore([]byte("super-duper-secret")))(c)

	return w, c
}

// TestBindInput tests the bindInput function for various input scenarios, including valid input,
// missing required fields, and parsing of decimal numbers.
func TestBindInput(t *testing.T) {
	type TestInput struct {
		Name   string      `json:"name" binding:"required"`
		Number json.Number `json:"number" binding:"required"`
	}

	t.Run("valid-input", func(t *testing.T) {
		_, c := newContext(t, `{"name": "test", "number": "123"}`)

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
		_, c := newContext(t, `{"number": "123"}`)

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
		_, c := newContext(t, `{"name": "test", "number": "123.45"}`)

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

// TestInitSession tests the initSession function to ensure it properly sets session data.
func TestInitSession(t *testing.T) {
	_, c := newContext(t, "")

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

// TestClearSession tests the clearSession function to ensure it properly clears session data.
func TestClearSession(t *testing.T) {
	_, c := newContext(t, "")

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

// TestGetSessionToken tests the getSessionToken function for valid and invalid session states.
func TestGetSessionToken(t *testing.T) {
	_, c := newContext(t, "")

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

// TestCreateSession tests the createSession function to ensure it creates a session and sets the
// correct cookie values.
func TestCreateSession(t *testing.T) {
	api := newAPI()
	defer api.db.Close()
	_, c := newContext(t, "")

	api.createSession("testuser", c)

	session := sessions.Default(c)
	if session.Get("username") != "testuser" {
		t.Errorf("Expected session username 'testuser', got '%s'", session.Get("username"))
	}

	if session.Get("token") == nil {
		t.Errorf("Expected session token to be set, got nil")
	}
}

// TestHandleSessionValidation tests the handleSessionValidation function for various session states
func TestHandleSessionValidation(t *testing.T) {

	t.Run("no-session", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		_, c := newContext(t, "")

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status %d for no session, got %d", http.StatusUnauthorized,
				c.Writer.Status())
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		_, c := newContext(t, "")

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

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		_, c := newContext(t, "")

		api.createSession("testuser123", c)

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status %d for valid session, got %d", http.StatusOK,
				c.Writer.Status())
		}
	})
}

// TestHandleRegistration tests the handleRegisterUser function for various input scenarios.
func TestHandleRegistration(t *testing.T) {

	t.Run("missing-field", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for missing field, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("valid-input", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser", "password": "password123"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid input, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("duplicate-username", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser", "password": "password123"}`)

		api.handleRegisterUser(c)
		w, c = newContext(t, `{"name": "newuser", "password": "other-password123"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d for duplicate username, got %d", http.StatusConflict,
				w.Code)
		}
	})
}

// TestHandleLoginUser tests the handleLoginUser function for valid and invalid credential scenarios
func TestHandleLoginUser(t *testing.T) {

	t.Run("invalid-credentials", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		w, c := newContext(t, `{"name": "nonexistentuser", "password": "wrongpassword"}`)

		api.handleLoginUser(c)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for invalid credentials, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("valid-credentials", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		body := `{"name": "testuser", "password": "password123"}`

		w, c := newContext(t, body)
		api.handleRegisterUser(c)

		w, c = newContext(t, body)
		api.handleLoginUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid credentials, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleLogoutUser(t *testing.T) {

	t.Run("no-session", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()
		w, c := newContext(t, "")

		api.handleLogoutUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for no session, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Body.String(), "No active session") {
			t.Errorf("Expected response body to contain 'No active session', got %s",
				w.Body.String())
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI()
		defer api.db.Close()

		s := Session{
			UserName: "testuser",
			Token:    "testtoken",
		}

		w, c := newContext(t, "")
		initSession(&s, c)

		api.handleLogoutUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid session, got %d", http.StatusOK, w.Code)
		}
	})
}
