//go:build prod

package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitEnv(r *gin.Engine) {
	defer Trace("InitEnv(db *sql.DB, r *gin.Engine)")()

	registerStaticAssetRoute(r)
	registerNoRoute(r)
}

//go:embed frontend/build
var frontendBuildFiles embed.FS

// registerNoRoute sets up a catch-all route for the Gin engine that handles requests to undefined
// routes.
func registerNoRoute(r *gin.Engine) {
	defer Trace("registerNoRoute(r *gin.Engine)")()

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		f, err := frontendBuildFiles.ReadFile("frontend/build/index.html")
		if err != nil {
			slog.Error("Failed to load index.html", "error", err)
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			return
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", f)
	})
}

// registerStaticAssetRoute sets up a route for serving static assets from the embedded filesystem.
func registerStaticAssetRoute(r *gin.Engine) {
	defer Trace("registerStaticAssetRoute(r *gin.Engine)")()

	f, err := fs.Sub(frontendBuildFiles, "frontend/build/assets")
	FatalOnError(err, "Failed to create sub filesystem for frontend assets")

	r.StaticFS("/assets", http.FS(f))
}
