package db2

import (
	"database/sql"
	"errors"
	"macro/errs"
)

// CreateFoodByToken creates a new food entry in the database.
func (db *Database) CreateFoodByToken(token string, in *FoodParams) (*Food, error) {
	food := &Food{
		Name:         in.Name,
		Brand:        in.Brand,
		Calories:     in.Calories,
		Carbs:        in.Carbs,
		Protein:      in.Protein,
		Fat:          in.Fat,
		ServingCount: in.ServingCount,
		ServingSize:  in.ServingSize,
	}

	p := args(in, token)
	if err := db.QueryRow(`
		INSERT INTO food
		  	(name, brand, calories, carbs, protein, fat, serving_count, serving_size, user_name)
	    VALUES
			(?, ?, ?, ?, ?, ?, ?, ?,
			(SELECT user_name FROM session WHERE token = ? AND expires > CURRENT_TIMESTAMP))
		RETURNING id, user_name, created;
	`, p...).Scan(&food.ID, &food.UserName, &food.Created); err != nil {
		if IsForeignKeyError(err) || errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotAuthorized(err, "No active session").With("args", args)
		}

		return nil, errs.Unexpected(err, "").With("args", args)
	}

	return food, nil
}
