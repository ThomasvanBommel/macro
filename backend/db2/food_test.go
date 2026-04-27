package db2

import (
	"macro/errs"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFoodByToken(t *testing.T) {
	db := newDatabase()
	defer db.Close()

	mustCreateUser(t, db, "valid_user", "valid_password")
	token, _, err := db.CreateSession("valid_user", "valid_password")
	require.NoError(t, err, "failed to create session: %v", err)

	tests := []struct {
		name     string
		token    string
		foodName string
		wantErr  bool
		httpCode int
	}{
		{"valid input", token, "Pizza", false, 0},
		{"invalid token", "invalid_token", "Burger", true, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		in := &FoodIn{
			Name:         tt.foodName,
			Brand:        "the bestest ever brand",
			Calories:     666,
			Carbs:        42,
			Protein:      42,
			Fat:          42,
			ServingCount: 1,
			ServingSize:  "dollop",
		}
		out, err := db.CreateFoodByToken(tt.token, in)

		if tt.wantErr {
			assert.Error(t, err, "expected error")

			var e *errs.Error
			require.ErrorAs(t, err, &e, "Expected errs.Error")
			assert.Equal(t, tt.httpCode, e.Code(), "want %d, got %d\n%v", tt.httpCode, e.Code(), err)
		} else {
			assert.NoError(t, err, "unexpected error")

			assert.Greater(t, out.ID, -1, "invalid food ID: %d", out.ID)
			assert.Equal(t, "valid_user", out.UserName, "unexpected username: %q", out.UserName)
		}
	}
}
