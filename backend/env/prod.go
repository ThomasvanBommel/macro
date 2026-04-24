//go:build prod

package env

import (
	"embed"
	"io/fs"
	"log/slog"
	"macro/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const SECURE = true
const SESSION_TIMEOUT_SEC = 3600

//go:embed build
var frontendBuildFiles embed.FS

//go:embed migrations/*.sql
var MigrationFiles embed.FS

func Init(r *gin.Engine) {
	defer util.Trace("Init(r *gin.Engine)")()

	registerStaticAssetRoute(r)
	registerNoRoute(r)
}

// registerNoRoute sets up a catch-all route for the Gin engine
func registerNoRoute(r *gin.Engine) {
	defer util.Trace("registerNoRoute(r *gin.Engine)")()

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "API endpoint not found"})
			return
		}

		f, err := frontendBuildFiles.ReadFile("build/index.html")
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
	defer util.Trace("registerStaticAssetRoute(r *gin.Engine)")()

	f, err := fs.Sub(frontendBuildFiles, "build/assets")
	util.FatalOnError(err, "Failed to create sub filesystem for frontend assets")

	r.StaticFS("/assets", http.FS(f))
}
