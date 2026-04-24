package api2

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustCreateUser(t *testing.T, db *DB, username, password string) {
	t.Helper()
	err := db.CreateUser(username, password)
	require.NoError(t, err, "failed to create user %q: %v", username, err)
}

func TestRegister(t *testing.T) {
	r, db := newEngine()
	defer db.Close()

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid input", `{"username": "testuser", "password": "testpass"}`, http.StatusCreated},
		{"missing username", `{"password": "testpass"}`, http.StatusUnprocessableEntity},
		{"empty username", `{"username": "", "password": "123"}`, http.StatusUnprocessableEntity},
		{"missing password", `{"username": "testuser"}`, http.StatusUnprocessableEntity},
		{"empty password", `{"username": "t", "password": ""}`, http.StatusUnprocessableEntity},
		{"duplicate user", `{"username": "testuser", "password": "testpass"}`, http.StatusConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req := newRequest("POST", "/api/auth/register", tt.body)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")
		})
	}
}

func TestLogin(t *testing.T) {
	r, db := newEngine()
	defer db.Close()

	mustCreateUser(t, db, "usr", "123")

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid credentials", `{"username": "usr", "password": "123"}`, http.StatusOK},
		{"invalid password", `{"username": "usr", "password": "111"}`, http.StatusUnauthorized},
		{"nonexistent user", `{"username": "nonexistent", "password": "any"}`, http.StatusNotFound},
		{"missing username", `{"password": "123"}`, http.StatusUnprocessableEntity},
		{"empty username", `{"username": "", "password": "123"}`, http.StatusUnprocessableEntity},
		{"missing password", `{"username": "usr"}`, http.StatusUnprocessableEntity},
		{"empty password", `{"username": "usr", "password": ""}`, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req := newRequest("POST", "/api/auth/login", tt.body)
			r.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")

			if tt.wantStatus == http.StatusOK {
				var out struct {
					Username string `json:"username"`
					Expires  string `json:"expires"`
				}

				require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out), "Unable to parse json")
				assert.Equal(t, "usr", out.Username, "Unexpected username in response")
				assert.NotEmpty(t, out.Expires, "Expires field should not be empty")
				assert.NotEmpty(t, w.Header().Get("Set-Cookie"), "Set-Cookie missing")
			}
		})
	}
}

func TestSessionInfo(t *testing.T) {
	r, db := newEngine()
	defer db.Close()

	mustCreateUser(t, db, "usr", "123")

	// No active session
	w, req := newRequest("GET", "/api/auth/session", "")
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Expected 401 for no active session")

	// Log in to create a session
	w, req = newRequest("POST", "/api/auth/login", `{"username": "usr", "password": "123"}`)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Login failed")

	cookie := w.Header().Get("Set-Cookie")
	require.NotEmpty(t, cookie, "Set-Cookie header missing after login")

	// Get session info with valid cookie
	w, req = newRequest("GET", "/api/auth/session", "")
	req.Header.Set("Cookie", cookie)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Failed to get session info")

	var out struct {
		Username string `json:"username"`
		Expires  string `json:"expires"`
	}

	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &out), "Unable to parse json")
	assert.Equal(t, "usr", out.Username, "Unexpected username in session info")
	assert.NotEmpty(t, out.Expires, "Expires field should not be empty")
}

func TestLogout(t *testing.T) {
	r, db := newEngine()
	defer db.Close()

	mustCreateUser(t, db, "usr", "123")

	// First, log in to set the session cookie
	w, req := newRequest("POST", "/api/auth/login", `{"username": "usr", "password": "123"}`)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Login failed")

	cookie := w.Header().Get("Set-Cookie")
	require.NotEmpty(t, cookie, "Set-Cookie header missing after login")

	// Now log out using the same cookie
	w, req = newRequest("POST", "/api/auth/logout", "")
	req.Header.Set("Cookie", cookie)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Logout failed")
	assert.Contains(t, w.Header().Get("Set-Cookie"), "Max-Age=0", "Cookie should be cleared")
}
