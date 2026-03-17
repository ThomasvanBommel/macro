package main

import (
	"database/sql"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

//go:embed frontend/build/*
var frontendFS embed.FS

//go:embed migrations/*.sql
var migrationsFS embed.FS

func init_db() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "/app/macro.db")
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	goose.SetBaseFS(migrationsFS)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	log.Println("Running database migrations...")
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := init_db()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	buildFS, err := fs.Sub(frontendFS, "frontend/build")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(buildFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.Error(w, "API endpoint not implemented", http.StatusNotImplemented)
			return
		}

		file, err := buildFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
		if err == nil {
			file.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		indexFile, err := buildFS.Open("index.html")
		defer indexFile.Close()
		stat, _ := indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
	})

	log.Println("macro backend is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
