package main

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed frontend/build/*
var frontendFS embed.FS

func main() {
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
