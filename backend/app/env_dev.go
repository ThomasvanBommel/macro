//go:build !prod

package main

import (
	"github.com/gin-gonic/gin"
)

func InitEnv(_ *gin.Engine) {
	// Do nothing
}
