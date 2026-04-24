package db2

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"macro/errs"
	"net/http"
	"unicode"

	"golang.org/x/crypto/bcrypt"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

const MAX_USERNAME_LENGTH = 25
const MAX_PASSWORD_LENGTH = 500

func ValidateUsername(username string) error {
	if len(username) == 0 {
		return errs.New(http.StatusUnprocessableEntity, "Username cannot be empty", nil)
	}

	n := len(username)
	if n > MAX_USERNAME_LENGTH {
		return errs.New(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("Username of %d characters exceeds the limit of %d",
				n, MAX_USERNAME_LENGTH), nil)
	}

	for _, c := range username {
		if !unicode.IsDigit(c) && !unicode.IsLetter(c) && c != '_' && c != '-' && c != '.' {
			return errs.New(http.StatusUnprocessableEntity,
				"Username may only contain letters, digits, underscores, hyphens, and periods", nil)
		}
	}

	return nil
}

func HashPassword(password string) ([]byte, error) {
	if len(password) == 0 {
		return nil, errs.New(http.StatusUnprocessableEntity, "Password cannot be empty", nil)
	}

	n := len(password)
	if n > MAX_PASSWORD_LENGTH {
		return nil, errs.New(
			http.StatusUnprocessableEntity,
			fmt.Sprintf("Password of %d characters exceeds the limit of %d",
				n, MAX_PASSWORD_LENGTH), nil)
	}

	h := sha256.New()
	_, err := h.Write([]byte(password))
	if err != nil {
		return nil, errs.Internal(err)
	}

	return h.Sum(nil), nil
}

func IsUniqueConstraintError(err error) bool {
	if e, ok := errors.AsType[*sqlite.Error](err); ok {
		return e.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE ||
			e.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY
	}

	return false
}

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
		return errs.Internal(err)
	}

	_, err = db.Exec("INSERT INTO user (name, password_hash) VALUES (?, ?)", username, hash)
	if err != nil {
		if IsUniqueConstraintError(err) {
			return errs.New(http.StatusConflict, "Username is already taken", err)
		}

		return errs.Internal(err)
	}

	return nil
}
