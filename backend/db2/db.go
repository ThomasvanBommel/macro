package db2

import (
	"context"
	"database/sql"
	"log/slog"
	"macro/env"
	"macro/util"
	"time"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type Database struct {
	*sql.DB
}

func Init(path string) *Database {
	db, err := sql.Open("sqlite", path)
	util.FatalOnError(err, "Failed to open database")
	db.SetMaxOpenConns(1)

	goose.SetBaseFS(env.MigrationFiles)
	util.FatalOnError(goose.SetDialect("sqlite3"), "Failed to set goose dialect")
	util.FatalOnError(goose.Up(db, "migrations"), "Failed to apply migrations")

	_, err = db.Exec(`
		PRAGMA cache_size = -10000;
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = OFF;
		PRAGMA foreign_keys = ON;
	`)
	util.FatalOnError(err, "Failed to set PRAGMA options")

	database := &Database{db}

	// cleanup sessions every 12h and on startup
	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		defer ticker.Stop()

		database.cleanUpSessions()
		for range ticker.C {
			database.cleanUpSessions()
		}
	}()

	return database
}

func (db *Database) cleanUpSessions() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	res, err := db.ExecContext(ctx, `DELETE FROM sessions WHERE expires < CURRENT_TIMESTAMP`)
	if err != nil {
		slog.Error("Failed to clear expired sessions", "error", err.Error())
		return
	}

	rows, _ := res.RowsAffected()
	slog.Info("Cleared expired sessions", "rows", rows)
}
