package main

import (
	"database/sql"
	"embed"
	"encoding/json"
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

// newAPI creates a new API instance with an in-memory database for testing.
func newAPI(t *testing.T) *API {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:?cache=shared")
	FatalOnError(err, "Failed to open DB in memory")

	goose.SetBaseFS(testMigrationFiles)
	err = goose.SetDialect("sqlite3")
	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	return &API{db: &Database{db}}
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

// TestBindInput tests the bindInput function for various input scenarios, including valid input,
// missing required fields, and parsing of decimal numbers.
func TestBindInput(t *testing.T) {
	type TestInput struct {
		Name   string      `json:"name" binding:"required"`
		Number json.Number `json:"number" binding:"required"`
	}

	t.Run("valid-input", func(t *testing.T) {
		_, c := newContext(t, `{"name": "test", "number": "123"}`)

		var input TestInput
		res := bindInput(c, &input)

		if !res {
			t.Errorf("Expected successful bind, got error: %v", c.Errors)
		}

		if input.Name != "test" {
			t.Errorf("Expected name 'test', got '%s'", input.Name)
		}

		if input.Number != "123" {
			t.Errorf("Expected number '123', got '%s'", input.Number)
		}
	})

	t.Run("missing-name", func(t *testing.T) {
		_, c := newContext(t, `{"number": "123"}`)

		var input TestInput
		res := bindInput(c, &input)

		if res {
			t.Errorf("Expected bind error for missing name, got success")
		}

		if len(c.Errors) == 0 {
			t.Errorf("Expected error for missing name, got none")
		}
	})

	t.Run("parse-decimal", func(t *testing.T) {
		_, c := newContext(t, `{"name": "test", "number": "123.45"}`)

		var input TestInput
		res := bindInput(c, &input)

		if !res {
			t.Errorf("Expected successful bind, got error: %v", c.Errors)
		}

		if input.Number != "123.45" {
			t.Errorf("Expected number '123.45', got '%s'", input.Number)
		}
	})
}

// TestInitSession tests the initSession function to ensure it properly sets session data.
func TestInitSession(t *testing.T) {
	_, c := newContext(t, "")

	s := Session{
		UserName: "testuser",
		Token:    "testtoken",
	}

	initSession(&s, c)

	session := sessions.Default(c)
	if session.Get("username") != s.UserName {
		t.Errorf("Expected session username '%s', got '%s'", s.UserName, session.Get("username"))
	}

	if session.Get("token") != s.Token {
		t.Errorf("Expected session token '%s', got '%s'", s.Token, session.Get("token"))
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
		t.Errorf("Expected session username to be cleared, got '%s'", session.Get("username"))
	}

	if session.Get("token") != nil {
		t.Errorf("Expected session token to be cleared, got '%s'", session.Get("token"))
	}
}

// TestGetSessionToken tests the getSessionToken function for valid and invalid session states.
func TestGetSessionToken(t *testing.T) {
	_, c := newContext(t, "")

	s := sessions.Default(c)
	s.Set("username", "testuser")
	s.Set("token", "testtoken")
	s.Save()

	token, success := getSessionToken(c)
	if !success {
		t.Errorf("Expected successful session token retrieval, got failure")
	}

	if token != "testtoken" {
		t.Errorf("Expected session token 'testtoken', got '%s'", token)
	}

	s.Clear()
	s.Save()

	_, success = getSessionToken(c)
	if success {
		t.Errorf("Expected failure retrieving token from cleared session, got success")
	}
}

// TestCreateSession tests the createSession function to ensure it creates a session and sets the
// correct cookie values.
func TestCreateSession(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	_, c := newContext(t, "")

	api.createSession("testuser", c)

	session := sessions.Default(c)
	if session.Get("username") != "testuser" {
		t.Errorf("Expected session username 'testuser', got '%s'", session.Get("username"))
	}

	if session.Get("token") == nil {
		t.Errorf("Expected session token to be set, got nil")
	}
}

// TestHandleSessionValidation tests the handleSessionValidation function for various session states
func TestHandleSessionValidation(t *testing.T) {

	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		_, c := newContext(t, "")

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status %d for no session, got %d", http.StatusUnauthorized,
				c.Writer.Status())
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		_, c := newContext(t, "")

		session := sessions.Default(c)
		session.Set("username", "testuser")
		session.Set("token", "invalidtoken")
		session.Save()

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusUnauthorized {
			t.Errorf("Expected status %d for invalid session, got %d", http.StatusUnauthorized,
				c.Writer.Status())
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		_, c := newContext(t, "")

		api.createSession("testuser123", c)

		api.handleSessionValidation(c)
		if c.Writer.Status() != http.StatusOK {
			t.Errorf("Expected status %d for valid session, got %d", http.StatusOK,
				c.Writer.Status())
		}
	})
}

