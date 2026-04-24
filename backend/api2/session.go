package api2

import (
	"macro/env"
	"macro/util"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// SessionMiddleware configures cookie-backed sessions
func SessionMiddleware() gin.HandlerFunc {
	s := []byte(os.Getenv("SESSION_SECRET"))
	if len(s) == 0 {
		s = util.GenerateRandomBytes(128)
	}

	store := cookie.NewStore(s)
	store.Options(sessions.Options{
		MaxAge:   env.SESSION_TIMEOUT_SEC,
		Secure:   env.SECURE,
		HttpOnly: true,
	})

	return sessions.Sessions("macro_session", store)
}
