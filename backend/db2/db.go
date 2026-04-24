package db2

import (
	"database/sql"
	"macro/env"
	"macro/util"

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

	return &Database{db}
}