// TestHandleRegistration tests the handleRegisterUser function for various input scenarios.
func TestHandleRegistration(t *testing.T) {

	t.Run("missing-field", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for missing field, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("valid-input", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser", "password": "password123"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid input, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("duplicate-username", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		w, c := newContext(t, `{"name": "newuser", "password": "password123"}`)

		api.handleRegisterUser(c)
		w, c = newContext(t, `{"name": "newuser", "password": "other-password123"}`)

		api.handleRegisterUser(c)
		if w.Code != http.StatusConflict {
			t.Errorf("Expected status %d for duplicate username, got %d", http.StatusConflict,
				w.Code)
		}
	})
}

// TestHandleLoginUser tests the handleLoginUser function for valid and invalid credential scenarios
func TestHandleLoginUser(t *testing.T) {

	t.Run("invalid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		w, c := newContext(t, `{"name": "nonexistentuser", "password": "wrongpassword"}`)

		api.handleLoginUser(c)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for invalid credentials, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("valid-credentials", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		body := `{"name": "testuser", "password": "password123"}`

		w, c := newContext(t, body)
		api.handleRegisterUser(c)

		w, c = newContext(t, body)
		api.handleLoginUser(c)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid credentials, got %d", http.StatusOK, w.Code)
		}
	})
}

// TestHandleLogoutUser tests the handleLogoutUser function for both no-session and valid-session
// scenarios.
func TestHandleLogoutUser(t *testing.T) {

	t.Run("no-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()
		w, c := newContext(t, "")

		api.handleLogoutUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for no session, got %d", http.StatusOK, w.Code)
		}

		if !strings.Contains(w.Body.String(), "No active session") {
			t.Errorf("Expected response body to contain 'No active session', got %s",
				w.Body.String())
		}
	})

	t.Run("valid-session", func(t *testing.T) {
		api := newAPI(t)
		defer api.db.Close()

		s := Session{
			UserName: "testuser",
			Token:    "testtoken",
		}

		w, c := newContext(t, "")
		initSession(&s, c)

		api.handleLogoutUser(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid session, got %d", http.StatusOK, w.Code)
		}
	})
}

