package api

import (
	"database/sql"
	"errors"
	"macro/db"
	"macro/env"
	"macro/util"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

// init disables Gin and Goose logging for cleaner test output.
func init() {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = nil

	goose.SetLogger(goose.NopLogger())
}

// newAPI creates a new API instance with an in-memory database for testing.
func newAPI(t *testing.T) *API {
	t.Helper()

	database, err := sql.Open("sqlite", ":memory:?cache=shared")
	util.FatalOnError(err, "Failed to open DB in memory")

	goose.SetBaseFS(env.MigrationFiles)
	err = goose.SetDialect("sqlite3")
	util.FatalOnError(err, "Failed to set goose dialect")
	util.FatalOnError(goose.Up(database, "migrations"), "Failed to apply migrations")

	// Enable foreign keys and set PRAGMAs for performance
	_, err = database.Exec(`
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = OFF;
		PRAGMA cache_size = -10000;
	`)
	util.FatalOnError(err, "Failed to set PRAGMAs")

	// Double check foreign keys are enabled
	var fkEnabled int
	err = database.QueryRow(`PRAGMA foreign_keys;`).Scan(&fkEnabled)
	util.FatalOnError(err, "Failed to check foreign_keys PRAGMA")
	util.FatalIf(fkEnabled != 1, "Foreign keys not enabled in SQLite")

	return Init(gin.New(), &db.Database{DB: database})
}

// newContext creates a test Gin context with the given JSON body and session middleware applied.
func newContext(t *testing.T, body string) (*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest("POST", "/", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	sessions.Sessions("macro_session", cookie.NewStore([]byte("super-duper-secret")))(c)

	return w, c
}

func executeHandler(t *testing.T, body string, session *db.Session, handler APIHandler) (
	*httptest.ResponseRecorder, *gin.Context) {
	t.Helper()

	w, c := newContext(t, body)
	if session != nil {
		initSession(session, c)
	}

	Wrap(handler)(c)
	c.Writer.WriteHeaderNow()
	return w, c
}

func TestDBError(t *testing.T) {
	t.Run("not-found-default-message", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		res := api.DBError(sql.ErrNoRows)
		if res.Status() != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, res.Status())
		}

		err, ok := res.(*ErrorResponse)
		if !ok {
			t.Fatalf("Expected ErrorResponse result, got %T", res)
		}
		if err.Error() != "Not found." {
			t.Errorf("Expected message %q, got %q", "Not found.", err.Error())
		}

		if res.Payload() != (JSONError{Error: "Not found."}) {
			t.Errorf("Expected payload %v, got %v", JSONError{Error: "Not found."}, res.Payload())
		}
	})

	t.Run("conflict-override-message", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser")

		err := api.db.CreateUser("testuser", "password123")
		if err == nil {
			t.Fatal("Expected duplicate CreateUser to return an error")
		}

		res := api.DBError(err, DBErrorMessages{Unique: "Username already exists."})
		if res.Status() != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, res.Status())
		}

		err, ok := res.(*ErrorResponse)
		if !ok {
			t.Fatalf("Expected ErrorResponse result, got %T", res)
		}
		if err.Error() != "Username already exists." {
			t.Errorf("Expected message %q, got %q", "Username already exists.", err.Error())
		}

		if res.Payload() != (JSONError{Error: "Username already exists."}) {
			t.Errorf("Expected payload %v, got %v", JSONError{Error: "Username already exists."}, res.Payload())
		}
	})

	t.Run("default-internal-error-message", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		res := api.DBError(errors.New("boom"))
		if res.Status() != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, res.Status())
		}

		err, ok := res.(*ErrorResponse)
		if !ok {
			t.Fatalf("Expected ErrorResponse result, got %T", res)
		}
		if err.Error() != "An unexpected database error occurred." {
			t.Errorf("Expected message %q, got %q", "An unexpected database error occurred.", err.Error())
		}

		if res.Payload() != (JSONError{Error: "Internal server error."}) {
			t.Errorf("Expected payload %v, got %v", JSONError{Error: "Internal server error."}, res.Payload())
		}
	})
}
