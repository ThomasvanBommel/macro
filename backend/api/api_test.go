package api

import (
	"database/sql"
	"errors"
	"macro/db"
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

	goose.SetBaseFS(db.MigrationFiles)
	err = goose.SetDialect("sqlite3")
	util.FatalOnError(err, "Failed to set goose dialect")
	util.FatalOnError(goose.Up(database, "migrations"), "Failed to apply migrations")

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

func mustCreateUser(t *testing.T, api *API, name string) {
	t.Helper()
	if err := api.db.CreateUser(name, "password123"); err != nil {
		t.Fatalf("CreateUser(%q) failed: %v", name, err)
	}
}

func mustCreateSession(t *testing.T, api *API, name string) *db.Session {
	t.Helper()
	s, err := api.db.CreateSession(name)
	if err != nil {
		t.Fatalf("CreateSession(%q) failed: %v", name, err)
	}
	return s
}

func mustCreateFood(t *testing.T, api *API, token string, name string) *db.Food {
	t.Helper()
	f, err := api.db.CreateFoodByToken(db.FoodParams{
		Name:         name,
		Brand:        "Test Brand",
		Calories:     100,
		Carbs:        10,
		Protein:      5,
		Fat:          2,
		ServingSize:  "g",
		ServingCount: 1,
	}, token)
	if err != nil {
		t.Fatalf("CreateFoodByToken(%q) failed: %v", name, err)
	}
	return f
}

func mustCreateEntry(t *testing.T, api *API, token string, foodID int, mealName string) *db.EntryWithFood {
	t.Helper()
	e, err := api.db.CreateEntryByToken(db.EntryParams{
		FoodId:   foodID,
		MealName: mealName,
		Date:     "1901-01-01",
		Servings: 1,
	}, token)
	if err != nil {
		t.Fatalf("CreateEntryByToken(%d) failed: %v", foodID, err)
	}
	return e
}

func TestDBError(t *testing.T) {
	t.Run("not-found-default-message", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		res := api.DBError(sql.ErrNoRows)
		if res.Status() != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, res.Status())
		}

		err, ok := res.Result().(error)
		if !ok {
			t.Fatalf("Expected error result, got %T", res.Result())
		}
		if err.Error() != "Not found." {
			t.Errorf("Expected message %q, got %q", "Not found.", err.Error())
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

		mappedErr, ok := res.Result().(error)
		if !ok {
			t.Fatalf("Expected error result, got %T", res.Result())
		}
		if mappedErr.Error() != "Username already exists." {
			t.Errorf("Expected message %q, got %q", "Username already exists.", mappedErr.Error())
		}
	})

	t.Run("default-internal-error-message", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		res := api.DBError(errors.New("boom"))
		if res.Status() != http.StatusInternalServerError {
			t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, res.Status())
		}

		err, ok := res.Result().(error)
		if !ok {
			t.Fatalf("Expected error result, got %T", res.Result())
		}
		if err.Error() != "An unexpected database error occurred." {
			t.Errorf("Expected message %q, got %q", "An unexpected database error occurred.", err.Error())
		}
	})
}
