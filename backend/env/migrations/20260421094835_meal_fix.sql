-- Migration to fix meal names because I forgot to enable foreign keys initially.

-- +goose Up
INSERT INTO meal (name)
VALUES ('breakfast'), ('lunch'), ('dinner'), ('snack')
ON CONFLICT (name) DO NOTHING;

UPDATE entry
SET meal_name = LOWER(meal_name);

UPDATE entry
SET meal_name = 'breakfast'
WHERE meal_name NOT IN ('breakfast', 'lunch', 'dinner', 'snack');

DELETE FROM meal
WHERE name NOT IN ('breakfast', 'lunch', 'dinner', 'snack');

-- +goose Down
INSERT INTO meal (name)
VALUES ('Breakfast'), ('Lunch'), ('Dinner'), ('Snack')
ON CONFLICT (name) DO NOTHING;

UPDATE entry
SET meal_name = UPPER(SUBSTR(meal_name, 1, 1)) || SUBSTR(meal_name, 2);

DELETE FROM meal
WHERE name NOT IN ('Breakfast', 'Lunch', 'Dinner', 'Snack');
