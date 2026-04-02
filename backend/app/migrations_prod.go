//go:build !dev

package main

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// applyMigrations applies database migrations using the goose library. It sets the base filesystem
// for migrations to the embedded migrations and configures the dialect to SQLite3. If any step
// fails, it panics, ensuring that the application does not run with an uninitialized database.
func ApplyMigrations(db *sql.DB) {
	defer Trace("applyMigrations(db)")()

	goose.SetBaseFS(migrationFiles)
	err := goose.SetDialect("sqlite3")

	FatalOnError(err, "Failed to set goose dialect")
	FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")
}
