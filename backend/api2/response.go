package api2

import (
	"macro/errs"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, err error) {
	code := 500
	msg := "An unexpected error occurred"

	if e, ok := err.(*errs.Error); ok {
		c.Error(e.Unwrap())
		code = e.Code()
		msg = e.Error()
	} else {
		c.Error(err)
	}

	c.AbortWithStatusJSON(code, map[string]string{"error": msg})
}

func DataResponse(c *gin.Context, code int, data any) { c.JSON(code, data) }

func Created(c *gin.Context, data any) { DataResponse(c, 201, data) }
