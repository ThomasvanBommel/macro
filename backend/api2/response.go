package api2

import (
	"macro/errs"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, err error) {
	code := 500
	msg := "An unexpected error occurred"

	if err != nil {
		c.Error(err)
	}

	if e, ok := err.(*errs.Error); ok {
		code = int(e.Code())
		msg = e.Public()
	}

	c.AbortWithStatusJSON(code, map[string]string{"error": msg})
}

func DataResponse(c *gin.Context, code int, data any) { c.JSON(code, data) }

func OK(c *gin.Context, data any)      { DataResponse(c, 200, data) }
func Created(c *gin.Context, data any) { DataResponse(c, 201, data) }

func NotAuthorized(c *gin.Context, msg string) { ErrorResponse(c, errs.NotAuthorized(nil, msg)) }
func BadInput(c *gin.Context, msg string)      { ErrorResponse(c, errs.BadInput(nil, msg)) }

// func InternalError(c *gin.Context, err error) { ErrorResponse(c, errs.Internal(err)) }
