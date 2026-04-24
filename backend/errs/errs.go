package errs

import "net/http"

type Error struct {
	code int    // http status code
	msg  string // public message
	err  error  // internal error
}

func (e *Error) Code() int     { return e.code }
func (e *Error) Error() string { return e.msg }
func (e *Error) Unwrap() error { return e.err }

func New(code int, msg string, err error) *Error { return &Error{code, msg, err} }
func Internal(err error) *Error {
	return New(http.StatusInternalServerError, "An internal error occurred", err)
}
