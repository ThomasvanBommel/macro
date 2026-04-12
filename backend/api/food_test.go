package api

import (
	"encoding/json"
	"macro/db"
	"net/http"
	"strconv"
	"testing"
)

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
		w, _ := executeHandler(t, body, &db.Session{UserName: "testuser", Token: "invalidtoken"}, api.handleCreateFood)
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
		_, err := api.db.CreateFoodByToken(db.FoodParams{
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
