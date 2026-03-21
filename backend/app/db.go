package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"

	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

type DatabaseWrapper struct {
	*sql.DB
}

func initDB() (*sql.DB, error) {
	// Open and configure database
	db, err := sql.Open("sqlite", "/app/data/macro.db")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	// Apply migrations
	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	log.Println("Running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	// ticker := time.NewTicker(10 * time.Minute)
	// go func() {
	// 	for range ticker.C {
	// 		log.Println("Cleaning up expired sessions...")
	// 		if err := cleanupExpiredSessions(db); err != nil {
	// 			log.Printf("Failed to clean up expired sessions: %v", err)
	// 		}
	// 	}
	// }()

	return db, nil
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

	// Insert into database and fetch session info
	var res ResponseSessionModel
	q := `INSERT INTO session (user_name, token, expires)
		  VALUES (?, ?, datetime('now', '+1 hour'))
		  RETURNING user_name, created, expires;`
	err := db.QueryRow(q, user.Name, token).Scan(&res.UserName, &res.Created, &res.Expires)
	if err != nil {
		return "", nil, err
	}

	return token, &res, nil
}

func (db *DatabaseWrapper) getSessionInfo(token string) (*ResponseSessionModel, error) {
	var s ResponseSessionModel
	q := `SELECT user_name, created, expires
		  FROM session
		  WHERE token = ?;`
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

// func (db *DatabaseWrapper) selectSession(token string) (*ResponseSessionModel, error) {
// 	var s ResponseSessionModel
// 	q := `SELECT user_name, created, expires
// 		  FROM session
// 		  WHERE token = ?
// 		    AND expires > datetime('now');`
// 	err := db.QueryRow(q, token).Scan(&s.UserName, &s.Created, &s.Expires)
// 	if err != nil { return nil, err }

// 	return &s, nil
// }

// func createUser(db *sql.DB, username string, password string) error {
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec(
// 		"INSERT INTO users (username, password_hash) VALUES (?, ?);", username, string(hash))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func createSession(db *sql.DB, user_id int) (string, error) {
// 	b := make([]byte, 32)
// 	if _, err := rand.Read(b); err != nil {
// 		return "", err
// 	}
// 	token := hex.EncodeToString(b)

// 	q := `
// 		INSERT INTO sessions (user_id, token, expires_at)
// 		VALUES (?, ?, datetime('now', '+1 hour'));
// 	`

// 	_, err := db.Exec(q, user_id, token)
// 	if err != nil {
// 		return "", err
// 	}

// 	return token, nil
// }

// func cleanupExpiredSessions(db *sql.DB) error {
// 	q := `
// 		DELETE
// 		FROM sessions
// 		WHERE expires_at <= datetime('now');
// 	`
// 	_, err := db.Exec(q)
// 	return err
// }

// func loginUser(db *sql.DB, username string, password string) (string, error) {
// 	var user_id int
// 	var hash string
// 	err := db.QueryRow(
// 		"SELECT id, password_hash FROM users WHERE username = ?", username).Scan(&user_id, &hash)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", nil
// 		}
// 		return "", err
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	if err != nil {
// 		return "", nil
// 	}

// 	return createSession(db, user_id)
// }

// func selectSession(db *sql.DB, token string) (*Session, error) {
// 	var s Session
// 	q := `
// 		SELECT user_id, username, expires_at
// 		FROM users
// 		JOIN sessions
// 		  ON users.id = user_id
// 		WHERE token = ?
// 		  AND expires_at > datetime('now');
// 	`
// 	err := db.QueryRow(q, token).Scan(&s.UserID, &s.Username, &s.ExpiresAt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &s, nil
// }

// func deleteSession(db *sql.DB, token string) error {
// 	q := `
// 		DELETE
// 		FROM sessions
// 		WHERE token = ?;
// 	`
// 	_, err := db.Exec(q, token)
// 	return err
// }

// func insertFood(db *sql.DB, req PutFoodRequest) error {
// 	q := `
// 		INSERT INTO food (name, brand, created_by, calories, carbs, protein, fat, serving_size)
// 		VALUES (?, ?, (
// 			SELECT user_id FROM sessions WHERE token = ?
// 		), ?, ?, ?, ?, ?);
// 	`
// 	_, err := db.Exec(q, req.Name, req.Brand, req.SessionToken, req.Calories, req.Carbs,
// 		req.Protein, req.Fat, req.ServingSize)

// 	return err
// }

// func selectFoods(db *sql.DB) ([]Food, error) {
// 	q := `
// 		SELECT id, name, brand, created_at, created_by, calories, carbs, protein, fat, serving_size
// 		FROM food;
// 	`
// 	rows, err := db.Query(q)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var foods []Food
// 	for rows.Next() {
// 		var f Food
// 		if err := rows.Scan(&f.ID, &f.Name, &f.Brand, &f.CreatedAt, &f.CreatedBy, &f.Calories,
// 			&f.Carbs, &f.Protein, &f.Fat, &f.ServingSize); err != nil {
// 			return nil, err
// 		}
// 		foods = append(foods, f)
// 	}
// 	return foods, nil
// }

// func selectMeals(db *sql.DB) ([]Meal, error) {
// 	q := "SELECT name FROM meals;"
// 	rows, err := db.Query(q)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var meals []Meal
// 	for rows.Next() {
// 		var m Meal
// 		if err := rows.Scan(&m.Name); err != nil {
// 			return nil, err
// 		}
// 		meals = append(meals, m)
// 	}
// 	return meals, nil
// }

// func insertEntry(db *sql.DB, req PutEntryRequest) error {
// 	q := `
// 		INSERT INTO entries (user_id, food_id, meal_id, servings, date)
// 		VALUES (
// 			(SELECT user_id FROM sessions WHERE token = ?),
// 			?, ?, ?, ?
// 		);
// 	`
// 	_, err := db.Exec(q, req.SessionToken, req.FoodID, req.Meal, req.Servings, req.Date)
// 	return err
// }

// func selectEntries(db *sql.DB, user_id int, date string) ([]EntryResponse, error) {
// 	q := `
// 		SELECT e.id, f.id, f.name, f.brand, f.calories, f.carbs, f.protein, f.fat, f.serving_size,
// 		       e.meal_id, e.servings, e.date, e.created_at
// 		FROM entries e
// 		JOIN food f ON e.food_id = f.id
// 		WHERE e.user_id = ?
// 		  AND e.date = ?;
// 	`
// 	rows, err := db.Query(q, user_id, date)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var entries []EntryResponse
// 	for rows.Next() {
// 		var er EntryResponse
// 		if err := rows.Scan(&er.ID, &er.Food.ID, &er.Food.Name, &er.Food.Brand,
// 			&er.Food.Calories, &er.Food.Carbs, &er.Food.Protein, &er.Food.Fat, &er.Food.ServingSize,
// 			&er.Meal.Name, &er.Servings, &er.Date, &er.CreatedAt); err != nil {
// 			return nil, err
// 		}
// 		entries = append(entries, er)
// 	}
// 	return entries, nil
// }

// func selectEntries(db *sql.DB, user_id int, date string) ([]Entry, error) {
// 	q := `
// 		SELECT id, food_id, meal_id, servings, date
// 		FROM entries
// 		WHERE user_id = ?
// 		  AND date = ?;
// 	`
// 	rows, err := db.Query(q, user_id, date)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var entries []Entry
// 	for rows.Next() {
// 		var e Entry
// 		if err := rows.Scan(&e.ID, &e.FoodID, &e.Meal.Name, &e.Servings, &e.Date); err != nil {
// 			return nil, err
// 		}
// 		entries = append(entries, e)
// 	}
// 	return entries, nil
// }
