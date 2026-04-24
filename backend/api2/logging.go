package api2

import (
	"log/slog"
	"macro/util"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LoggingMiddleware returns a Gin middleware function that logs HTTP requests and responses
func LoggingMiddleware(level slog.Level) gin.HandlerFunc {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
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
			id = util.GenerateRandomHexString(10)
			session.Set("id", id)
			session.Save()
		}

		args := []any{
			"method", method,
			"path", path,
			"status", status,
			"duration_ns", duration.Nanoseconds(),
			"duration_ms", duration.Milliseconds(),
			"client_ip", clientIP,
			"session_id", id,
			"session_token", token,
			"username", session.Get("username"),
		}

		lvl := slog.LevelInfo
		if len(c.Errors) > 0 {
			lvl = slog.LevelError
			args = append(args, "error_count", len(c.Errors))
			args = append(args, "errors", c.Errors.ByType(gin.ErrorTypePrivate).String())
		}

		slog.Log(c, lvl, "HTTP request", args...)
	}
}
