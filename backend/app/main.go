package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	InitLogger()

	db := InitDatabase()
	defer db.Close()

	r := gin.Default()
	InitAPI(r, db)
	InitStaticRoutes(r)

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
