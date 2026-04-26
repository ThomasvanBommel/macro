package errs

import (
	"fmt"
	"log/slog"
	"maps"
	"net/http"
)

type ErrCode int

const (
	ErrBadCredentials ErrCode = iota
	ErrNotAuthorized
	ErrUnexpected
	ErrBadInput
	ErrConflict
	ErrUnknown
)

var statusMap = map[ErrCode]int{
	ErrBadCredentials: http.StatusUnauthorized,
	ErrNotAuthorized:  http.StatusUnauthorized,
	ErrUnexpected:     http.StatusInternalServerError,
	ErrBadInput:       http.StatusUnprocessableEntity,
	ErrConflict:       http.StatusConflict,
	ErrUnknown:        http.StatusInternalServerError,
}

type Error struct {
	private  error
	public   string
	code     ErrCode
	metadata map[string]any
}

func (e *Error) Error() string            { return e.private.Error() }
func (e *Error) Unwrap() error            { return e.private }
func (e *Error) Public() string           { return e.public }
func (e *Error) Code() ErrCode            { return e.code }
func (e *Error) Metadata() map[string]any { return e.metadata }

// HTTPStatus returns the recommended HTTP status code for this error.
func (e *Error) HTTPStatus() int {
	if status, ok := statusMap[e.code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// With adds a key-value pair to the error's metadata.
func (e *Error) With(key string, value any) *Error {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}
	e.metadata[key] = value
	return e
}

// New creates a new error with the given message and code.
//
// err is the underlying error and is considered private, not exposed to the client.
//
// msg is the public message that could be sent to the client.
//
// metadata is merged in the order they are provided, overwriting old values. See [maps.Copy]
//
//   - If err is nil, a new error is created with the given message and code.
//   - If err is not nil, it is wrapped with the given message and code.
//   - If metadata is provided, it is added to the error's metadata.
//   - If the code is not recognized, ErrUnknown is used.
func New(err error, msg string, code ErrCode, metadata ...map[string]any) *Error {
	if _, ok := statusMap[code]; !ok {
		slog.Error("Unknown error code", "code", code)
		code = ErrUnknown
	}

	// Create error if not provided, wrap otherwise
	if err == nil {
		err = fmt.Errorf("%s", msg)
	} else {
		err = fmt.Errorf("%v: %w", msg, err)
	}

	// Initialize metadata if provided
	e := &Error{err, msg, code, nil}
	if len(metadata) > 0 {
		e.metadata = metadata[0]

		// Merge additional metadata
		for _, m := range metadata[1:] {
			maps.Copy(e.metadata, m)
		}
	}

	return e
}

// BadCredentials is a shortcut for [New] with ErrBadCredentials.
func BadCredentials(err error, msg string, metadata ...map[string]any) *Error {
	return New(err, msg, ErrBadCredentials, metadata...)
}

// NotAuthorized is a shortcut for [New] with ErrNotAuthorized.
func NotAuthorized(err error, msg string, metadata ...map[string]any) *Error {
	return New(err, msg, ErrNotAuthorized, metadata...)
}

// Unexpected is a shortcut for [New] with ErrUnexpected.
//   - If msg is empty, "Unexpected error" is used.
func Unexpected(err error, msg string, metadata ...map[string]any) *Error {
	if msg == "" {
		msg = "Unexpected error"
	}

	return New(err, msg, ErrUnexpected, metadata...)
}

// BadInput is a shortcut for [New] with ErrBadInput.
func BadInput(err error, msg string, metadata ...map[string]any) *Error {
	return New(err, msg, ErrBadInput, metadata...)
}

// Conflict is a shortcut for [New] with ErrConflict.
func Conflict(err error, msg string, metadata ...map[string]any) *Error {
	return New(err, msg, ErrConflict, metadata...)
}
