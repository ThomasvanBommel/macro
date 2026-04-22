package api

import (
	"encoding/json"
	"macro/db"
	"net/http"
	"strconv"
	"testing"
)

func mustCreateEntryByParam(t *testing.T, api *API, token string, params db.EntryParams) *db.EntryWithFood {
	t.Helper()
	e, err := api.db.CreateEntryByToken(params, token)
	if err != nil {
		t.Fatalf("CreateEntryByToken(%v) failed: %v", params, err)
	}
	return e
}

func mustCreateEntryByVar(t *testing.T, api *API, token string, foodID int, mealName string,
	date string, servings int) *db.EntryWithFood {
	t.Helper()
	return mustCreateEntryByParam(t, api, token, db.EntryParams{
		FoodId:   foodID,
		MealName: mealName,
		Date:     date,
		Servings: servings,
	})
}

func mustCreateEntry(t *testing.T, api *API, token string, foodID int, mealName string) *db.EntryWithFood {
	t.Helper()
	return mustCreateEntryByParam(t, api, token, db.EntryParams{
		FoodId:   foodID,
		MealName: mealName,
		Date:     "1901-01-01",
		Servings: 1,
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
		body := `{"food_id": 1, "meal_name": "lunch", "date": "1901-01-01", "servings": 0}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-session", func(t *testing.T) {
		body := `{"food_id": 1, "meal_name": "lunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, &db.Session{UserName: "testuser", Token: "invalidtoken"}, api.handleCreateEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-food-id", func(t *testing.T) {
		body := `{"food_id": 99999, "meal_name": "lunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}
	})

	// invalid meal name
	t.Run("invalid-meal-name", func(t *testing.T) {
		body := `{"food_id": ` + strconv.Itoa(f.ID) + `, "meal_name": "brunch", "date": "1901-01-01", "servings": 1}`
		w, _ := executeHandler(t, body, s, api.handleCreateEntry)
		if w.Code != http.StatusUnprocessableEntity {
			t.Errorf("Expected status %d, got %d", http.StatusUnprocessableEntity, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"food_id": ` + strconv.Itoa(f.ID) + `, "meal_name": "lunch", "date": "1901-01-01", "servings": "1"}`
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
			_, err := api.db.CreateEntryByToken(db.EntryParams{FoodId: f.ID, MealName: "breakfast", Date: "1901-01-01", Servings: 1}, s.Token)
			if err != nil {
				t.Fatalf("CreateEntryByToken failed: %v", err)
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
func TestHandleEditEntry(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")
	f1 := mustCreateFood(t, api, s.Token, "Food One")
	f2 := mustCreateFood(t, api, s.Token, "Food Two")
	e := mustCreateEntry(t, api, s.Token, f1.ID, "breakfast")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "lunch", "date": "1901-01-01", "servings": "1"}`
		w, _ := executeHandler(t, body, nil, api.handleEditEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-servings", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "lunch", "date": "1901-01-01", "servings": "0"}`
		w, _ := executeHandler(t, body, s, api.handleEditEntry)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("missing-entry", func(t *testing.T) {
		body := `{"id": 99999, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "lunch", "date": "1901-01-01", "servings": "1"}`
		w, _ := executeHandler(t, body, s, api.handleEditEntry)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `, "food_id": ` + strconv.Itoa(f2.ID) + `, "meal_name": "lunch", "date": "1901-01-02", "servings": "2"}`
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
		if entry.MealName != "lunch" {
			t.Errorf("Expected meal name %q, got %q", "lunch", entry.MealName)
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
	e := mustCreateEntry(t, api, s.Token, f.ID, "breakfast")

	t.Run("unauthorized", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `}`
		w, _ := executeHandler(t, body, nil, api.handleDeleteEntry)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status %d, got %d", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("invalid-entry", func(t *testing.T) {
		body := `{"id": 99999}`
		w, _ := executeHandler(t, body, s, api.handleDeleteEntry)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})

	t.Run("valid", func(t *testing.T) {
		body := `{"id": ` + strconv.Itoa(e.ID) + `}`
		w, _ := executeHandler(t, body, s, api.handleDeleteEntry)
		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status %d, got %d", http.StatusNoContent, w.Code)
		}

		entries, err := api.db.ListUserEntriesWithFoodByNameAndDate("testuser", "1901-01-01")
		if err != nil {
			t.Fatalf("ListUserEntriesWithFoodByNameAndDate failed: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("Expected deleted entry to be removed, got %d remaining entries", len(entries))
		}
	})
}
