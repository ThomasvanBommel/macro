//go:build !prod

package env

import (
	"github.com/gin-gonic/gin"
)

var Secure = false

func Init(_ *gin.Engine) {
	// Do nothing
}
