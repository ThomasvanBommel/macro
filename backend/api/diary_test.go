package api

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestHandleGetDiary(t *testing.T) {
	api := newAPI(t)
	defer api.db.Close()
	mustCreateUser(t, api, "testuser")
	s := mustCreateSession(t, api, "testuser")

	f1 := mustCreateFoodN(t, api, s.Token, 1)
	f10 := mustCreateFoodN(t, api, s.Token, 10)
	f100 := mustCreateFoodN(t, api, s.Token, 100)
	f1000 := mustCreateFoodN(t, api, s.Token, 1000)

	mustCreateEntryByVar(t, api, s.Token, f1.ID, "breakfast", "1901-01-01", 1)
	mustCreateEntryByVar(t, api, s.Token, f10.ID, "lunch", "1901-01-01", 10)
	mustCreateEntryByVar(t, api, s.Token, f100.ID, "dinner", "1901-01-01", 100)
	mustCreateEntryByVar(t, api, s.Token, f1000.ID, "snack", "1901-01-01", 1000)

	t.Run("missing-date", func(t *testing.T) {
		body := `{"name": "testuser"}`
		w, _ := executeHandler(t, body, s, api.handleGetDiary)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})

	t.Run("invalid-user", func(t *testing.T) {
		body := `{"name": "invaliduser", "date": "1901-01-01"}`
		w, _ := executeHandler(t, body, s, api.handleGetDiary)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	t.Run("valid-request", func(t *testing.T) {
		body := `{"name": "testuser", "date": "1901-01-01"}`
		w, _ := executeHandler(t, body, s, api.handleGetDiary)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var res DiaryResponse
		if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
			t.Errorf("Failed to unmarshal response: %v", err)
		}

		expectedTotals := Totals{
			Calories: 11.11,
			Carbs:    11.11,
			Protein:  11.11,
			Fat:      11.11,
		}

		if res.Totals != expectedTotals {
			t.Errorf("Expected totals %+v, got %+v", expectedTotals, res.Totals)
		}

		if len(res.Breakfast.Entries) != 1 || res.Breakfast.Entries[0].Food.ID != f1.ID {
			t.Errorf("Expected breakfast entry with food ID %d, got %+v", f1.ID, res.Breakfast.Entries)
		}

		if len(res.Lunch.Entries) != 1 || res.Lunch.Entries[0].Food.ID != f10.ID {
			t.Errorf("Expected lunch entry with food ID %d, got %+v", f10.ID, res.Lunch.Entries)
		}

		if len(res.Dinner.Entries) != 1 || res.Dinner.Entries[0].Food.ID != f100.ID {
			t.Errorf("Expected dinner entry with food ID %d, got %+v", f100.ID, res.Dinner.Entries)
		}

		if len(res.Snack.Entries) != 1 || res.Snack.Entries[0].Food.ID != f1000.ID {
			t.Errorf("Expected snack entry with food ID %d, got %+v", f1000.ID, res.Snack.Entries)
		}

	})
}
