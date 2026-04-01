package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/hex"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Database is a wrapper around sql.DB that provides additional methods for our application logic.
// It embeds *sql.DB, so it inherits all of its methods, but also allows us to define our own
// methods that are specific to our application's needs.
type Database struct {
	*sql.DB
}

// InitDatabase initializes the database by opening a connection to the SQLite database file,
// applying any necessary migrations, and returning a Database instance. It panics if any step
// fails, ensuring that the application does not run with an uninitialized database.
func InitDatabase() *Database {
	defer Trace("InitDatabase()")()

	db := openDatabase("/app/data/macro.db")

	applyMigrations(db)

	// TODO: start session cleanup ticker

	return &Database{db}
}

// openDatabase opens a SQLite database at the specified path, applies migrations, and returns a
// Database instance. It panics if any step fails, ensuring that the application does not run with
// an uninitialized database.
func openDatabase(path string) *sql.DB {
	defer Trace("openDatabase(path)", "path", path)()

	db, err := sql.Open("sqlite", path)
	FatalOnError(err, "Failed to open database")

	db.SetMaxOpenConns(1)
	return db
}

// applyMigrations applies database migrations using the goose library. It sets the base filesystem
// for migrations to the embedded migrations and configures the dialect to SQLite3. If any step
// fails, it panics, ensuring that the application does not run with an uninitialized database.
func applyMigrations(db *sql.DB) {
	defer Trace("applyMigrations(db)")()

	goose.SetBaseFS(migrationFiles)
	err := goose.SetDialect("sqlite3")

	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")
}

// createUser creates a new user in the database with the given name and password. It hashes the
// password using bcrypt before storing it. If any step fails, it returns an error.
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

// createSession creates a new session for the user with the given name. It generates a random
// token, inserts it into the session table with an expiration time of 1 hour, and returns the
// created Session struct or an error if the operation fails.
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

// deleteSession deletes a session from the database based on the provided token. It executes a
// DELETE statement and returns any error encountered.
func (db *Database) deleteSession(token string) error {
	defer Trace("deleteSession(token)", "token", token)()

	q := `DELETE FROM session WHERE token = ?;`
	_, err := db.Exec(q, token)
	return err
}

// getUserByNameAndPassword retrieves a user from the database by their name and password. It first
// queries the database for a user with the specified name, then compares the provided password with
// the stored password hash using bcrypt. If the user is found and the password matches, it returns
// a pointer to the User struct; otherwise, it returns an error.
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

// rowScanner is an interface that abstracts the Scan method used by both sql.Row and sql.Rows.
// This allows the same scanning logic to be used for both single-row and multi-row queries.
type rowScanner interface {
	Scan(dest ...any) error
}

// collectRows is a helper function that iterates over sql.Rows and applies a scanning function to
// each row. It takes a pointer to sql.Rows and a scanning function as arguments, and returns a
// slice of the scanned items or an error if any occurs during iteration or scanning. The function
// ensures that the rows are properly closed after processing.
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

// scanEntryWithFood is a helper function that scans a single row from a query that joins the entry
// and food tables, and returns an EntryWithFood struct. It takes a rowScanner (which can be either
// sql.Row or sql.Rows) as an argument and scans the fields into the appropriate struct fields. If
// the scanning is successful, it returns the EntryWithFood struct and a nil error; otherwise, it
// returns an empty EntryWithFood struct and the encountered error.
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

// listUserEntriesWithFoodByNameAndDate retrieves a list of entries for a specific user and date,
// along with the associated food information. It executes a SQL query that joins the entry and food
// tables, filtering by the user's name and the specified date. The results are collected into a
// slice of EntryWithFood structs using the collectRows helper function. If any error occurs during
// the query execution or row scanning, it returns an error.
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

//------------------------------ OLD:

type DatabaseWrapper struct {
	*sql.DB
}

// func initDB() (*sql.DB, error) {
// 	// Open and configure database
// 	db, err := sql.Open("sqlite", "/app/data/macro.db")
// 	if err != nil {
// 		return nil, err
// 	}
// 	db.SetMaxOpenConns(1)

