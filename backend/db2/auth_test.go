package db2

import (
	"errors"
	"fmt"
	"macro/errs"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

func randUser(length int, seed int) string {
	username := make([]byte, length)
	for i := range username {
		switch seed * i * 7 % 5 {
		case 0:
			username[i] = 'a' + byte(i%26)
		case 1:
			username[i] = '0' + byte(i%10)
		case 2:
			username[i] = '_'
		case 3:
			username[i] = '-'
		case 4:
			username[i] = '.'
		}
	}

	return string(username)
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"", true},
		{"a", false},
		{"valid_user...---", false},
		{"invalid user", true},
		{"usernameiswayyyyytofriggenlong", true},
		{randUser(MAX_USERNAME_LENGTH, 1), false},
		{randUser(MAX_USERNAME_LENGTH+1, 2), true},
		{"user!", true},
	}

	for _, tt := range tests {
		err := ValidateUsername(tt.name)
		assert.Equal(t, tt.wantErr, err != nil, "input: %q", tt.name)
	}
}

func TestHashPassword(t *testing.T) {
	tests := []struct {
		password string
		wantErr  bool
	}{
		{"", true},
		{"short", false},
		{string(make([]byte, MAX_PASSWORD_LENGTH)), false},
		{string(make([]byte, MAX_PASSWORD_LENGTH+1)), true},
	}

	for _, tt := range tests {
		hash, err := HashPassword(tt.password)
		assert.Equal(t, tt.wantErr, err != nil, "input: %q", tt.password)
		if err == nil {
			assert.Len(t, hash, 32, "input: %q", tt.password)
		}
	}
}

type mockSQLiteError struct {
	msg  string
	code int
}

func (e *mockSQLiteError) Error() string { return e.msg }
func (e *mockSQLiteError) Code() int     { return e.code }

func TestIsUniqueConstraintError(t *testing.T) {
	tests := []struct {
		code     int
		expected bool
	}{
		{sqlite3.SQLITE_CONSTRAINT_UNIQUE, true},
		{sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY, true},
		{123, false},
	}

	for _, tt := range tests {
		if e, ok := errors.AsType[*sqlite.Error](&mockSQLiteError{code: tt.code}); ok {
			assert.Equal(t, tt.expected, IsUniqueConstraintError(e), "code: %d", tt.code)
		}
	}

	b := IsUniqueConstraintError(fmt.Errorf("some other error"))
	assert.False(t, b, "returned true for non-sqlite error")
}

func TestCreateUser(t *testing.T) {
	db := newDatabase()
	defer db.Close()

	tests := []struct {
		username string
		password string
		wantErr  bool
		httpCode int
	}{
		{"valid_user", "valid_password", false, 0},
		{"valid_user", "another_password", true, http.StatusConflict},
		{"", "password", true, http.StatusUnprocessableEntity},
		{randUser(MAX_USERNAME_LENGTH, 3), "password", false, 0},
		{randUser(MAX_USERNAME_LENGTH+1, 4), "password", true, http.StatusUnprocessableEntity},
		{"user_with_long_password", string(make([]byte, MAX_PASSWORD_LENGTH)), false, 0},
		{"user_with_longer_password", string(make([]byte, MAX_PASSWORD_LENGTH+1)), true, http.StatusUnprocessableEntity},
	}

	for _, tt := range tests {
		err := db.CreateUser(tt.username, tt.password)
		assert.Equal(t, tt.wantErr, err != nil, "input: %q / %q", tt.username, tt.password)
		if err != nil {
			e, ok := err.(*errs.Error)
			require.True(t, ok, "expected *errs.Error, got %T", err)
			assert.Equal(t, tt.httpCode, e.Code(), "input: %q / %q", tt.username, tt.password)
		}
	}
}
