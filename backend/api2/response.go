package api2

import (
	"log/slog"
	"macro/errs"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, err error) {
	code := 500
	msg := "An unexpected error occurred"

	slog.Error("ErrorResponse", "error", err.Error())

	if err != nil {
		c.Error(err)
	}

	if e, ok := err.(*errs.Error); ok {
		// c.Error(e)
		code = e.Code()
		msg = e.Error()
	}

	c.AbortWithStatusJSON(code, map[string]string{"error": msg})
}

func DataResponse(c *gin.Context, code int, data any) { c.JSON(code, data) }

func OK(c *gin.Context, data any)      { DataResponse(c, 200, data) }
func Created(c *gin.Context, data any) { DataResponse(c, 201, data) }

func Unauthorized(c *gin.Context, msg string)  { ErrorResponse(c, errs.New(401, msg, nil)) }
func Unprocessable(c *gin.Context, msg string) { ErrorResponse(c, errs.New(422, msg, nil)) }

func InternalError(c *gin.Context, err error) { ErrorResponse(c, errs.Internal(err)) }