// 	// Apply migrations
// 	goose.SetBaseFS(migrationsFS)
// 	if err := goose.SetDialect("sqlite3"); err != nil {
// 		return nil, err
// 	}

// 	log.Println("Running database migrations...")
// 	if err := goose.Up(db, "migrations"); err != nil {
// 		return nil, err
// 	}

// 	// Start session cleanup ticker
// 	ticker := time.NewTicker(24 * time.Hour)
// 	go func() {
// 		for range ticker.C {
// 			log.Println("Cleaning up expired sessions...")
// 			if err := cleanupExpiredSessions(db); err != nil {
// 				log.Printf("Failed to clean up expired sessions: %v", err)
// 			}
// 		}
// 	}()

// 	return db, nil
// }

func cleanupExpiredSessions(db *sql.DB) error {
	q := `DELETE FROM session WHERE expires <= datetime('now');`
	_, err := db.Exec(q)
	return err
}

func (db *DatabaseWrapper) register(user RequestUserModel) (string, *ResponseSessionModel, error) {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil, err
	}

	// Insert into database
	q := "INSERT INTO user (name, password_hash) VALUES (?, ?);"
	_, err = db.Exec(q, user.Name, string(hash))
	if err != nil {
		return "", nil, err
	}

	// Login
	return db.login(user)
}

func (db *DatabaseWrapper) login(user RequestUserModel) (string, *ResponseSessionModel, error) {
	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	token := hex.EncodeToString(b)

	// Check password
	var hash string
	q := "SELECT password_hash FROM user WHERE name = ?;"
	err := db.QueryRow(q, user.Name).Scan(&hash)
	if err != nil {
		return "", nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password)); err != nil {
		return "", nil, err
	}

	// Insert into database and fetch session info
	var res ResponseSessionModel
	q = `INSERT INTO session (user_name, token, expires)
		  VALUES (?, ?, datetime('now', '+1 hour'))
		  RETURNING user_name, created, expires;`
	err = db.QueryRow(q, user.Name, token).Scan(&res.UserName, &res.Created, &res.Expires)
	if err != nil {
		return "", nil, err
	}

	return token, &res, nil
}

