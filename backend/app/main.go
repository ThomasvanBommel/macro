package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

//go:embed frontend/build/*
var frontendFiles embed.FS

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database initialized successfully")

	assetFS, err := fs.Sub(frontendFiles, "frontend/build/assets")
	if err != nil {
		log.Fatal(err)
	}

	dbWrapper := &DatabaseWrapper{db}

	router := gin.Default()
	router.Use(func(c *gin.Context) { c.Set("db", dbWrapper); c.Next() })

	initAPI(router)

	if os.Getenv("DEVMODE") == "true" {
		router.NoRoute(func(c *gin.Context) {
			c.String(http.StatusTeapot, "This is dev mode, goof. Use the Vite server :5173")
		})
	} else {
		router.StaticFS("/assets", http.FS(assetFS))

		router.NoRoute(func(c *gin.Context) {
			// reject requests for paths starting with /api
			if strings.HasPrefix(c.Request.URL.Path, "/api") {
				c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
				return
			}

			file, err := frontendFiles.ReadFile("frontend/build/index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "Failed to load index.html")
				return
			}

			c.Data(http.StatusOK, "text/html; charset=utf-8", file)
		})
	}

	log.Println("macro backend is running on http://localhost:8080")
	router.Run()
}
