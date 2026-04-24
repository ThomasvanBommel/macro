//go:build !prod

package env

import (
	"embed"

	"github.com/gin-gonic/gin"
)

const SECURE = false
const SESSION_TIMEOUT_SEC = 600

//go:embed migrations/*.sql
var MigrationFiles embed.FS

func Init(_ *gin.Engine) {
	// Do nothing
}