func (db *DatabaseWrapper) getSessionInfo(token string) (*ResponseSessionModel, error) {
	var s ResponseSessionModel
	q := `SELECT user_name, created, expires
		  FROM session
		  WHERE token = ?
		    AND expires > datetime('now');`
	err := db.QueryRow(q, token).Scan(&s.UserName, &s.Created, &s.Expires)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (db *DatabaseWrapper) deleteSession(token string) error {
	q := `DELETE FROM session WHERE token = ?;`
	_, err := db.Exec(q, token)
	return err
}

func (db *DatabaseWrapper) getEntries(req RequestEntriesModel) ([]ResponseEntryModel, error) {
	q := `SELECT e.id, e.user_name, e.meal_name, e.date, e.servings, e.created,
				 f.id, f.name, f.brand, f.created, f.user_name, f.calories, f.carbs, f.protein, 
				 f.fat, f.serving_size, f.serving_count
		  FROM entry e
		  JOIN food f ON e.food_id = f.id
		  WHERE e.user_name = ? AND e.date = ?;`

	rows, err := db.Query(q, req.UserName, req.Date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// var entries []ResponseEntryModel
	entries := make([]ResponseEntryModel, 0)
	for rows.Next() {
		var e ResponseEntryModel
		var f ResponseFoodModel
		err := rows.Scan(&e.ID, &e.UserName, &e.MealName, &e.Date, &e.Servings, &e.Created,
			&f.ID, &f.Name, &f.Brand, &f.Created, &f.UserName, &f.Calories, &f.Carbs,
			&f.Protein, &f.Fat, &f.ServingSize, &f.ServingCount)
		if err != nil {
			return nil, err
		}
		e.Food = f
		entries = append(entries, e)
	}

	return entries, nil
}

func (db *DatabaseWrapper) addFood(req *RequestFoodModel, token string) (int, error) {
	q := `INSERT INTO food (name, brand, user_name, calories, carbs, protein, fat, serving_size, 
		  	serving_count)
		  VALUES (?, ?, (
			SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?, ?, ?)
		  RETURNING id;`

	var id int
	err := db.QueryRow(q, req.Name, req.Brand, token, req.Calories, req.Carbs, req.Protein,
		req.Fat, req.ServingSize, req.ServingCount).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// insert entry from RequestAddEntryModel, returning ResponseEntryModel with food info populated
func (db *DatabaseWrapper) addEntry(req RequestAddEntryModel, token string) (*ResponseEntryModel,
	error) {
	q := `INSERT INTO entry (user_name, food_id, meal_name, date, servings)
		  VALUES ((
		  	SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?)
		  RETURNING id, user_name, meal_name, date, servings, created;`

	var e ResponseEntryModel
	err := db.QueryRow(q, token, req.FoodId, req.MealName, req.Date, req.Servings).
		Scan(&e.ID, &e.UserName, &e.MealName, &e.Date, &e.Servings, &e.Created)
	if err != nil {
		return nil, err
	}

	// Get food info
	q = `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food
		 WHERE id = ?;`
	err = db.QueryRow(q, req.FoodId).Scan(&e.Food.ID, &e.Food.Name, &e.Food.Brand, &e.Food.Created,
		&e.Food.UserName, &e.Food.Calories, &e.Food.Carbs, &e.Food.Protein, &e.Food.Fat,
		&e.Food.ServingSize, &e.Food.ServingCount)
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (db *DatabaseWrapper) getFoods() ([]ResponseFoodModel, error) {
	q := `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food;`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// entries := make([]ResponseEntryModel, 0)
	// var foods []ResponseFoodModel
	foods := make([]ResponseFoodModel, 0)
	for rows.Next() {
		var f ResponseFoodModel
		err := rows.Scan(&f.ID, &f.Name, &f.Brand, &f.Created, &f.UserName, &f.Calories,
			&f.Carbs, &f.Protein, &f.Fat, &f.ServingSize, &f.ServingCount)
		if err != nil {
			return nil, err
		}
		foods = append(foods, f)
	}

	return foods, nil
}

func (db *DatabaseWrapper) createFood(req *RequestCreateFoodModel, token string) (ResponseFoodModel, error) {
	q := `INSERT INTO food (name, brand, user_name, calories, carbs, protein, fat, serving_size, 
		  	serving_count)
		  VALUES (?, ?, (
			SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?, ?, ?)
		  RETURNING id, name, brand, created, user_name, calories, carbs, protein, fat, 
		    serving_size, serving_count;`

	var f ResponseFoodModel
	err := db.QueryRow(q, req.Name, req.Brand, token, req.Calories, req.Carbs, req.Protein,
		req.Fat, req.ServingSize, req.ServingCount).Scan(&f.ID, &f.Name, &f.Brand, &f.Created,
		&f.UserName, &f.Calories, &f.Carbs, &f.Protein, &f.Fat, &f.ServingSize,
		&f.ServingCount)
	if err != nil {
		return ResponseFoodModel{}, err
	}

	return f, nil
}

func (db *DatabaseWrapper) createEntry(req *RequestAddEntryModel, token string) (*ResponseEntryModel, error) {
	q := `INSERT INTO entry (user_name, food_id, meal_name, date, servings)
		  VALUES ((
		  	SELECT user_name FROM session WHERE token = ?
		  ), ?, ?, ?, ?)
		  RETURNING id, user_name, meal_name, date, servings, created;`

	var e ResponseEntryModel
	err := db.QueryRow(q, token, req.FoodId, req.MealName, req.Date, req.Servings).
		Scan(&e.ID, &e.UserName, &e.MealName, &e.Date, &e.Servings, &e.Created)
	if err != nil {
		return nil, err
	}

	// Get food info
	q = `SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, 
	     	serving_count
		 FROM food
		 WHERE id = ?;`
	err = db.QueryRow(q, req.FoodId).Scan(&e.Food.ID, &e.Food.Name, &e.Food.Brand, &e.Food.Created,
		&e.Food.UserName, &e.Food.Calories, &e.Food.Carbs, &e.Food.Protein, &e.Food.Fat,
		&e.Food.ServingSize, &e.Food.ServingCount)
	if err != nil {
		return nil, err
	}

	return &e, nil
}
