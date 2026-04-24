package api2

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
