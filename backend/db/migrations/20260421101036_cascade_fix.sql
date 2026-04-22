-- Migration to add ON DELETE CASCADE to foreign keys in entry and food tables.

-- +goose Up
CREATE TABLE entry_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT NOT NULL,
    food_id INTEGER NOT NULL,
    meal_name TEXT NOT NULL,
    date DATE NOT NULL,
    servings INTEGER NOT NULL,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_name) REFERENCES user(name) ON DELETE CASCADE,
    FOREIGN KEY (food_id) REFERENCES food(id) ON DELETE CASCADE,
    FOREIGN KEY (meal_name) REFERENCES meal(name) ON DELETE CASCADE
);

INSERT INTO entry_new (id, user_name, food_id, meal_name, date, servings, created)
SELECT id, user_name, food_id, meal_name, date, servings, created
FROM entry;

DROP TABLE entry;
ALTER TABLE entry_new RENAME TO entry;

CREATE TABLE food_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brand TEXT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_name TEXT NOT NULL,
    calories INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    serving_size TEXT NOT NULL,
    serving_count INTEGER NOT NULL,
    FOREIGN KEY (user_name) REFERENCES user(name) ON DELETE CASCADE
);

INSERT INTO food_new (id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, serving_count)
SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, serving_count
FROM food;

DROP TABLE food;
ALTER TABLE food_new RENAME TO food;

-- +goose Down
CREATE TABLE entry_old (
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

INSERT INTO entry_old (id, user_name, food_id, meal_name, date, servings, created)
SELECT id, user_name, food_id, meal_name, date, servings, created
FROM entry;

DROP TABLE entry;
ALTER TABLE entry_old RENAME TO entry;

CREATE TABLE food_old (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    brand TEXT,
    created DATETIME DEFAULT CURRENT_TIMESTAMP,
    user_name TEXT NOT NULL, -- user who created the food entry
    calories INTEGER NOT NULL,
    carbs INTEGER NOT NULL,
    protein INTEGER NOT NULL,
    fat INTEGER NOT NULL,
    serving_size TEXT NOT NULL,
    serving_count INTEGER NOT NULL,
    FOREIGN KEY (user_name) REFERENCES user(name)
);

INSERT INTO food_old (id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, serving_count)
SELECT id, name, brand, created, user_name, calories, carbs, protein, fat, serving_size, serving_count
FROM food;

DROP TABLE food;
ALTER TABLE food_old RENAME TO food;