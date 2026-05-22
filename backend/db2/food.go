package db2

import (
	"macro/errs"

	"github.com/gin-gonic/gin/binding"
)

// CreateFood creates a new food entry in the database.
func (db *Database) CreateFood(in *NewFood) (*Food, error) {
	err := binding.Validator.ValidateStruct(in)
	if err != nil {
		return nil, errs.BadInput(err, err.Error()).With("input", in)
	}

	f := &Food{
		Name:         in.Name,
		Brand:        in.Brand,
		UserName:     in.UserName,
		Calories:     in.Calories,
		Carbs:        in.Carbs,
		Protein:      in.Protein,
		Fat:          in.Fat,
		ServingCount: in.ServingCount,
		ServingSize:  in.ServingSize,
	}

	q := `	INSERT INTO food
			(name,brand,user_name,calories,carbs,protein,fat,serving_count,serving_size)
			VALUES (?,?,?,?,?,?,?,?,?) RETURNING id, created;`

	args := ToSlice(in)
	if err := db.QueryRow(q, args...).Scan(&f.ID, &f.Created); err != nil {
		return nil, errs.Unexpected(err, "").With("args", args)
	}

	return f, nil
}

func (db *Database) SearchFoodsByName(in string) ([]*Food, error) {
	rows, err := db.Query(`
		SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_count, 
			serving_size
		FROM food
		WHERE LOWER(name) LIKE LOWER(?)
		ORDER BY name ASC;
	`, "%"+in+"%")
	if err != nil {
		return nil, errs.Unexpected(err, "")
	}
	defer rows.Close()

	var foods []*Food
	for rows.Next() {
		f := &Food{}
		var calories, carbs, protein, fat, servingCount int
		if err := rows.Scan(
			&f.ID, &f.Name, &f.Brand, &f.Created, &f.UserName, &calories, &carbs, &protein, &fat,
			&servingCount, &f.ServingSize,
		); err != nil {
			return nil, errs.Unexpected(err, "")
		}

		f.Calories = ToDecimal(calories)
		f.Carbs = ToDecimal(carbs)
		f.Protein = ToDecimal(protein)
		f.Fat = ToDecimal(fat)
		f.ServingCount = ToDecimal(servingCount)

		foods = append(foods, f)
	}

	if err := rows.Err(); err != nil {
		return nil, errs.Unexpected(err, "")
	}

	return foods, nil
}
