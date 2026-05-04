package api2

import (
	"log/slog"
	"macro/db2"
	"macro/errs"

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
	auth.POST("/login", x.Login)
	auth.POST("/logout", x.Logout)
	auth.GET("/session", x.SessionInfo)

	food := api.Group("/food")
	food.POST("/new", x.NewFood)

	return r, x
}

func BindJSON(c *gin.Context, input any) error {
	if err := c.ShouldBindJSON(input); err != nil {
		return errs.BadInput(err, err.Error())
	}
	return nil
}
