package db2

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"macro/env"
	"macro/errs"
	"macro/util"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const MAX_USERNAME_LENGTH = 25
const MAX_PASSWORD_LENGTH = 500

var USERNAME_REGEX = regexp.MustCompile("^[\\w-.]+$")

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return errs.BadInput(nil, "Username cannot be empty")
	}

	n := len(username)
	if n > MAX_USERNAME_LENGTH {
		msg := fmt.Sprintf("Username cannot exceed %d characters", MAX_USERNAME_LENGTH)
		return errs.BadInput(nil, msg).With("username", username)
	}

	if !USERNAME_REGEX.MatchString(username) {
		return errs.BadInput(nil, "Invalid username. Use: A-Z, a-z, 0-9, _, -, .")
	}

	return nil
}

func HashPassword(password string) (hash []byte, err error) {
	n := len(password)
	if n == 0 {
		return nil, errs.BadInput(nil, "Password cannot be empty")
	}

	if n > MAX_PASSWORD_LENGTH {
		return nil, errs.BadInput(nil, "Password cannot exceed 500 characters").With("length", n)
	}

	h := sha256.New()
	if _, err := h.Write([]byte(password)); err != nil {
		return nil, errs.Unexpected(err, "")
	}

	return h.Sum(nil), nil
}

// CreateUser creates a new user with the given username and password
func (db *Database) CreateUser(username, password string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return err
	}

	hash, err = bcrypt.GenerateFromPassword(hash, -1)
	if err != nil {
		return errs.Unexpected(err, "")
	}

	_, err = db.Exec("INSERT INTO user (name, password_hash) VALUES (?, ?)", username, hash)
	if err != nil {
		if IsUniqueConstraintError(err) {
			return errs.Conflict(err, "Username is already taken").With("name", username)
		}

		return errs.Unexpected(err, "")
	}

	return nil
}

// CreateSession creates a new session for the user with the given username and password
// It returns an session token, expiry time, and a possible error
func (db *Database) CreateSession(username, password string) (tok, exp string, err error) {
	if err := ValidateUsername(username); err != nil {
		return "", "", err
	}

	hash, err := HashPassword(password)
	if err != nil {
		return "", "", err
	}

	var db_hash string
	q := `SELECT password_hash FROM user WHERE name = ?`
	if err := db.QueryRow(q, username).Scan(&db_hash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errs.BadCredentials(err, "Invalid username or password").
				With("name", username)
		}

		return "", "", errs.Unexpected(err, "")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(db_hash), hash); err != nil {
		return "", "", errs.BadCredentials(err, "Invalid username or password").
			With("name", username)
	}

	token := util.GenerateRandomHexString(32)
	var expiry string
	q = `INSERT INTO session (token, user_name, expires) 
	     VALUES (?, ?, datetime('now', '+' || ? || ' seconds'))
		 RETURNING expires`
	if err := db.QueryRow(q, token, username, env.SESSION_TIMEOUT_SEC).Scan(&expiry); err != nil {
		return "", "", errs.Unexpected(err, "")
	}

	return token, expiry, nil
}

// GetSessionInfo retrieves the username and expiry time for the session with the given token
func (db *Database) GetSessionInfo(token string) (usr, exp string, err error) {
	var username, expiry string
	q := `SELECT user_name, expires
	      FROM session
		  WHERE token = ?
		  AND expires > datetime('now')`
	if err := db.QueryRow(q, token).Scan(&username, &expiry); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errs.NotAuthorized(err, "No active session")
		}

		return "", "", errs.Unexpected(err, "")
	}

	return username, expiry, nil
}
