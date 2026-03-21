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

	var foods []ResponseFoodModel
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
