//go:build !dev

package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/build/*
var frontendBuildFiles embed.FS

// InitStaticRoutes sets up the routes for serving static files and handling undefined routes. It
// registers the necessary routes for serving static assets and handling requests to undefined
// routes. This function is called during the initialization of the server to ensure that static
// files are served correctly and that undefined routes are handled gracefully.
func InitStaticRoutes(r *gin.Engine) {
	defer Trace("InitStaticRoutes(r *gin.Engine)")()

	registerStaticAssetRoute(r)
	registerNoRoute(r)
}

// registerNoRoute sets up a catch-all route for the Gin engine that handles requests to undefined
// routes. It checks if the request path starts with "/api" and returns a 404 JSON response if it
// does, indicating that the API endpoint was not found. For all other paths, it serves the
// "index.html" file from the embedded frontend build files, allowing the frontend application to
// handle client-side routing for non-API requests. This ensures that API requests receive
// appropriate error responses, while still allowing the frontend to function correctly.
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
// It creates a sub-filesystem from the embedded frontend build files, and configures the Gin engine
// to serve files from this sub-filesystem at the "/assets" URL path. This allows the application to
// serve static assets like JavaScript, CSS, and image files. Panics if there is an error creating
// the sub-filesystem, ensuring that the application does not run with an invalid configuration.
func registerStaticAssetRoute(r *gin.Engine) {
	defer Trace("registerStaticAssetRoute(r *gin.Engine)")()

	f, err := fs.Sub(frontendBuildFiles, "frontend/build/assets")
	FatalOnError(err, "Failed to create sub filesystem for frontend assets")

	r.StaticFS("/assets", http.FS(f))
}
