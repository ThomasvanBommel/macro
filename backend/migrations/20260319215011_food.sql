-- +goose Up
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
    FOREIGN KEY (created_by) REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS food;
