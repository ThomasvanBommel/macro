package main

import (
	"log/slog"

	"github.com/gin-gonic/gin"
)

func main() {
	db := InitDatabase()
	defer db.Close()

	router := gin.Default()
	InitAPI(router, db)

	slog.Info("Starting server: http://localhost:8080/")
	router.Run()
}
