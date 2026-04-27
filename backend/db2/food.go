package db2

import (
	"database/sql"
	"errors"
	"macro/errs"
)

type FoodIn struct {
	Name         string
	Brand        string
	ServingSize  string
	Calories     int
	Carbs        int
	Protein      int
	Fat          int
	ServingCount int
}

type FoodOut struct {
	ID       int
	UserName string
	Created  string
}

func (f *FoodIn) Args() []any {
	return []any{
		f.Name,
		f.Brand,
		f.Calories,
		f.Carbs,
		f.Protein,
		f.Fat,
		f.ServingCount,
		f.ServingSize,
	}
}

// CreateFoodByToken creates a new food entry in the database.
func (db *Database) CreateFoodByToken(token string, in *FoodIn) (out *FoodOut, err error) {
	q := `INSERT INTO food (user_name, name, brand, calories, carbs, protein, fat, 
			serving_count, serving_size)
	      VALUES ((
		  	SELECT user_name FROM session WHERE token = ? AND expires > CURRENT_TIMESTAMP
		  ), ?, ?, ?, ?, ?, ?, ?, ?)
		  RETURNING id, user_name, created`

	var id int
	var username, created string
	args := append([]any{token}, in.Args()...)

	if err := db.QueryRow(q, args...).Scan(&id, &username, &created); err != nil {
		if IsForeignKeyError(err) || errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotAuthorized(err, "No active session").With("args", args)
		}

		return nil, errs.Unexpected(err, "").With("args", args)
	}

	return &FoodOut{id, username, created}, nil
}
