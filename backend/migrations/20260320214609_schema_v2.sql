-- +goose Up
CREATE TABLE user (
    name TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE session (
    token TEXT PRIMARY KEY,
    user_name TEXT NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires DATETIME NOT NULL,
    FOREIGN KEY (user_name) REFERENCES user(name) ON DELETE CASCADE
);

/* Food:
   calories, carbs, protein, fat, and serving_count are stored as integer representing the value 
   multiplied by 100 to avoid floating point issues. */
CREATE TABLE food (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brand TEXT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_by_user_name TEXT NOT NULL,
    calories INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    serving_size TEXT NOT NULL,
    serving_count INTEGER NOT NULL,
    FOREIGN KEY (created_by_user_name) REFERENCES user(name)
);

CREATE TABLE meal (name TEXT PRIMARY KEY);
INSERT INTO meal (name) VALUES ('Breakfast'), ('Lunch'), ('Dinner'), ('Snack');

/* Entry:
   servings is stored as integer representing the value multiplied by 100 to avoid floating point 
   issues. */
CREATE TABLE entry (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    food_id INTEGER NOT NULL,
    meal_name TEXT NOT NULL,
    date DATE NOT NULL,
    servings INTEGER NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_name) REFERENCES user(name),
    FOREIGN KEY (food_id) REFERENCES food(id),
    FOREIGN KEY (meal_name) REFERENCES meal(name)
);

-- +goose Down
DROP TABLE IF EXISTS session;
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS food;
DROP TABLE IF EXISTS meal;
DROP TABLE IF EXISTS entry;