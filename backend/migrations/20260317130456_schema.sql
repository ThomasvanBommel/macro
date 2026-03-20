-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE food (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brand TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by INTEGER NOT NULL,
    calories INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    serving_size REAL NOT NULL,
    FOREIGN KEY (created_by) REFERENCES users(id)
);

CREATE TABLE meals (
    name TEXT PRIMARY KEY
);

INSERT INTO meals (name) VALUES ('Breakfast'), ('Lunch'), ('Dinner'), ('Snack');

CREATE TABLE entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    food_id INTEGER NOT NULL,
    meal_id INTEGER NOT NULL,
    date DATE NOT NULL,
    servings REAL NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (food_id) REFERENCES food(id),
    FOREIGN KEY (meal_id) REFERENCES meals(id)
);

-- +goose Down
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS meals;
DROP TABLE IF EXISTS food;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;