package api2

import (
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
			w := newRecorder(r, "POST", "/api/auth/register", tt.body)
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
			w := newRecorder(r, "POST", "/api/auth/login", tt.body)
			assert.Equal(t, tt.wantStatus, w.Code, "Unexpected status code")
		})
	}
}
