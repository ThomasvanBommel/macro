package main

import (
	"log/slog"
	"os"

	"macro/api"
	"macro/db"
	"macro/env"
	"macro/util"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	log := util.Logger{}
	r.Use(log.Middleware())

	db := db.NewDatabase()
	defer db.Close()

	api.Init(r, db)
	env.Init(r)

	// If running in a trusted environment (e.g. GCP Cloud Run)
	if os.Getenv("K_SERVICE") != "" {
		r.ForwardedByClientIP = true
		r.SetTrustedProxies(nil)
	} else {
		// In development, trust localhost and Docker bridge IPs
		r.SetTrustedProxies([]string{"127.0.0.1", "172.17.0.1"})
	}

	slog.Info("Starting server: http://localhost:8080/")
	r.Run()
}
