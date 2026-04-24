package api2

import (
	"log/slog"
	"macro/db2"
	"macro/errs"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DB struct{ *db2.Database }

func Init(db *db2.Database, level slog.Level) (*gin.Engine, *DB) {
	r := gin.New()
	r.Use(LoggingMiddleware(level))
	r.Use(SessionMiddleware())
	r.Use(gin.Recovery())

	x := &DB{db}

	api := r.Group("/api")
	api.GET("/health", func(c *gin.Context) { c.JSON(200, nil) })

	auth := api.Group("/auth")
	auth.POST("/register", x.Register)
	return r, x
}

func BindJSON(c *gin.Context, input any) error {
	if err := c.ShouldBindJSON(input); err != nil {
		return errs.New(http.StatusUnprocessableEntity, err.Error(), err)
	}
	return nil
}
