package db2

import (
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
)

func newDatabase() *Database {
	goose.SetLogger(goose.NopLogger())
	return Init(":memory:?cache=shared")
}

func TestInit(t *testing.T) {
	db := newDatabase()
	defer db.Close()

	_, err := db.Exec("SELECT 1")
	require.NoError(t, err, "Failed to query database: %v", err)

	var fkEnabled int
	err = db.QueryRow("PRAGMA foreign_keys").Scan(&fkEnabled)
	require.NoError(t, err, "Failed to query foreign_keys PRAGMA: %v", err)
	require.Equal(t, 1, fkEnabled, "Expected foreign_keys PRAGMA to be 1, got %d", fkEnabled)

	tables := []string{"user", "session", "food", "meal", "entry"}
	for _, table := range tables {
		_, err := db.Exec("SELECT 1 FROM " + table + " LIMIT 1")
		require.NoError(t, err, "Failed to query %s table: %v", table, err)
	}
}
