package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var testMigrationFiles embed.FS

// init disables Gin and Goose logging for cleaner test output.
func init() {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = nil

	goose.SetLogger(goose.NopLogger())
}

// newAPI creates a new API instance with an in-memory database for testing.
func newAPI(t *testing.T) *API {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:?cache=shared")
	FatalOnError(err, "Failed to open DB in memory")

	goose.SetBaseFS(testMigrationFiles)
	err = goose.SetDialect("sqlite3")
	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	return InitAPI(gin.New(), &Database{db})
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

func executeHandler(t *testing.T, body string, session *Session, handler APIHandler) (*httptest.ResponseRecorder, *gin.Context) {
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
	if err := api.db.createUser(name, "password123"); err != nil {
		t.Fatalf("createUser(%q) failed: %v", name, err)
	}
}

func mustCreateSession(t *testing.T, api *API, name string) *Session {
	t.Helper()
	s, err := api.db.createSession(name)
	if err != nil {
		t.Fatalf("createSession(%q) failed: %v", name, err)
	}
	return s
}

func mustCreateFood(t *testing.T, api *API, token string, name string) *Food {
	t.Helper()
	f, err := api.db.createFoodByToken(CreateFoodParams{
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
		t.Fatalf("createFoodByToken(%q) failed: %v", name, err)
	}
	return f
}

func mustCreateEntry(t *testing.T, api *API, token string, foodID int, mealName string) *EntryWithFood {
	t.Helper()
	e, err := api.db.createEntryByToken(EntryParams{
		FoodId:   foodID,
		MealName: mealName,
		Date:     "1901-01-01",
		Servings: 1,
	}, token)
	if err != nil {
		t.Fatalf("createEntryByToken(%d) failed: %v", foodID, err)
	}
	return e
}

// TestInitSession tests the initSession function to ensure it properly sets session data.
func TestInitSession(t *testing.T) {
	_, c := newContext(t, "")

	s := Session{UserName: "testuser", Token: "testtoken"}
	initSession(&s, c)

	session := sessions.Default(c)
	if session.Get("username") != s.UserName {
		t.Errorf("Expected session username %q, got %q", s.UserName, session.Get("username"))
	}
	if session.Get("token") != s.Token {
		t.Errorf("Expected session token %q, got %q", s.Token, session.Get("token"))
	}
}

// TestClearSession tests the clearSession function to ensure it properly clears session data.
func TestClearSession(t *testing.T) {
	_, c := newContext(t, "")

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	clearSession(c)

	session := sessions.Default(c)
	if session.Get("username") != nil {
		t.Errorf("Expected session username to be cleared, got %q", session.Get("username"))
	}
	if session.Get("token") != nil {
		t.Errorf("Expected session token to be cleared, got %q", session.Get("token"))
	}
}

// TestGetSessionToken tests the getSessionToken function for valid and invalid session states.
func TestGetSessionToken(t *testing.T) {
	_, c := newContext(t, "")

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	token, ok := getSessionToken(c)
	if !ok {
		t.Fatal("Expected successful session token retrieval, got failure")
	}
	if token != "testtoken" {
		t.Errorf("Expected session token %q, got %q", "testtoken", token)
	}

	s.Clear()
	s.Save()

	_, ok = getSessionToken(c)
	if ok {
		t.Fatal("Expected failure retrieving token from cleared session, got success")
	}
}

// TestCreateSession tests the createSession function to ensure it creates a session and sets the correct cookie values.
func TestCreateSession(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")

	_, c := newContext(t, "")
	if err := api.createSession("testuser", c); err != nil {
		t.Fatalf("createSession returned error: %v", err)
	}

	session := sessions.Default(c)
	if session.Get("username") != "testuser" {
		t.Errorf("Expected session username %q, got %q", "testuser", session.Get("username"))
	}
	if session.Get("token") == nil {
		t.Fatal("Expected session token to be set, got nil")
	}
}

func TestHandleSessionValidation(t *testing.T) {
	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", nil, api.handleSessionValidation)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", &Session{UserName: "testuser", Token: "invalidtoken"}, api.handleSessionValidation)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser123")
		s := mustCreateSession(t, api, "testuser123")

		w, _ := executeHandler(t, "", s, api.handleSessionValidation)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleRegistration(t *testing.T) {
	t.Run("missing-field", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "newuser"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("valid-input", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "newuser", "password": "password123"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	t.Run("duplicate-username", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		executeHandler(t, `{"name": "newuser", "password": "password123"}`, nil, api.handleRegisterUser)
		w, _ := executeHandler(t, `{"name": "newuser", "password": "other-password123"}`, nil, api.handleRegisterUser)
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
		}
	})
}

func TestHandleLoginUser(t *testing.T) {
	t.Run("invalid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, `{"name": "nonexistentuser", "password": "wrongpassword"}`, nil, api.handleLoginUser)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser")

		w, c := executeHandler(t, `{"name": "testuser", "password": "password123"}`, nil, api.handleLoginUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if sessions.Default(c).Get("token") == nil {
			t.Fatal("Expected login to initialize a session token")
		}
	})
}

func TestHandleLogoutUser(t *testing.T) {
	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		w, _ := executeHandler(t, "", nil, api.handleLogoutUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if w.Body.Len() != 0 {
			t.Errorf("Expected empty response body, got %s", w.Body.String())
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		mustCreateUser(t, api, "testuser")
		s := mustCreateSession(t, api, "testuser")

		w, c := executeHandler(t, "", s, api.handleLogoutUser)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
		if _, ok := getSessionToken(c); ok {
			t.Fatal("Expected logout to clear the session token")
		}
	})
}

func TestHandleCreateEntry(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")
	f := mustCreateFood(t, api, s.Token, "Test Food")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, nil, api.handleCreateEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-servings", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 0}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, &Session{UserName: "testuser", Token: "invalidtoken"}, api.handleCreateEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-food-id", func(t *testing.T) {
		body := `{"food_id": 99999, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"food_id": ` + strconv.Itoa(f.ID) + `, "meal_name": "Lunch", "date": "1901-01-01", "servings": "1"}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleListUserEntries(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")
	f := mustCreateFood(t, api, s.Token, "Test Food")

	t.Run("invalid-user", func(t *testing.T) {
		w, _ := executeHandler(t, `{"name": "billybobbyjoetthorton", "date": "1901-01-01"}`, nil, api.handleListUserEntries)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("valid-user-no-entries", func(t *testing.T) {
		w, _ := executeHandler(t, `{"name": "testuser", "date": "1901-01-01"}`, nil, api.handleListUserEntries)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected 0 entries, got %d", len(entries))
		}
	})

	t.Run("valid-user-with-entries", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			_, err := api.db.createEntryByToken(EntryParams{FoodId: f.ID, MealName: "Breakfast", Date: "1901-01-01", Servings: 1}, s.Token)
			if err != nil {
				t.Fatalf("createEntryByToken failed: %v", err)
			}
		}

		w, _ := executeHandler(t, `{"name": "testuser", "date": "1901-01-01"}`, nil, api.handleListUserEntries)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(entries) != 3 {
			t.Errorf("Expected 3 entries, got %d", len(entries))
		}
	})
}

func TestHandleCreateFood(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", "protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, _ := executeHandler(t, body, nil, api.handleCreateFood)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-serving-count", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", "protein": "5", "fat": "2", "serving_size": "g", "serving_count": "0"}`
		w, _ := executeHandler(t, body, s, api.handleCreateFood)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("negative-macros", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "-100", "carbs": "-10", "protein": "-5", "fat": "-2", "serving_size": "g", "serving_count": "1"}`
		w, _ := executeHandler(t, body, s, api.handleCreateFood)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-token", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", "protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, _ := executeHandler(t, body, &Session{UserName: "testuser", Token: "invalidtoken"}, api.handleCreateFood)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", "protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, _ := executeHandler(t, body, s, api.handleCreateFood)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleListFoods(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")

	t.Run("no-foods", func(t *testing.T) {
		w, _ := executeHandler(t, "", nil, api.handleListFoods)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var foods []FoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &foods); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(foods) != 0 {
			t.Errorf("Expected 0 foods, got %d", len(foods))
		}
	})

	for i := 1; i <= 47; i++ {
		_, err := api.db.createFoodByToken(CreateFoodParams{
			Name:         "Test Food " + strconv.Itoa(i),
			Brand:        "Test Brand",
			Calories:     100 * i,
			Carbs:        10 * i,
			Protein:      5 * i,
			Fat:          2 * i,
			ServingSize:  "g",
			ServingCount: 1,
		}, s.Token)
		if err != nil {
			t.Fatalf("createFoodByToken failed: %v", err)
		}
	}

	t.Run("with-foods", func(t *testing.T) {
		w, _ := executeHandler(t, "", nil, api.handleListFoods)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var foods []FoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &foods); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if len(foods) != 47 {
			t.Errorf("Expected 47 foods, got %d", len(foods))
		}
	})
}

func TestHandleEditEntry(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")
	f1 := mustCreateFood(t, api, s.Token, "Food One")
	f2 := mustCreateFood(t, api, s.Token, "Food Two")
	e := mustCreateEntry(t, api, s.Token, f1.ID, "Breakfast")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "Lunch", "date": "1901-01-01", "servings": "1"}`
		w, _ := executeHandler(t, body, nil, api.handleEditEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-servings", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "Lunch", "date": "1901-01-01", "servings": "0"}`
		w, _ := executeHandler(t, body, s, api.handleEditEntry)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing-entry", func(t *testing.T) {
		body := `{"id": 99999, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "Lunch", "date": "1901-01-01", "servings": "1"}`
		w, _ := executeHandler(t, body, s, api.handleEditEntry)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "Lunch", "date": "1901-01-02", "servings": "2"}`
		w, _ := executeHandler(t, body, s, api.handleEditEntry)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var entry EntryWithFoodResponse
		if err := json.Unmarshal(w.Body.Bytes(), &entry); err != nil {
			t.Fatalf("Failed to parse response body: %v", err)
		}
		if entry.Food.ID != f2.ID {
			t.Errorf("Expected food ID %d, got %d", f2.ID, entry.Food.ID)
		}
		if entry.MealName != "Lunch" {
			t.Errorf("Expected meal name %q, got %q", "Lunch", entry.MealName)
		}
		if entry.Date != "1901-01-02T00:00:00Z" {
			t.Errorf("Expected date %q, got %q", "1901-01-02T00:00:00Z", entry.Date)
		}
	})
}

func TestHandleDeleteEntry(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")
	f := mustCreateFood(t, api, s.Token, "Food One")
	e := mustCreateEntry(t, api, s.Token, f.ID, "Breakfast")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `}`
		w, _ := executeHandler(t, body, nil, api.handleDeleteEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `}`
		w, _ := executeHandler(t, body, s, api.handleDeleteEntry)
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		entries, err := api.db.listUserEntriesWithFoodByNameAndDate("testuser", "1901-01-01")
		if err != nil {
			t.Fatalf("listUserEntriesWithFoodByNameAndDate failed: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected deleted entry to be removed, got %d remaining entries", len(entries))
		}
	})
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

		err := api.db.createUser("testuser", "password123")
		if err == nil {
			t.Fatal("Expected duplicate createUser to return an error")
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
