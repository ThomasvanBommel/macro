package main

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// Database wraps sql.DB with app-specific query helpers.
type Database struct {
	*sql.DB
}

//go:embed migrations/*.sql
var migrationFiles embed.FS

// InitDatabase opens the primary DB connection and applies migrations.
func InitDatabase() *Database {
	db := openDatabase("/data/macro.db")

	// Apply migrations
	goose.SetBaseFS(migrationFiles)
	err := goose.SetDialect("sqlite3")
	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	// TODO: start session cleanup ticker

	return &Database{db}
}

// openDatabase opens a SQLite connection and sets connection limits.
func openDatabase(path string) *sql.DB {
	defer Trace("openDatabase(path)", "path", path)()

	db, err := sql.Open("sqlite", path)
	FatalOnError(err, "Failed to open database")

	db.SetMaxOpenConns(1)
	return db
}

// createUser hashes the password and inserts a new user.
func (db *Database) createUser(name string, password string) error {
	defer Trace("createUser(name, password)", "name", name)()

	// Password can be max 72 bytes
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		return err
	}

	q := "INSERT INTO user (name, password_hash) VALUES (?, ?);"
	_, err = db.Exec(q, name, hash)
	return err
}

// createSession creates a 1-hour session and returns its record.
func (db *Database) createSession(name string) (*Session, error) {
	defer Trace("createSession(userName)", "userName", name)()

	t := GenerateRandomHexString(32)

	q := `INSERT INTO session (user_name, token, expires)
		  VALUES (?, ?, datetime('now', '+1 hour'))
		  RETURNING user_name, token, created, expires;`

	var s Session
	err := db.QueryRow(q, name, t).Scan(&s.UserName, &s.Token, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// getSessionByToken returns the session for a token.
func (db *Database) getSessionByToken(token string) (*Session, error) {
	defer Trace("getSessionByToken(token)", "token", token)()

	var s Session
	q := `SELECT user_name, token, created, expires 
	      FROM session 
		  WHERE token = ?
		    AND expires > datetime('now');`
	err := db.QueryRow(q, token).Scan(&s.UserName, &s.Token, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// deleteSession deletes a session by token.
func (db *Database) deleteSession(token string) error {
	defer Trace("deleteSession(token)", "token", token)()

	q := `DELETE FROM session WHERE token = ?;`
	_, err := db.Exec(q, token)
	return err
}

// getUserByNameAndPassword authenticates a user via bcrypt hash comparison.
func (db *Database) getUserByNameAndPassword(name string, password string) (*User, error) {
	defer Trace("getUserByNameAndPassword(name, password)", "name", name)()

	var u User
	q := "SELECT name, password_hash, created FROM user WHERE name = ?;"
	err := db.QueryRow(q, name).Scan(&u.Name, &u.PasswordHash, &u.Created)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// rowScanner abstracts Scan for both sql.Row and sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

// collectRows scans all rows with the provided scanner function.
func collectRows[T any](rows *sql.Rows, scan func(rowScanner) (T, error)) ([]T, error) {
	defer rows.Close()

	items := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// scanEntryWithFood maps a joined entry+food row to EntryWithFood.
func scanEntryWithFood(scanner rowScanner) (EntryWithFood, error) {
	var e EntryWithFood
	err := scanner.Scan(
		&e.ID, &e.UserName, &e.MealName, &e.Date, &e.Servings, &e.Created,
		&e.Food.ID, &e.Food.Name, &e.Food.Brand, &e.Food.Created, &e.Food.UserName,
		&e.Food.Calories, &e.Food.Carbs, &e.Food.Protein, &e.Food.Fat,
		&e.Food.ServingSize, &e.Food.ServingCount,
	)
	return e, err
}

// listUserEntriesWithFoodByNameAndDate returns joined entries and foods for a user/date.
func (db *Database) listUserEntriesWithFoodByNameAndDate(name string, date string) (
	[]EntryWithFood, error) {
	defer Trace("listUserEntriesWithFoodByNameAndDate(name, date)", "name", name, "date", date)()

	q := `SELECT e.id, e.user_name, e.meal_name, e.date, e.servings, e.created,
				 f.id, f.name, f.brand, f.created, f.user_name, f.calories, f.carbs, f.protein, 
				 f.fat, f.serving_size, f.serving_count
		  FROM entry e
		  JOIN food f ON e.food_id = f.id
		  WHERE e.user_name = ? AND e.date = ?;`

	rows, err := db.Query(q, name, date)
	if err != nil {
		return nil, err
	}

	return collectRows(rows, scanEntryWithFood)
}

// createFoodByToken inserts a food tied to the user resolved from session token.
func (db *Database) createFoodByToken(food CreateFoodParams, token string) (*Food, error) {
	defer Trace("createFoodByToken(food, token)", "food", food, "token", token)()

	q := `INSERT INTO food (name, brand, user_name, calories, carbs, protein, fat, serving_size, 
		  	serving_count)
		  VALUES (?, ?, (
			SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?, ?, ?)
		  RETURNING id, name, brand, created, user_name, calories, carbs, protein, fat, 
		    serving_size, serving_count;`

	var f Food
	err := db.QueryRow(q, food.Name, food.Brand, token, food.Calories, food.Carbs,
		food.Protein, food.Fat, food.ServingSize, food.ServingCount).Scan(&f.ID, &f.Name,
		&f.Brand, &f.Created, &f.UserName, &f.Calories, &f.Carbs, &f.Protein,
		&f.Fat, &f.ServingSize, &f.ServingCount)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// scanFoods maps a food row to Food.
func scanFoods(scanner rowScanner) (Food, error) {
	defer Trace("scanFoods(scanner)")()

	var f Food
	err := scanner.Scan(&f.ID, &f.Name, &f.Brand, &f.Created, &f.UserName, &f.Calories, &f.Carbs,
		&f.Protein, &f.Fat, &f.ServingSize, &f.ServingCount)
	return f, err
}

// listFoods returns all foods.
func (db *Database) listFoods() ([]Food, error) {
	defer Trace("listFoods()")()

	q := `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food;`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}

	return collectRows(rows, scanFoods)
}

// getEntryWithFoodById returns one joined entry+food record by entry ID.
func (db *Database) getEntryWithFoodById(id int) (*EntryWithFood, error) {
	defer Trace("getEntryWithFoodById(id)", "id", id)()

	q := `SELECT e.id, e.user_name, e.meal_name, e.date, e.servings, e.created,
				 f.id, f.name, f.brand, f.created, f.user_name, f.calories, f.carbs, f.protein, 
				 f.fat, f.serving_size, f.serving_count
		  FROM entry e
		  JOIN food f ON e.food_id = f.id
		  WHERE e.id = ?;`

	var e EntryWithFood
	err := db.QueryRow(q, id).Scan(&e.ID, &e.UserName, &e.MealName, &e.Date, &e.Servings,
		&e.Created, &e.Food.ID, &e.Food.Name, &e.Food.Brand, &e.Food.Created,
		&e.Food.UserName, &e.Food.Calories, &e.Food.Carbs, &e.Food.Protein, &e.Food.Fat,
		&e.Food.ServingSize, &e.Food.ServingCount)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

// createEntryByToken inserts an entry for the session user and returns the hydrated record.
func (db *Database) createEntryByToken(entry CreateEntryParams, token string) (*EntryWithFood,
	error) {
	defer Trace("createEntryByToken(entry, token)", "entry", entry, "token", token)()

	q := `INSERT INTO entry (user_name, food_id, meal_name, date, servings)
		  VALUES ((
		  	SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?)
		  RETURNING id;`

	var id int
	err := db.QueryRow(q, token, entry.FoodId, entry.MealName, entry.Date, entry.Servings).Scan(&id)
	if err != nil {
		return nil, err
	}

	return db.getEntryWithFoodById(id)
}

// searchFoodsByName returns foods with names partially matching the query.
func (db *Database) searchFoodsByName(query string) ([]Food, error) {
	defer Trace("searchFoodsByName(query)", "query", query)()

	q := `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food
		 WHERE LOWER(name) LIKE LOWER(?)
		 ORDER BY name ASC;`

	rows, err := db.Query(q, "%"+query+"%")
	if err != nil {
		return nil, err
	}

	return collectRows(rows, scanFoods)
}

func (db *Database) searchFoodsByNameSortedUserFromTokenFirst(query string, token string) ([]Food, error) {
	defer Trace("searchFoodsByNameSortedUserFromTokenFirst(query, token)", "query", query, "token", token)()

	q := `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food
		 WHERE LOWER(name) LIKE LOWER(?)
		 ORDER BY (user_name = (SELECT user_name FROM session WHERE token = ?)) DESC, name ASC;`

	rows, err := db.Query(q, "%"+query+"%", token)
	if err != nil {
		return nil, err
	}

	return collectRows(rows, scanFoods)
}
