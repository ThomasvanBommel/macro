package main

import (
	"log/slog"
)

// log.go defines logging utilities and conventions for the application.
func InitLogger() {
	// slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// })))
}

// Trace logs function entry and exit; defer its return value at call sites.
func Trace(fn string, args ...any) func() {
	slog.Debug(">> "+fn, args...)

	return func() {
		slog.Debug("<< "+fn, args...)
	}
}

// FatalOnError logs and panics when err is non-nil.
func FatalOnError(err error, msg string) {
	if err != nil {
		slog.Error(msg, "error", err)
		panic(err)
	}
}

// FatalIf logs and panics when cond is true.
func FatalIf(cond bool, msg string) {
	if cond {
		slog.Error(msg)
		panic(msg)
	}
}
