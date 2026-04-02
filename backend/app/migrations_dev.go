//go:build dev

package main

import "database/sql"

// ApplyMigrations is a no-op function in development builds. In development mode, it is assumed \
// that migrations are applied manually or are not needed. This function is defined to maintain a
// consistent interface with the production version of the ApplyMigrations function. It takes a
// *sql.DB parameter but does not use it, and simply returns without performing any actions.
func ApplyMigrations(db *sql.DB) {
	// Do nothing
}
