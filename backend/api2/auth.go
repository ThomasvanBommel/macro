package api2

import "github.com/gin-gonic/gin"

type AuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (db *DB) Register(c *gin.Context) {
	var input AuthInput
	if err := BindJSON(c, &input); err != nil {
		ErrorResponse(c, err)
		return
	}

	if err := db.CreateUser(input.Username, input.Password); err != nil {
		ErrorResponse(c, err)
		return
	}

	Created(c, nil)
}
