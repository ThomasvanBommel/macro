package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed frontend/build/*
var frontendFS embed.FS

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	buildFS, err := fs.Sub(frontendFS, "frontend/build")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/register" && r.Method == "POST" {
			handleRegisterUser(w, r, db)
			return
		}

		if strings.HasPrefix(r.URL.Path, "/api/") {
			handleAPI(w, r)
			return
		}

		if os.Getenv("DEVMODE") == "true" {
			http.NotFound(w, r)
			return
		}

		handleStatic(w, r, buildFS)
	})

	log.Println("macro backend is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
