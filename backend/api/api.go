package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"macro/db"
	"macro/env"
	"macro/util"
	"math"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

// Wrap converts an APIHandler to a gin.HandlerFunc, handling APIResponse and errors uniformly.
func Wrap(fn APIHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		r := fn(c)
		if r == nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		if err, ok := r.(error); ok && err != nil {
			c.Error(err)
		}

		payload := r.Payload()
		if payload == nil {
			c.Status(r.Status())
			return
		}

		c.JSON(r.Status(), payload)
	}
}

func Init(r *gin.Engine, db *db.Database) *API {
	// Configure cookie-backed sessions
	s := []byte(os.Getenv("SESSION_SECRET"))
	if len(s) == 0 {
		s = util.GenerateRandomBytes(128)
	}

	store := cookie.NewStore(s)
	store.Options(sessions.Options{
		MaxAge:   SESSION_TIMEOUT_SEC,
		HttpOnly: true,
		Secure:   env.SECURE,
	})

	r.Use(sessions.Sessions("macro_session", store))

	err := func(code int, ovr ...string) func(error, ...string) APIResponse {
		return func(e error, msg ...string) APIResponse {
			err := e

			if len(msg) == 0 && len(ovr) > 0 {
				msg = ovr
			}

			if err == nil && len(msg) > 0 {
				err = errors.New(msg[0])
			} else if err == nil {
				err = errors.New("An unexpected error occurred.")
			}

			return &ErrorResponse{code, err, msg}
		}
	}

	// Initialize API with DB connection and response helpers
	api := &API{
		db: db,

		OK:             func(d any) APIResponse { return &DataResponse{200, d} },
		Message:        func(m string) APIResponse { return &DataResponse{200, gin.H{"message": m}} },
		Created:        func(d any) APIResponse { return &DataResponse{201, d} },
		Accepted:       func() APIResponse { return &DataResponse{202, nil} },
		NoContent:      func() APIResponse { return &DataResponse{204, nil} },
		PartialContent: func(d any) APIResponse { return &DataResponse{206, d} },

		BadRequest:      err(400),
		Unauthorized:    err(401),
		Forbidden:       err(403),
		NotFound:        err(404),
		Conflict:        err(409),
		Unprocessable:   err(422),
		TooManyRequests: err(429),

		InternalServerError: err(500, "Internal server error."),
		NotImplemented:      err(501, "Not implemented."),
	}

	// Register API routes
	// auth := r.Group("/api/auth")
	// auth.POST("/register", Wrap(api.handleRegisterUser))
	// auth.POST("/login", Wrap(api.handleLoginUser))
	// auth.POST("/validate-session", Wrap(api.handleSessionValidation))
	// auth.POST("/logout", Wrap(api.handleLogoutUser))

	a := r.Group("/api")
	a.POST("/register", Wrap(api.handleRegisterUser))
	a.POST("/login", Wrap(api.handleLoginUser))
	a.POST("/logout", Wrap(api.handleLogoutUser))
	a.POST("/validate-session", Wrap(api.handleSessionValidation))
	a.POST("/food", Wrap(api.handleCreateFood))
	a.POST("/foods", Wrap(api.handleListFoods))
	a.POST("/entry", Wrap(api.handleCreateEntry))
	a.POST("/entries", Wrap(api.handleListUserEntries))
	a.POST("/diary", Wrap(api.handleGetDiary))
	a.POST("/food/search", Wrap(api.handleFoodSearch))

	a.POST("/entry/edit", Wrap(api.handleEditEntry))
	a.POST("/entry/delete", Wrap(api.handleDeleteEntry))

	a.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, nil) })

	return api
}

// DBErrorMessages holds overridable response messages for known DB error conditions.
// Zero-value fields fall back to their default.
type DBErrorMessages struct{ Unique, NotFound, Busy, ForeignKey, Default string }

// DBError maps a database error to an APIResponse, with optional message overrides.
func (a *API) DBError(err error, msgs ...DBErrorMessages) APIResponse {
	m := DBErrorMessages{
		"Already exists.",
		"Not found.",
		"Database is busy, please try again.",
		"Invalid foreign key.",
		"An unexpected database error occurred.",
	}

	if len(msgs) > 0 {
		o := msgs[0]

		if o.Unique != "" {
			m.Unique = o.Unique
		}
		if o.NotFound != "" {
			m.NotFound = o.NotFound
		}
		if o.Busy != "" {
			m.Busy = o.Busy
		}
		if o.ForeignKey != "" {
			m.ForeignKey = o.ForeignKey
		}
		if o.Default != "" {
			m.Default = o.Default
		}
	}

	if errors.Is(err, sql.ErrNoRows) {
		return a.NotFound(errors.New(m.NotFound))
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		switch sqliteErr.Code() {
		case sqlite3.SQLITE_CONSTRAINT_UNIQUE, sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY:
			return a.Conflict(errors.New(m.Unique))
		case sqlite3.SQLITE_BUSY:
			return a.TooManyRequests(errors.New(m.Busy))
		case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
			return a.Unprocessable(errors.New(m.ForeignKey))
		}
	}

	return a.InternalServerError(errors.New(m.Default))
}

// roundToTwoDecimalPlaces rounds a float to two decimal places.
func roundToTwoDecimalPlaces(f float64) float64 {
	return math.Round(f*100) / 100
}

// scaleForClient converts an integer DB value to a float for client responses.
func scaleForClient(n int) float64 {
	return roundToTwoDecimalPlaces(float64(n) / 100)
}

// scaleForStorage converts a JSON number to an integer for DB storage.
func scaleForStorage(n json.Number) int {
	f, err := n.Float64()
	if err != nil {
		return 0
	}

	return int(roundToTwoDecimalPlaces(f) * 100)
}
