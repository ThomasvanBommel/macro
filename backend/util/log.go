package util

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Logger is a struct that holds the Gin router and provides methods for initializing logging.
type Logger struct{}

// Middleware returns a Gin middleware function that logs HTTP requests and responses using slog.
func (l *Logger) Middleware() gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		// Level: slog.LevelDebug,
	}))

	slog.SetDefault(logger)

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		session := sessions.Default(c)
		id := session.Get("id")
		token := session.Get("token")

		if id == nil {
			id = GenerateRandomHexString(10)
			session.Set("id", id)
			session.Save()
		}

		error, _ := c.Get("error")

		slog.Info("HTTP request",
			"method", method,
			"path", path,
			"status", status,
			"duration_ns", duration.Nanoseconds(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", clientIP,
			"session_id", id,
			"session_token", token,
			"username", session.Get("username"),
			"error", error,
		)
	}
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
