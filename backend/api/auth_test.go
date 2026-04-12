package api

import (
	"encoding/json"
	"macro/db"
	"net/http"
	"testing"

	"github.com/gin-contrib/sessions"
)

// TestInitSession tests the initSession function to ensure it properly sets session data.
func TestInitSession(t *testing.T) {
	_, c := newContext(t, "")

	s := db.Session{UserName: "testuser", Token: "testtoken"}
	initSession(&s, c)

	session := sessions.Default(c)
	if session.Get("username") != s.UserName {
		t.Errorf("Expected session username %q, got %q", s.UserName, session.Get("username"))
	}
	if session.Get("token") != s.Token {
		t.Errorf("Expected session token %q, got %q", s.Token, session.Get("token"))
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
		t.Errorf("Expected session username to be cleared, got %q", session.Get("username"))
	}
	if session.Get("token") != nil {
		t.Errorf("Expected session token to be cleared, got %q", session.Get("token"))
	}
}

// TestGetSessionToken tests the getSessionToken function for valid and invalid session states.
func TestGetSessionToken(t *testing.T) {
	_, c := newContext(t, "")

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	token, ok := getSessionToken(c)
	if !ok {
		t.Fatal("Expected successful session token retrieval, got failure")
	}
	if token != "testtoken" {
		t.Errorf("Expected session token %q, got %q", "testtoken", token)
	}

	s.Clear()
	s.Save()

	_, ok = getSessionToken(c)
	if ok {
		t.Fatal("Expected failure retrieving token from cleared session, got success")
	}
}

// TestCreateSession tests the createSession function to ensure it creates a session and sets the correct cookie values.
func TestCreateSession(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")

	_, c := newContext(t, "")
	if err := api.createSession("testuser", c); err != nil {
		t.Fatalf("createSession returned error: %v", err)
	}

	session := sessions.Default(c)
	if session.Get("username") != "testuser" {
		t.Errorf("Expected session username %q, got %q", "testuser", session.Get("username"))
	}
	if session.Get("token") == nil {
		t.Fatal("Expected session token to be set, got nil")
	}
}

func TestHandleSessionValidation(t *testing.T) {
	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", nil, api.handleSessionValidation)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", &db.Session{UserName: "testuser", Token: "invalidtoken"}, api.handleSessionValidation)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser123")
		s := mustCreateSession(t, api, "testuser123")

		w, _ := executeHandler(t, "", s, api.handleSessionValidation)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var res struct {
			Username string `json:"username" binding:"required"`
			Expires  string `json:"expires" binding:"required"`
		}

		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if res.Username != "testuser123" {
			t.Errorf("Expected username %q in response, got %q", "testuser123", res.Username)
		}

		if res.Expires != s.Expires {
			t.Errorf("Expected expires %q in response, got %q", s.Expires, res.Expires)
		}
	})
}

func TestHandleRegistration(t *testing.T) {
	t.Run("missing-field", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "newuser"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("valid-input", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "newuser", "password": "password123"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("duplicate-username", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		executeHandler(t, `{"name": "newuser", "password": "password123"}`, nil, api.handleRegisterUser)
		w, _ := executeHandler(t, `{"name": "newuser", "password": "other-password123"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})
}

func TestHandleLoginUser(t *testing.T) {
	t.Run("invalid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "nonexistentuser", "password": "wrongpassword"}`, nil, api.handleLoginUser)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser")

		w, c := executeHandler(t, `{"name": "testuser", "password": "password123"}`, nil, api.handleLoginUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if sessions.Default(c).Get("token") == nil {
			t.Fatal("Expected login to initialize a session token")
		}
	})
}

func TestHandleLogoutUser(t *testing.T) {
	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", nil, api.handleLogoutUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.Len() != 0 {
			t.Errorf("Expected empty response body, got %s", w.Body.String())
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser")
		s := mustCreateSession(t, api, "testuser")

		w, c := executeHandler(t, "", s, api.handleLogoutUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if _, ok := getSessionToken(c); ok {
			t.Fatal("Expected logout to clear the session token")
		}
	})
}
