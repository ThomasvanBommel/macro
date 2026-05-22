package db2

import (
	"macro/errs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFood(t *testing.T) {
	db := newDatabase()
	defer db.Close()

	mustCreateUser(t, db, "bob", "supersecure")

	new := func() *NewFood {
		return &NewFood{
			Name:         "apple",
			Brand:        "apple",
			UserName:     "bob",
			Calories:     1337.5,
			Carbs:        69.7,
			Protein:      5.1,
			Fat:          5.44,
			ServingCount: 5,
			ServingSize:  "ml",
		}
	}

	expectError := func(in *NewFood, exp errs.ErrCode) {
		t.Helper()
		_, err := db.CreateFood(in)

		var e errs.Error
		require.ErrorAs(t, err, &e)
		assert.Equal(t, exp, e.Code())
	}

	// name too short
	f := new()
	f.Name = "a"
	expectError(f, errs.ErrBadInput)

	// name too long
	f = new()
	f.Name = "abcdefghijklmnopqrstuvwxyz01234"
	expectError(f, errs.ErrBadInput)

	tests := []struct {
		name string
		in   *NewFood
		modify func(*NewFood)
	}{
		{
			name: "name-too-short",
			in: new(),
			modify: func(in *Newfood) { in.Name = "a" },
		}
	}
}
