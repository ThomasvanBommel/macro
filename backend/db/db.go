package db

import (
	"database/sql"
	"embed"
	"macro/util"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

// Database wraps sql.DB with app-specific query helpers.
type Database struct {
	*sql.DB
}

//go:embed migrations/*.sql
var MigrationFiles embed.FS

// NewDatabase opens the primary DB connection and applies migrations.
func NewDatabase() *Database {
	db := openDatabase("/data/macro.db")

	// Apply migrations
	goose.SetBaseFS(MigrationFiles)
	err := goose.SetDialect("sqlite3")
	util.FatalOnError(err, "Failed to set goose dialect")
	util.FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	// TODO: start session cleanup ticker

	return &Database{db}
}

// openDatabase opens a SQLite connection and sets connection limits.
func openDatabase(path string) *sql.DB {
	defer util.Trace("openDatabase(path)", "path", path)()

	db, err := sql.Open("sqlite", path)
	util.FatalOnError(err, "Failed to open database")

	// Set PRAGMAs for performance. WAL allows concurrent reads/writes, the others speed up writes.
	_, err = db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = OFF;
		PRAGMA cache_size = -10000;
	`)
	util.FatalOnError(err, "Failed to set PRAGMAs")

	db.SetMaxOpenConns(1)
	return db
}

// CreateUser hashes the password and inserts a new user.
func (db *Database) CreateUser(name string, password string) error {
	defer util.Trace("CreateUser(name, password)", "name", name)()

	// Password can be max 72 bytes
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		return err
	}

	q := "INSERT INTO user (name, password_hash) VALUES (?, ?);"
	_, err = db.Exec(q, name, hash)
	return err
}

// CreateSession creates a session for the user and returns the session record.
func (db *Database) CreateSession(name string, timeout_sec int) (*Session, error) {
	defer util.Trace("CreateSession(userName)", "userName", name)()

	t := util.GenerateRandomHexString(32)

	q := `INSERT INTO session (user_name, token, expires)
		  VALUES (?, ?, datetime('now', '+' || ? || ' seconds'))
		  RETURNING user_name, token, created, expires;`

	var s Session
	err := db.QueryRow(q, name, t, timeout_sec).Scan(&s.UserName, &s.Token, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// GetSessionByToken returns the session for a token.
func (db *Database) GetSessionByToken(token string) (*Session, error) {
	defer util.Trace("GetSessionByToken(token)", "token", token)()

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

// DeleteSession deletes a session by token.
func (db *Database) DeleteSession(token string) error {
	defer util.Trace("DeleteSession(token)", "token", token)()

	q := `DELETE FROM session WHERE token = ?;`
	_, err := db.Exec(q, token)
	return err
}

// GetUserByNameAndPassword authenticates a user via bcrypt hash comparison.
func (db *Database) GetUserByNameAndPassword(name string, password string) (*User, error) {
	defer util.Trace("GetUserByNameAndPassword(name, password)", "name", name)()

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

// ListUserEntriesWithFoodByNameAndDate returns joined entries and foods for a user/date.
func (db *Database) ListUserEntriesWithFoodByNameAndDate(name string, date string) (
	[]EntryWithFood, error) {
	defer util.Trace("ListUserEntriesWithFoodByNameAndDate(name, date)", "name", name, "date", date)()

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

// CreateFoodByToken inserts a food tied to the user resolved from session token.
func (db *Database) CreateFoodByToken(food FoodParams, token string) (*Food, error) {
	defer util.Trace("CreateFoodByToken(food, token)", "food", food, "token", token)()

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
	defer util.Trace("scanFoods(scanner)")()

	var f Food
	err := scanner.Scan(&f.ID, &f.Name, &f.Brand, &f.Created, &f.UserName, &f.Calories, &f.Carbs,
		&f.Protein, &f.Fat, &f.ServingSize, &f.ServingCount)
	return f, err
}

// ListFoods returns all foods.
func (db *Database) ListFoods() ([]Food, error) {
	defer util.Trace("ListFoods()")()

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
	defer util.Trace("getEntryWithFoodById(id)", "id", id)()

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

// CreateEntryByToken inserts an entry for the session user and returns the hydrated record.
func (db *Database) CreateEntryByToken(entry EntryParams, token string) (*EntryWithFood, error) {
	defer util.Trace("CreateEntryByToken(entry, token)", "entry", entry, "token", token)()

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

// SearchFoodsByName returns foods with names partially matching the query.
func (db *Database) SearchFoodsByName(query string) ([]Food, error) {
	defer util.Trace("SearchFoodsByName(query)", "query", query)()

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

// SearchFoodsByNameSortedUserFromTokenFirst returns foods matching the query, sorted with the
// session user's foods first.
func (db *Database) SearchFoodsByNameSortedUserFromTokenFirst(query string, token string) ([]Food, error) {
	defer util.Trace("SearchFoodsByNameSortedUserFromTokenFirst(query, token)", "query", query, "token", token)()

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

// EditEntryAuthByToken updates an entry if it belongs to the session user and returns the record.
func (db *Database) EditEntryAuthByToken(id int, p EntryParams, t string) (*EntryWithFood, error) {
	defer util.Trace("EditEntryAuthByToken(id, params, token)", "id", id, "params", p, "token", t)()

	q := `UPDATE entry
		  SET food_id = ?, meal_name = ?, date = ?, servings = ?
		  WHERE id = ? AND user_name = (SELECT user_name FROM session WHERE token = ?)
		  RETURNING id;`

	var entryId int
	err := db.QueryRow(q, p.FoodId, p.MealName, p.Date, p.Servings, id, t).Scan(&entryId)
	if err != nil {
		return nil, err
	}

	return db.getEntryWithFoodById(entryId)
}

// DeleteEntryAuthByToken deletes an entry if it belongs to the session user.
func (db *Database) DeleteEntryAuthByToken(id int, token string) error {
	defer util.Trace("DeleteEntryAuthByToken(id, token)", "id", id, "token", token)()

	q := `DELETE 
	      FROM entry 
		  WHERE id = ? AND user_name = (SELECT user_name FROM session WHERE token = ?);`
	_, err := db.Exec(q, id, token)
	return err
}
