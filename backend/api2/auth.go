package api2

import (
	"fmt"
	"time"

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
	session.Set("expiry", expiry)
	session.Save()

	OK(c, gin.H{"username": input.Username, "expires": expiry})
}

func (db *DB) SessionInfo(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	expiry := session.Get("expiry")

	if username == nil || expiry == nil {
		Unauthorized(c, "No active session")
		return
	}

	ex, ok := expiry.(string)
	if !ok {
		InternalError(c, fmt.Errorf("Failed to cast session expiry to string: %v", expiry))
		return
	}

	t, err := time.Parse(time.RFC3339, ex)
	if err != nil {
		InternalError(c, fmt.Errorf("Invalid session expiry date format: %v", ex))
		return
	}

	if t.Before(time.Now().UTC()) {
		Unauthorized(c, "Session has expired")
		return
	}

	OK(c, gin.H{"username": username, "expires": expiry})
}

func (db *DB) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: -1})
	session.Clear()
	session.Save()

	OK(c, nil)
}