// TestHandleCreateEntry tests the handleCreateEntry function for various scenarios
func TestHandleCreateEntry(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()

	// test user & session
	api.db.createUser("testuser", "password123")
	s, _ := api.db.createSession("testuser")

	// test food
	f, _ := api.db.createFoodByToken(
		CreateFoodParams{
			Name:         "Test Food",
			Brand:        "Test Brand",
			Calories:     100,
			Carbs:        10,
			Protein:      5,
			Fat:          2,
			ServingSize:  "g",
			ServingCount: 1,
		}, s.Token)

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, c := newContext(t, body)

		api.handleCreateEntry(c)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for unauthorized request, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("invalid-servings", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 0}`
		w, c := newContext(t, body)
		initSession(s, c)

		api.handleCreateEntry(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid servings, got %d", http.StatusBadRequest,
				w.Code)
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		s := Session{
			UserName: "testuser",
			Token:    "invalidtoken",
		}

		body := `{"food_id": 1, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, c := newContext(t, body)
		initSession(&s, c)

		api.handleCreateEntry(c)
		if w.Code != http.StatusUnauthorized {
			t.Error(c.Get("error"))
			t.Errorf("Expected status %d for invalid session, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("invalid-food-id", func(t *testing.T) {
		body := `{"food_id": 99999, "meal_name": "Lunch", "date": "1901-01-01", "servings": 1}`
		w, c := newContext(t, body)
		api.createSession("testuser", c)

		api.handleCreateEntry(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid food id, got %d", http.StatusBadRequest,
				w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"food_id": ` + strconv.Itoa(f.ID) + `, "meal_name": "Lunch", "date": ` +
			`"1901-01-01", "servings": "1"}`
		w, c := newContext(t, body)
		api.createSession("testuser", c)

		api.handleCreateEntry(c)
		if w.Code != http.StatusOK {
			t.Error(c.Get("error"))
			t.Errorf("Expected status %d for valid request, got %d", http.StatusOK, w.Code)
		}
	})
}

// TestHandleListUserEntries tests the handleListUserEntries function for various user and entry
// scenarios.
func TestHandleListUserEntries(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()

	// test user & session
	api.db.createUser("testuser", "password123")
	s, _ := api.db.createSession("testuser")

	// test food
	f, _ := api.db.createFoodByToken(
		CreateFoodParams{
			Name:         "Test Food",
			Brand:        "Test Brand",
			Calories:     100,
			Carbs:        10,
			Protein:      5,
			Fat:          2,
			ServingSize:  "g",
			ServingCount: 1,
		}, s.Token)

	t.Run("invalid-user", func(t *testing.T) {
		body := `{"name": "billybobbyjoetthorton","date": "1901-01-01"}`
		w, c := newContext(t, body)

		api.handleListUserEntries(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for invalid user request, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		err := json.Unmarshal(w.Body.Bytes(), &entries)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("Expected 0 entries for new user, got %d", len(entries))
		}
	})

	t.Run("valid-user-no-entries", func(t *testing.T) {
		body := `{"name": "testuser","date": "1901-01-01"}`
		w, c := newContext(t, body)

		api.handleListUserEntries(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid user request, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		err := json.Unmarshal(w.Body.Bytes(), &entries)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		if len(entries) != 0 {
			t.Errorf("Expected 0 entries for new user, got %d", len(entries))
		}
	})

	t.Run("valid-user-with-entries", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			api.db.createEntryByToken(CreateEntryParams{
				FoodId:   f.ID,
				MealName: "Breakfast",
				Date:     "1901-01-01",
				Servings: 1,
			}, s.Token)
		}

		body := `{"name": "testuser","date": "1901-01-01"}`
		w, c := newContext(t, body)

		api.handleListUserEntries(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for valid user request, got %d", http.StatusOK, w.Code)
		}

		var entries []EntryWithFoodResponse
		err := json.Unmarshal(w.Body.Bytes(), &entries)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		if len(entries) != 3 {
			t.Errorf("Expected 3 entries for user, got %d", len(entries))
		}
	})
}

// TestHandleCreateFood tests the handleCreateFood function for various scenarios
func TestHandleCreateFood(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()

	// test user & session
	api.db.createUser("testuser", "password123")
	s, _ := api.db.createSession("testuser")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", ` +
			`"protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, c := newContext(t, body)

		api.handleCreateFood(c)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d for unauthorized request, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("invalid-serving-count", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", ` +
			`"protein": "5", "fat": "2", "serving_size": "g", "serving_count": "0"}`
		w, c := newContext(t, body)
		initSession(s, c)

		api.handleCreateFood(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid serving count, got %d", http.StatusBadRequest,
				w.Code)
		}
	})

	t.Run("negative-macros", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "-100", "carbs": "-10", ` +
			`"protein": "-5", "fat": "-2", "serving_size": "g", "serving_count": "1"}`
		w, c := newContext(t, body)
		initSession(s, c)

		api.handleCreateFood(c)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for negative macros, got %d", http.StatusBadRequest,
				w.Code)
		}
	})

	t.Run("invalid-token", func(t *testing.T) {
		s := Session{
			UserName: "testuser",
			Token:    "invalidtoken",
		}

		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", ` +
			`"protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, c := newContext(t, body)
		initSession(&s, c)

		api.handleCreateFood(c)
		if w.Code != http.StatusUnauthorized {
			t.Error(c.Get("error"))
			t.Errorf("Expected status %d for invalid session token, got %d", http.StatusUnauthorized,
				w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"name": "Test Food", "brand": "Test Brand", "calories": "100", "carbs": "10", ` +
			`"protein": "5", "fat": "2", "serving_size": "g", "serving_count": "1"}`
		w, c := newContext(t, body)
		initSession(s, c)

		api.handleCreateFood(c)
		if w.Code != http.StatusOK {
			t.Error(c.Get("error"))
			t.Errorf("Expected status %d for valid request, got %d", http.StatusOK, w.Code)
		}
	})
}

func TestHandleListFoods(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()

	// test user & session
	api.db.createUser("testuser", "password123")
	s, _ := api.db.createSession("testuser")

	t.Run("no-foods", func(t *testing.T) {
		w, c := newContext(t, "")
		api.handleListFoods(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for no foods, got %d", http.StatusOK, w.Code)
		}

		var foods []FoodResponse
		err := json.Unmarshal(w.Body.Bytes(), &foods)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		if len(foods) != 0 {
			t.Errorf("Expected 0 foods, got %d", len(foods))
		}
	})

	// test foods
	for i := 1; i <= 47; i++ {
		api.db.createFoodByToken(CreateFoodParams{
			Name:         "Test Food " + strconv.Itoa(i),
			Brand:        "Test Brand",
			Calories:     100 * i,
			Carbs:        10 * i,
			Protein:      5 * i,
			Fat:          2 * i,
			ServingSize:  "g",
			ServingCount: 1,
		}, s.Token)
	}

	t.Run("with-foods", func(t *testing.T) {
		w, c := newContext(t, "")
		api.handleListFoods(c)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d for foods, got %d", http.StatusOK, w.Code)
		}

		var foods []FoodResponse
		err := json.Unmarshal(w.Body.Bytes(), &foods)
		if err != nil {
			t.Errorf("Failed to parse response body: %v", err)
		}

		if len(foods) != 47 {
			t.Errorf("Expected 47 foods, got %d", len(foods))
		}
	})
}
