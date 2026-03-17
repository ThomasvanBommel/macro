package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

func handleRegisterUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}
	
	resp, err := registerUser(db, req)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func handleAPI(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "API endpoint not implemented", http.StatusNotImplemented)
}

func handleStatic(w http.ResponseWriter, r *http.Request, buildFS fs.FS) {
	file, err := buildFS.Open(strings.TrimPrefix(r.URL.Path, "/"))
	if err == nil {
		defer file.Close()
		stat, _ := file.Stat()
		http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file.(io.ReadSeeker))
		return
	}

	indexFile, err := buildFS.Open("index.html")
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	defer indexFile.Close()
	stat, _ := indexFile.Stat()
	http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
}
