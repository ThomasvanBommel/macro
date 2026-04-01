package main

import (
	"log/slog"
)

// Trace is a helper function for logging the entry and exit of a function. It takes the function
// name and optional key-value pairs as arguments, and returns a function that should be deferred
// at the beginning of the traced function. When the returned function is called (at the end of the
// traced function), it logs the exit of the function along with any provided context.
func Trace(fn string, args ...any) func() {
	slog.Debug(">> "+fn, args...)

	return func() {
		slog.Debug("<< "+fn, args...)
	}
}

// FatalOnError is a helper function that checks if an error is non-nil and, if so, logs the
// provided message along with the error details, and then panics. This is useful for quickly
// handling errors that should not be ignored and need to be addressed immediately.
func FatalOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		panic(err)
	}
}

// FatalIf is a helper function that checks if a given condition is true and, if so, logs the
// provided message and panics. This is useful for asserting conditions that must hold true for the
// application to function correctly, and provides a way to quickly identify and handle unexpected
// states in the code.
func FatalIf(cond bool, msg string) {
	if cond {
		slog.Error(msg)
		panic(msg)
	}
}
