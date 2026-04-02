//go:build dev

package main

import "github.com/gin-gonic/gin"

// InitStaticRoutes is a placeholder function for development builds that does not set up any routes
// for serving static files. In development mode, it is assumed that the frontend assets are served
// separately (e.g., by a development server like Vite or Webpack Dev Server), so this function does
// not perform any actions. It is defined to maintain a consistent interface with the production
// version of the function, allowing the main application code to call InitStaticRoutes without
// needing to check the build mode.
func InitStaticRoutes(router *gin.Engine) {
	// Do nothing
}
