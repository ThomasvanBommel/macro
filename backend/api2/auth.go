package api2

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

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

func (db *DB) Login(c *gin.Context) {
	var input AuthInput
	if err := BindJSON(c, &input); err != nil {
		ErrorResponse(c, err)
		return
	}

	err, token, expiry := db.CreateSession(input.Username, input.Password)
	if err != nil {
		ErrorResponse(c, err)
		return
	}

	session := sessions.Default(c)
	session.Set("username", input.Username)
	session.Set("token", token)
	session.Save()

	OK(c, gin.H{"username": input.Username, "expires": expiry})
}

func (db *DB) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: -1})
	session.Clear()
	session.Save()

	OK(c, nil)
}
